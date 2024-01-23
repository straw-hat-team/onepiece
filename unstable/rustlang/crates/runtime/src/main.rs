fn main() {
  let file_path = std::env::current_dir()
    .unwrap()
    .join("target/wasm32-unknown-unknown/debug/monitoring_wasm.wasm");

  let url = extism::Wasm::file(file_path.as_path());
  let manifest = extism::Manifest::new([url]);
  let mut plugin = extism::Plugin::new(&manifest, [], true).unwrap();

  let input = r#"{"id": "1", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}"#;
  let stream_id = plugin.call::<&str, &str>("stream_id", input).unwrap();
  let state = plugin.call::<&str, &str>("initial_state", "remove_me").unwrap();
  println!("{}", state);
}
