package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/models"
)

func HandleCreateChip(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req, err := parseChirpRequest(r)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := validateChirpBody(req.Body); err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		userID, err := parseUserID(req.UserId)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		chirp := &models.Chirp{
			Body:   cleanProfanity(req.Body),
			UserID: userID,
		}

		if err := cfg.DB.Create(chirp).Error; err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
			return
		}

		response := buildChirpResponse(chirp)
		RespondWithJSON(w, http.StatusCreated, response)
	}
}

func parseChirpRequest(r *http.Request) (*models.ChirpRequest, error) {
	var req models.ChirpRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON format")
	}

	return &req, nil
}

func validateChirpBody(body string) error {
	const maxChirpLength = 140

	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("chirp body cannot be empty")
	}

	if len(body) > maxChirpLength {
		return fmt.Errorf("chirp is too long (max %d characters)", maxChirpLength)
	}

	return nil
}

func parseUserID(userIDStr string) (uint, error) {
	if userIDStr == "" {
		return 0, fmt.Errorf("user ID is required")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format")
	}

	if userID == 0 {
		return 0, fmt.Errorf("user ID must be greater than 0")
	}

	return uint(userID), nil
}

func buildChirpResponse(chirp *models.Chirp) models.ChirpResponse {
	return models.ChirpResponse{
		Id:        strconv.FormatUint(uint64(chirp.ID), 10),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserId:    strconv.FormatUint(uint64(chirp.UserID), 10),
	}
}

func cleanProfanity(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(text)

	for i, word := range words {
		for _, profane := range profaneWords {
			if strings.EqualFold(word, profane) {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}
