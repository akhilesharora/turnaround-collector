package camera

import (
	"net/http"
	"net/http/httptest"
	"testing"

	testutil "github.com/akhilesharora/turnaround-collector/pkg/testutils"
)

func TestCameraServer(t *testing.T) {
	tests := []struct {
		name           string
		cameraID       string
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "valid request",
			cameraID:       "camera_1",
			expectedStatus: http.StatusOK,
			expectedType:   "image/jpeg",
		},
		{
			name:           "missing camera ID",
			cameraID:       "",
			expectedStatus: http.StatusBadRequest,
			expectedType:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &testutil.MockLogger{}
			server := NewServer(logger)

			req := httptest.NewRequest(http.MethodGet, "/snap.jpg", nil)
			if tt.cameraID != "" {
				req.Header.Set("X-Camera-ID", tt.cameraID)
			}
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedType != "" && w.Header().Get("Content-Type") != tt.expectedType {
				t.Errorf("expected content-type %s, got %s", tt.expectedType, w.Header().Get("Content-Type"))
			}
		})
	}
}
