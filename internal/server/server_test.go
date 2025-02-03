package server

import (
	"context"
	"github.com/grokkos/go-ledger/internal/handler"
	"github.com/grokkos/go-ledger/internal/service"
	"net/http"
	"testing"
	"time"
)

func TestServer_Start(t *testing.T) {
	// Create a test server
	ledgerService := service.NewLedgerService()
	ledgerHandler := handler.NewHandler(ledgerService)
	srv := NewServer(ledgerHandler)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Server.Start() error = %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Make a test request
	resp, err := http.Get("http://localhost:8080/api/v1/balance")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Shutdown server
	if err := srv.server.Shutdown(context.Background()); err != nil {
		t.Errorf("Server.Shutdown() error = %v", err)
	}
}
