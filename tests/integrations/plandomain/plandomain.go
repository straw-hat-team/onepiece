package plandomain

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanExists = errors.New("plan already exists")
var ErrPlanNotFound = errors.New("plan not found")
var ErrPlanArchived = errors.New("plan already archived")

var Decider = onepiece.NewDecider(decide, evolve)

type state struct {
	planId     *string
	isArchived bool
}

func decide(state state, command *planproto.Command) ([]*planproto.Event, error) {
	switch c := command.Command.(type) {
	case *planproto.Command_CreatePlan:
		if state.planId != nil {
			return nil, ErrPlanExists
		}

		return []*planproto.Event{
			{
				Context: command.Context,
				Event: &planproto.Event_PlanCreated{
					PlanCreated: &planproto.PlanCreated{
						PlanId:           c.CreatePlan.PlanId,
						Title:            c.CreatePlan.Title,
						Color:            c.CreatePlan.Color,
						GoalAmount:       c.CreatePlan.GoalAmount,
						Description:      c.CreatePlan.Description,
						Icon:             c.CreatePlan.Icon,
						CreatedAt:        c.CreatePlan.CreatedAt,
						DepositAccountId: c.CreatePlan.DepositAccountId,
					},
				},
			},
		}, nil

	case *planproto.Command_ArchivePlan:
		if state.planId == nil {
			return nil, ErrPlanNotFound
		}
		if state.isArchived {
			return nil, ErrPlanArchived
		}

		return []*planproto.Event{
			{
				Context: command.Context,
				Event: &planproto.Event_PlanArchived{
					PlanArchived: &planproto.PlanArchived{
						PlanId:     c.ArchivePlan.PlanId,
						ArchivedBy: c.ArchivePlan.ArchivedBy,
						ArchivedAt: c.ArchivePlan.ArchivedAt,
					},
				},
			},
		}, nil

	case *planproto.Command_UpdatePlan:
		if state.planId == nil {
			return nil, ErrPlanNotFound
		}
		if state.isArchived {
			return nil, ErrPlanArchived
		}

		return []*planproto.Event{
			{
				Context: command.Context,
				Event: &planproto.Event_PlanUpdated{
					PlanUpdated: &planproto.PlanUpdated{
						PlanId:      c.UpdatePlan.PlanId,
						Title:       c.UpdatePlan.Title,
						Color:       c.UpdatePlan.Color,
						GoalAmount:  c.UpdatePlan.GoalAmount,
						Description: c.UpdatePlan.Description,
						Icon:        c.UpdatePlan.Icon,
						UpdatedAt:   c.UpdatePlan.UpdatedAt,
					},
				},
			},
		}, nil
	default:
		return nil, onepiece.ErrUnknownCommand
	}
}

func evolve(state state, event *planproto.Event) state {
	switch event.Event.(type) {
	case *planproto.Event_PlanCreated:
		state.planId = &event.GetPlanCreated().PlanId
		return state
	case *planproto.Event_PlanArchived:
		state.isArchived = true
		return state
	default:
		return state
	}
}
