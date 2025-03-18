package middleware

import (
	"net/http"
	"strings"

	"cyber-rbac/pkg/config"
	"cyber-rbac/pkg/logger"
)

// AuthMiddleware validates Bearer token from the request header
func AuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Check if Authorization header is missing
		if authHeader == "" {
			logger.Logger.Warn("Missing Authorization header")
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// Validate Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Logger.Warn("Invalid Authorization format")
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Verify if the token matches the configured Bearer token
		if tokenParts[1] != cfg.BearerToken {
			logger.Logger.Warn("Unauthorized access attempt with invalid token")
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
