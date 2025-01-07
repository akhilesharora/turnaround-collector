package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akhilesharora/turnaround-collector/internal/camera"
)

func main() {
	logger := log.New(os.Stdout, "CAMERA: ", log.LstdFlags)
	server := camera.NewServer(logger)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	// Start camera server
	go func() {
		logger.Printf("Starting camera server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("Error during shutdown: %v", err)
	}
}
