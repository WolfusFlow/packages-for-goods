package main

import (
	"log"
	"net/http"
	"pfg/internal/config"
	"pfg/internal/db"
	"pfg/internal/handler"
	"pfg/internal/html"
	"pfg/internal/pack"

	htmltpl "html/template" // alias to avoid naming conflict

	"github.com/go-chi/chi/v5"
)

func main() {
	// Load environment config
	cfg := config.Load()

	// Connect to Postgres
	conn, err := db.Connect(cfg.GetPostgresURL())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Initialize repository, service, handlers
	repo := db.NewRepository(conn)
	service := pack.NewService(repo)
	jsonHandler := handler.NewHandler(service)

	// Load standard Go templates (*.html)
	templates, err := htmltpl.ParseGlob("internal/templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	htmlHandler := html.NewHTMLHandler(service, templates)

	// Setup router
	r := chi.NewRouter()

	// API endpoints (JSON)
	r.Route("/api", func(r chi.Router) {
		r.Post("/calculate", jsonHandler.CalculatePacks)
		r.Get("/packs", jsonHandler.ListPackSizes)
		r.Post("/packs", jsonHandler.AddPackSize)
		r.Delete("/packs", jsonHandler.DeletePackSize)
	})

	// Root info message
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Packs for Goods API! Use /api for JSON endpoints or /packs for HTML UI."))
	})

	// HTML UI endpoints
	r.Get("/packs", htmlHandler.RenderPackList)
	r.Post("/packs/add", htmlHandler.HandleAddPack)
	r.Post("/packs/delete", htmlHandler.HandleDeletePack)
	r.Get("/calculate", htmlHandler.RenderCalculateForm)
	r.Post("/calculate", htmlHandler.RenderCalculateForm)

	// Start server
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
