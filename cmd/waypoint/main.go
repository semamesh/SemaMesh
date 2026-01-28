package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"semamesh.io/internal/proxy"
)

func main() {
	// 1. Define the "Final Destination" Handler (The Reverse Proxy)
	// This function actually sends the traffic to the Internet
	proxyHandler := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// In a real Transparent Proxy, we use the Host header
			// For this MVP, I added OpenAI as the default target but can be modified via env var
			// if the agent didn't specify a host (e.g. transparent redirect)
			targetAddr := os.Getenv("SEMA_TARGET_URL")
            if targetAddr == "" {
                targetAddr = "https://api.openai.com" // Default fallback
            }
            target, _ := url.Parse(targetAddr)

			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			// PRO TIP: OpenAI requires an API Key.
			// Usually the Agent sends it, but the Mesh can inject it here too!
		},
	}

	// 2. Build the Middleware Chain
	// Order: Logger -> TokenGuard -> [The Internet]
	stack := proxy.BuildChain(
		proxyHandler,
		proxy.WithLogging,
		proxy.WithTokenGuard,
	)

	// 3. Start the Server
	// We listen on 15001 because that's where our eBPF redirects traffic
	addr := ":15001"
	log.Printf("SemaMesh Waypoint listening on %s...", addr)

	if err := http.ListenAndServe(addr, stack); err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
}