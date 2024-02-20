package eventsourcing

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/gofrs/uuid"
	"github.com/straw-hat-team/onepiece/go/onepiece"
	"github.com/straw-hat-team/onepiece/go/onepiece/eventsourcing/onepiecemessage"
	"io"
)

type Result[Event any] struct {
	NextExpectedVersion uint64
	Events              []Event
}

type ContentType = esdb.ContentType

const (
	ContentTypeBinary = esdb.ContentTypeBinary
	ContentTypeJson   = esdb.ContentTypeJson
)

var (
	ErrOptimisticConcurrency = errors.New("optimistic concurrency error")
)

type Any = esdb.Any
type StreamExists = esdb.StreamExists
type NoStream = esdb.NoStream
type StreamRevision = esdb.StreamRevision
type ExpectedRevision = esdb.ExpectedRevision

type CorrelationId string
type CausationId string

// CommandHandler is a function that receives a command and returns a list of events.
// TODO: figure out how to avoid passing the db here. Take into consideration multi-tenancy.
type CommandHandler[Command any, Event any] func(context context.Context, db *esdb.Client, command Command, opts *Options) (*Result[Event], error)

type StreamId[Command any] func(command Command) (string, error)

type Metadata map[string]any

type Options struct {
	ExpectedRevision ExpectedRevision
	Metadata         Metadata
	CorrelationId    *CorrelationId
	CausationId      *CausationId
}

type UnmarshalEvent[Event any] func(eventType string, data []byte) (Event, error)
type MarshalEvent[Event any] func(event Event) (ContentType, []byte, error)

type GetEventType[Event any] func(event Event) (*onepiecemessage.MessageType, error)

var (
	maxReadSize = ^uint64(0)
)

func NewDecider[State any, Command any, Event any](
	decider *onepiece.Decider[State, Command, Event],
	getStreamId StreamId[Command],
	marshalEvent MarshalEvent[Event],
	unmarshalEvent UnmarshalEvent[Event],
	getEventType GetEventType[Event],
) CommandHandler[Command, Event] {
	return func(context context.Context, db *esdb.Client, command Command, opts *Options) (*Result[Event], error) {
		streamID, err := getStreamId(command)
		if err != nil {
			return nil, err
		}
		stream, err := db.ReadStream(context, streamID, esdb.ReadStreamOptions{
			Direction: esdb.Forwards,
			From:      esdb.Start{},
		}, maxReadSize)

		if err != nil {
			return nil, err
		}

		defer stream.Close()

		state := decider.InitialState()

		var lastResolvedEvent *esdb.ResolvedEvent

		for {
			resolvedEvent, err := stream.Recv()

			if err, ok := esdb.FromError(err); !ok {
				if err.Code() == esdb.ErrorCodeResourceNotFound {
					break
				} else if errors.Is(err, io.EOF) {
					break
				} else {
					return nil, err
				}
			}

			event, err := unmarshalEvent(
				resolvedEvent.Event.EventType,
				resolvedEvent.Event.Data,
			)
			if err != nil {
				return nil, err
			}

			state = decider.Evolve(state, event)
			lastResolvedEvent = resolvedEvent
		}

		if decider.IsTerminal(state) {
			return nil, onepiece.ErrTerminalState
		}

		events, err := decider.Decide(state, command)

		if err != nil {
			return nil, err
		}

		metadata, err := getEventMetadata(opts)
		if err != nil {
			return nil, err
		}

		eventData := make([]esdb.EventData, len(events))
		for i, event := range events {
			eventType, err := getEventType(event)
			if err != nil {
				return nil, err
			}

			contentType, data, err := marshalEvent(event)
			if err != nil {
				return nil, err
			}

			eventData[i] = esdb.EventData{
				EventType:   eventType.String(),
				ContentType: contentType,
				Data:        data,
				Metadata:    metadata,
			}
		}

		writeResult, err := db.AppendToStream(context, streamID, esdb.AppendToStreamOptions{
			ExpectedRevision: getExpectedRevision(opts, lastResolvedEvent),
		}, eventData...)

		if err, ok := esdb.FromError(err); !ok {
			if err.Code() == esdb.ErrorCodeWrongExpectedVersion {
				return nil, ErrOptimisticConcurrency
			}
			return nil, err
		}

		return &Result[Event]{
			Events:              events,
			NextExpectedVersion: writeResult.NextExpectedVersion,
		}, nil
	}
}

func getExpectedRevision(opts *Options, lastResolvedEvent *esdb.ResolvedEvent) ExpectedRevision {
	if opts != nil && opts.ExpectedRevision != nil {
		return opts.ExpectedRevision
	} else if lastResolvedEvent == nil {
		return NoStream{}
	} else {
		return Revision(lastResolvedEvent.OriginalEvent().EventNumber)
	}
}

func getEventMetadata(opts *Options) ([]byte, error) {
	var metadata map[string]any

	if opts != nil && opts.Metadata != nil {
		metadata = opts.Metadata
	} else {
		metadata = make(map[string]any)
	}

	metadata["$correlationId"] = getCorrelation(opts)
	metadata["$causationId"] = getCausationId(opts)

	bytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func getCausationId(opts *Options) *CausationId {
	if opts != nil && opts.CausationId != nil {
		return opts.CausationId
	} else {
		return NewCausationId()
	}
}

func getCorrelation(opts *Options) *CorrelationId {
	if opts != nil && opts.CorrelationId != nil {
		return opts.CorrelationId
	} else {
		return NewCorrelationId()
	}
}

func NewCausationId() *CausationId {
	id := CausationId(uuid.Must(uuid.NewV4()).String())
	return &id
}

func NewCorrelationId() *CorrelationId {
	id := CorrelationId(uuid.Must(uuid.NewV4()).String())
	return &id
}

func Revision(value uint64) StreamRevision {
	return esdb.Revision(value)
}
