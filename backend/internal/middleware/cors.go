package middleware

import (
	"net/http"

	chiCors "github.com/go-chi/cors"
)

// CORS returns configured CORS middleware.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return chiCors.Handler(chiCors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true, // Required for httpOnly cookies
		MaxAge:           300,
	})
}
