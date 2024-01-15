package planinfra

import (
	"encoding/json"
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain"
)

const domain = "plan"
const version = "v1"

var SendCommand = onepiece.NewEventSourcingDecider(
	onepiece.NewDecider(
		plandomain.Decide,
		plandomain.Evolve,
	),
	streamID,
	marshalEvent,
	unmarshalEvent,
	eventTypeProvider,
)

func streamID(command *plandomain.Command) (string, error) {
	switch c := command.Command.(type) {
	case *plandomain.Command_CreatePlan:
		return onepiece.StreamID(domain, c.CreatePlan.PlanId), nil
	default:
		return "", onepiece.ErrUnknownCommand
	}
}

func unmarshalEvent(eventType string, data []byte) (*plandomain.Event, error) {
	event := &plandomain.Event{}
	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func marshalEvent(event *plandomain.Event) (onepiece.ContentType, []byte, error) {
	bytes, err := json.Marshal(event)
	return onepiece.ContentTypeJson, bytes, err
}

func eventTypeProvider(event *plandomain.Event) (string, error) {
	switch event.Event.(type) {
	case *plandomain.Event_PlanCreated:
		return onepiece.EventType(domain, version, "plan-created"), nil
	default:
		return "", onepiece.ErrUnknownEvent
	}
}
