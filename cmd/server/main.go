package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pfg/internal/auth"
	"pfg/internal/config"
	"pfg/internal/db"
	"pfg/internal/handler"
	"pfg/internal/pack"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.Load()

	pool, err := db.Connect(cfg.GetPostgresURL())
	if err != nil {
		log.Fatal(err)
	}

	repo := db.NewRepository(pool)
	service := pack.NewService(repo)

	// Setup HTTP router
	r := chi.NewRouter()
	h := handler.NewHandler(service)

	r.Post("/pack", h.CalculatePacks)
	r.Get("/packs", h.ListPackSizes)

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.RequireAdmin) // Middleware to protect admin routes
		r.Post("/packs", h.AddPackSize)
		r.Delete("/packs", h.DeletePackSize)
	})

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start HTTP server in background
	go func() {
		log.Printf("Server running on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown on interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
