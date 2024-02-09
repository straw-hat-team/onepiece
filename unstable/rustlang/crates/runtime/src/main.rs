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
        .join("target/wasm32-wasi/debug/monitoring_wasm.wasm");
        // .join("target/wasm32-unknown-unknown/debug/monitoring_wasm.wasm");

    let url = extism::Wasm::file(file_path.as_path());
    let manifest = extism::Manifest::new([url]);
    let mut plugin = extism::Plugin::new(&manifest, [], true).unwrap();

    let mut decider = WasmEventSourcingDecider::new(&mut plugin);

    let result = decider
        .dispatch_command(
            client,
            serde_json::json!({
              "CreateMonitoring": {"id": "2", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}
            }).to_string(),
            None,
        )
        .await?;

    println!("result: {:?}", result);

    Ok(())
}
