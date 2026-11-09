package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	handler := setupRoutes()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", body["status"])
	}
	if body["message"] != "TurfZone API is running" {
		t.Errorf("expected message 'TurfZone API is running', got %v", body["message"])
	}
}

func TestRootRoute(t *testing.T) {
	handler := setupRoutes()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["name"] != "TurfZone API" {
		t.Errorf("expected name 'TurfZone API', got %v", body["name"])
	}
	if body["version"] != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %v", body["version"])
	}
}

func TestNotFoundRoute(t *testing.T) {
	handler := setupRoutes()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/non-existent", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestCORSIntegration(t *testing.T) {
	t.Run("Preflight OPTIONS request with wildcard", func(t *testing.T) {
		handler := setupRoutes()

		req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/health", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
			t.Fatalf("expected 200 or 204 for OPTIONS, got %d", rr.Code)
		}

		if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
			t.Errorf("expected Access-Control-Allow-Origin '*', got %q", got)
		}
	})

	t.Run("Specific origin allowed with credentials", func(t *testing.T) {
		t.Setenv("CORS_ALLOWED_ORIGIN", "http://localhost:3000")
		handler := setupRoutes()

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
			t.Errorf("expected origin 'http://localhost:3000', got %q", got)
		}

		if got := rr.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Errorf("expected credentials 'true', got %q", got)
		}
	})
}
