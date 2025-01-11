package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	stdLog "log"
	"os"
	"url-shortener/internal/config"
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
	_ = log
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

	// TODO: server http
}
