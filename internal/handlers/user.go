package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/G0SU19O2/Chirpy/internal/auth"
	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/models"
	"gorm.io/gorm"
)

func createUser(db *gorm.DB, email string, password string) (*models.User, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{Email: email, HashedPassword: hash}
	result := db.Create(user)
	return user, result.Error
}

func userToResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:        strconv.FormatUint(uint64(user.ID), 10),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	}
}

func HandleCreateUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.UserRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		user, err := createUser(cfg.DB, req.Email, req.Password)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Something wrong")
			return
		}
		resp := userToResponse(user)
		RespondWithJSON(w, http.StatusCreated, resp)
	}
}

func HandleLoginUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.UserRequest
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		var user models.User
		result := cfg.DB.Where("email = ?", req.Email).First(&user)
		if result.Error != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		resp := userToResponse(&user)
		RespondWithJSON(w, http.StatusOK, resp)
	}
}
