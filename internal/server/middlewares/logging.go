package middlewares

import (
	"github.com/sparfenov/httpmux/pkg/logger"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	logger logger.Interface
}

func NewLoggingMiddleware(l logger.Interface) *LoggingMiddleware {
	return &LoggingMiddleware{logger: l}
}

func (l *LoggingMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		l.logger.Infof(
			"Request: %s\t[%s]\t%s\t%s\t%s\t%d\t%s",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL,
			r.UserAgent(),
			r.RemoteAddr,
			r.Response.StatusCode,
			time.Since(start).String(),
		)
	})
}
