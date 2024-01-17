package createplan

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanExists = errors.New("plan already exists")

var Decider = onepiece.NewDecider(decide, evolve)

type State struct {
	PlanId *string
}

func decide(state State, command *planproto.CreatePlan) ([]*planproto.Event, error) {
	if state.PlanId != nil {
		return nil, ErrPlanExists
	}

	return []*planproto.Event{
		{
			Event: &planproto.Event_PlanCreated{
				PlanCreated: &planproto.PlanCreated{
					PlanId:           command.PlanId,
					Title:            command.Title,
					Color:            command.Color,
					GoalAmount:       command.GoalAmount,
					Description:      command.Description,
					Icon:             command.Icon,
					CreatedAt:        command.CreatedAt,
					DepositAccountId: command.DepositAccountId,
				},
			},
		},
	}, nil
}

func evolve(state State, event *planproto.Event) State {
	switch e := event.Event.(type) {
	case *planproto.Event_PlanCreated:
		state.PlanId = &e.PlanCreated.PlanId
		return state
	default:
		return state
	}
}
