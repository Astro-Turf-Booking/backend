package health

import (
	"net/http"
)

// RegisterRoutes registers the root and health check routes on the provided ServeMux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", HandleRoot)
	mux.HandleFunc("GET /health", HandleHealth)
}
