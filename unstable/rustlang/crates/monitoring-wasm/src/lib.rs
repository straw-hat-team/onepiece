use extism_pdk::*;
use anyhow::anyhow;

#[plugin_fn]
pub fn stream_id(Json(command): Json<monitoring::Command>) -> FnResult<String> {
    let id = match command {
        monitoring::Command::CreateMonitoring(command) => command.id,
        monitoring::Command::PauseMonitoring(command) => command.id,
        monitoring::Command::ResumeMonitoring(command) => command.id,
    };

    Ok(format!("monitoring.{}", id))
}

#[plugin_fn]
pub fn initial_state(_input: ()) -> FnResult<Json<monitoring::State>> {
    Ok(Json(monitoring::initial_state()))
}

#[derive(Debug, serde::Deserialize)]
struct EvolveCommand {
    pub event: String,
    pub state: String,
}

#[plugin_fn]
pub fn evolve(Json(evolve_command): Json<EvolveCommand>) -> FnResult<Json<monitoring::State>> {
    let state: monitoring::State = serde_json::from_str(&evolve_command.state)?;
    let event: monitoring::Event = serde_json::from_str(&evolve_command.event)?;

    Ok(Json(monitoring::evolve(
        &state,
        &event,
    )))
}

#[plugin_fn]
pub fn is_terminal(Json(state): Json<monitoring::State>) -> FnResult<Json<bool>> {
    Ok(Json(monitoring::is_terminal(&state)))
}

#[plugin_fn]
pub fn event_type(Json(event): Json<monitoring::Event>) -> FnResult<String> {
    let event_type = match event {
        monitoring::Event::MonitoringStarted { .. } => "MonitoringStarted",
        monitoring::Event::MonitoringPaused { .. } => "MonitoringPaused",
        monitoring::Event::MonitoringResumed { .. } => "MonitoringResumed",
    };
    Ok(event_type.to_string())
}

#[plugin_fn]
pub fn marshal_event(Json(event): Json<monitoring::Event>) -> FnResult<Vec<u8>> {
    Ok(serde_json::to_vec(&event)?)
}

#[derive(Debug, serde::Deserialize)]
struct UnmarshalEventCommand {
    // pub event_type: String,
    pub payload: String,
}

#[plugin_fn]
pub fn unmarshal_event(
    Json(command): Json<UnmarshalEventCommand>,
) -> FnResult<Json<monitoring::Event>> {
    let u: monitoring::Event = serde_json::from_str(&command.payload)?;
    Ok(Json(u))
}

#[derive(Debug, serde::Deserialize)]
struct DecideCommand {
    pub state: String,
    pub command: String,
}

#[plugin_fn]
pub fn decide(Json(decide_command): Json<DecideCommand>) -> FnResult<Json<Vec<monitoring::Event>>> {
  let command: monitoring::Command = serde_json::from_str(&decide_command.command)?;
  let state: monitoring::State = serde_json::from_str(&decide_command.state)?;

  match monitoring::decide(&state, &command) {
    Ok(events) => {
      Ok(Json(events))
    }
    Err(err) => {
      Err(WithReturnCode::from(Error(err)))
    }
  }
}

struct Error(monitoring::Error);

impl Into<extism_pdk::Error> for Error {
  fn into(self) -> extism_pdk::Error {
    return anyhow!("Missing attribute: {:?}", self.0)
  }
}
