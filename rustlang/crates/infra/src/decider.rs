use std::result::Result;
use std::fmt::Debug;
use std::cmp::PartialEq;

pub type InitialState<State> = fn() -> State;
pub type IsTerminal<State> = fn(state: &State) -> bool;
pub type Decide<State, Command, Event, Error> = fn(state: &State, command: Command) -> Result<Vec<Event>, Error>;
pub type Evolve<State, Event> = fn(state: &State, event: Event) -> State;

pub struct Decider<State, Command, Event, Error> {
  initial_state: InitialState<State>,
  decide: Decide<State, Command, Event, Error>,
  evolve: Evolve<State, Event>,
  is_terminal: IsTerminal<State>,
}

impl<State, Command, Event, Error> Decider<State, Command, Event, Error>
  where
    Event: PartialEq + Debug,
    Error: PartialEq + Debug,
    State: PartialEq + Debug,
{
  pub fn new(
    decide: Decide<State, Command, Event, Error>,
    evolve: Evolve<State, Event>,
    initial_state: InitialState<State>,
    is_terminal: Option<IsTerminal<State>>,
  ) -> Self {
    Decider {
      decide,
      evolve,
      initial_state,
      is_terminal: is_terminal.unwrap_or(never_terminal),
    }
  }

  pub fn initial_state(&self) -> State {
    (self.initial_state)()
  }

  pub fn decide(&self, state: &State, command: Command) -> Result<Vec<Event>, Error> {
    (self.decide)(state, command)
  }

  pub fn evolve(&self, state: &State, event: Event) -> State {
    (self.evolve)(state, event)
  }

  pub fn is_terminal(&self, state: &State) -> bool {
    (self.is_terminal)(state)
  }
}

fn never_terminal<State>(_state: &State) -> bool {
  false
}
