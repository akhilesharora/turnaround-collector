package target

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akhilesharora/turnaround-collector/pkg/testutils"
)

// Mocks
type mockProcessor struct {
	processFunc func([]byte) error
}

func (m *mockProcessor) Process(imageData []byte) error {
	return m.processFunc(imageData)
}

func TestTargetServer(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		processError   error
		expectedStatus int
		checkLogs      func(t *testing.T, logs []string)
	}{
		{
			name:           "valid request",
			method:         http.MethodPost,
			path:           "/image",
			processError:   nil,
			expectedStatus: http.StatusOK,
			checkLogs: func(t *testing.T, logs []string) {
				found := false
				for _, log := range logs {
					if found = strings.Contains(log, "successfully processed"); found {
						break
					}
				}
				if !found {
					t.Error("expected success log message")
				}
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			path:           "/image",
			processError:   nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid path",
			method:         http.MethodPost,
			path:           "/invalid",
			processError:   nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "processing error",
			method:         http.MethodPost,
			path:           "/image",
			processError:   errors.New("process error"),
			expectedStatus: http.StatusInternalServerError,
			checkLogs: func(t *testing.T, logs []string) {
				found := false
				for _, log := range logs {
					if found = strings.Contains(log, "error processing"); found {
						break
					}
				}
				if !found {
					t.Error("expected error log message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &testutils.MockLogger{}
			processor := &mockProcessor{
				processFunc: func([]byte) error {
					return tt.processError
				},
			}

			server := NewServer(logger, processor)

			body := bytes.NewReader([]byte("test image"))
			req := httptest.NewRequest(tt.method, tt.path, body)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkLogs != nil {
				tt.checkLogs(t, logger.Logs)
			}

			if t.Failed() {
				t.Logf("Logs:\n%s", strings.Join(logger.Logs, "\n"))
			}
		})
	}
}
