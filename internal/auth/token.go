package auth

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func InitTokenAuth(secret string) {
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

// GenerateDevToken returns a hardcoded token for development purposes.
func GenerateDevToken() string {
	_, tokenStr, _ := TokenAuth.Encode(map[string]interface{}{
		"user":    "admin",
		"isAdmin": true,
		"exp":     jwtauth.ExpireIn(30 * time.Minute),
	})
	return tokenStr
}
