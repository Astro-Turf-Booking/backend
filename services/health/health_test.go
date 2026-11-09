package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"turfzone-backend/services/health"
)

func TestHealthEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	health.RegisterRoutes(mux)

	t.Run("GET / returns API metadata", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var body map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to parse json response: %v", err)
		}

		if body["name"] != "TurfZone API" || body["version"] != "1.0.0" {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("GET /health returns health status", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var body map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to parse json response: %v", err)
		}

		if body["status"] != "ok" || body["message"] != "TurfZone API is running" {
			t.Errorf("unexpected body: %+v", body)
		}
	})

	t.Run("GET /unknown returns 404", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/unknown", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}
