package main

import (
	"context"
	"errors"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/straw-hat-team/onepiece/go/onepiece/eventsourcing"
	"google.golang.org/protobuf/encoding/protojson"
	"strconv"
	golang "unstable"
	"unstable/plandomain/planproto"
	"unstable/planinfra"
)

type GetPlanLimitReached func(ctx context.Context, depositAccountId string) (bool, error)
type GenerateId func() string

type HandlerOptions struct {
	EventStore          *esdb.Client
	GenerateId          GenerateId
	GetPlanLimitReached GetPlanLimitReached
}

func NewHandler(o HandlerOptions) func(command *planproto.CreatePlan) (*golang.CommandHandlerResponse, error) {
	return func(command *planproto.CreatePlan) (*golang.CommandHandlerResponse, error) {
		// NOTE: this could be the side effect.
		// I said could be because what makes it a side effect is depending upon
		// the runtime environment dependency injection.
		limitReached, err := o.GetPlanLimitReached(context.Background(), command.DepositAccountId)
		if err != nil {
			return nil, err
		}
		if limitReached {
			return nil, errors.New("plan Limit Reached")
		}

		// NOTE: Just to make a point, context depending on if you want to use the existing
		// id or generate a new one. Or if you want to use the exact same Command message contract.
		if command.PlanId == "" {
			// NOTE: this could be the side effect.
			// I said could be because what makes it a side effect is depending upon
			// the runtime environment dependency injection.
			command.PlanId = o.GenerateId()
		}

		result, err := planinfra.DispatchCreatePlan(
			context.Background(),
			o.EventStore,
			command,
			&eventsourcing.Options{
				ExpectedRevision: eventsourcing.NoStream{},
				Metadata:         nil,
				CorrelationId:    eventsourcing.NewCorrelationId(),
				CausationId:      eventsourcing.NewCausationId(),
			},
		)
		if err != nil {
			return nil, err
		}

		return &golang.CommandHandlerResponse{
			NextExpectedVersion: result.NextExpectedVersion,
		}, nil
	}
}

func main() {
	ctx := context.Background()
	nc, js := golang.NewNats()
	eventStore := golang.MustNewEventStore()
	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "plan-counts"})
	golang.Must(err)
	getPlanLimitReachedService, err := NewGetPlanLimitReached(ctx, kv, 50)
	golang.Must(err)

	// NOTE: Ignore how the service works, this is just to make a point.
	_, err = golang.NewService[*planproto.CreatePlan](
		nc,
		"PlanCreator",
		"create-plan",
		unmarshalCommand,
		NewHandler(
			HandlerOptions{
				EventStore:          eventStore,
				GenerateId:          generateUuid,
				GetPlanLimitReached: getPlanLimitReachedService,
			},
		),
	)
	golang.Must(err)

	<-ctx.Done()
}

func NewGetPlanLimitReached(ctx context.Context, kv jetstream.KeyValue, limit uint) (GetPlanLimitReached, error) {
	if limit == 0 || limit > 50 {
		return nil, errors.New("invalid limit")
	}

	return func(context context.Context, depositAccountId string) (bool, error) {
		entry, err := kv.Get(ctx, depositAccountId)
		if err != nil {
			if errors.Is(err, jetstream.ErrKeyNotFound) {
				return false, nil
			}
			return false, err
		}

		count, err := strconv.ParseUint(string(entry.Value()), 10, 32)
		if err != nil {
			return false, err
		}

		return count >= uint64(limit), nil
	}, nil
}

func generateUuid() string {
	return uuid.Must(uuid.NewV4()).String()
}

func unmarshalCommand(data []byte) (*planproto.CreatePlan, error) {
	command := &planproto.CreatePlan{}
	err := protojson.Unmarshal(data, command)
	if err != nil {
		return nil, err
	}
	return command, err
}
