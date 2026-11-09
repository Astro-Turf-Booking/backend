package main

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"turfzone-backend/cmd/internal/database"
	"turfzone-backend/cmd/internal/migrations"
	"turfzone-backend/middleware"
	"turfzone-backend/services/auth"
	"turfzone-backend/services/health"

	"github.com/rs/cors"
	"gorm.io/gorm"
)

func setupRoutes(db ...*gorm.DB) http.Handler {
	mux := http.NewServeMux()

	var database *gorm.DB
	if len(db) > 0 {
		database = db[0]
	}

	// Dedicated sub-mux for v1 API routes
	apiV1Mux := http.NewServeMux()
	health.RegisterRoutes(apiV1Mux) // Registers routes like /health onto the sub-mux
	auth.RegisterRoutes(apiV1Mux, database)

	// Mount the sub-mux under /api/v1/ and strip the prefix before handing over
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1Mux))

	allowedOriginsStr := cmp.Or(os.Getenv("CORS_ALLOWED_ORIGIN"), "*")
	var allowedOrigins []string
	if allowedOriginsStr == "*" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: allowedOriginsStr != "*",
		MaxAge:           300,
	})

	return c.Handler(middleware.Logging(mux))
}

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
	_ = db //

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	port := cmp.Or(os.Getenv("PORT"), "8080")
	addr := fmt.Sprintf(":%s", port)

	handler := setupRoutes(db)

	log.Printf("Starting dev server on http://localhost:%s", port)
	if err := http.ListenAndServe(addr, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed to start: %v", err)
	}
}
