package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orchestra-mcp/plugin-storage-sqlite/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

func main() {
	workspace := flag.String("workspace", ".", "Root workspace directory")

	storage := internal.NewStoragePlugin(*workspace)

	p := plugin.New("storage.sqlite").
		Version("1.0.0").
		Description("SQLite storage with typed tables, CAS versioning, and markdown dual-write").
		Author("Orchestra").
		Binary("storage-sqlite").
		ProvidesStorage("sqlite").
		SetStorageHandler(storage).
		BuildWithTools()

	// Parse standard plugin flags (--orchestrator-addr, --listen-addr, --certs-dir).
	p.ParseFlags()

	// Re-read workspace after flag.Parse has been called.
	storage = internal.NewStoragePlugin(*workspace)
	p.Server().SetStorageHandler(storage)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("storage.sqlite: %v", err)
	}
}
