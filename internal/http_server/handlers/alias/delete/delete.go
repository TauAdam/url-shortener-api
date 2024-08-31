package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/tauadam/url-shortener-api/internal/storage"
	"github.com/tauadam/url-shortener-api/lib/api/response"
	"github.com/tauadam/url-shortener-api/lib/logger/sl"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2 --name=ShortcutDeleter
type ShortcutDeleter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, shortcutDeleter ShortcutDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.shortcut.delete.New"

		logger := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if len(alias) == 0 {
			logger.Info("invalid alias")
			render.JSON(w, r, response.Fail("redirect: invalid alias"))
			return
		}

		err := shortcutDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				logger.Info("URL for given alias not found", "alias", alias)
				render.JSON(w, r, response.Fail(storage.ErrNotFound.Error()))
				return
			}
			logger.Error("failed to delete shortcut", sl.Err(err))
			render.JSON(w, r, response.Fail("internal error"))
			return
		}

		logger.Info("successfully deleted shortcut", "alias", alias)
		response.Success()
	}
}
