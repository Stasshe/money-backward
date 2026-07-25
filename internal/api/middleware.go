package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoggingMiddleware struct {
	logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (lm *LoggingMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		lm.logger.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// Log response
		duration := time.Since(start)
		lm.logger.Printf("[%s] %s %s - Status: %d (%v)", r.Method, r.RequestURI, r.RemoteAddr, wrapped.statusCode, duration)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

type AuthMiddleware struct {
	// Stub for authentication validation
	// In a real implementation, this would validate tokens, API keys, etc.
	apiKeys map[string]bool
}

func NewAuthMiddleware(validKeys []string) *AuthMiddleware {
	keys := make(map[string]bool)
	for _, key := range validKeys {
		keys[key] = true
	}
	return &AuthMiddleware{apiKeys: keys}
}

func (am *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Support Bearer token format
		var token string
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = authHeader
		}

		if !am.isValidToken(token) {
			http.Error(w, "invalid authorization token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) isValidToken(token string) bool {
	if len(am.apiKeys) == 0 {
		// If no keys configured, allow all (stub mode)
		return true
	}
	return am.apiKeys[token]
}

type RateLimitMiddleware struct {
	requestsPerMinute int
	tracker           map[string][]time.Time
}

func NewRateLimitMiddleware(requestsPerMinute int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		requestsPerMinute: requestsPerMinute,
		tracker:           make(map[string][]time.Time),
	}
}

func (rlm *RateLimitMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		now := time.Now()
		oneMinuteAgo := now.Add(-time.Minute)

		// Clean old requests
		if times, exists := rlm.tracker[ip]; exists {
			var recent []time.Time
			for _, t := range times {
				if t.After(oneMinuteAgo) {
					recent = append(recent, t)
				}
			}
			rlm.tracker[ip] = recent
		}

		requests := rlm.tracker[ip]
		if len(requests) >= rlm.requestsPerMinute {
			http.Error(w, fmt.Sprintf("rate limit exceeded: %d requests per minute", rlm.requestsPerMinute), http.StatusTooManyRequests)
			return
		}

		rlm.tracker[ip] = append(requests, now)
		next.ServeHTTP(w, r)
	})
}

func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
