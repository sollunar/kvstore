package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapper := &WrapperWriter{
				ResponseWriter: w,
				StatusCode:     http.StatusOK,
			}
			next.ServeHTTP(wrapper, r)

			logger.Info("request completed",
				zap.Int("status", wrapper.StatusCode),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
