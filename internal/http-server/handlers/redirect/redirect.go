package redirect

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
	resp "url-shortener/internal/http-server/api/response"
	"url-shortener/pkg/logger"
)

//go:generate go run github.com/vektra/mockery/v2 --name=URLGetter
type URLGetter interface {
	GetUrlByAlias(alias string) (string, error)
}

func New(log *slog.Logger, storage URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", "handlers.url.redirect.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		url, err := storage.GetUrlByAlias(alias)
		if err != nil {
			log.Error("error getting url by alias(url not found)", logger.ErrAttr(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("error getting url by alias(url not found)"))
			return
		}
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		log.Info("url redirected", slog.String("alias", alias), slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
