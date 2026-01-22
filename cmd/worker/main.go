package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudsweep/cloudsweep/internal/infrastructure/config"
	"github.com/cloudsweep/cloudsweep/internal/infrastructure/database"
	"github.com/cloudsweep/cloudsweep/internal/infrastructure/queue"
)

var version = "dev"

func main() {
	log.Printf("Starting CloudSweep Worker %s", version)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create worker server
	worker, err := queue.NewWorkerServer(cfg.Redis, db)
	if err != nil {
		log.Fatalf("Failed to create worker server: %v", err)
	}

	// Create task handlers
	mux := queue.NewServeMux(db)

	// Start worker in goroutine
	go func() {
		log.Println("Worker started, waiting for tasks...")
		if err := worker.Run(mux); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	worker.Shutdown()

	log.Println("Worker exited properly")
}
