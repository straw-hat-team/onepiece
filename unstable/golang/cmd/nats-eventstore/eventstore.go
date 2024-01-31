package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"strings"
	"time"
)

type appendOpts struct {
	expSeq ExpectedRevision
}

type appendOptFn func(o *appendOpts) error

func (f appendOptFn) appendOpt(o *appendOpts) error {
	return f(o)
}

// AppendOption is an option for the event store AppendToStream operation.
type AppendOption interface {
	appendOpt(o *appendOpts) error
}

// ExpectSequence indicates that the expected sequence of the subject sequence should
// be the value provided. If not, a conflict is indicated.
func ExpectSequence(seq ExpectedRevision) AppendOption {
	return appendOptFn(func(o *appendOpts) error {
		o.expSeq = seq
		return nil
	})
}

type ProposedMessage struct {
	ID       uuid.UUID
	Data     []byte
	Metadata []byte

	Created        time.Time
	SystemMetadata map[string]string
	ContentType    string
	EventType      string
}

func (p ProposedMessage) GetData() []byte {
	if p.Data == nil {
		return []byte("")
	}

	return p.Data
}

func (p ProposedMessage) GetMetadata() []byte {
	if p.Metadata == nil {
		return []byte("")
	}

	return p.Metadata
}

// Any means the write should not conflict with anything and should always succeed.
type Any struct{}

// StreamExists means the stream should exist.
type StreamExists struct{}

// NoStream means the stream being written to should not yet exist.
type NoStream struct{}

// ExpectedRevision the use of expected revision can be a bit tricky especially when discussing guaranties given by
// EventStoreDB server. The EventStoreDB server will assure idempotency for all requests using any value in
// ExpectedRevision except Any. When using Any, the EventStoreDB server will do its best to assure idempotency but
// will not guarantee it.
type ExpectedRevision interface {
	isExpectedRevision()
}

func (r Any) isExpectedRevision() {
}

func (r StreamExists) isExpectedRevision() {
}

func (r NoStream) isExpectedRevision() {
}

func (r StreamRevision) isExpectedRevision() {
}

// StreamRevision returns a stream position at a specific event revision.
type StreamRevision struct {
	Value uint64
}

// Revision returns a stream position at a specific event revision.
func Revision(value uint64) StreamRevision {
	return StreamRevision{
		Value: value,
	}
}

type EventStore struct {
	nc         *nats.Conn
	js         jetstream.JetStream
	streamName string
}

type RecordedEvent struct {
	EventID        string
	Data           []byte
	UserMetadata   []byte
	Created        time.Time
	SystemMetadata map[string]string
	ContentType    string
	EventType      string
	StreamID       string
	EventNumber    uint64
}

type MessagePayload struct {
	Data     []byte `json:"data"`
	Metadata []byte `json:"metadata"`
}

func (es *EventStore) streamSubject(streamID string) string {
	return fmt.Sprintf("%s.%s", es.streamName, streamID)
}

func (es *EventStore) ReadStream(ctx context.Context, streamID string) ([]*RecordedEvent, error) {
	cons, err := es.js.OrderedConsumer(ctx, es.streamName, jetstream.OrderedConsumerConfig{
		FilterSubjects: []string{es.streamSubject(streamID)},
		DeliverPolicy:  jetstream.DeliverAllPolicy,
	})
	if err != nil {
		return nil, err
	}

	expectedMsgCount := cons.CachedInfo().NumPending

	var events []*RecordedEvent
	for {
		msgs, err := cons.Fetch(500)
		if err != nil {
			return nil, err
		}

		for msg := range msgs.Messages() {
			event, err := fromNatsMsg(msg)
			if err != nil {
				return nil, err
			}

			events = append(events, event)

			err = msg.Ack()
			if err != nil {
				return nil, err
			}

			expectedMsgCount--
		}

		if expectedMsgCount == 0 {
			break
		}
	}

	return events, nil
}

func fromNatsMsg(msg jetstream.Msg) (*RecordedEvent, error) {
	headers := msg.Headers()
	md, err := msg.Metadata()
	if err != nil {
		return nil, err
	}
	eventTime, err := time.Parse(eventTimeFormat, headers.Get(eventCreatedHdr))
	if err != nil {
		return nil, fmt.Errorf("unpack: failed to parse event time: %s", err)
	}

	meta := make(map[string]string)
	for h := range headers {
		if strings.HasPrefix(h, eventMetaPrefixHdr) {
			key := h[len(eventMetaPrefixHdr):]
			meta[key] = headers.Get(h)
		}
	}

	data := &MessagePayload{}
	err = json.Unmarshal(msg.Data(), data)
	if err != nil {
		return nil, err
	}

	return &RecordedEvent{
		EventID:        headers.Get(nats.MsgIdHdr),
		EventType:      headers.Get(eventTypeHdr),
		ContentType:    headers.Get(eventContentTypeHdr),
		EventNumber:    md.Sequence.Stream,
		StreamID:       msg.Subject(),
		SystemMetadata: meta,
		Created:        eventTime,
		Data:           data.Data,
		UserMetadata:   data.Metadata,
	}, nil
}

type WriteResult struct {
	NextExpectedVersion uint64
}

func (es *EventStore) AppendToStream(
	ctx context.Context,
	streamID string,
	events []*ProposedMessage,
	opts ...AppendOption,
) (*WriteResult, error) {
	var o appendOpts
	for _, opt := range opts {
		if err := opt.appendOpt(&o); err != nil {
			return nil, err
		}
	}

	var ack *jetstream.PubAck

	for i, event := range events {
		popts := []jetstream.PublishOpt{
			jetstream.WithExpectStream(es.streamName),
		}

		if i == 0 && o.expSeq != nil {
			switch e := o.expSeq.(type) {
			case StreamRevision:
				popts = append(popts, jetstream.WithExpectLastSequencePerSubject(e.Value))
			case NoStream:
				popts = append(popts, jetstream.WithExpectLastSequencePerSubject(0))
			}
		}

		event.ID = uuid.Must(uuid.NewV4())
		event.Created = time.Now()

		msg, err := toNatsMsg(es.streamSubject(streamID), event)
		if err != nil {
			return nil, err
		}

		// TODO: Handle atomic write aka. Batch Publish
		ack, err = es.js.PublishMsg(ctx, msg, popts...)
		if err != nil {
			// TODO: this is a red-flag to be checking for a string in an error message
			if strings.Contains(err.Error(), "wrong last sequence") {
				return nil, ErrSequenceConflict
			}
			return nil, err
		}
	}

	// TODO: figure out how to handle this, ack.Sequence could be global instead of per subject
	return &WriteResult{
		NextExpectedVersion: ack.Sequence,
	}, nil
}

const (
	eventTypeHdr        = "onepiece-type"
	eventCreatedHdr     = "onepiece-created"
	eventContentTypeHdr = "onepiece-content-type"
	eventMetaPrefixHdr  = "onepiece-meta-"
	eventTimeFormat     = time.RFC3339Nano
)

var (
	ErrSequenceConflict = errors.New("onepiece: sequence conflict")
)

func toNatsMsg(streamID string, event *ProposedMessage) (*nats.Msg, error) {
	msg := nats.NewMsg(streamID)
	data, err := json.Marshal(&MessagePayload{Data: event.GetData(), Metadata: event.GetMetadata()})
	if err != nil {
		return nil, err
	}
	msg.Data = data
	msg.Header.Set(nats.MsgIdHdr, event.ID.String())
	msg.Header.Set(eventTypeHdr, event.EventType)
	msg.Header.Set(eventCreatedHdr, event.Created.Format(eventTimeFormat))
	msg.Header.Set(eventContentTypeHdr, event.ContentType)

	for k, v := range event.SystemMetadata {
		msg.Header.Set(fmt.Sprintf("%s%s", eventMetaPrefixHdr, k), v)
	}

	return msg, nil
}

func Create(ctx context.Context, name string, nc *nats.Conn, js jetstream.JetStream, config *jetstream.StreamConfig) (*EventStore, error) {
	es := &EventStore{
		streamName: name,
		nc:         nc,
		js:         js,
	}

	if config == nil {
		config = &jetstream.StreamConfig{}
	}
	config.Name = es.streamName

	if len(config.Subjects) == 0 {
		config.Subjects = []string{fmt.Sprintf("%s.>", es.streamName)}
	}

	_, err := es.js.CreateStream(ctx, *config)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func (es *EventStore) Update(ctx context.Context, config *jetstream.StreamConfig) error {
	if config == nil {
		config = &jetstream.StreamConfig{}
	}
	config.Name = es.streamName
	_, err := es.js.UpdateStream(ctx, *config)
	return err
}

func (es *EventStore) Delete(ctx context.Context) error {
	return es.js.DeleteStream(ctx, es.streamName)
}
