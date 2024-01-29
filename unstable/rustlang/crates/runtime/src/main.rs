use eventstore::Client;
use infra::wasmeventsourcing::WasmEventSourcingDecider;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let settings = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"
        .parse()
        .unwrap();
    let client = Client::new(settings).unwrap();

    let file_path = std::env::current_dir()
        .unwrap()
        .join("target/wasm32-unknown-unknown/debug/monitoring_wasm.wasm");

    let url = extism::Wasm::file(file_path.as_path());
    let manifest = extism::Manifest::new([url]);
    let mut plugin = extism::Plugin::new(&manifest, [], true).unwrap();

    let mut decider = WasmEventSourcingDecider::new(&mut plugin);

    let result = decider
        .dispatch_command(
            client,
            serde_json::json!({
              "CreateMonitoring": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
            }).to_string(),
            None,
        )
        .await?;

    println!("result: {:?}", result);

    let input = r#"
        {
          "CreateMonitoring": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
        }
      "#;
    let stream_id: String = plugin.call("stream_id", input).unwrap();
    println!("stream_id: {}", stream_id);
    let initial_state: String = plugin.call("initial_state", "").unwrap();
    println!("initial_state: {}", initial_state);

    let input = r#"
        {
          "state": {"id":null,"status":"Paused"},
          "event": {
            "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
          }
        }
      "#;
    let evolve: String = plugin.call("evolve", input).unwrap();
    println!("evolve: {}", evolve);

    let input = r#"
        {"id":null,"status":"Paused"}
      "#;
    let is_terminal: String = plugin.call("is_terminal", input).unwrap();
    println!("is_terminal: {}", is_terminal);

    let input = r#"
        {
            "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
          }
      "#;
    let event_type: String = plugin.call("event_type", input).unwrap();
    println!("event_type: {}", event_type);

    let marshal_event: String = plugin
        .call(
            "marshal_event",
            r#"
        {
            "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
          }
      "#,
        )
        .unwrap();
    println!("marshal_event: {}", marshal_event);

    let payload = serde_json::json!({
        "event_type": "MonitoringStarted",
        "payload": serde_json::json!({
          "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
        }).to_string()
    })
    .to_string();

    let unmarshal_event: String = plugin.call("unmarshal_event", payload).unwrap();
    println!("unmarshal_event: {}", unmarshal_event);

    let payload = serde_json::json!({
        "state": serde_json::json!({"id":null,"status":"Paused"}),
        "command": serde_json::json!({
          "CreateMonitoring": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
        })
    })
    .to_string();

    let decide: String = plugin.call("decide", payload).unwrap();
    println!("decide: {}", decide);

    Ok(())
}
