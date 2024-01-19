use bytes::Bytes;
use std::collections::HashMap;
use std::fmt::{Debug, Display};
use crate::decider::Decider;

type GetEventType<Event> = fn(&Event) -> String;
type GetStreamId<Command> = fn(&Command) -> String;
type UnmarshalEvent<Event> = fn(String, Bytes) -> Result<Event, Box<dyn std::error::Error>>;
type MarshalEvent<Event> = fn(&Event) -> Result<(String, Bytes), Box<dyn std::error::Error>>;

pub struct EventSourcingDecider<State, Command, Event, Error> {
  decider: Decider<State, Command, Event, Error>,
  get_event_type: GetEventType<Event>,
  get_stream_id: GetStreamId<Command>,
  unmarshal_event: UnmarshalEvent<Event>,
  marshal_event: MarshalEvent<Event>,
}

impl<State, Command, Event, Error> EventSourcingDecider<State, Command, Event, Error>
  where
    Event: PartialEq + Debug,
    Error: PartialEq + Debug,
    State: PartialEq + Debug,
{
  pub fn new(
    decider: Decider<State, Command, Event, Error>,
    get_event_type: GetEventType<Event>,
    get_stream_id: GetStreamId<Command>,
    unmarshal_event: UnmarshalEvent<Event>,
    marshal_event: MarshalEvent<Event>,
  ) -> Self {
    EventSourcingDecider {
      decider,
      get_event_type,
      get_stream_id,
      unmarshal_event,
      marshal_event,
    }
  }
  pub fn get_event_type(&self, event: &Event) -> String {
    (self.get_event_type)(event)
  }
  pub fn get_stream_id(&self, command: &Command) -> String {
    (self.get_stream_id)(command)
  }
  pub fn unmarshal_event(&self, event_type: String, data: Bytes) -> Result<Event, Box<dyn std::error::Error>> {
    (self.unmarshal_event)(event_type, data)
  }
  pub fn marshal_event(&self, event: &Event) -> Result<(String, Bytes), Box<dyn std::error::Error>> {
    (self.marshal_event)(event)
  }

  pub fn initial_state(&self) -> State {
    self.decider.initial_state()
  }

  pub fn decide(&self, state: &State, command: Command) -> Result<Vec<Event>, Error> {
    self.decider.decide(state, command)
  }

  pub fn evolve(&self, state: &State, event: Event) -> State {
    self.decider.evolve(state, event)
  }

  pub fn is_terminal(&self, state: &State) -> bool {
    self.decider.is_terminal(state)
  }
}


#[derive(Debug)]
pub struct DecisionResult<Event> {
  pub next_expected_version: u64,
  pub events: Vec<Event>,
}

#[derive(Debug)]
pub struct CorrelationId(String);

impl Display for CorrelationId {
  fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
    write!(f, "{}", self.0.clone())
  }
}

impl Default for CorrelationId {
  fn default() -> Self {
    CorrelationId(uuid::Uuid::new_v4().to_string())
  }
}

#[derive(Debug)]
pub struct CausationId(String);

impl Display for CausationId {
  fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
    write!(f, "{}", self.0.clone())
  }
}

impl Default for CausationId {
  fn default() -> Self {
    CausationId(uuid::Uuid::new_v4().to_string())
  }
}

pub type ExpectedRevision = eventstore::ExpectedRevision;

pub type Metadata = HashMap<String, String>;

#[derive(Debug)]
pub struct Options {
  pub metadata: Option<Metadata>,
  pub expected_revision: Option<ExpectedRevision>,
  pub correlation_id: Option<CorrelationId>,
  pub causation_id: Option<CausationId>,
}

impl Default for Options {
  fn default() -> Self {
    Options {
      expected_revision: None,
      metadata: None,
      correlation_id: None,
      causation_id: None,
    }
  }
}

pub enum MarshalError<Error> {
  UnmarshalEvent(Error),
  MarshalEvent(Error),
}

#[derive(Debug)]
pub enum CommandHandlerError<Error, MarshallingError> {
  StateIsTerminal,
  Domain(Error),
  MarshalError(MarshallingError),
  InvalidConnection,
  EventStore(eventstore::Error),
}
