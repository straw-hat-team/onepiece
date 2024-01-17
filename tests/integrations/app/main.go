package main

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/gofrs/uuid"
	"github.com/straw-hat-team/onepiece-go/onepiece/eventsourcing"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/planinfra"
)

func main() {
	err, db := newEsDb()

	if err != nil {
		panic(err)
	}

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
		db,
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

func newEsDb() (error, *esdb.Client) {
	settings, err := esdb.ParseConnectionString("esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000")

	if err != nil {
		panic(err)
	}

	db, err := esdb.NewClient(settings)
	return err, db
}
