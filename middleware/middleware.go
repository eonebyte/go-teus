package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const (
	RequestIDKey contextKey = "requestID"
	UserIDKey    contextKey = "userID"
)

func isDev() bool {
	env := os.Getenv("APP_ENV")

	switch env {
	case "", "dev", "development", "local":
		return true
	default:
		return false
	}
}

// Logger is a structured request logger middleware.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info("request",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.Status()),
					slog.Duration("latency", time.Since(start)),
					slog.String("request_id", middleware.GetReqID(r.Context())),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

// Recoverer catches panics and returns 500.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {

				if rec := recover(); rec != nil {

					requestID := middleware.GetReqID(
						r.Context(),
					)

					if isDev() {

						logger.Error(
							"panic recovered",
							slog.Any("panic", rec),
							slog.String("method", r.Method),
							slog.String("path", r.URL.Path),
							slog.String("request_id", requestID),
							slog.String(
								"stack_trace",
								string(debug.Stack()),
							),
						)

					} else {

						logger.Error(
							"panic recovered",
							slog.Any("panic", rec),
							slog.String("method", r.Method),
							slog.String("path", r.URL.Path),
							slog.String("request_id", requestID),
						)
					}

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusInternalServerError,
					)

					_, _ = w.Write([]byte(
						`{"success":false,"error":{"code":"INTERNAL_ERROR","message":"internal server error"}}`,
					))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// SetUserID injects a userID into the request context (call after auth validation).
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID retrieves userID from context.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}
