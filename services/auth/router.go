package auth

import (
	"net/http"

	"gorm.io/gorm"
)

// RegisterRoutes registers the authentication routes on the provided ServeMux.
func RegisterRoutes(mux *http.ServeMux, db *gorm.DB) {
	h := &Handler{DB: db}
	h.RegisterRoutes(mux)
}

// RegisterRoutes registers the authentication routes for this handler instance.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/logout", h.Logout)
}
