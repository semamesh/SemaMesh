package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// IntentMiddleware intercepts requests to check for policy violations
func IntentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Read the body to inspect the prompt
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// IMPORTANT: Restore the body so the next handler (the Mock LLM) can read it too
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Normalize text for analysis
		bodyStr := strings.ToLower(string(bodyBytes))

		// --- 2. LOGIC: Intent Detection ---
		// If the prompt contains "delete", we trigger a policy violation
		if strings.Contains(bodyStr, "delete") {
			log.Println("INTENT_ANALYSIS: Destructive intent detected! Triggering PAUSE.")

			// Create a file to simulate the stateful pause (for the smoke test to see)
			// In a real scenario, this would call the Controller API
			_ = os.WriteFile("/tmp/semamesh-violation", []byte("violation"), 0644)

			http.Error(w, "SemaMesh Policy Violation: Agent Paused", http.StatusForbidden)
			return
		}

		// --- 3. LOGIC: Token Counting ---
		// Simple approximation: count words as tokens
		tokenCount := len(strings.Fields(bodyStr))

		// Get limit from environment or use default
		limitStr := os.Getenv("TOKEN_LIMIT")
		if limitStr == "" {
			limitStr = "1000"
		}

		limit, _ := strconv.Atoi(limitStr)

		// Check Quota
		if tokenCount > limit {
			log.Printf("QUOTA_VIOLATION: Request size %d exceeds limit %d", tokenCount, limit)
			http.Error(w, "SemaMesh: Token Quota Exceeded", http.StatusTooManyRequests)
			return
		}

		// Log success for the smoke test to grep
		log.Printf("INTENT_ANALYSIS: Safe. Tokens: %d", tokenCount)

		// 4. Pass request to the actual LLM
		next.ServeHTTP(w, r)
	})
}