package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// HandleRoot handles requests to the root endpoint ("/").
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"name":    "TurfZone API",
		"version": "1.0.0",
	})
}

// HandleHealth handles requests to the health check endpoint ("/health").
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"message": "TurfZone API is running",
		"time":    time.Now().UTC(),
	})
}
