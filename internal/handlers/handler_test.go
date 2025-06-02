package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/G0SU19O2/Chirpy/internal/models"
)

func TestHandleCreateChipWithMocks(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "Empty body",
			payload:        `{"user_id":"123","body":""}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Error != "chirp body cannot be empty" {
					t.Errorf("expected error message 'chirp body cannot be empty', got '%s'", errorResp.Error)
				}
			},
		},
		{
			name:           "Missing user_id",
			payload:        `{"body":"This is a chirp"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user_id",
			payload:        `{"user_id":"abc","body":"This is a chirp"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Zero user_id",
			payload:        `{"user_id":"0","body":"This is a chirp"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Too long chirp",
			payload:        `{"user_id":"123","body":"` + strings.Repeat("a", 141) + `"}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Error != "chirp is too long (max 140 characters)" {
					t.Errorf("expected error message about chirp length, got '%s'", errorResp.Error)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/chirps", bytes.NewBufferString(tc.payload))
			rec := httptest.NewRecorder()

			chirpReq, err := parseChirpRequest(req)
			if err == nil {
				err = validateChirpBody(chirpReq.Body)
				if err == nil {
					_, err = parseUserID(chirpReq.UserId)
				}
			}

			if err != nil {
				RespondWithError(rec, http.StatusBadRequest, err.Error())
			}

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, resp)
			}
		})
	}
}
