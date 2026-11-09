package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser     UserRole = "USER"
	RoleProvider UserRole = "PROVIDER"
)

type User struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email               string     `gorm:"uniqueIndex;not null" json:"email"`
	FullName            string     `gorm:"not null" json:"full_name"`
	Password            string     `gorm:"not null" json:"-"`
	ProfilePicture      string     `json:"profile_picture"`
	PhoneNumber         string     `gorm:"uniqueIndex" json:"phone_number"`
	ResetPasswordToken  string     `json:"-"`
	ResetPasswordExpiry *time.Time `json:"-"`
	Role                UserRole   `gorm:"type:varchar(20);not null;default:'USER';check:role IN ('USER', 'PROVIDER')" json:"role"`
	LastLogin           *time.Time `json:"last_login"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	IsEmailVerified     bool       `gorm:"default:false" json:"is_email_verified"`
}
