package updateplan

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanNotFound = errors.New("plan not found")
var ErrPlanArchived = errors.New("plan already archived")

var Decider = onepiece.NewDecider(decide, evolve)

type State struct {
	PlanId     *string
	IsArchived bool
}

func decide(state State, command *planproto.UpdatePlan) ([]*planproto.Event, error) {
	if state.PlanId == nil {
		return nil, ErrPlanNotFound
	}
	if state.IsArchived {
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

func evolve(state State, event *planproto.Event) State {
	switch e := event.Event.(type) {
	case *planproto.Event_PlanCreated:
		state.PlanId = &e.PlanCreated.PlanId
		return state
	case *planproto.Event_PlanArchived:
		state.IsArchived = true
		return state
	default:
		return state
	}
}
