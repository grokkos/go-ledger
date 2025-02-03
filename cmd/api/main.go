package main

import (
	"github.com/grokkos/go-ledger/internal/handler"
	"github.com/grokkos/go-ledger/internal/server"
	"github.com/grokkos/go-ledger/internal/service"
	"log"
)

func main() {
	// Initialize services and handlers
	ledgerService := service.NewLedgerService()
	ledgerHandler := handler.NewHandler(ledgerService)

	// Create and start server
	srv := server.NewServer(ledgerHandler)
	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
