package createplan

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanExists = errors.New("plan already exists")

var Decider = onepiece.NewDecider(decide, evolve)

type state struct {
	planId *string
}

func decide(state state, command *planproto.CreatePlan) ([]*planproto.Event, error) {
	if state.planId != nil {
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

func evolve(state state, event *planproto.Event) state {
	switch e := event.Event.(type) {
	case *planproto.Event_PlanCreated:
		state.planId = &e.PlanCreated.PlanId
		return state
	default:
		return state
	}
}
