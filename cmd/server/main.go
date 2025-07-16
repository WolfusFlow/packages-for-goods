package main

import (
	"log"
	"net/http"

	"pfg/internal/config"
	"pfg/internal/db"
	"pfg/internal/handler"
	"pfg/internal/html"
	"pfg/internal/pack"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg.GetPostgresURL())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	repo := db.NewRepository(conn)
	service := pack.NewService(repo)
	jsonHandler := handler.NewHandler(service)

	tmpls, err := html.ParseTemplates()

	htmlHandler := html.NewHTMLHandler(service, tmpls)

	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/html/static"))))

	r.Route("/api", func(r chi.Router) {
		r.Post("/calculate", jsonHandler.CalculatePacks)
		r.Get("/packs", jsonHandler.ListPackSizes)
		r.Post("/packs", jsonHandler.AddPackSize)
		r.Delete("/packs", jsonHandler.DeletePackSize)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Packs for Goods API! Use /api for JSON endpoints or /packs for HTML UI."))
	})

	r.Get("/packs", htmlHandler.RenderPackList)
	r.Post("/packs/add", htmlHandler.HandleAddPack)
	r.Post("/packs/delete", htmlHandler.HandleDeletePack)
	r.Get("/calculate", htmlHandler.RenderCalculateForm)
	r.Post("/calculate", htmlHandler.RenderCalculateForm)

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
