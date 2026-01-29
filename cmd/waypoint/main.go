package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/semamesh/semamesh/internal/proxy"
)

func main() {
	// 1. Configuration
	target := os.Getenv("TARGET_LLM")
	if target == "" {
		// Default fallback for local testing
		target = "http://mock-llm-service.default.svc.cluster.local:8080"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	// 2. Setup Reverse Proxy with Debugging & Timeout
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	// FIX: Log errors so we know if the backend is unreachable
	reverseProxy.ErrorLog = log.New(os.Stderr, "PROXY_ERR: ", 0)

	// FIX: Use a custom Transport to disable Keep-Alives and enforce timeouts
	reverseProxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Don't hang forever
			KeepAlive: 30 * time.Second,
		}).DialContext,
		DisableKeepAlives: true, // Prevents "hanging" connections
	}

	// FIX: Rewrite the Host header to match the target (Critical for Nginx/Cloud endpoints)
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
		// Debug Log: Prove we are sending the request
		log.Printf("PROXY_DEBUG: Forwarding to %s%s", targetURL.Host, req.URL.Path)
	}

	// 3. Apply Middleware (Intent Analysis)
	finalHandler := proxy.IntentMiddleware(reverseProxy)

	// 4. Start Server
	log.Printf("🚀 Waypoint Proxy starting on :%s forwarding to %s", port, target)
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}