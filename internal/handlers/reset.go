package handlers

import (
	"net/http"

	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/models"
	"gorm.io/gorm"
)

func HandleReset(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.Platform != "DEV" {
			RespondWithError(w, http.StatusForbidden, "Forbidden")
			return
		}
		result := cfg.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
		if result.Error != nil {
			RespondWithError(w, http.StatusInternalServerError, "Cannot remove users")
			return
		}
		RespondWithJSON(w, http.StatusOK, "Remove all users")
	}
}
