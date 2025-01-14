package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/http-server/api/response"
	storage2 "url-shortener/internal/storage"
	"url-shortener/pkg/logger"
)

//go:generate go run github.com/vektra/mockery/v2 --name=URLDeleter
type URLDeleter interface {
	DeleteAlias(alias string) error
}

func New(log *slog.Logger, storage URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", "handlers.url.delete.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		err := storage.DeleteAlias(alias)
		if err != nil {
			if errors.Is(err, storage2.ErrUrlNotFound) {
				log.Error("url not found", logger.ErrAttr(err))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("url not found"))
				return
			}
			log.Error("error deleting url", logger.ErrAttr(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("error deleting url"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp.Success())

	}
}
