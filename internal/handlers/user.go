package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/models"
	"gorm.io/gorm"
)

func HandleCreateUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.UserRequest
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		resp, err := createUserResponse(cfg.DB, req.Email)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Something wrong")
			return
		}
		RespondWithJSON(w, http.StatusCreated, resp)
	}
}

func createUserResponse(db *gorm.DB, email string) (models.UserResponse, error) {
	user := &models.User{Email: email}
	result := db.Create(user)
	return models.UserResponse{
		ID:        strconv.FormatUint(uint64(user.ID), 10),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	}, result.Error
}
