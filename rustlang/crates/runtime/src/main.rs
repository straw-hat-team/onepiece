use bytes::Bytes;
use infra::decider::Decider;
use eventstore::{Client};
use infra::eventsourcing::{DecisionResult, Options, EventSourcingDecider, CommandHandlerError, MarshalError};

fn get_stream_id(command: &monitoring::Command) -> String {
  match command {
    monitoring::Command::CreateMonitoring(command) => format!("monitoring:{}", command.id),
    monitoring::Command::PauseMonitoring(command) => format!("monitoring:{}", command.id),
    monitoring::Command::ResumeMonitoring(command) => format!("monitoring:{}", command.id),
  }
}

fn marshal_event(event: &monitoring::Event) -> Result<Bytes, serde_json::Error> {
  serde_json::to_vec(&event).map(Bytes::from)
}

fn unmarshal_event(event_type: &str, data: Bytes) -> Result<monitoring::Event, MarshalError<serde_json::Error>> {
  match event_type {
    "MonitoringStarted" | "MonitoringPaused" | "MonitoringResumed" => {
      return match serde_json::from_slice(&data) {
        Ok(event) => Ok(event),
        Err(err) => Err(MarshalError::UnmarshalEvent(err)),
      };
    }
    _ => Err(MarshalError::UnknownEventType),
  }
}

fn event_type(event: &monitoring::Event) -> &str {
  match event {
    monitoring::Event::MonitoringStarted { .. } => "MonitoringStarted",
    monitoring::Event::MonitoringPaused { .. } => "MonitoringPaused",
    monitoring::Event::MonitoringResumed { .. } => "MonitoringResumed",
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

  return command_handler.dispatch_command(client, &command, opts).await;
}

#[tokio::main]
async fn main() -> Result<(), CommandHandlerError<monitoring::Error, serde_json::Error>> {
  let result = run(None).await?;
  println!("{:?}", result);
  Ok(())
}
