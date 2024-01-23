// wit_bindgen::generate!({
//     world: "host",
//     exports: {
//         world: Decider,
//     },
// });
//
// struct Decider;
//
// impl Guest for Decider {
//   fn run() {
//     print("Hello, world!");
//   }
// }

use extism_pdk::*;


struct Command(monitoring::Command);
struct State(monitoring::State);
struct Event(monitoring::Event);
struct Error(monitoring::Error);

//
// impl extism_pdk::ToBytes<'_> for State {
//   type Bytes = dyn AsRef<[u8]>;
//
//   fn to_bytes(&self) -> Result<Self::Bytes, Error> {
//     Ok(
//       serde_json::to_string(&self.0)
//         .map_err(|e| Error::msg(e.to_string()))?
//         .into_bytes(),
//     )
//   }
// }
//
// impl FromBytesOwned for Command {
//   fn from_bytes_owned(data: &[u8]) -> Result<Self, Error> {
//     serde_json::from_slice(data).map_err(|e| Error::msg(e.to_string()))
//   }
// }

#[plugin_fn]
pub fn stream_id(command: Command) -> FnResult<String> {
  let id = match command.0 {
    monitoring::Command::CreateMonitoring(command) => command.id,
    monitoring::Command::PauseMonitoring(command) => command.id,
    monitoring::Command::ResumeMonitoring(command) => command.id,
  };

  Ok(format!("monitoring:{}", id))
}

#[plugin_fn]
pub fn initial_state() -> FnResult<State> {
  let state = monitoring::initial_state();
  Ok(State(state))
}

#[plugin_fn]
pub fn evolve(state: State, event: Event) -> FnResult<State> {
  Ok(State(monitoring::evolve(&state.0, &event.0)))
}
