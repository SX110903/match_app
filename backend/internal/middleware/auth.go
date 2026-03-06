package middleware

import (
	"net/http"
	"strings"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/pkg/response"
)

// RequireAuth validates the JWT access token and injects claims into context.
func RequireAuth(jwtSvc auth.IJWTService, blacklist auth.ITokenBlacklist) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Unauthorized(w, "missing authorization token")
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwtSvc.ValidateAccessToken(tokenStr)
			if err != nil {
				response.Unauthorized(w, "invalid token")
				return
			}

			// Check blacklist
			blacklisted, err := blacklist.IsBlacklisted(r.Context(), claims.ID)
			if err != nil || blacklisted {
				response.Unauthorized(w, "token has been revoked")
				return
			}

			ctx := r.Context()
			ctx = auth.ContextWithClaims(ctx, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
