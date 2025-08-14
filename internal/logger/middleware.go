package logger

import (
	"net/http"
	"time"
)

func Middleware(logger *Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
	
		logger.Info("Request started", "", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		})

		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		logger.Info("Request completed", "", map[string]interface{}{
			"status":   recorder.status,
			"duration": duration.Milliseconds(),
			"bytes":    recorder.bytesWritten,
		})
	})
}