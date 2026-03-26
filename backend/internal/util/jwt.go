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

	// parse duration like "24h", "15m"
	d, err := time.ParseDuration(expiresIn)
	if err != nil {
		jwtExpiry = 24 * time.Hour
		return
	}
	jwtExpiry = d
}

// GenerateJWT creates a signed token for a given user ID
func GenerateJWT(userID string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("jwt not initialized")
	}

	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(jwtExpiry).Unix(),
		"iat":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT verifies the token and returns the claims
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("jwt not initialized")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// STRICT: only allow HS256
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Explicit expiration check
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if int64(exp) < time.Now().Unix() {
		return nil, errors.New("token expired")
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
