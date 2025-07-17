package auth

import (
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

// RequireAdmin checks the "isAdmin" claim in the JWT and blocks for non-authorized.
var RedirectToUnauthorized http.HandlerFunc

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API: Reject with JSON
		if strings.HasPrefix(r.URL.Path, "/api/") {
			_, claims, _ := jwtauth.FromContext(r.Context())
			if claims == nil || claims["isAdmin"] != true {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// HTML: Check cookie
		cookie, err := TokenAuthCookie(r)
		if err != nil {
			RedirectToUnauthorized(w, r)
			return
		}

		token, err := TokenAuth.Decode(cookie.Value)
		if err != nil {
			RedirectToUnauthorized(w, r)
			return
		}

		claims := token.PrivateClaims()
		if claims["isAdmin"] != true {
			RedirectToUnauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Extracts the cookie safely
func TokenAuthCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("admin_token")
}
