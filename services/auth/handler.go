package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"turfzone-backend/models"
	"turfzone-backend/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Handler handles authentication-related HTTP requests.
type Handler struct {
	DB *gorm.DB // database connection used by all auth handlers.
}

// RegisterRequest defines the expected payload for user registration.
type RegisterRequest struct {
	Email       string          `json:"email"`
	FullName    string          `json:"full_name"`
	Password    string          `json:"password"`
	PhoneNumber string          `json:"phone_number"`
	Role        models.UserRole `json:"role"`
}

// LoginRequest defines the expected payload for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the data payload returned on successful authentication.
type AuthResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token,omitempty"`
}

// Register handles user registration requests.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	fullName := strings.TrimSpace(req.FullName)
	phoneNumber := strings.TrimSpace(req.PhoneNumber)
	password := req.Password

	if email == "" || fullName == "" || password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email, full name, and password are required")
		return
	}

	if !strings.Contains(email, "@") {
		utils.RespondError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	if len(password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	if h.DB != nil {
		var existingUser models.User
		if err := h.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
			utils.RespondError(w, http.StatusConflict, "User with this email already exists")
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusInternalServerError, "Database error checking existing user")
			return
		}

		if phoneNumber != "" {
			var phoneUser models.User
			if err := h.DB.Where("phone_number = ?", phoneNumber).First(&phoneUser).Error; err == nil {
				utils.RespondError(w, http.StatusConflict, "User with this phone number already exists")
				return
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.RespondError(w, http.StatusInternalServerError, "Database error checking phone number")
				return
			}
		}
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	now := time.Now().UTC()
	user := models.User{
		ID:          uuid.New(),
		Email:       email,
		FullName:    fullName,
		Password:    hashedPassword,
		PhoneNumber: phoneNumber,
		Role:        models.RoleUser,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if h.DB != nil {
		if err := h.DB.Create(&user).Error; err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.RespondSuccess(w, http.StatusCreated, "User registered successfully", AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login handles user login requests.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	if email == "" || password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if h.DB == nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database connection unavailable")
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	now := time.Now().UTC()
	user.LastLogin = &now
	_ = h.DB.Model(&user).Update("last_login", now)

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.RespondSuccess(w, http.StatusOK, "Login successful", AuthResponse{
		User:  user,
		Token: token,
	})
}

// Logout handles user logout requests by clearing the auth token cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	utils.RespondSuccess(w, http.StatusOK, "Logout successful", nil)
}
