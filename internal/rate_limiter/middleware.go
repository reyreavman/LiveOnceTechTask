package ratelimiter

import "net/http"

func Middleware(rateLimiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := getClientIP(r)
		if !rateLimiter.Allow(clientID) {
			w.Header().Set("Retry-After", rateLimiter.window.String())
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}