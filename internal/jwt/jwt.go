package jwt

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var Auth *jwtauth.JWTAuth

func InitTokenAuth(secret string) {
	Auth = jwtauth.New("HS256", []byte(secret), nil)
}

// GenerateDevToken returns an admin token string for development or testing purposes.
func GenerateDevToken() string {
	if Auth == nil {
		panic("TokenAuth is not initialized. Call InitTokenAuth(secret) before using JWT operations.")
	}

	_, tokenStr, _ := Auth.Encode(map[string]interface{}{
		"user":    "admin",
		"isAdmin": true,
		"exp":     jwtauth.ExpireIn(30 * time.Minute),
	})
	return tokenStr
}
