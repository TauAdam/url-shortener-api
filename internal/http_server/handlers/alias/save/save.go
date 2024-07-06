package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"vigilant-octo-spoon/lib/api/response"
	"vigilant-octo-spoon/lib/logger/sl"
)

type AliasSaver interface {
	SaveAlias(alias, urlText string) (int64, error)
}
type Request struct {
	Url   string `json:"url" validate:"required,alias"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	Alias string `json:"alias,omitempty"`
	response.Response
}

func New(logger *slog.Logger, aliasSaver AliasSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.alias.save.New"
		logger.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to parse request", sl.Err(err))
			render.JSON(w, r, response.Error("failed to parse request"))
			return
		}
	}
}
