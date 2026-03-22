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

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// no token -> proceed without user (public access)
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
			util.BadRequest(w, r, err)
			return
		}

		userID, err := util.ExtractUserID(claims)
		if err != nil {
			util.BadRequest(w, r, err)
			return
		}

		// Add userID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
