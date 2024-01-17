package planactor

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

var ErrPlanExists = errors.New("plan already exists")
var ErrPlanNotFound = errors.New("plan not found")
var ErrPlanArchived = errors.New("plan already archived")
var ErrPlanUnarchived = errors.New("plan must be archived")
var ErrPlanDrained = errors.New("plan already drained")

var Decider = onepiece.NewDecider(decide, evolve)

type state struct {
	planId     *string
	isArchived bool
	isDrained  bool
}

func decide(state state, command *planproto.Command) ([]*planproto.Event, error) {
	switch c := command.Command.(type) {
	case *planproto.Command_CreatePlan:
		if state.planId != nil {
			return nil, ErrPlanExists
		}

		return []*planproto.Event{
			{
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
	case *planproto.Command_DrainPlan:
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
						PlanId:     c.DrainPlan.PlanId,
						TransferId: c.DrainPlan.TransferId,
						DrainedAt:  c.DrainPlan.DrainedAt,
					},
				},
			},
		}, nil
	case *planproto.Command_FailDrainPlan:
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
				Event: &planproto.Event_PlanDrainFailed{
					PlanDrainFailed: &planproto.PlanDrainFailed{
						PlanId:     c.FailDrainPlan.PlanId,
						TransferId: c.FailDrainPlan.TransferId,
						FailedAt:   c.FailDrainPlan.FailedAt,
					},
				},
			},
		}, nil

	default:
		return nil, onepiece.ErrUnknownCommand
	}
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
