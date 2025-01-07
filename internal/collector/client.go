package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	client    *http.Client
	baseURL   string
	targetURL string
}

func NewClient(timeout time.Duration, baseURL, targetURL string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL:   baseURL,
		targetURL: targetURL,
	}
}

func (c *Client) FetchImage(ctx context.Context, cameraID int) ([]byte, error) {
	baseURL := strings.TrimRight(c.baseURL, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/snap.jpg", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set camera ID in a custom header
	req.Header.Set("X-Camera-ID", fmt.Sprintf("camera_%d", cameraID))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) SendImage(ctx context.Context, imageData []byte) error {
	targetURL := c.targetURL

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedURL.String(),
		bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "image/jpeg")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("send image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
