package main

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/straw-hat-team/onepiece/go/onepiece/eventsourcing"
	golang "unstable"
	"unstable/plandomain/planproto"
	"unstable/planinfra"
)

func main() {
	eventStore := golang.MustNewEventStore()
	planID := uuid.Must(uuid.NewV4()).String()
	command := &planproto.CreatePlan{
		PlanId: planID,
		Title:  "Vacation",
		Color:  "#FF0000",
		GoalAmount: &planproto.Amount{
			Amount:       1000,
			Denomination: "USD",
		},
		Description:      "Plan for a vacation",
		Icon:             "https://some-url.com/icon.png",
		CreatedAt:        nil,
		DepositAccountId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
	}

	result, err := planinfra.DispatchCommand(
		context.Background(),
		eventStore,
		&planproto.Command{
			Command: &planproto.Command_CreatePlan{CreatePlan: command},
		},
		&eventsourcing.Options{
			ExpectedRevision: eventsourcing.Any{},
			Metadata:         nil,
			CorrelationId:    eventsourcing.NewCorrelationId(),
			CausationId:      eventsourcing.NewCausationId(),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", result)
}
