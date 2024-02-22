use async_nats::jetstream;
use futures::StreamExt;
use serde_json::{Map, Value};
use std::env;
use std::str::from_utf8;
use tokio::io::AsyncReadExt;

#[derive(Debug, serde::Deserialize)]
struct ServiceCommand {
    pub metadata: Option<Map<String, serde_json::Value>>,
    pub payload: serde_json::Value,
}

#[tokio::main]
async fn main() {
    let nats_url = env::var("NATS_URL").unwrap_or_else(|_| "nats://localhost:4222".to_string());
    let client = async_nats::connect(nats_url).await.unwrap();
    let jetstream = jetstream::new(client);

    let store = match jetstream.get_object_store("profiles").await {
        Ok(store) => store,
        Err(_) => {
            jetstream
                .create_object_store(async_nats::jetstream::object_store::Config {
                    bucket: "profiles".to_string(),
                    ..Default::default()
                })
                .await.unwrap()
        }
    };

    let binding = std::env::current_dir()
        .unwrap()
        .join("target/wasm32-wasi/debug/monitoring_wasm.wasm");
    let file_path = binding.as_path();

    let file_key = "monitoring.create-monitoring";
    let mut file = tokio::fs::File::open(file_path).await.unwrap();
    store.put(file_key, &mut file).await.unwrap();

    let mut entry = store.get("monitoring.create-monitoring").await.unwrap();
    let mut data = Vec::new();
    let r = entry.read_to_end(&mut data);
    println!("{:?} ---> {:?}", entry.info.name, data);
}
