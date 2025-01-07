package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/akhilesharora/turnaround-collector/internal/collector"
)

func main() {
	logger := log.New(os.Stdout, "COLLECTOR: ", log.LstdFlags)

	// Get configuration
	cameraCount, err := strconv.Atoi(os.Getenv("CAMERA_COUNT"))
	if err != nil || cameraCount < 1 {
		cameraCount = 3
	}

	maxConcurrent, err := strconv.Atoi(os.Getenv("MAX_CONCURRENT"))
	if err != nil || maxConcurrent < 1 {
		maxConcurrent = (cameraCount + 1) / 2
	}

	config := collector.Config{
		CameraCount:   getEnvInt("CAMERA_COUNT", 3),
		PollInterval:  getEnvDuration("POLL_INTERVAL", 5*time.Second),
		MaxConcurrent: getEnvInt("MAX_CONCURRENT", 0),
		CameraBaseURL: getEnv("CAMERA_BASE_URL", "http://camera"),
		TargetURL:     getEnv("TARGET_URL", "http://target:8080/image"),
	}

	httpClient := collector.NewClient(
		5*time.Second,
		config.CameraBaseURL,
		config.TargetURL,
	)

	c := collector.NewCollector(
		config,
		httpClient,
		httpClient,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())

	// Start collector
	errCh := make(chan error, 1)
	go func() {
		errCh <- c.Start(ctx)
	}()

	// Wait for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Println("Shutting down collector...")
	cancel()

	// Wait for collector to finish with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled {
			logger.Printf("Error during shutdown: %v", err)
		}
	case <-shutdownCtx.Done():
		logger.Printf("Shutdown timed out")
	}
}

// Helper functions
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if str, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(str); err == nil {
			return value
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if str, exists := os.LookupEnv(key); exists {
		if value, err := time.ParseDuration(str); err == nil {
			return value
		}
	}
	return fallback
}
