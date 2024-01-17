package planinfra

import (
	"encoding/json"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/archiveplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/createplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/drainplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/commands/faildrainplan"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planactor"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

const domain = "plan"
const version = "v1"

var DispatchCommand = onepiece.NewEventSourcingDecider(
	planactor.Decider,
	streamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

var DispatchCreatePlan = onepiece.NewEventSourcingDecider(
	createplan.Decider,
	createPlanStreamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

var DispatchArchivePlan = onepiece.NewEventSourcingDecider(
	archiveplan.Decider,
	archivePlanStreamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

var DispatchDrainPlan = onepiece.NewEventSourcingDecider(
	drainplan.Decider,
	drainPlanStreamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

var DispatchFailDrainPlan = onepiece.NewEventSourcingDecider(
	faildrainplan.Decider,
	failPlanStreamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

func streamID(command *planproto.Command) (string, error) {
	switch c := command.Command.(type) {
	case *planproto.Command_CreatePlan:
		return createPlanStreamID(c.CreatePlan)
	case *planproto.Command_ArchivePlan:
		return archivePlanStreamID(c.ArchivePlan)
	case *planproto.Command_DrainPlan:
		return drainPlanStreamID(c.DrainPlan)
	case *planproto.Command_FailDrainPlan:
		return failPlanStreamID(c.FailDrainPlan)
	default:
		return "", onepiece.ErrUnknownCommand
	}
}

func failPlanStreamID(command *planproto.FailDrainPlan) (string, error) {
	return onepiece.StreamID(domain, command.PlanId), nil
}

func drainPlanStreamID(command *planproto.DrainPlan) (string, error) {
	return onepiece.StreamID(domain, command.PlanId), nil
}

func archivePlanStreamID(command *planproto.ArchivePlan) (string, error) {
	return onepiece.StreamID(domain, command.PlanId), nil
}

func createPlanStreamID(command *planproto.CreatePlan) (string, error) {
	return onepiece.StreamID(domain, command.PlanId), nil
}

func unmarshalEvent(eventType string, data []byte) (*planproto.Event, error) {
	event := &planproto.Event{}
	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func marshalEvent(event *planproto.Event) (onepiece.ContentType, []byte, error) {
	bytes, err := json.Marshal(event)
	return onepiece.ContentTypeJson, bytes, err
}

func eventTypeProvider(event *planproto.Event) (string, error) {
	switch event.Event.(type) {
	case *planproto.Event_PlanCreated:
		return onepiece.EventType(domain, version, "plan-created"), nil
	default:
		return "", onepiece.ErrUnknownEvent
	}
}
