package planactor

import (
	"errors"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/archiveplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/createplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/drainplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/faildrainplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/updateplan"
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
		return createplan.Decider.Decide(createplan.State{
			PlanId: state.planId,
		}, c.CreatePlan)

	case *planproto.Command_ArchivePlan:
		return archiveplan.Decider.Decide(archiveplan.State{
			PlanId:     state.planId,
			IsArchived: state.isArchived,
		}, c.ArchivePlan)

	case *planproto.Command_UpdatePlan:
		return updateplan.Decider.Decide(updateplan.State{
			PlanId:     state.planId,
			IsArchived: state.isArchived,
		}, c.UpdatePlan)

	case *planproto.Command_DrainPlan:
		return drainplan.Decider.Decide(drainplan.State{
			PlanId:     state.planId,
			IsArchived: state.isArchived,
			IsDrained:  state.isDrained,
		}, c.DrainPlan)

	case *planproto.Command_FailDrainPlan:
		return faildrainplan.Decider.Decide(faildrainplan.State{
			PlanId:     state.planId,
			IsArchived: state.isArchived,
			IsDrained:  state.isDrained,
		}, c.FailDrainPlan)

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
