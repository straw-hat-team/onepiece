use extism_pdk::*;


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
    pub event: monitoring::Event,
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
    pub payload: String
}

#[plugin_fn]
pub fn unmarshal_event(Json(command): Json<UnmarshalEventCommand>) -> FnResult<Json<monitoring::Event>> {
    debug!("unmarshal_event: {:?}", command);
    let u: monitoring::Event = serde_json::from_str(&command.payload.as_str())?;
    println!("{:#?}", u);

    Ok(Json(u))
}
