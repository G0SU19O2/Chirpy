package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerValidateChirp(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     string
		expectedStatus  int
		expectedCleaned string
		expectedError   string
	}{
		{
			name:            "Valid chirp without profanity",
			requestBody:     `{"body":"This is a valid chirp"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "This is a valid chirp",
		},
		{
			name:            "Valid chirp at exactly 140 characters",
			requestBody:     `{"body":"` + strings.Repeat("a", 140) + `"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: strings.Repeat("a", 140),
		},
		{
			name:           "Chirp too long (141 characters)",
			requestBody:    `{"body":"` + strings.Repeat("a", 141) + `"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Chirp is too long",
		},
		{
			name:            "Empty chirp body",
			requestBody:     `{"body":""}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "",
		},
		{
			name:           "Invalid JSON - missing quotes",
			requestBody:    `{body:test}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:           "Invalid JSON - malformed",
			requestBody:    `{"body":"test"`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:           "Empty request body",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:            "Missing body field",
			requestBody:     `{"message":"test"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "",
		},
		{
			name:            "Chirp with single profane word",
			requestBody:     `{"body":"This is a kerfuffle opinion"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "This is a **** opinion",
		},
		{
			name:            "Chirp with multiple profane words",
			requestBody:     `{"body":"What a kerfuffle! I love sharbert and fornax"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "What a kerfuffle! I love **** and ****",
		},
		{
			name:            "Chirp with uppercase profane words",
			requestBody:     `{"body":"KERFUFFLE and Sharbert are FORNAX"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "**** and **** are ****",
		},
		{
			name:            "Chirp with mixed case profane words",
			requestBody:     `{"body":"KerfUffle ShArbert FoRnAx"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "**** **** ****",
		},
		{
			name:            "Chirp with profane words with punctuation (should not be replaced)",
			requestBody:     `{"body":"Sharbert! kerfuffle? fornax."}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "Sharbert! kerfuffle? fornax.",
		},
		{
			name:            "Chirp with words containing profane substrings (should not be replaced)",
			requestBody:     `{"body":"sharberted kerfuffled fornaxed"}`,
			expectedStatus:  http.StatusOK,
			expectedCleaned: "sharberted kerfuffled fornaxed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/validate_chirp", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handlerValidateChirp(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			expectedContentType := "application/json"
			if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
			}

			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			if tt.expectedCleaned != "" {
				if cleanedBody, ok := responseBody["cleaned_body"].(string); !ok || cleanedBody != tt.expectedCleaned {
					t.Errorf("Expected cleaned_body: %q, got: %v", tt.expectedCleaned, responseBody["cleaned_body"])
				}
				if _, hasError := responseBody["error"]; hasError {
					t.Errorf("Expected no error field, but got one: %v", responseBody["error"])
				}
			}

			if tt.expectedError != "" {
				if errorMsg, ok := responseBody["error"].(string); !ok || errorMsg != tt.expectedError {
					t.Errorf("Expected error: %s, got: %v", tt.expectedError, responseBody["error"])
				}
				if _, hasCleaned := responseBody["cleaned_body"]; hasCleaned {
					t.Errorf("Expected no cleaned_body field, but got one: %v", responseBody["cleaned_body"])
				}
			}
		})
	}
}
