package archiveplan

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

func decide(state state, command *planproto.ArchivePlan) ([]*planproto.Event, error) {
	if state.planId == nil {
		return nil, ErrPlanNotFound
	}
	if state.isArchived {
		return nil, ErrPlanArchived
	}

	return []*planproto.Event{
		{
			Event: &planproto.Event_PlanArchived{
				PlanArchived: &planproto.PlanArchived{
					PlanId:     command.PlanId,
					ArchivedBy: command.ArchivedBy,
					ArchivedAt: command.ArchivedAt,
				},
			},
		},
	}, nil
}

// WIT Variants
// Today Compose, Top level, two components, call fast, type switching between langs*
// Future, memcopy if same lang*,

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
