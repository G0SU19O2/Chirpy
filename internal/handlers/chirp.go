package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/G0SU19O2/Chirpy/internal/models"
)

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.ChirpRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	RespondWithJSON(w, http.StatusOK, models.CleanResponse{CleanedBody: cleanProfanity(req.Body)})
}

func cleanProfanity(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(text)

	for i, word := range words {
		for _, profane := range profaneWords {
			if strings.ToLower(word) == profane {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}