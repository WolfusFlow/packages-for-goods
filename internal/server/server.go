package server

import (
	"net/http"
	"time"

	"pfg/internal/auth"
	"pfg/internal/handler"
	"pfg/internal/html"
	"pfg/internal/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

func NewRouter(jsonHandler *handler.Handler, htmlHandler *html.HTMLHandler, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(httprate.LimitByIP(100, 5*time.Minute))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Request", zap.String("method", r.Method), zap.String("url", r.URL.Path))
			next.ServeHTTP(w, r)
		})
	})

	r.Handle("/static/*", http.StripPrefix("/static/", html.StaticFileServer()))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt.Auth))
		r.Use(jwtauth.Authenticator(jwt.Auth))
		r.Use(auth.RequireToken)

		r.Group(func(r chi.Router) {
			r.Get("/packs", jsonHandler.ListPackSizes)
			r.Post("/packs", jsonHandler.AddPackSize)
			r.Delete("/packs", jsonHandler.DeletePackSize)
		})

		r.Post("/calculate", jsonHandler.CalculatePacks)
	})

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !auth.RequireAdminOnly(r) {
					htmlHandler.RenderUnauthorized(w, r)
					return
				}
				next.ServeHTTP(w, r)
			})
		})

		r.Get("/packs", htmlHandler.RenderPackList)
		r.Post("/packs/add", htmlHandler.HandleAddPack)
		r.Post("/packs/delete", htmlHandler.HandleDeletePack)
	})

	// Public routes
	r.Get("/", htmlHandler.RenderWelcomePage)
	r.Get("/calculate", htmlHandler.RenderCalculateForm)
	r.Post("/calculate", htmlHandler.RenderCalculateForm)
	r.Get("/login", htmlHandler.RenderLoginForm)
	r.Post("/login", htmlHandler.HandleLoginPost)

	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	return r
}
