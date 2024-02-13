use async_nats::jetstream;
use async_nats::jetstream::kv::Entry;
use async_nats::service::ServiceExt;
use eventstore::Client;
use futures::StreamExt;
use infra::wasmeventsourcing::WasmEventSourcingDecider;
use serde_json::{Map, Value};
use std::env;
use std::str::from_utf8;

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

    let kv = jetstream
        .create_key_value(async_nats::jetstream::kv::Config {
            bucket: "profiles".to_string(),
            max_value_size: 1024 * 1024 * 10, // 10mg, I think
            ..Default::default()
        })
        .await
        .unwrap();

    let file_path = std::env::current_dir()
        .unwrap()
        .join("target/wasm32-wasi/debug/monitoring_wasm.wasm");

    match std::fs::read(file_path) {
        Ok(contents) => {
            // convert the len to megabytes
            let megabytes = contents.len() / 1024 / 1024;
            println!("len: {:?} megabytes: {:?}", contents.len(), megabytes);
            let status = kv
                .put("sue.color", contents.into())
                .await
                .unwrap();
            println!("status: {:?}", status);
        }
        Err(e) => println!("Failed to read file: {}", e),
    }

    let entry = kv.entry("sue.color").await.unwrap();
    if let Some(entry) = entry {
        println!(
            "{} @ {} -> {}",
            entry.key,
            entry.revision,
            from_utf8(&entry.value).unwrap()
        );
    }

    kv.put("sue.color", "green".into()).await.unwrap();
    let entry = kv.entry("sue.color").await.unwrap();
    if let Some(entry) = entry {
        println!(
            "{} @ {} -> {}",
            entry.key,
            entry.revision,
            from_utf8(&entry.value).unwrap()
        );
    }

    kv.update("sue.color", "red".into(), 1)
        .await
        .expect_err("expected error");

    kv.update("sue.color", "red".into(), 2).await.unwrap();
    let entry = kv.entry("sue.color").await.unwrap();
    if let Some(entry) = entry {
        println!(
            "{} @ {} -> {}",
            entry.key,
            entry.revision,
            from_utf8(&entry.value).unwrap()
        );
    }

    let name = jetstream.stream_names().next().await.unwrap().unwrap();
    println!("KV stream name: {name}");

    // let nats_url =
    //       std::env::var("NATS_URL").unwrap_or_else(|_| "nats://localhost:4222".to_string());
    //
    //   let nats = async_nats::connect(nats_url).await.unwrap();
    //
    //   let jetstream = jetstream::new(nats.clone());
    //
    //   let kv = jetstream
    //       .create_key_value(async_nats::jetstream::kv::Config {
    //           bucket: "wasm-mods".to_string(),
    //           ..Default::default()
    //       })
    //       .await
    //       .unwrap();
    //
    //   match kv.entry("monitoring.create-monitoring").await.unwrap() {
    //       None => {
    //           let file_path = std::env::current_dir()
    //               .unwrap()
    //               .join("target/wasm32-wasi/debug/monitoring_wasm.wasm");
    //
    //           match std::fs::read(file_path) {
    //               Ok(contents) => {
    //                   let status = kv
    //                       .put("monitoring.create-monitoring", contents.into())
    //                       .await
    //                       .unwrap();
    //                   println!("status: {:?}", status);
    //               }
    //               Err(e) => println!("Failed to read file: {}", e),
    //           }
    //       }
    //       Some(_) => {
    //
    //       }
    //   }
}
