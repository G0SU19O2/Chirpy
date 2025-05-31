package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type UserRequest struct {
	Email string `json:"email"`
}
type UserResponse struct {
	Id         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Email      string `json:"email"`
}

type User struct {
	gorm.Model
	Email string `gorm:"size:255;uniqueIndex"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UserRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	resp, err := createUserResponse(cfg.db, req.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something wrong")
		return
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func createUserResponse(db *gorm.DB, email string) (UserResponse, error) {
	user := &User{Email: email}
	result := db.Create(user)
	return UserResponse{Id: strconv.FormatUint(uint64(user.ID), 10), Created_at: user.CreatedAt.String(), Updated_at: user.UpdatedAt.String(), Email: user.Email}, result.Error
}
