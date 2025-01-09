package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"os"
	"url-shortener/internal/config"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading config: ", err)
	}

	logger := loggerSetup.InitLogger(cfg.Env)
	if logger == nil {
		log.Fatal("error initializing logger")
	}
	_ = logger
	logger.Debug("logger initialized")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("error initializing storage: ", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer storage.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger.Logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: server http
}
