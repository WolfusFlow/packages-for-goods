package app

import (
	"context"
	"fmt"
	"net/http"

	"pfg/internal/config"
	"pfg/internal/db"
	"pfg/internal/handler"
	"pfg/internal/html"
	"pfg/internal/jwt"
	"pfg/internal/pack"
	"pfg/internal/server"

	"go.uber.org/zap"
)

type App struct {
	cfg     *config.Config
	dbConn  db.Conn
	httpSrv *http.Server
	logger  *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	logger.Info("Initializing application")

	jwt.InitTokenAuth(cfg.JWTSecret)
	logger.Info("JWT auth initialized")

	conn, err := db.Connect(cfg.GetPostgresURL())
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}
	logger.Info("Database connection established")

	repo := db.NewRepository(conn)
	service := pack.NewService(repo)

	jsonHandler := handler.NewHandler(service, logger)

	tmpls, err := html.ParseTemplates()
	if err != nil {
		logger.Error("Failed to parse templates", zap.Error(err))
		return nil, err
	}
	logger.Info("Templates parsed successfully")
	for _, tmpl := range tmpls.Templates() {
		fmt.Println("Loaded template:", tmpl.Name())
	}

	htmlHandler := html.NewHTMLHandler(service, tmpls, cfg, logger)

	router := server.NewRouter(jsonHandler, htmlHandler, logger)

	app := &App{
		cfg:    cfg,
		dbConn: conn,
		logger: logger,
		httpSrv: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: router,
		},
	}

	logger.Info("Application initialized", zap.String("port", cfg.Port))
	return app, nil
}

func (a *App) Start() error {
	a.logger.Info("Starting HTTP server", zap.String("addr", a.httpSrv.Addr))
	return a.httpSrv.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down server gracefully")

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		a.logger.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	a.logger.Info("Closing database connection")
	if err := a.dbConn.Close(); err != nil {
		a.logger.Error("Failed to close DB", zap.Error(err))
		return err
	}

	a.logger.Info("Shutdown complete")
	return nil
}
