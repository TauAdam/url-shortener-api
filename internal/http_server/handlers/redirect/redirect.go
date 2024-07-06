package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"vigilant-octo-spoon/internal/storage"
	"vigilant-octo-spoon/lib/api/response"
	"vigilant-octo-spoon/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.New"

		logger := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if len(alias) == 0 {
			logger.Info("invalid alias")
			render.JSON(w, r, response.Error("redirect: invalid alias"))
			return
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrNotFound) {
			logger.Info("URL for given alias not found", "alias", alias)
			render.JSON(w, r, response.Error(storage.ErrNotFound.Error()))
			return
		}
		if err != nil {
			logger.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		logger.Info("successfully got URL", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
