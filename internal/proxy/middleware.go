package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Middleware function signature
type Middleware func(http.Handler) http.Handler

// 1. Logger: Simply prints what is happening
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Waypoint] Request: %s %s", r.Method, r.URL.Host)
		next.ServeHTTP(w, r)
	})
}

// 2. TokenGuard: The "Bank Teller"
// In a real app, this would check Redis. Here, we mock a 1000-token limit.
func WithTokenGuard(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Fetch limit from Environment, default to 1000 if not set
        limitStr := os.Getenv("SEMA_DEFAULT_TOKEN_LIMIT")
        limit := 1000
        if val, err := strconv.Atoi(limitStr); err == nil {
            limit = val
        }

        if r.Method == "POST" {
            // ... (body reading logic) ...
            if tokenCount > limit {
                // NEW: Call our utility function
                NotifyViolation(r.Header.Get("X-Agent-Name"), "Token Quota Exceeded")

                log.Printf("[TokenGuard] BLOCKED: %d tokens exceeds limit", tokenCount)
                http.Error(w, "SemaMesh: Token Quota Exceeded", http.StatusTooManyRequests)
                return
            }
        }
        next.ServeHTTP(w, r)
    })
}

// 3. Chain Helper: Wraps the final handler with all middlewares
func BuildChain(finalHandler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}
	return finalHandler
}