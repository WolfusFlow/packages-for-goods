package server

import (
	"net/http"

	"pfg/internal/auth"
	"pfg/internal/handler"
	"pfg/internal/html"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter(jsonHandler *handler.Handler, htmlHandler *html.HTMLHandler) http.Handler {
	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", html.StaticFileServer()))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Use(auth.RequireToken)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAdmin)

			r.Get("/packs", jsonHandler.ListPackSizes)
			r.Post("/packs", jsonHandler.AddPackSize)
			r.Delete("/packs", jsonHandler.DeletePackSize)
		})

		r.Post("/calculate", jsonHandler.CalculatePacks)
	})

	// HTML pages (some protected)
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAdmin)
		r.Get("/packs", htmlHandler.RenderPackList)
		r.Post("/packs/add", htmlHandler.HandleAddPack)
		r.Post("/packs/delete", htmlHandler.HandleDeletePack)
	})

	// Public pages
	r.Get("/", htmlHandler.RenderWelcomePage)
	r.Get("/calculate", htmlHandler.RenderCalculateForm)
	r.Post("/calculate", htmlHandler.RenderCalculateForm)
	r.Get("/login", htmlHandler.RenderLoginForm)
	r.Post("/login", htmlHandler.HandleLoginPost)

	// Logout
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
