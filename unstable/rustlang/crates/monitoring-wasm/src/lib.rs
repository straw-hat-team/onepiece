use extism_pdk::*;
use anyhow::anyhow;

#[plugin_fn]
pub fn stream_id(Json(command): Json<monitoring::Command>) -> FnResult<String> {
    let id = match command {
        monitoring::Command::CreateMonitoring(command) => command.id,
        monitoring::Command::PauseMonitoring(command) => command.id,
        monitoring::Command::ResumeMonitoring(command) => command.id,
    };

    Ok(format!("monitoring:{}", id))
}

#[plugin_fn]
pub fn initial_state(_input: ()) -> FnResult<Json<monitoring::State>> {
    Ok(Json(monitoring::initial_state()))
}

#[derive(serde::Deserialize)]
struct EvolveCommand {
    #[serde(with = "serde_json::json")]
    pub event: monitoring::Event,
    #[serde(with = "serde_json::json")]
    pub state: monitoring::State,
}

#[plugin_fn]
pub fn evolve(Json(event_state): Json<EvolveCommand>) -> FnResult<Json<monitoring::State>> {
    Ok(Json(monitoring::evolve(
        &event_state.state,
        &event_state.event,
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
    pub event_type: String,
    pub payload: String,
}

#[plugin_fn]
pub fn unmarshal_event(
    Json(command): Json<UnmarshalEventCommand>,
) -> FnResult<Json<monitoring::Event>> {
    debug!("unmarshal_event: {:?}", command);
    let u: monitoring::Event = serde_json::from_str(&command.payload.as_str())?;
    println!("{:#?}", u);

    Ok(Json(u))
}

// // Custom deserialization for the Command enum
// fn deserialize_command<'de, D>(deserializer: D) -> Result<monitoring::Command, D::Error>
//   where
//     D: serde::Deserializer<'de>,
// {
//   let s = String::deserialize(deserializer)?;
//   serde_json::from_str(&s).map_err(serde::de::Error::custom)
// }

#[derive(serde::Deserialize)]
struct DecideCommand {
    #[serde(with = "serde_json::json")]
    // #[serde(deserialize_with = "deserialize_command")]
    pub state: monitoring::State,
    #[serde(with = "serde_json::json")]
    pub command: monitoring::Command,

  // pub state: Json<monitoring::State>,
    // pub command: Json<monitoring::Command>,
}

#[plugin_fn]
pub fn decide(Json(command): Json<DecideCommand>) -> FnResult<Json<Vec<monitoring::Event>>> {
  match monitoring::decide(&command.state, &command.command) {
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
