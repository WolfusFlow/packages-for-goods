package auth

import (
	"net/http"
	"strings"

	"pfg/internal/jwt"
)

func RequireAdminOnly(r *http.Request) bool {
	isAdmin := false

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tok, err := jwt.Auth.Decode(tokenStr); err == nil {
			claims := tok.PrivateClaims()
			isAdmin, _ = claims["isAdmin"].(bool)
		}
	}

	if !isAdmin {
		if cookie, err := r.Cookie("admin_token"); err == nil {
			if tok, err := jwt.Auth.Decode(cookie.Value); err == nil {
				claims := tok.PrivateClaims()
				isAdmin, _ = claims["isAdmin"].(bool)
			}
		}
	}

	return isAdmin
}

func RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdminOnly(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
