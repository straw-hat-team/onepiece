build:
	cargo build --target wasm32-wasi
	wasm-tools component new ../../target/wasm32-wasi/debug/monitoring_wasm.wasm \
      -o my-component.wasm --adapt ./wasi_snapshot_preview1.wasm
