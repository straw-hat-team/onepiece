package main

import (
	"context"
	"fmt"
	golang "unstable"
)

func main() {
	nc, js := golang.NewNats()
	es, err := Create(context.Background(), "eventstore", nc, js, nil)
	golang.Must(err)
	result, err := es.AppendToStream(
		context.Background(),
		"plans.1",
		[]*ProposedMessage{
			{
				Data:        []byte(`{"name":"plan 1","plan_id":"1"}`),
				ContentType: "application/json",
				EventType:   "plan.created",
			},
		},
	)
	golang.Must(err)
	fmt.Sprintf("%v", result)

	events, err := es.ReadStream(
		context.Background(),
		"plans.1",
	)
	golang.Must(err)
	fmt.Sprintf("%v", events)
}
