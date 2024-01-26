use extism_pdk::*;

#[derive(serde::Deserialize)]
struct EventState {
    pub event: monitoring::Event,
    pub state: monitoring::State,
}

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

#[plugin_fn]
pub fn evolve(Json(event_state): Json<EventState>) -> FnResult<Json<monitoring::State>> {
    Ok(Json(monitoring::evolve(
        &event_state.state,
        &event_state.event,
    )))
}
