package drainplan

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanNotFound = errors.New("plan not found")
var ErrPlanUnarchived = errors.New("plan must be archived")
var ErrPlanDrained = errors.New("plan already drained")

var Decider = onepiece.NewDecider(decide, evolve)

type State struct {
	PlanId     *string
	IsArchived bool
	IsDrained  bool
}

func decide(state State, command *planproto.DrainPlan) ([]*planproto.Event, error) {
	if state.PlanId == nil {
		return nil, ErrPlanNotFound
	}
	if state.IsArchived == false {
		return nil, ErrPlanUnarchived
	}
	if state.IsDrained {
		return nil, ErrPlanDrained
	}

	return []*planproto.Event{
		{
			Event: &planproto.Event_PlanDrained{
				PlanDrained: &planproto.PlanDrained{
					PlanId:     command.PlanId,
					TransferId: command.TransferId,
					DrainedAt:  command.DrainedAt,
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
	case *planproto.Event_PlanDrained:
		state.IsDrained = true
		return state
	default:
		return state
	}
}
