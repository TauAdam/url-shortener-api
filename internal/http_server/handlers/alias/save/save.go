package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"vigilant-octo-spoon/lib/api/response"
	"vigilant-octo-spoon/lib/logger/sl"
	"vigilant-octo-spoon/lib/random"
)

type AliasSaver interface {
	SaveAlias(alias, urlText string) (int64, error)
}
type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	Alias string `json:"alias,omitempty"`
	response.Response
}

const AliasLength = 5

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
		logger.Info("successfully parsed request", slog.Any("req", req))
		if err := validator.New().Struct(req); err != nil {
			logger.Error("failed to validate request", sl.Err(err))

			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			render.JSON(w, r, response.ValidationError(validationErrors))
			return
		}

		alias := req.Alias
		if len(alias) == 0 {
			alias = random.NewString(AliasLength)
		}

	}
}
