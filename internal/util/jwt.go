package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret []byte
	jwtExpiry time.Duration
)

// InitJWT must be called once from main.go
func InitJWT(secret string, expiresIn string) {
	jwtSecret = []byte(secret)

	// parse duration like "24h", "15m", "7d"
	d, err := time.ParseDuration(expiresIn)
	if err != nil {
		jwtExpiry = 24 * time.Hour
		return
	}
	jwtExpiry = d
}

// GenerateJWT creates a signed token for a given user ID
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT verifies the token and returns the claims
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractUserID extracts user_id from JWT claims
func ExtractUserID(claims jwt.MapClaims) (string, error) {
	uid, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("missing user_id in token")
	}
	return uid, nil
}
