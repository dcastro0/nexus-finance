package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dcastro0/nexus-finance/pkg/security"
)

type AuthMiddleware struct {
	TokenService *security.TokenService
}

func NewAuthMiddleware(tokenService *security.TokenService) *AuthMiddleware {
	return &AuthMiddleware{TokenService: tokenService}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := m.TokenService.ValidateToken(parts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "account_id", claims.AccountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
