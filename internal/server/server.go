package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/grokkos/go-ledger/internal/handler"
)

type Server struct {
	router  *mux.Router
	server  *http.Server
	handler *handler.Handler
}

func NewServer(handler *handler.Handler) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		handler: handler,
	}

	s.setupRouter()
	s.setupServer()

	return s
}

func (s *Server) setupRouter() {
	// Register routes with API versioning
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/transactions", s.handler.RecordTransaction).Methods("POST")
	api.HandleFunc("/transactions", s.handler.GetTransactions).Methods("GET")
	api.HandleFunc("/balance", s.handler.GetBalance).Methods("GET")
}

func (s *Server) setupServer() {
	s.server = &http.Server{
		Addr:         ":8080",
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		fmt.Println("Server starting on port 8080...")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v\n", err)
		}
	}()

	return s.waitForShutdown()
}

func (s *Server) waitForShutdown() error {
	// Create channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-quit
	fmt.Println("\nServer is shutting down...")

	// Create deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	fmt.Println("Server exited properly")
	return nil
}
