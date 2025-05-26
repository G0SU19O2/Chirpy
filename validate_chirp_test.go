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
		name           string
		requestBody    string
		expectedStatus int
		expectedValid  *bool
		expectedError  string
	}{
		{
			name:           "Valid chirp",
			requestBody:    `{"body":"This is a valid chirp"}`,
			expectedStatus: http.StatusOK,
			expectedValid:  boolPtr(true),
		},
		{
			name:           "Valid chirp at exactly 140 characters",
			requestBody:    `{"body":"` + strings.Repeat("a", 140) + `"}`,
			expectedStatus: http.StatusOK,
			expectedValid:  boolPtr(true),
		},
		{
			name:           "Chirp too long (141 characters)",
			requestBody:    `{"body":"` + strings.Repeat("a", 141) + `"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Chirp is too long",
		},
		{
			name:           "Empty chirp body",
			requestBody:    `{"body":""}`,
			expectedStatus: http.StatusOK,
			expectedValid:  boolPtr(true),
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
			name:           "Missing body field",
			requestBody:    `{"message":"test"}`,
			expectedStatus: http.StatusOK,
			expectedValid:  boolPtr(true),
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

			if tt.expectedValid != nil {
				if valid, ok := responseBody["valid"].(bool); !ok || valid != *tt.expectedValid {
					t.Errorf("Expected valid: %v, got: %v", *tt.expectedValid, responseBody["valid"])
				}
				if _, hasError := responseBody["error"]; hasError {
					t.Errorf("Expected no error field, but got one: %v", responseBody["error"])
				}
			}

			if tt.expectedError != "" {
				if errorMsg, ok := responseBody["error"].(string); !ok || errorMsg != tt.expectedError {
					t.Errorf("Expected error: %s, got: %v", tt.expectedError, responseBody["error"])
				}
				if _, hasValid := responseBody["valid"]; hasValid {
					t.Errorf("Expected no valid field, but got one: %v", responseBody["valid"])
				}
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
