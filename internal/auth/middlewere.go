package auth

import (
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

// RequireAdmin checks the "isAdmin" claim in the JWT and blocks for non authorise.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			_, claims, _ := jwtauth.FromContext(r.Context())
			if claims == nil || claims["admin"] != true {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Query().Get("admin") != "1" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please authenticate to access admin area."))
			return
		}

		next.ServeHTTP(w, r)
	})
}
