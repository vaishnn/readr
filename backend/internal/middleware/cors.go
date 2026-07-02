package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CORS returns a middleware with permissive settings for development and
// locked-down settings for production.
func CORS(env string) func(http.Handler) http.Handler {
	origins := []string{"http://localhost:4200"}
	if env == "production" {
		// In production, set this to the actual frontend domain via env config.
		origins = []string{"https://readr.example.com"}
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
