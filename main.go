package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultPort        = "8080"
	DefaultWorkerCount = 3
	DefaultQueueSize   = 100
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	fmt.Println("===========================================")
	fmt.Println("    Go Concurrent Job Processor")
	fmt.Println("===========================================")
	fmt.Printf("  Workers : %d\n", DefaultWorkerCount)
	fmt.Printf("  Queue   : %d slots\n", DefaultQueueSize)
	fmt.Printf("  Port    : %s\n", port)
	fmt.Println("===========================================")

	queue := NewJobQueue(DefaultQueueSize)

	pool := NewWorkerPool(DefaultWorkerCount, queue)
	pool.Start()

	handlers := NewHandlers(queue)

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", LoggingMiddleware(handlers.JobsRouter))
	mux.HandleFunc("/jobs/", LoggingMiddleware(handlers.JobsRouter))
	mux.HandleFunc("/health", LoggingMiddleware(handlers.HealthCheck))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("\n[Server] Listening on http://localhost:%s\n", port)
		fmt.Println("[Server] Endpoints:")
		fmt.Println("  POST /jobs        - Submit a new job")
		fmt.Println("  GET  /jobs        - List all jobs")
		fmt.Println("  GET  /jobs/{id}   - Get job status")
		fmt.Println("  GET  /health      - Health check")
		fmt.Println()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] Failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n[Server] Received shutdown signal...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[Server] Forced shutdown: %v", err)
	}

	pool.Stop()

	fmt.Println("[Server] Graceful shutdown complete")
}
