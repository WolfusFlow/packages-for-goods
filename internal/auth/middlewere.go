package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// Should be an end var
var TokenAuth = jwtauth.New("HS256", []byte("your-secret"), nil)

// RequireAdmin checks the "isAdmin" claim in the JWT and blocks for non authorise.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		if claims["isAdmin"] != true {
			http.Error(w, "admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GenerateDevToken returns a hardcoded token for development purposes.
func GenerateDevToken() string {
	_, tokenStr, _ := TokenAuth.Encode(map[string]interface{}{
		"user":    "admin",
		"isAdmin": true,
	})
	return tokenStr
}
