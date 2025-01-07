package interfaces

import (
	"net/http"
)

type Server interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
