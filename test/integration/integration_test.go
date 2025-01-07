package integration

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/akhilesharora/turnaround-collector/internal/camera"
	"github.com/akhilesharora/turnaround-collector/internal/collector"
	"github.com/akhilesharora/turnaround-collector/internal/target"
)

func TestIntegration(t *testing.T) {
	cameraLogger := log.New(os.Stdout, "CAMERA: ", log.LstdFlags)
	targetLogger := log.New(os.Stdout, "TARGET: ", log.LstdFlags)
	collectorLogger := log.New(os.Stdout, "COLLECTOR: ", log.LstdFlags)

	// Start test servers
	cameraServer := camera.NewServer(cameraLogger)
	targetServer := target.NewServer(targetLogger, nil)

	camera := httptest.NewServer(cameraServer)
	defer camera.Close()

	target := httptest.NewServer(targetServer)
	defer target.Close()

	// Configuring collector
	config := collector.Config{
		CameraCount:   1,
		PollInterval:  100 * time.Millisecond, // Short interval for quick test
		MaxConcurrent: 1,
		CameraBaseURL: camera.URL,
		TargetURL:     target.URL + "/image",
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
		collectorLogger,
	)

	// Run collector with 2 second context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.Start(ctx)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled {
			t.Errorf("collector failed: %v", err)
		}
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			// This is expected, as the test will cancel after a short duration
			t.Log("Test completed within expected time")
		} else {
			t.Error("test timed out")
		}
	}
}
