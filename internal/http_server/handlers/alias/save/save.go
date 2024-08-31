package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/tauadam/url-shortener-api/internal/storage"
	"github.com/tauadam/url-shortener-api/lib/api/response"
	"github.com/tauadam/url-shortener-api/lib/logger/sl"
	"github.com/tauadam/url-shortener-api/lib/random"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2 --name=ShortcutSaver
type ShortcutSaver interface {
	SaveShortcut(urlText, alias string) (int64, error)
}
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	Alias string `json:"alias,omitempty"`
	response.Response
}

const AliasLength = 5

func New(logger *slog.Logger, shortcutSaver ShortcutSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.alias.save.New"
		logger := logger.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to parse request", sl.Err(err))
			render.JSON(w, r, response.Fail("failed to parse request"))
			return
		}

		logger.Info("successfully parsed request", slog.Any("req", req))
		if err := validator.New().Struct(req); err != nil {
			logger.Error("failed to validate request", sl.Err(err))

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				render.JSON(w, r, response.ValidationError(validationErrors))
			}
			return
		}

		alias := req.Alias
		if len(alias) == 0 {
			alias = random.NewString(AliasLength)
		}

		id, err := shortcutSaver.SaveShortcut(req.URL, req.Alias)
		if errors.Is(err, storage.ErrAlreadyExists) {
			logger.Info("URL already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Fail("URL already exists"))
			return
		}
		if err != nil {
			logger.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, response.Fail("failed to save url"))
			return
		}
		logger.Info("successfully saved url", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Alias:    alias,
			Response: response.Success(),
		})
	}
}
