package planinfra

import (
	"encoding/json"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
)

const domain = "plan"
const version = "v1"

var SendCommand = onepiece.NewEventSourcingDecider(
	plandomain.Decider,
	streamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

func streamID(command *planproto.Command) (string, error) {
	switch c := command.Command.(type) {
	case *planproto.Command_CreatePlan:
		return onepiece.StreamID(domain, c.CreatePlan.PlanId), nil
	default:
		return "", onepiece.ErrUnknownCommand
	}
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
