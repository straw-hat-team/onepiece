use std::fmt::Debug;
use std::cmp::PartialEq;
use crate::decider::Decider;

pub enum SpecResult<Event, Error> {
  Event{ events: Vec<Event> },
  Error{ error: Error },
}

pub struct Spec<'a, Aggregate, Command, Event, Error> {
  decider: Decider<Aggregate, Command, Event, Error>,
  given: Vec<&'a Event>,
  when: Option<&'a Command>
}

impl<'a, Aggregate, Command, Event, Error> Spec<'a, Aggregate, Command, Event, Error>
  where
    Aggregate: PartialEq + Debug,
    Event: PartialEq + Debug,
    Error: PartialEq + Debug,
{
  pub fn new(decider: Decider<Aggregate, Command, Event, Error>) -> Self {
    Spec {
      decider,
      given: Vec::new(),
      when: None,
    }
  }

  pub fn given(mut self, events: Vec<&'a Event>) -> Self {
    self.given = events;
    self
  }

  pub fn when(mut self, command: &'a Command) -> Self {
    self.when = Some(command);
    self
  }

  pub fn then(self, then: SpecResult<Event, Error>) {
    let when = self.when.expect("when is not set");
    let mut aggregate = self.decider.initial_state();

    for event in self.given {
      aggregate = self.decider.evolve(&aggregate, event);
    }

    let result = self.decider.decide(&aggregate, when);

    match then {
      SpecResult::Event { events } => {
        assert_eq!(result, Ok(events));
      },
      SpecResult::Error { error } => {
        assert_eq!(result, Err(error));
      },
    }
  }
}
