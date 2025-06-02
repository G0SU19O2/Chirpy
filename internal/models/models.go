package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"size:255;uniqueIndex"`
}
type Chirp struct {
	gorm.Model
	Body string `gorm:"not null"`
	UserID uint `gorm:"constraint:OnDelete:CASCADE;"`
	User User `gorm:"foreignKey:UserID"`
}

type ChirpResponse struct {
	Id string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body string `json:"body"`
	UserId string `json:"user_id"`
}
type UserRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

type ChirpRequest struct {
	UserId string `json:"user_id"`
	Body string `json:"body"`
}

type CleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}