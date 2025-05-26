package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ChirpRequest struct {
	Body string `json:"body"`
}

type CleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ChirpRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	respondWithJSON(w, http.StatusOK, CleanResponse{cleanProfanity(req.Body)})
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

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}
