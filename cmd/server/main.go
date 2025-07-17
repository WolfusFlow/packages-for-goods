package main

import (
	"log"
	"net/http"

	"pfg/internal/auth"
	"pfg/internal/config"
	"pfg/internal/db"
	"pfg/internal/handler"
	"pfg/internal/html"
	"pfg/internal/pack"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	cfg := config.Load()
	auth.InitTokenAuth(cfg.JWTSecret)

	conn, err := db.Connect(cfg.GetPostgresURL())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	repo := db.NewRepository(conn)
	service := pack.NewService(repo)
	jsonHandler := handler.NewHandler(service)

	tmpls, err := html.ParseTemplates()
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	htmlHandler := html.NewHTMLHandler(service, tmpls, cfg)

	auth.RedirectToUnauthorized = htmlHandler.RenderUnauthorized

	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", html.StaticFileServer()))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAdmin)
			r.Get("/packs", jsonHandler.ListPackSizes)
			r.Post("/packs", jsonHandler.AddPackSize)
			r.Delete("/packs", jsonHandler.DeletePackSize)
		})

		r.Post("/calculate", jsonHandler.CalculatePacks)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAdmin)
		r.Get("/packs", htmlHandler.RenderPackList)
		r.Post("/packs/add", htmlHandler.HandleAddPack)
		r.Post("/packs/delete", htmlHandler.HandleDeletePack)
	})

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

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
