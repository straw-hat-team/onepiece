package onepiecetesting

import (
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

type Message struct {
	Type    string    `json:"type"`
	Payload yaml.Node `json:"payload"`
}

type Case struct {
	Given     []Message `json:"given" yaml:"given,omitempty"`
	When      Message   `json:"when" yaml:"when"`
	Then      []Message `json:"then" yaml:"then,omitempty"`
	Exception *Message  `json:"Exception,omitempty" yaml:"exception,omitempty"`
}

type UseCase struct {
	Description string `json:"description" yaml:"description"`
	Case        Case   `json:"case" yaml:"case"`
}

type TestingFile struct {
	UseCases []UseCase `json:"useCases" yaml:"useCases"`
}

func NewTestingFile(t *testing.T, fileName string) *TestingFile {
	file, err := os.ReadFile(fileName)
	require.NoError(t, err, "error reading file %s", fileName)

	var tf TestingFile
	err = yaml.Unmarshal(file, &tf)
	require.NoError(t, err, "error unmarshalling file %s", fileName)

	return &tf
}

type UnmarshalMessage[Message any] func(eventType string, payload yaml.Node) (Message, error)

func RunTestingFile[State any, Command any, Event any](
	t *testing.T,
	fileName string,
	decider *onepiece.Decider[State, Command, Event],
	unmarshalCommand UnmarshalMessage[Command],
	unmarshalEvent UnmarshalMessage[Event],
) {
	tf := NewTestingFile(t, fileName)

	for _, useCase := range tf.UseCases {
		t.Run(useCase.Description, func(t *testing.T) {
			givens := make([]Event, len(useCase.Case.Given))
			for i, message := range useCase.Case.Given {
				event, err := unmarshalEvent(message.Type, message.Payload)
				require.NoError(t, err, "error decoding payload")
				givens[i] = event
			}

			command, err := unmarshalCommand(useCase.Case.When.Type, useCase.Case.When.Payload)
			require.NoError(t, err, "error decoding payload")

			thens := make([]Event, len(useCase.Case.Then))
			for i, message := range useCase.Case.Then {
				event, err := unmarshalEvent(message.Type, message.Payload)
				require.NoError(t, err, "error decoding payload")
				thens[i] = event
			}

			NewTestCase(t, decider).
				Given(givens...).
				When(command).
				Then(thens...)
		})
	}
}
