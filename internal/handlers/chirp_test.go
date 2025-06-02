package handlers

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/G0SU19O2/Chirpy/internal/models"
	"gorm.io/gorm"
)

func TestValidateChirpBody(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid Chirp",
			body:        "This is a valid chirp with valid length",
			expectError: false,
		},
		{
			name:        "Empty Chirp",
			body:        "",
			expectError: true,
			errorMsg:    "chirp body cannot be empty",
		},
		{
			name:        "Whitespace Only Chirp",
			body:        "   ",
			expectError: true,
			errorMsg:    "chirp body cannot be empty",
		},
		{
			name:        "Too Long Chirp",
			body:        strings.Repeat("a", 141), // 141 chars
			expectError: true,
			errorMsg:    "chirp is too long (max 140 characters)",
		},
		{
			name:        "Max Length Chirp",
			body:        strings.Repeat("a", 140), // 140 chars
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateChirpBody(tc.body)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if err.Error() != tc.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestParseUserID(t *testing.T) {
	tests := []struct {
		name        string
		userIDStr   string
		expectedID  uint
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid User ID",
			userIDStr:   "123",
			expectedID:  123,
			expectError: false,
		},
		{
			name:        "Empty User ID",
			userIDStr:   "",
			expectedID:  0,
			expectError: true,
			errorMsg:    "user ID is required",
		},
		{
			name:        "Non-numeric User ID",
			userIDStr:   "abc",
			expectedID:  0,
			expectError: true,
			errorMsg:    "invalid user ID format",
		},
		{
			name:        "Zero User ID",
			userIDStr:   "0",
			expectedID:  0,
			expectError: true,
			errorMsg:    "user ID must be greater than 0",
		},
		{
			name:        "Very large User ID",
			userIDStr:   "4294967295", // Max uint32
			expectedID:  4294967295,
			expectError: false,
		},
		{
			name:        "User ID exceeding uint32",
			userIDStr:   "4294967296", // Max uint32 + 1
			expectedID:  0,
			expectError: true,
			errorMsg:    "invalid user ID format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := parseUserID(tc.userIDStr)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if err.Error() != tc.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if id != tc.expectedID {
					t.Errorf("expected user ID %d, got %d", tc.expectedID, id)
				}
			}
		})
	}
}

func TestCleanProfanity(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "No profanity",
			text:     "This is a clean text",
			expected: "This is a clean text",
		},
		{
			name:     "Contains profanity",
			text:     "This is kerfuffle and sharbert",
			expected: "This is **** and ****",
		},
		{
			name:     "Case insensitive",
			text:     "This is KERFUFFLE and Sharbert",
			expected: "This is **** and ****",
		},
		{
			name:     "Mixed with other words",
			text:     "kerfuffle start sharbert middle fornax end",
			expected: "**** start **** middle **** end",
		},
		{
			name:     "Only profanity",
			text:     "kerfuffle",
			expected: "****",
		},
		{
			name:     "Empty string",
			text:     "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanProfanity(tc.text)
			if result != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestBuildChirpResponse(t *testing.T) {
	now := time.Now()
	chirp := &models.Chirp{
		Model: gorm.Model{
			ID:        123,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Body:   "Test chirp body",
		UserID: 456,
	}

	response := buildChirpResponse(chirp)

	if response.Id != "123" {
		t.Errorf("expected ID '123', got '%s'", response.Id)
	}

	if response.Body != "Test chirp body" {
		t.Errorf("expected body 'Test chirp body', got '%s'", response.Body)
	}

	if response.UserId != "456" {
		t.Errorf("expected user ID '456', got '%s'", response.UserId)
	}

	expectedCreatedAt := now.Format(time.RFC3339)
	if response.CreatedAt != expectedCreatedAt {
		t.Errorf("expected created_at '%s', got '%s'", expectedCreatedAt, response.CreatedAt)
	}

	expectedUpdatedAt := now.Format(time.RFC3339)
	if response.UpdatedAt != expectedUpdatedAt {
		t.Errorf("expected updated_at '%s', got '%s'", expectedUpdatedAt, response.UpdatedAt)
	}
}

func TestParseChirpRequest(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectedReq *models.ChirpRequest
		expectError bool
	}{
		{
			name:    "Valid request",
			payload: `{"user_id":"123","body":"Test chirp"}`,
			expectedReq: &models.ChirpRequest{
				UserId: "123",
				Body:   "Test chirp",
			},
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			payload:     `{"user_id":123,"body"`,
			expectError: true,
		},
		{
			name:    "Missing body field",
			payload: `{"user_id":"123"}`,
			expectedReq: &models.ChirpRequest{
				UserId: "123",
			},
			expectError: false,
		},
		{
			name:    "Missing user_id field",
			payload: `{"body":"Test chirp"}`,
			expectedReq: &models.ChirpRequest{
				Body: "Test chirp",
			},
			expectError: false,
		},
		{
			name:        "Empty payload",
			payload:     `{}`,
			expectedReq: &models.ChirpRequest{},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/chirps", bytes.NewBufferString(tc.payload))
			result, err := parseChirpRequest(req)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if result.UserId != tc.expectedReq.UserId {
					t.Errorf("expected user_id '%s', got '%s'", tc.expectedReq.UserId, result.UserId)
				}
				if result.Body != tc.expectedReq.Body {
					t.Errorf("expected body '%s', got '%s'", tc.expectedReq.Body, result.Body)
				}
			}
		})
	}
}
