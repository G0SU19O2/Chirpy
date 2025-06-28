package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/G0SU19O2/Chirpy/internal/auth"
	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/models"
)

func HandleWebHook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil || apiKey != cfg.PolkaAPIKey {
			RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}
		var req models.WebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if req.Event != "user.upgraded" {
			RespondWithError(w, http.StatusBadRequest, "Unsupported event type")
			return
		}

		userID, err := strconv.ParseUint(req.Data.UserID, 10, 32)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var user models.User
		if err := cfg.DB.First(&user, uint(userID)).Error; err != nil {
			RespondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		user.IsChirpyRed = true
		if err := cfg.DB.Save(&user).Error; err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
