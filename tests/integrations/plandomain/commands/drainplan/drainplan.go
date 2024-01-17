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

type state struct {
	planId     *string
	isArchived bool
	isDrained  bool
}

func decide(state state, command *planproto.DrainPlan) ([]*planproto.Event, error) {
	if state.planId == nil {
		return nil, ErrPlanNotFound
	}
	if state.isArchived == false {
		return nil, ErrPlanUnarchived
	}
	if state.isDrained {
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

func evolve(state state, event *planproto.Event) state {
	switch e := event.Event.(type) {
	case *planproto.Event_PlanCreated:
		state.planId = &e.PlanCreated.PlanId
		return state
	case *planproto.Event_PlanArchived:
		state.isArchived = true
		return state
	case *planproto.Event_PlanDrained:
		state.isDrained = true
		return state
	default:
		return state
	}
}
