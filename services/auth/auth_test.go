package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister_Validation(t *testing.T) {
	h := &Handler{DB: nil}

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "wrong method GET",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantMsg:    "Method not allowed",
		},
		{
			name:       "invalid JSON",
			method:     http.MethodPost,
			body:       "{invalid",
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid request body",
		},
		{
			name:       "missing email",
			method:     http.MethodPost,
			body:       `{"full_name":"Test User","password":"password123"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Email, full name, and password are required",
		},
		{
			name:       "invalid email",
			method:     http.MethodPost,
			body:       `{"email":"invalid-email","full_name":"Test User","password":"password123"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid email address",
		},
		{
			name:       "password too short",
			method:     http.MethodPost,
			body:       `{"email":"test@example.com","full_name":"Test User","password":"123"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Password must be at least 6 characters",
		},
		{
			name:       "successful registration defaults to USER without role field",
			method:     http.MethodPost,
			body:       `{"email":"test@example.com","full_name":"Test User","password":"password123"}`,
			wantStatus: http.StatusCreated,
			wantMsg:    "User registered successfully",
		},
		{
			name:       "successful registration with role defaults to USER",
			method:     http.MethodPost,
			body:       `{"email":"test@example.com","full_name":"Test User","password":"password123","role":"PROVIDER"}`,
			wantStatus: http.StatusCreated,
			wantMsg:    "User registered successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/auth/register", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			h.Register(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d (body: %s)", tt.wantStatus, rr.Code, rr.Body.String())
			}

			var resp map[string]any
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if msg, ok := resp["message"].(string); !ok || !strings.Contains(msg, tt.wantMsg) {
				t.Errorf("expected message containing %q, got %q", tt.wantMsg, msg)
			}

			if tt.wantStatus == http.StatusCreated {
				cookies := rr.Result().Cookies()
				var foundCookie bool
				for _, c := range cookies {
					if c.Name == "token" && c.Value != "" {
						foundCookie = true
						break
					}
				}
				if !foundCookie {
					t.Error("expected auth cookie 'token' to be set")
				}

				if data, ok := resp["data"].(map[string]any); ok {
					if user, ok := data["user"].(map[string]any); ok {
						if role, ok := user["role"].(string); !ok || role != "USER" {
							t.Errorf("expected user role to be 'USER', got %v", user["role"])
						}
					} else {
						t.Error("expected user in response data")
					}
				} else {
					t.Error("expected data in response")
				}
			}
		})
	}
}

func TestLogin_Validation(t *testing.T) {
	h := &Handler{DB: nil}

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "wrong method GET",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantMsg:    "Method not allowed",
		},
		{
			name:       "invalid JSON",
			method:     http.MethodPost,
			body:       "{invalid",
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid request body",
		},
		{
			name:       "missing password",
			method:     http.MethodPost,
			body:       `{"email":"test@example.com"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Email and password are required",
		},
		{
			name:       "missing email",
			method:     http.MethodPost,
			body:       `{"password":"password123"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Email and password are required",
		},
		{
			name:       "database unavailable",
			method:     http.MethodPost,
			body:       `{"email":"test@example.com","password":"password123"}`,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "Database connection unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/auth/login", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			h.Login(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d (body: %s)", tt.wantStatus, rr.Code, rr.Body.String())
			}

			var resp map[string]any
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if msg, ok := resp["message"].(string); !ok || !strings.Contains(msg, tt.wantMsg) {
				t.Errorf("expected message containing %q, got %q", tt.wantMsg, msg)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	h := &Handler{DB: nil}

	t.Run("wrong method GET", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/auth/logout", nil)
		rr := httptest.NewRecorder()

		h.Logout(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("successful logout", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", nil)
		rr := httptest.NewRecorder()

		h.Logout(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		cookies := rr.Result().Cookies()
		var foundClearedCookie bool
		for _, c := range cookies {
			if c.Name == "token" && c.Value == "" && c.MaxAge == -1 {
				foundClearedCookie = true
				break
			}
		}
		if !foundClearedCookie {
			t.Error("expected cleared auth cookie 'token'")
		}
	})
}

func TestRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, nil)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d on registered /auth/logout, got %d", http.StatusOK, rr.Code)
	}
}
