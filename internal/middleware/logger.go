package middleware

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"my-calendar/internal/logger"
	"my-calendar/internal/metrics"
	"my-calendar/internal/clickhouse"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(rw, r)

		metrics.RequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(rw.statusCode),
		).Inc()

		duration := time.Since(start).Seconds()

		metrics.RequestDuration.
		    WithLabelValues(
		        r.Method,
		        r.URL.Path,
		    ).Observe(duration)

		go func() {
		    err := clickhouse.SaveLog(
		        clickhouse.RequestLog{
		            Timestamp:  time.Now(),
		            Method:     r.Method,
		            Path:       r.URL.Path,
		            Status:     uint16(rw.statusCode),
		            DurationMs: uint32(time.Since(start).Milliseconds()),
		        },
		    )		
		    if err != nil {
		        logger.L().Error(
		            "failed to save request log",
		            zap.Error(err),
		        )
		    }
		}()
		
		logger.L().Info(
			"http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", time.Since(start)),
		)
	})
}