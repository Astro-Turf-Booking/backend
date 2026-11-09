package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"turfzone-backend/middleware"
)

func TestLogging(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	testHandler := middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/test", nil)
	rr := httptest.NewRecorder()

	testHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	logOutput := logBuf.String()
	if !bytes.Contains([]byte(logOutput), []byte("[POST] /test 201")) {
		t.Errorf("expected log output to contain '[POST] /test 201', got: %s", logOutput)
	}
}
