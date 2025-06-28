package models

import (
	"time"

	"gorm.io/gorm"
)

type WebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

type RefreshToken struct {
	Token     string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint       `gorm:"constraint:OnDelete:CASCADE;"`
	User      User       `gorm:"foreignKey:UserID"`
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time `gorm:"default:NULL"`
}

type User struct {
	gorm.Model
	Email          string `gorm:"size:255;uniqueIndex"`
	HashedPassword string `gorm:"not null"`
	IsChirpyRed    bool   `gorm:"default:false"`
}

type UserUpdateRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Chirp struct {
	gorm.Model
	Body   string `gorm:"not null"`
	UserID uint   `gorm:"constraint:OnDelete:CASCADE;"`
	User   User   `gorm:"foreignKey:UserID"`
}

type ChirpResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
}
type UserRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserResponse struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Email        string `json:"email"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsChirpyRed  bool   `json:"is_chirpy_red,omitempty"`
}

type ChirpRequest struct {
	UserId string `json:"user_id"`
	Body   string `json:"body"`
}

type CleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
