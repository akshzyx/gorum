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

		if !strings.HasPrefix(authHeader, "Bearer ") {
			util.BadRequest(w, r, errors.New("invalid authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := util.ValidateJWT(tokenStr)
		if err != nil {
			// Token exists but is invalid → unauthorized
			util.Unauthorized(w, r, err)
			return
		}

		userID, err := util.ExtractUserID(claims)
		if err != nil {
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
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Missing token → reject immediately
			util.Unauthorized(w, r, errors.New("missing authorization token"))
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			util.BadRequest(w, r, errors.New("invalid authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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
