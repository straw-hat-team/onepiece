use bytes::Bytes;
use std::error::Error;
use infra::decider::Decider;
use eventstore::{AppendToStreamOptions, Client, EventData};
use infra::eventsourcing::{DecisionResult, Options, EventSourcingDecider, CommandHandlerError};

fn get_stream_id(command: &monitoring::Command) -> String {
  match command {
    monitoring::Command::CreateMonitoring(command) => format!("monitoring:{}", command.id),
    monitoring::Command::PauseMonitoring(command) => format!("monitoring:{}", command.id),
    monitoring::Command::ResumeMonitoring(command) => format!("monitoring:{}", command.id)
  }
}

fn marshal_event(event: &monitoring::Event) -> Result<(String, Bytes), Box<dyn Error>> {
  match event {
    monitoring::Event::MonitoringStarted { .. } => {
      let event_type = "MonitoringStarted".to_string();
      Ok((event_type, Bytes::from(serde_json::to_vec(&event)?)))
    }
    monitoring::Event::MonitoringPaused { .. } => {
      let event_type = "MonitoringPaused".to_string();
      Ok((event_type, Bytes::from(serde_json::to_vec(&event)?)))
    }
    monitoring::Event::MonitoringResumed { .. } => {
      let event_type = "MonitoringResumed".to_string();
      Ok((event_type, Bytes::from(serde_json::to_vec(&event)?)))
    }
  }
}

fn unmarshal_event(event_type: String, data: Bytes) -> Result<monitoring::Event, Box<dyn Error>> {
  match event_type.as_str() {
    "MonitoringStarted" => {
      let event: monitoring::Event = serde_json::from_slice(&data)?;
      Ok(event)
    }
    "MonitoringPaused" => {
      let event: monitoring::Event = serde_json::from_slice(&data)?;
      Ok(event)
    }
    "MonitoringResumed" => {
      let event: monitoring::Event = serde_json::from_slice(&data)?;
      Ok(event)
    }
    _ => Err("unknown event type".into()),
  }
}

fn event_type(event: &monitoring::Event) -> String {
  match event {
    monitoring::Event::MonitoringStarted { .. } => "MonitoringStarted".to_string(),
    monitoring::Event::MonitoringPaused { .. } => "MonitoringPaused".to_string(),
    monitoring::Event::MonitoringResumed { .. } => "MonitoringResumed".to_string(),
  }
}

async fn run(opts: Option<Options>) -> Result<DecisionResult<monitoring::Event>, CommandHandlerError<monitoring::Error, serde_json::Error>> {
  let settings = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000".parse().unwrap();
  let client = Client::new(settings).unwrap();
  let command = monitoring::Command::CreateMonitoring(monitoring::CreateMonitoring {
    id: uuid::Uuid::new_v4().to_string(),
    url: "https://www.google.com".to_string(),
  });

  let command_handler = EventSourcingDecider::new(
    Decider::new(
      monitoring::decide,
      monitoring::evolve,
      monitoring::initial_state,
      Some(monitoring::is_terminal),
    ),
    event_type,
    get_stream_id,
    unmarshal_event,
    marshal_event,
  );

  let stream_id = command_handler.get_stream_id(&command);


  let mut state = command_handler.initial_state();
  let mut stream = client
    .read_stream(stream_id.clone(), &Default::default())
    .await.unwrap();

  let mut last_event_expected_version = None;

  loop {
    match stream.next().await {
      Ok(Some(event)) => {
        let resolved_event = event.get_original_event();

        let event = command_handler.unmarshal_event(
          resolved_event.event_type.clone(),
          resolved_event.data.clone(),
        ).unwrap();

        state = command_handler.evolve(&state, event);
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

  if command_handler.is_terminal(&state) {
    return Err(CommandHandlerError::StateIsTerminal);
  }

  let events = command_handler.decide(&state, command).unwrap();
  let mut record_events: Vec<EventData> = vec![];

  let opts = opts.unwrap_or_default();
  let mut metadata = opts.metadata.unwrap_or_default();

  metadata.insert("$correlationId".to_string(), opts.correlation_id.unwrap_or_default().to_string());
  metadata.insert("$causationId".to_string(), opts.causation_id.unwrap_or_default().to_string());

  for event in &events {
    let (event_type, data) = command_handler.marshal_event(event).unwrap();
    match EventData::binary(event_type, data).metadata_as_json("{}") {
      Ok(record_event) => {
        record_events.push(record_event);
      }
      Err(err) => {
        return Err(CommandHandlerError::MarshalError(err));
      }
    }

  }

  let expected_version = last_event_expected_version.unwrap_or(eventstore::ExpectedRevision::NoStream);

  let options = AppendToStreamOptions::default().
    expected_revision(expected_version);

  let append_result = client.append_to_stream(stream_id.clone(), &options, record_events).await.unwrap();

  Ok(DecisionResult {
    next_expected_version: append_result.next_expected_version,
    events,
  })
}

#[tokio::main]
async fn main() -> Result<(), CommandHandlerError<monitoring::Error, serde_json::Error>> {
  let result = run(None).await?;
  println!("{:?}", result);
  Ok(())
}
