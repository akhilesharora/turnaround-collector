package target

import (
	"io"
	"net/http"

	"github.com/akhilesharora/turnaround-collector/pkg/interfaces"
)

type Server struct {
	logger    interfaces.Logger
	processor interfaces.ImageProcessor
}

func NewServer(logger interfaces.Logger, processor interfaces.ImageProcessor) *Server {
	if processor == nil {
		processor = &defaultProcessor{logger: logger}
	}
	return &Server{
		logger:    logger,
		processor: processor,
	}
}

type defaultProcessor struct {
	logger interfaces.Logger
}

func (p *defaultProcessor) Process(imageData []byte) error {
	p.logger.Printf("processed image of size %d bytes", len(imageData))
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/image" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("error reading request body: %v", err)
		http.Error(w, "Failed to read image", http.StatusBadRequest)
		return
	}

	if err := s.processor.Process(imageData); err != nil {
		s.logger.Printf("error processing image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	s.logger.Printf("successfully processed image from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
