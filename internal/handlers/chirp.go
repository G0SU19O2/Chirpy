package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/G0SU19O2/Chirpy/internal/auth"
	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/models"
	"gorm.io/gorm"
)

func HandleDeleteChirp(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpID, err := parseChirpIDFromPath(r)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

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

		chirp, err := findChirpByID(cfg.DB, chirpID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				RespondWithError(w, http.StatusNotFound, "Chirp not found")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp")
			return
		}

		if !isChirpOwner(tokenUserId, chirp.UserID) {
			RespondWithError(w, http.StatusUnauthorized, "You are not authorized to delete this chirp")
			return
		}

		if err := cfg.DB.Delete(&chirp).Error; err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func parseChirpIDFromPath(r *http.Request) (uint, error) {
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		return 0, fmt.Errorf("Chirp ID is required")
	}
	chirpID, err := strconv.ParseUint(chirpIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Invalid chirp ID format")
	}
	return uint(chirpID), nil
}

func findChirpByID(db *gorm.DB, chirpID uint) (*models.Chirp, error) {
	var chirp models.Chirp
	result := db.First(&chirp, chirpID)
	return &chirp, result.Error
}

func isChirpOwner(tokenUserId string, chirpUserID uint) bool {
	return tokenUserId == strconv.FormatUint(uint64(chirpUserID), 10)
}

func HandleGetChirpById(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpIDStr := r.PathValue("chirpID")
		if chirpIDStr == "" {
			RespondWithError(w, http.StatusBadRequest, "Chirp ID is required")
			return
		}
		chirpID, err := strconv.ParseUint(chirpIDStr, 10, 32)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID format")
			return
		}
		var chirp models.Chirp
		result := cfg.DB.First(&chirp, uint(chirpID))
		if result.Error != nil {
			RespondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		response := buildChirpResponse(&chirp)
		RespondWithJSON(w, http.StatusOK, response)
	}
}

func HandleGetAllChirps(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorIDStr := r.URL.Query().Get("author_id")
		sortOrder := r.URL.Query().Get("sort")

		query := cfg.DB
		if authorIDStr != "" {
			authorID, err := strconv.ParseUint(authorIDStr, 10, 32)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, "Invalid author_id format")
				return
			}
			query = query.Where("user_id = ?", uint(authorID))
		}

		var chirps []models.Chirp
		if err := query.Find(&chirps).Error; err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
			return
		}

		responses := make([]models.ChirpResponse, len(chirps))
		for i, chirp := range chirps {
			responses[i] = buildChirpResponse(&chirp)
		}

		sortResponses(responses, sortOrder)

		RespondWithJSON(w, http.StatusOK, responses)
	}
}

func sortResponses(responses []models.ChirpResponse, sortOrder string) {
	if sortOrder == "desc" {
		sort.Slice(responses, func(i, j int) bool {
			return responses[i].CreatedAt > responses[j].CreatedAt
		})
	} else {
		sort.Slice(responses, func(i, j int) bool {
			return responses[i].CreatedAt < responses[j].CreatedAt
		})
	}
}

func HandleCreateChirp(cfg *config.Config) http.HandlerFunc {
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

		if tokenUserId != strconv.FormatUint(uint64(userID), 10) {
			RespondWithError(w, http.StatusUnauthorized, "Invalid token")
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
