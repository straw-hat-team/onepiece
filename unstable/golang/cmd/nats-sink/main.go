package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	golang "unstable"
)

const groupName = "nats-sink"
const streamName = "EVENT_STORE_DB"

func main() {
	client := golang.MustNewEventStore()
	nc, js := golang.NewNats()
	defer nc.Drain()

	_, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"eventstoredb.>"},
	})
	golang.Must(err)

	err = client.CreatePersistentSubscriptionToAll(context.Background(), groupName, esdb.PersistentAllSubscriptionOptions{
		// TODO: find a way to start from the last position of the stream reading nats jetsream messages
		StartFrom: esdb.Start{},
	})
	if err, ok := esdb.FromError(err); !ok {
		if esdb.ErrorCodeResourceAlreadyExists != err.Code() {
			golang.Must(err)
		}
	} else {
		golang.Must(err)
	}

	sub, err := client.SubscribeToPersistentSubscriptionToAll(context.Background(), groupName, esdb.SubscribeToPersistentSubscriptionOptions{})
	golang.Must(err)

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			ev := event.EventAppeared.Event.OriginalEvent()
			subject := fmt.Sprintf("eventstoredb.%s.%s.%s", ev.StreamID, ev.EventType, ev.EventID.String())
			fmt.Println(subject)

			marshal, err := json.Marshal(ev)
			golang.Must(err)

			msg := &nats.Msg{
				Subject: subject,
				Header:  nil,
				Data:    marshal,
			}
			_, err = js.PublishMsg(context.Background(), msg)
			if err == nil {
				sub.Ack(event.EventAppeared.Event)
				golang.Must(err)
			} else {
				err := sub.Nack(err.Error(), esdb.NackActionRetry, event.EventAppeared.Event)
				golang.Must(err)
			}
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
}
