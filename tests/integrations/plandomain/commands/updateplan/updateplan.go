package updateplan

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanNotFound = errors.New("plan not found")
var ErrPlanArchived = errors.New("plan already archived")

var Decider = onepiece.NewDecider(decide, evolve)

type state struct {
	planId     *string
	isArchived bool
}

func decide(state state, command *planproto.UpdatePlan) ([]*planproto.Event, error) {
	if state.planId == nil {
		return nil, ErrPlanNotFound
	}
	if state.isArchived {
		return nil, ErrPlanArchived
	}

	return []*planproto.Event{
		{
			Event: &planproto.Event_PlanUpdated{
				PlanUpdated: &planproto.PlanUpdated{
					PlanId:      command.PlanId,
					Title:       command.Title,
					Color:       command.Color,
					GoalAmount:  command.GoalAmount,
					Description: command.Description,
					Icon:        command.Icon,
					UpdatedAt:   command.UpdatedAt,
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
	case *planproto.Event_PlanArchived:
		state.isArchived = true
		return state
	default:
		return state
	}
}
