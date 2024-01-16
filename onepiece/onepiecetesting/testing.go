package onepiecetesting

import (
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestCase[State any, Command any, Event any] struct {
	t *testing.T

	decider        *onepiece.Decider[State, Command, Event]
	previousEvents []Event
	command        Command
	expectedEvents []Event
	expectedError  error
}

func (tc *TestCase[State, Command, Event]) Given(events ...Event) *TestCase[State, Command, Event] {
	tc.previousEvents = append(tc.previousEvents, events...)
	return tc
}

func (tc *TestCase[State, Command, Event]) When(command Command) *TestCase[State, Command, Event] {
	tc.command = command
	return tc
}

func (tc *TestCase[State, Command, Event]) Then(event ...Event) *TestCase[State, Command, Event] {
	tc.expectedEvents = append(tc.expectedEvents, event...)
	return tc
}

func (tc *TestCase[State, Command, Event]) Catch(err error) *TestCase[State, Command, Event] {
	tc.expectedError = err
	return tc
}

func (tc *TestCase[State, Command, Event]) Assert() {
	state := tc.decider.InitialState()

	for _, event := range tc.previousEvents {
		state = tc.decider.Evolve(state, event)
	}

	if tc.decider.IsTerminal(state) {
		require.ErrorAs(tc.t, tc.expectedError, onepiece.ErrTerminalState)
	}

	events, err := tc.decider.Decide(state, tc.command)

	require.Equal(tc.t, tc.expectedEvents, events)
	require.Equal(tc.t, tc.expectedError, err)
}

func NewTestCase[State any, Command any, Event any](t *testing.T, decider *onepiece.Decider[State, Command, Event]) *TestCase[State, Command, Event] {
	require.NotNil(t, decider, "decider should not be nil")
	return &TestCase[State, Command, Event]{
		t:              t,
		previousEvents: []Event{},
		decider:        decider,
	}
}
