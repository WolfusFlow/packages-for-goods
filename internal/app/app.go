package app

import (
	"context"
	"log"
	"net/http"

	"pfg/internal/auth"
	"pfg/internal/config"
	"pfg/internal/db"
	"pfg/internal/handler"
	"pfg/internal/html"
	"pfg/internal/pack"
	"pfg/internal/server"
)

type App struct {
	cfg     *config.Config
	dbConn  db.Conn
	httpSrv *http.Server
}

func New(cfg *config.Config) (*App, error) {
	auth.InitTokenAuth(cfg.JWTSecret)

	conn, err := db.Connect(cfg.GetPostgresURL())
	if err != nil {
		return nil, err
	}

	repo := db.NewRepository(conn)
	service := pack.NewService(repo)
	jsonHandler := handler.NewHandler(service)

	tmpls, err := html.ParseTemplates()
	if err != nil {
		return nil, err
	}

	htmlHandler := html.NewHTMLHandler(service, tmpls, cfg)
	auth.RedirectToUnauthorized = htmlHandler.RenderUnauthorized

	router := server.NewRouter(jsonHandler, htmlHandler)

	return &App{
		cfg:    cfg,
		dbConn: conn,
		httpSrv: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: router,
		},
	}, nil
}

func (a *App) Start() error {
	log.Printf("Listening on %s...", a.httpSrv.Addr)
	return a.httpSrv.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if err := a.httpSrv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Closing database connection...")
	return a.dbConn.Close()
}
