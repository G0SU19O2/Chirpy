package main

import (
	"net/http"

	"gorm.io/gorm"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "DEV" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	result := cfg.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot remove users")
		return
	}
	respondWithJSON(w, http.StatusOK, "Remove all users")
}
