fn main() {
    let file_path = std::env::current_dir()
        .unwrap()
        .join("target/wasm32-unknown-unknown/debug/monitoring_wasm.wasm");

    let url = extism::Wasm::file(file_path.as_path());
    let manifest = extism::Manifest::new([url]);
    let mut plugin = extism::Plugin::new(&manifest, [], true).unwrap();

    let input = r#"
        {
          "CreateMonitoring": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
        }
      "#;
    let stream_id: String = plugin.call("stream_id", input).unwrap();
    println!("{}", stream_id);
    let state: String = plugin.call("initial_state", "").unwrap();
    println!("{}", state);

    let input = r#"
        {
          "state": {"id":null,"status":"Paused"},
          "event": {
            "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
          }
        }
      "#;
    let state: String = plugin.call("evolve", input).unwrap();
    println!("{}", state);

    let input = r#"
        {"id":null,"status":"Paused"}
      "#;
    let is_terminal: String = plugin.call("is_terminal", input).unwrap();
    println!("{}", is_terminal);

    let input = r#"
        {
            "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
          }
      "#;
    let event_type: String = plugin.call("event_type", input).unwrap();
    println!("{}", event_type);

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
    println!("{}", marshal_event);

    let payload = serde_json::json!({
        "event_type": "MonitoringStarted",
        "payload": serde_json::json!({
          "MonitoringStarted": {"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
        }).to_string()
    })
    .to_string();

    let unmarshal_event: String = plugin.call("unmarshal_event", payload).unwrap();
    println!("{}", unmarshal_event);
}
