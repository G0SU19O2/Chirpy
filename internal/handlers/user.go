package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

func userToResponse(user *models.User, token string, refreshToken string) models.UserResponse {
	return models.UserResponse{
		ID:           strconv.FormatUint(uint64(user.ID), 10),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
}

func createRefreshToken(db *gorm.DB, userID uint) (*models.RefreshToken, error) {
	token, err := auth.MakeRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshToken := &models.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(60 * time.Second),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return refreshToken, db.Create(refreshToken).Error
}

func HandleCreateUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		resp := userToResponse(user, "", "")
		RespondWithJSON(w, http.StatusCreated, resp)
	}
}

func HandleLoginUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		expiresIn := req.ExpiresInSeconds
		if expiresIn <= 0 || expiresIn > 3600 {
			expiresIn = 3600
		}
		token, err := auth.MakeJWT(strconv.FormatUint(uint64(user.ID), 10), cfg.JWTSecret)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Could not create token")
			return
		}
		refreshToken, err := createRefreshToken(cfg.DB, user.ID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Could not create refresh token")
			return
		}
		resp := userToResponse(&user, token, refreshToken.Token)
		RespondWithJSON(w, http.StatusOK, resp)
	}
}

func HandleRefreshToken(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		refreshToken, err := auth.ValidateRefreshToken(cfg.DB, token)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var user models.User
		result := cfg.DB.First(&user, refreshToken.UserID)
		if result.Error != nil {
			RespondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		newToken, err := auth.MakeJWT(strconv.FormatUint(uint64(user.ID), 10), cfg.JWTSecret)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Could not create new token")
			return
		}

		resp := map[string]string{"token": newToken}
		RespondWithJSON(w, http.StatusOK, resp)
	}
}

func HandleRevokeToken(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		refreshToken, err := auth.ValidateRefreshToken(cfg.DB, token)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err := auth.RevokeRefreshToken(cfg.DB, refreshToken); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleUpdateUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		tokenUserId, err := auth.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		var user models.User
		if err := cfg.DB.First(&user, tokenUserId).Error; err != nil {
			RespondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		var req models.UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if req.Email != "" {
			user.Email = req.Email
		}

		if req.Password != "" {
			hashedPassword, err := auth.HashPassword(req.Password)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
				return
			}
			user.HashedPassword = hashedPassword
		}

		if err := cfg.DB.Save(&user).Error; err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		resp := userToResponse(&user, "", "")
		RespondWithJSON(w, http.StatusOK, resp)
	}
}