package interfaces

import (
	"context"
)

// Collector interfaces
type ImageFetcher interface {
	FetchImage(ctx context.Context, cameraID int) ([]byte, error)
}

type ImageSender interface {
	SendImage(ctx context.Context, imageData []byte) error
}

// Target server interfaces
type ImageProcessor interface {
	Process(imageData []byte) error
}
