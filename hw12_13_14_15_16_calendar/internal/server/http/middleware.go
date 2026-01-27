package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

func loggingMiddleware(logger app.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			remoteIP := strings.Split(r.RemoteAddr, ":")[0]
			logger.Info(fmt.Sprintf(
				"%s [%s] %s %s %s %d %d %v",
				remoteIP,
				start.Format("01/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.RequestURI,
				r.Proto,
				r.Response.StatusCode,
				time.Since(start)/time.Microsecond,
				r.UserAgent(),
			))
		}()
		next.ServeHTTP(w, r)
	})
}
