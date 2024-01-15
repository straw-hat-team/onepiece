package plandomain

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
)

type Plan struct {
	PlanId *string
}

var ErrAlreadyExists = errors.New("plan already exists")

func Decide(state Plan, command *Command) ([]*Event, error) {
	switch c := command.Command.(type) {
	case *Command_CreatePlan:
		if state.PlanId != nil {
			return nil, ErrAlreadyExists
		}

		return []*Event{
			{
				Context: command.Context,
				Event: &Event_PlanCreated{
					PlanCreated: &PlanCreated{
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
	default:
		return nil, onepiece.ErrUnknownCommand
	}
}

func Evolve(state Plan, event *Event) Plan {
	switch event.Event.(type) {
	case *Event_PlanCreated:
		state.PlanId = &event.GetPlanCreated().PlanId
		return state
	default:
		return state
	}
}
