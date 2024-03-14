use async_nats::jetstream;
use async_nats::jetstream::kv::Entry;
use async_nats::service::ServiceExt;
use eventstore::Client;
use extism::ToBytes;
use futures::StreamExt;
use infra::wasmeventsourcing::{Options, WasmEventSourcingDecider};
use serde_json::{Map, Value};

#[derive(Debug)]
struct CommandHandler<'a> {
    service: &'a str,
    method: &'a str,
}

fn parse_command(input: &str) -> Option<CommandHandler> {
    let parts: Vec<&str> = input.split('.').collect();

    if parts.len() == 4 && parts[0] == "srv" && parts[1] == "command" {
        Some(CommandHandler {
            service: parts[2],
            method: parts[3],
        })
    } else {
        None
    }
}

#[derive(Debug, serde::Deserialize)]
struct ServiceCommand {
    pub metadata: Option<Map<String, serde_json::Value>>,
    pub payload: serde_json::Value,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let nats_url =
        std::env::var("NATS_URL").unwrap_or_else(|_| "nats://localhost:4222".to_string());

    let nats = async_nats::ConnectOptions::new()
      .name("decider-wasm")
      .connect(nats_url).await.unwrap();

    let _jetstream = jetstream::new(nats.clone());

    let settings = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"
        .parse()
        .unwrap();

    let client = Client::new(settings).unwrap();

    let service = nats
        .service_builder()
        .description("Event Sourcing WASM Decider")
        .stats_handler(|endpoint, _stats| serde_json::json!({ "endpoint": endpoint }))
        .start("decider-wasm", "0.0.1")
        .await
        .unwrap();

    let mut endpoint = service.endpoint("srv.command.*.*").await.unwrap();

    loop {
      if let Some(request) = endpoint.next().await {
        let command_handler = parse_command(request.message.subject.as_str()).unwrap();
        let service_command: ServiceCommand =
          serde_json::from_slice(&request.message.payload).unwrap();

        println!("{:?} ---> {:?}", command_handler, service_command);

        let file_path = std::env::current_dir()
          .unwrap()
          .join("target/wasm32-wasi/debug/monitoring_wasm.wasm");
        let url = extism::Wasm::file(file_path.as_path());

        let manifest = extism::Manifest::new([url])
          .with_timeout(std::time::Duration::from_secs(1))
          .with_memory_max(1024 * 1024 * 50);

        let mut plugin = extism::Plugin::new(&manifest, [], true).unwrap();

        // metadata.insert(
        //   "$correlationId".to_string(),
        //   opts.correlation_id.unwrap_or_default().to_string(),
        // );
        // metadata.insert(
        //   "$causationId".to_string(),
        //   opts.causation_id.unwrap_or_default().to_string(),
        // );

        let opts = Options::default();

        let result = WasmEventSourcingDecider::new(&mut plugin)
          .dispatch_command(client.clone(), service_command.payload.to_string(), Some(opts))
          .await
          .unwrap();

        println!("POG CRAZY");
        println!("result: {:?}", result);

        let resp = serde_json::json!({"next_expected_version": result.next_expected_version}).to_bytes().unwrap();
        request.respond(Ok(resp.into())).await.unwrap();
      }
    }
}

// $by_category example $ce-monitoring
// $by_event_type example $et-MonitoringStarted
// $by_correlation_id example $bc-<correlation id>
// $stream_by_category example $category-monitoring
// $streams example
//
// {
// "correlationIdProperty": "$myCorrelationId"
// }


// first
// -

// curl "https://4iykoi7jk2kp5hhd5irhbdprn40yxest.lambda-url.us-west-2.on.aws/?myCustomParameter=squirrel"
// $OP.REQ.SRV.<account id>.<service>.<method>
