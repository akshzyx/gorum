package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/akshzyx/gorum/internal/util"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// helper: extracts token from Authorization header OR cookie
func extractToken(r *http.Request) (string, error) {
	// 1. Try Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return "", errors.New("invalid authorization header")
		}
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	// 2. Try cookie
	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", errors.New("missing auth token")
}

// helper: extracts userID from token
func extractUser(r *http.Request) (string, error) {
	tokenStr, err := extractToken(r)
	if err != nil {
		return "", err
	}

	claims, err := util.ValidateJWT(tokenStr)
	if err != nil {
		return "", err
	}

	userID, err := util.ExtractUserID(claims)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// OptionalAuth attaches user to context if token is present
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := extractToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := util.ValidateJWT(tokenStr)
		if err != nil {
			util.Unauthorized(w, r, err)
			return
		}

		userID, err := util.ExtractUserID(claims)
		if err != nil {
			util.Unauthorized(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth enforces authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := extractUser(r)
		if err != nil {
			util.Unauthorized(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts user_id from context safely
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
