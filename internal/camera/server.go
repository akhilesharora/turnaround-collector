package camera

import (
	"net/http"

	"github.com/akhilesharora/turnaround-collector/pkg/interfaces"
)

var mockImage = []byte{0xFF, 0xD8, 0xFF, 0xE0}

type Server struct {
	logger interfaces.Logger
}

func NewServer(logger interfaces.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cameraID := r.Header.Get("X-Camera-ID")

	if cameraID == "" {
		http.Error(w, "Camera ID required", http.StatusBadRequest)
		return
	}

	s.logger.Printf("serving image request from camera %s", cameraID)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(mockImage)
}
