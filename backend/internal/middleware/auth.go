package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/readr/api/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

// Authenticate validates the Bearer token in the Authorization header.
// Rejects the request with 401 if the token is missing, invalid, or not an access token.
func Authenticate(accessSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateToken(strings.TrimPrefix(header, "Bearer "), accessSecret)
			if err != nil || claims.Type != auth.TokenTypeAccess {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromCtx extracts the authenticated user's ID from the request context.
// Returns an empty string and false if not set (i.e. on unprotected routes).
func UserIDFromCtx(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok && id != ""
}
