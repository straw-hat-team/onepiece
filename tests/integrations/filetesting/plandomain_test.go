package filetesting_test

import (
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/onepiece/onepiecetesting"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planactor"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestCreatePLan(t *testing.T) {
	onepiecetesting.RunTestingFile(
		t,
		"testing.yaml",
		planactor.Decider,
		func(eventType string, payload yaml.Node) (*planproto.Command, error) {
			switch eventType {
			case "CreatePlan":
				var command planproto.CreatePlan
				err := payload.Decode(&command)
				if err != nil {
					return nil, err
				}
				return &planproto.Command{Command: &planproto.Command_CreatePlan{CreatePlan: &command}}, nil
			}
			return nil, onepiece.ErrUnknownCommand
		},
		func(eventType string, payload yaml.Node) (*planproto.Event, error) {
			switch eventType {
			case "PlanCreated":
				var event planproto.PlanCreated
				err := payload.Decode(&event)
				if err != nil {
					return nil, err
				}
				return &planproto.Event{Event: &planproto.Event_PlanCreated{PlanCreated: &event}}, nil
			}
			return nil, onepiece.ErrUnknownEvent
		},
	)
}
