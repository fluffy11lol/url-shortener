package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	stdLog "log"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	deleteh "url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/storage/sqlite"
	"url-shortener/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		stdLog.Fatal("error loading config: ", err)
	}

	log := logger.InitLogger(cfg.Env)
	if log == nil {
		stdLog.Fatal("error initializing log")
	}
	log.Info("starting url-shortener app", slog.String("env", cfg.Env))
	log.Debug("logger initialized")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error initializing storage: ", log.ErrAttr(err))
		os.Exit(1)
	}
	defer storage.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log.Logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log.Logger, storage))
	router.Get("/{alias}", redirect.New(log.Logger, storage))
	router.Delete("/url/{alias}", deleteh.New(log.Logger, storage))
	server := http.Server{
		Addr:         cfg.HttpServer.Host + ":" + cfg.HttpServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	log.Info("server started", slog.String("address", cfg.HttpServer.Host+":"+cfg.HttpServer.Port))
	if err = server.ListenAndServe(); err != nil {
		log.Error("error starting server: ", log.ErrAttr(err))
		os.Exit(1)
	}
	log.Info("server stopped")
}
