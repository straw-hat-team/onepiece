package onepiece

type Decide[State any, Command any, Event any] func(state State, command Command) ([]Event, error)

type Evolve[State any, Event any] func(state State, event Event) State
type InitialState[State any] func() State

type IsTerminal[State any] func(state State) bool

type Decider[State any, Command any, Event any] struct {
	decide       Decide[State, Command, Event]
	evolve       Evolve[State, Event]
	initialState InitialState[State]
	isTerminal   IsTerminal[State]
}

func (d *Decider[State, Command, Event]) Decide(state State, command Command) ([]Event, error) {
	return d.decide(state, command)
}

func (d *Decider[State, Command, Event]) Evolve(state State, event Event) State {
	return d.evolve(state, event)
}

func (d *Decider[State, Command, Event]) InitialState() State {
	return d.initialState()
}

func (d *Decider[State, Command, Event]) IsTerminal(state State) bool {
	return d.isTerminal(state)
}

func NewDecider[State any, Command any, Event any](
	decide Decide[State, Command, Event],
	evolve Evolve[State, Event],
) *Decider[State, Command, Event] {
	decider := &Decider[State, Command, Event]{
		decide:       decide,
		evolve:       evolve,
		initialState: EmptyInitialState[State],
		isTerminal:   NeverTerminal[State],
	}

	return decider
}

func (d *Decider[State, Command, Event]) WithInitialState(initialState InitialState[State]) *Decider[State, Command, Event] {
	d.initialState = initialState
	return d
}

func (d *Decider[State, Command, Event]) WithIsTerminal(isTerminal IsTerminal[State]) *Decider[State, Command, Event] {
	d.isTerminal = isTerminal
	return d
}

func EmptyInitialState[State any]() State {
	var state State
	return state
}

func NeverTerminal[State any](_state State) bool {
	return false
}
