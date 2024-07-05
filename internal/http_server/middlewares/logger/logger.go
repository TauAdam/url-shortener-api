package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("Custom logger middleware initialized")

		f := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("user_agent", r.UserAgent()),
				slog.String("remote_ip", r.RemoteAddr),
				slog.String("origin", r.Header.Get("Origin")),
				slog.String("referer", r.Referer()),
				slog.String("host", r.Host),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			now := time.Now()
			defer func() {
				entry.Info("Request served",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(now).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(f)
	}
}
