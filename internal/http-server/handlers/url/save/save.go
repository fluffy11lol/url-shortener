package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/http-server/api/response"
	storage2 "url-shortener/internal/storage"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/random"
)

const aliasLength = 5

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, storage URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("error decoding request body: ", logger.ErrAttr(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request body"))
			return
		}

		log.Info("request received", slog.Any("request", req))

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request", logger.ErrAttr(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		if req.Alias == "" {
			req.Alias = random.GetRandomAlias(aliasLength)
		}

		id, err := storage.SaveURL(req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, storage2.ErrUrlAlreadyExist) {
				log.Error("url already exist",
					slog.String("alias", req.Alias), logger.ErrAttr(err))
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, resp.Error("url already exist"))
				return
			}
			log.Error("error saving url: ", logger.ErrAttr(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url saved", slog.Int64("id", id))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Response: resp.Success(),
			Alias:    req.Alias,
		})
	}
}
