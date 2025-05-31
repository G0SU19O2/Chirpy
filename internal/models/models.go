package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"size:255;uniqueIndex"`
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
	Body string `json:"body"`
}

type CleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}