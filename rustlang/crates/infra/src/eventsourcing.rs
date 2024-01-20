use std::collections::HashMap;
use std::fmt::{Debug, Display};
use crate::decider::Decider;

type GetEventType<Event> = fn(&Event) -> &str;
type GetStreamId<Command> = fn(&Command) -> String;
type UnmarshalEvent<Event, MarshallingError> = fn(&str, bytes::Bytes) -> Result<Event, MarshalError<MarshallingError>>;
type MarshalEvent<Event, MarshallingError> = fn(&Event) -> Result<bytes::Bytes, MarshallingError>;

pub struct EventSourcingDecider<State, Command, Event, Error, MarshallingError> {
  decider: Decider<State, Command, Event, Error>,
  get_event_type: GetEventType<Event>,
  get_stream_id: GetStreamId<Command>,
  unmarshal_event: UnmarshalEvent<Event, MarshallingError>,
  marshal_event: MarshalEvent<Event, MarshallingError>,
}

impl<State, Command, Event, Error, MarshallingError> EventSourcingDecider<State, Command, Event, Error, MarshallingError>
  where
    State: PartialEq + Debug,
    Event: PartialEq + Debug,
    Error: PartialEq + Debug,
    MarshallingError: Debug,
{
  pub fn new(
    decider: Decider<State, Command, Event, Error>,
    get_event_type: GetEventType<Event>,
    get_stream_id: GetStreamId<Command>,
    unmarshal_event: UnmarshalEvent<Event, MarshallingError>,
    marshal_event: MarshalEvent<Event, MarshallingError>,
  ) -> Self {
    EventSourcingDecider {
      decider,
      get_event_type,
      get_stream_id,
      unmarshal_event,
      marshal_event,
    }
  }

  fn unmarshal_event(&self, event_type: &str, data: bytes::Bytes) -> Result<Event, MarshalError<MarshallingError>> {
    (self.unmarshal_event)(event_type, data)
  }
  fn marshal_event<'a>(&'a self, event: &'a Event) -> Result<bytes::Bytes, MarshalError<MarshallingError>> {
    (self.marshal_event)(event).map_err(MarshalError::MarshalEvent)
  }

  pub async fn dispatch_command(&self, client: eventstore::Client, command: &Command, opts: Option<Options>) -> Result<DecisionResult<Event>, CommandHandlerError<Error, MarshallingError>> {
    let stream_id = (self.get_stream_id)(command);
    let mut state = self.decider.initial_state();
    let mut stream = client
      .read_stream(stream_id.as_str(), &Default::default())
      .await.unwrap();

    let mut last_event_expected_version = None;

    loop {
      match stream.next().await {
        Ok(Some(event)) => {
          let resolved_event = event.get_original_event();

          let event = self.unmarshal_event(
            resolved_event.event_type.as_str(),
            resolved_event.data.clone(),
          ).unwrap();

          state = self.decider.evolve(&state, &event);
          last_event_expected_version = Some(eventstore::ExpectedRevision::Exact(resolved_event.revision));
        }
        Ok(None) => {
          break;
        }
        Err(eventstore::Error::ResourceNotFound) => {
          break;
        }
        Err(err) => {
          return Err(CommandHandlerError::EventStore(err));
        }
      }
    }

    if self.decider.is_terminal(&state) {
      return Err(CommandHandlerError::StateIsTerminal);
    }

    let events = self.decider.decide(&state, command).unwrap();
    let mut record_events: Vec<eventstore::EventData> = vec![];

    let opts = opts.unwrap_or_default();
    let mut metadata = opts.metadata.unwrap_or_default();

    metadata.insert("$correlationId".to_string(), opts.correlation_id.unwrap_or_default().to_string());
    metadata.insert("$causationId".to_string(), opts.causation_id.unwrap_or_default().to_string());

    for event in &events {
      let event_type = (self.get_event_type)(event);
      let data = self.marshal_event(event).unwrap();
      match eventstore::EventData::binary(event_type, data).metadata_as_json("{}") {
        Ok(record_event) => {
          record_events.push(record_event);
        }
        Err(err) => {
          return Err(CommandHandlerError::MarshalMetadata(err));
        }
      }
    }

    let expected_version = last_event_expected_version.unwrap_or(eventstore::ExpectedRevision::NoStream);

    let options = eventstore::AppendToStreamOptions::default().
      expected_revision(expected_version);

    let append_result = client.append_to_stream(stream_id.as_str(), &options, record_events).await.unwrap();

    Ok(DecisionResult {
      next_expected_version: append_result.next_expected_version,
      events,
    })
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

#[derive(Debug)]
pub enum MarshalError<Error> {
  UnmarshalEvent(Error),
  MarshalEvent(Error),
  MarshalMetadata(Error),
  UnknownEventType,
}

#[derive(Debug)]
pub enum CommandHandlerError<Error, MarshallingError> {
  StateIsTerminal,
  Domain(Error),
  MarshalError(MarshallingError),
  EventStore(eventstore::Error),
  MarshalMetadata(serde_json::Error),
}
