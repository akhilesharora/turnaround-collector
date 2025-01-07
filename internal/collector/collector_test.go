package collector

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	testutil "github.com/akhilesharora/turnaround-collector/pkg/testutils"
)

// Mock fetcher
type mockFetcher struct {
	fetchFunc func(ctx context.Context, cameraID int) ([]byte, error)
}

func (m *mockFetcher) FetchImage(ctx context.Context, cameraID int) ([]byte, error) {
	return m.fetchFunc(ctx, cameraID)
}

// Mock sender
type mockSender struct {
	sendFunc func(ctx context.Context, imageData []byte) error
}

func (m *mockSender) SendImage(ctx context.Context, imageData []byte) error {
	return m.sendFunc(ctx, imageData)
}

func TestCollector(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		fetchError     error
		sendError      error
		expectSuccess  bool
		expectedErrLog string
	}{
		{
			name: "successful image processing",
			config: Config{
				CameraCount:   1,
				PollInterval:  10 * time.Millisecond,
				MaxConcurrent: 1,
				CameraBaseURL: "http://camera",
				TargetURL:     "http://target:8080/image",
			},
			expectSuccess:  true,
			expectedErrLog: "",
		},
		{
			name: "fetch error",
			config: Config{
				CameraCount:   1,
				PollInterval:  10 * time.Millisecond,
				MaxConcurrent: 1,
			},
			fetchError:     errors.New("fetch failed"),
			expectSuccess:  false,
			expectedErrLog: "fetch failed",
		},
		{
			name: "send error",
			config: Config{
				CameraCount:   1,
				PollInterval:  10 * time.Millisecond,
				MaxConcurrent: 1,
			},
			sendError:      errors.New("send failed"),
			expectSuccess:  false,
			expectedErrLog: "send failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := &testutil.MockLogger{}

			// Set up fetcher
			fetcher := &mockFetcher{
				fetchFunc: func(ctx context.Context, cameraID int) ([]byte, error) {
					if tt.fetchError != nil {
						return nil, tt.fetchError
					}
					return []byte("test image"), nil
				},
			}

			// Set up sender
			sender := &mockSender{
				sendFunc: func(ctx context.Context, imageData []byte) error {
					if tt.sendError != nil {
						return tt.sendError
					}
					return nil
				},
			}

			collector := NewCollector(tt.config, fetcher, sender, logger)

			// Cancel after ~100ms
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Start the collector
			errCh := make(chan error, 1)
			go func() {
				errCh <- collector.Start(ctx)
			}()

			// After 50ms cancel context, enough for at least one poll cycle
			time.Sleep(50 * time.Millisecond)
			cancel()

			<-errCh

			hasSuccessLog := false
			hasErrorLog := false
			for _, logMsg := range logger.Logs {
				if strings.Contains(logMsg, "Successfully processed") {
					hasSuccessLog = true
				}
				if tt.expectedErrLog != "" && strings.Contains(logMsg, tt.expectedErrLog) {
					hasErrorLog = true
				}
			}

			if tt.expectSuccess {
				// Expect at least one successful process
				if !hasSuccessLog {
					t.Errorf("expected success in logs but found none\nAll logs:\n%s",
						strings.Join(logger.Logs, "\n"))
				}
			} else {
				// Expect no success
				if hasSuccessLog {
					t.Errorf("did NOT expect success log, but found one\nAll logs:\n%s",
						strings.Join(logger.Logs, "\n"))
				}
			}

			// Expect error
			if tt.expectedErrLog != "" && !hasErrorLog {
				t.Errorf("expected to find '%s' in logs but did not.\nAll logs:\n%s",
					tt.expectedErrLog, strings.Join(logger.Logs, "\n"))
			}
		})
	}
}
