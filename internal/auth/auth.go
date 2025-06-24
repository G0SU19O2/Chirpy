package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/G0SU19O2/Chirpy/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	hashed, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), error
}

func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID, tokenSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject:   userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("authorization header is missing")
	}
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", errors.New("authorization header must start with 'Bearer '")
	}
	return token[7:], nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ValidateRefreshToken(db *gorm.DB, tokenStr string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	result := db.Where("token = ?", tokenStr).First(&refreshToken)
	if result.Error != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.RevokedAt != nil || time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token is invalid or expired")
	}

	return &refreshToken, nil
}

func RevokeRefreshToken(db *gorm.DB, token *models.RefreshToken) error {
	now := time.Now()
	token.RevokedAt = &now
	token.UpdatedAt = now

	result := db.Save(token)
	return result.Error
}
