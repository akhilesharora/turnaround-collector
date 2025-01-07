package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/akhilesharora/turnaround-collector/pkg/interfaces"
)

type Config struct {
	CameraCount   int
	PollInterval  time.Duration
	MaxConcurrent int
	CameraBaseURL string
	TargetURL     string
}

type Collector struct {
	config  Config
	fetcher interfaces.ImageFetcher
	sender  interfaces.ImageSender
	logger  interfaces.Logger
	sem     chan struct{}
}

func NewCollector(config Config, fetcher interfaces.ImageFetcher, sender interfaces.ImageSender, logger interfaces.Logger) *Collector {
	if config.MaxConcurrent <= 0 {
		config.MaxConcurrent = config.CameraCount
	}
	if config.PollInterval == 0 {
		config.PollInterval = 5 * time.Second
	}

	return &Collector{
		config:  config,
		fetcher: fetcher,
		sender:  sender,
		logger:  logger,
		sem:     make(chan struct{}, config.MaxConcurrent),
	}
}

func (c *Collector) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, c.config.CameraCount)

	// Start a goroutine for each camera
	for i := 0; i < c.config.CameraCount; i++ {
		wg.Add(1)
		go func(cameraID int) {
			defer wg.Done()
			// pollCamera runs indefinitely until ctx is canceled
			// or it gets to a fatal situation
			err := c.pollCamera(ctx, cameraID)
			if err != nil && !errors.Is(err, context.Canceled) {
				select {
				case errCh <- fmt.Errorf("camera %d error: %w", cameraID, err):
				case <-ctx.Done():
				}
			}
		}(i + 1)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Handle errors of camera polling
	go func() {
		for err := range errCh {
			c.logger.Printf("Collector caught error: %v", err)
		}
	}()

	<-ctx.Done()
	c.logger.Printf("Collector context canceled; shutting down.")

	return nil
}

func (c *Collector) pollCamera(ctx context.Context, cameraID int) error {
	ticker := time.NewTicker(c.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			select {
			case c.sem <- struct{}{}: // Acquire semaphore
				err := c.processCameraImage(ctx, cameraID)
				<-c.sem // Release semaphore
				if err != nil {
					c.logger.Printf("Camera %d error: %v", cameraID, err)
					continue
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (c *Collector) processCameraImage(ctx context.Context, cameraID int) error {
	// Fetch image
	imageData, err := c.fetcher.FetchImage(ctx, cameraID)
	if err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}

	// Send image
	if err := c.sender.SendImage(ctx, imageData); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	c.logger.Printf("Successfully processed image from camera %d", cameraID)
	return nil
}
