package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/G0SU19O2/Chirpy/internal/models"
)

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, models.ErrorResponse{Error: message})
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}