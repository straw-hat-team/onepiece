package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	golang "unstable"

	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, js := golang.NewNats()
	defer nc.Drain()

	streamName := "EVENTS"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	})

	js.Publish(ctx, "events.1", nil)
	js.Publish(ctx, "events.2", nil)
	js.Publish(ctx, "events.3", nil)

	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{})

	wg := sync.WaitGroup{}
	wg.Add(3)

	cc, _ := cons.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		fmt.Println("received msg on", msg.Subject())
		wg.Done()
	})
	wg.Wait()

	cc.Stop()

	js.Publish(ctx, "events.1", nil)
	js.Publish(ctx, "events.2", nil)
	js.Publish(ctx, "events.3", nil)

	msgs, _ := cons.Fetch(2)
	var i int
	for msg := range msgs.Messages() {

		msg.Ack()
		i++
	}
	fmt.Printf("got %d messages\n", i)

	msgs, _ = cons.FetchNoWait(100)
	i = 0
	for msg := range msgs.Messages() {
		msg.Ack()
		i++
	}
	fmt.Printf("got %d messages\n", i)

	fetchStart := time.Now()
	msgs, _ = cons.Fetch(1, jetstream.FetchMaxWait(time.Second))
	i = 0
	for msg := range msgs.Messages() {
		msg.Ack()
		i++
	}

	fmt.Printf("got %d messages in %v\n", i, time.Since(fetchStart))

	dur, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "processor",
	})

	msgs, _ = dur.Fetch(1)
	msg := <-msgs.Messages()
	fmt.Printf("received %q from durable consumer\n", msg.Subject())

	stream.DeleteConsumer(ctx, "processor")

	_, err := stream.Consumer(ctx, "processor")

	fmt.Println("consumer deleted:", errors.Is(err, jetstream.ErrConsumerNotFound))
}
