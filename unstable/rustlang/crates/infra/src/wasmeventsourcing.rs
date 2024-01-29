// use std::collections::HashMap;
// use std::fmt::{Debug, Display};
//
// pub struct WasmEventSourcingDecider {
//     plugin: extism::Plugin,
// }
//
// impl WasmEventSourcingDecider {
//     pub fn new(plugin: extism::Plugin) -> Self {
//         WasmEventSourcingDecider { plugin }
//     }
//
//     pub async fn dispatch_command(
//         &self,
//         client: eventstore::Client,
//         command: Vec<u8>,
//         opts: Option<Options>,
//     ) -> Result<DecisionResult, ()> {
//         let stream_id: String = self.plugin.clone().call("stream_id", command)?;
//         let mut state: String = self.plugin.clone().call("initial_state", "")?;
//         let mut stream = client
//             .read_stream(stream_id.as_str(), &Default::default())
//             .await
//             .unwrap();
//
//         let mut last_event_expected_version = None;
//
//         loop {
//             match stream.next().await {
//                 Ok(Some(event)) => {
//                     let resolved_event = event.get_original_event();
//
//                     let event: String = self.plugin.clone().call(
//                         "unmarshal_event",
//                         serde_json::json!({
//                             "event_type": resolved_event.event_type.as_str(),
//                             "payload": resolved_event.data.clone()
//                         })
//                         .to_string(),
//                     )?;
//
//                     state = self
//                         .plugin
//                         .clone()
//                         .call(
//                             "evolve",
//                             serde_json::json!({
//                               "state": state,
//                               "event": event,
//                             })
//                             .to_string()
//                             .as_str(),
//                         )?;
//                     last_event_expected_version =
//                         Some(eventstore::ExpectedRevision::Exact(resolved_event.revision));
//                 }
//                 Ok(None) => {
//                     break;
//                 }
//                 Err(eventstore::Error::ResourceNotFound) => {
//                     break;
//                 }
//                 Err(err) => {
//                     return Err(CommandHandlerError::EventStore(err));
//                 }
//             }
//         }
//
//         let is_terminal: bool = self.plugin.clone().call("is_terminal", state.as_str())?;
//
//         if is_terminal {
//             return Err(CommandHandlerError::StateIsTerminal);
//         }
//
//         let events: Vec<String> = self.plugin.clone().call("decide", state.as_str(), command.clone())?;
//         let mut record_events: Vec<eventstore::EventData> = vec![];
//
//         let opts = opts.unwrap_or_default();
//         let mut metadata = opts.metadata.unwrap_or_default();
//
//         metadata.insert(
//             "$correlationId".to_string(),
//             opts.correlation_id.unwrap_or_default().to_string(),
//         );
//         metadata.insert(
//             "$causationId".to_string(),
//             opts.causation_id.unwrap_or_default().to_string(),
//         );
//
//         for event in &events {
//             let event_type: String = self.plugin.clone().call("event_type", event)?;
//             let data: bytes::Bytes = self.plugin.clone().call("marshal_event", event)?;
//             match eventstore::EventData::binary(event_type, data).metadata_as_json("{}") {
//                 Ok(record_event) => {
//                     record_events.push(record_event);
//                 }
//                 Err(err) => {
//                     return Err(CommandHandlerError::MarshalMetadata(err));
//                 }
//             }
//         }
//
//         let expected_version =
//             last_event_expected_version.unwrap_or(eventstore::ExpectedRevision::NoStream);
//
//         let options =
//             eventstore::AppendToStreamOptions::default().expected_revision(expected_version);
//
//         let append_result = client
//             .append_to_stream(stream_id.as_str(), &options, record_events)
//             .await
//             .unwrap();
//
//         Ok(DecisionResult {
//             next_expected_version: append_result.next_expected_version,
//             events,
//         })
//     }
// }
//
// #[derive(Debug)]
// pub struct DecisionResult {
//     pub next_expected_version: u64,
//     pub events: Vec<String>,
// }
//
// #[derive(Debug)]
// pub struct CorrelationId(String);
//
// impl Display for CorrelationId {
//     fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
//         write!(f, "{}", self.0.clone())
//     }
// }
//
// impl Default for CorrelationId {
//     fn default() -> Self {
//         CorrelationId(uuid::Uuid::new_v4().to_string())
//     }
// }
//
// #[derive(Debug)]
// pub struct CausationId(String);
//
// impl Display for CausationId {
//     fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
//         write!(f, "{}", self.0.clone())
//     }
// }
//
// impl Default for CausationId {
//     fn default() -> Self {
//         CausationId(uuid::Uuid::new_v4().to_string())
//     }
// }
//
// pub type ExpectedRevision = eventstore::ExpectedRevision;
//
// pub type Metadata = HashMap<String, String>;
//
// #[derive(Debug)]
// pub struct Options {
//     pub metadata: Option<Metadata>,
//     pub expected_revision: Option<ExpectedRevision>,
//     pub correlation_id: Option<CorrelationId>,
//     pub causation_id: Option<CausationId>,
// }
//
// impl Default for Options {
//     fn default() -> Self {
//         Options {
//             expected_revision: None,
//             metadata: None,
//             correlation_id: None,
//             causation_id: None,
//         }
//     }
// }
//
// #[derive(Debug)]
// pub enum MarshalError<Error> {
//     UnmarshalEvent(Error),
//     MarshalEvent(Error),
//     MarshalMetadata(Error),
//     UnknownEventType,
// }
//
// #[derive(Debug)]
// pub enum CommandHandlerError<Error, MarshallingError> {
//     StateIsTerminal,
//     Domain(Error),
//     MarshalError(MarshallingError),
//     EventStore(eventstore::Error),
//     MarshalMetadata(serde_json::Error),
// }
