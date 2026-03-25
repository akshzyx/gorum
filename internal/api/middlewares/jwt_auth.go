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

// helper: extracts userID from Authorization header
func extractUserFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization token")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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

// OptionalAuth attaches user to context if token is present.
// If no token is provided, request continues as unauthenticated.
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token → treat as guest request
			next.ServeHTTP(w, r)
			return
		}

		userID, err := extractUserFromHeader(r)
		if err != nil {
			// Token exists but is invalid → unauthorized
			util.Unauthorized(w, r, err)
			return
		}

		// Attach user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth enforces authentication.
// Request is rejected if token is missing or invalid.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := extractUserFromHeader(r)
		if err != nil {
			util.Unauthorized(w, r, err)
			return
		}

		// At this point user is guaranteed to be authenticated
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
