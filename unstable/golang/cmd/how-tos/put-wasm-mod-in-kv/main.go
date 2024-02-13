package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"os"
	golang "unstable"
)

func main() {
	ctx := context.Background()
	nc, js := golang.NewNats()
	defer nc.Drain()

	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "wasm-pepehands"})
	golang.Must(err)

	fileBytes, err := os.ReadFile("/Users/ubi/Developer/github.com/straw-hat-team/onepiece/unstable/rustlang/target/wasm32-wasi/debug/monitoring_wasm.wasm")
	golang.Must(err)

	_, err = kv.Put(ctx, "sue.color", fileBytes)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}
	entry, _ := kv.Get(ctx, "sue.color")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))
}
