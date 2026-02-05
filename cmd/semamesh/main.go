package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/semamesh/SemaMesh/pkg/audit"
	"github.com/semamesh/SemaMesh/pkg/identity"
	"github.com/semamesh/SemaMesh/pkg/proxy"
)

func main() {
	// 1. Config Flags
	// Default to empty string (""). This tells the K8s client to use "In-Cluster Config" (Service Account).
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

	devMode := flag.Bool("dev", true, "Run in local dev mode (mock identity)")
	flag.Parse()

	// 2. Initialize Audit Logging
	audit.Init()

	// 3. Initialize Identity System
	idManager := identity.NewManager(*devMode)
	if !*devMode {
		log.Println("🔌 Attempting to connect to Kubernetes Cluster...")
		if err := idManager.StartWatcher(*kubeconfig); err != nil {
			log.Fatalf("Failed to start K8s Watcher: %v", err)
		}
	} else {
		log.Println("⚠️ Running in DEV MODE. Kubernetes connection skipped.")
	}

	// 4. Initialize Handler
	semaHandler, err := proxy.NewSemaHandler("https://api.openai.com", idManager)
	if err != nil {
		log.Fatalf("Failed to initialize SemaHandler: %v", err)
	}

    // 5. Start Metrics
    go func() {
       log.Println("📊 Starting Metrics Server on :9090/metrics")
       http.Handle("/metrics", promhttp.Handler())
       // REMOVE THE DUPLICATE LINE HERE
       if err := http.ListenAndServe(":9090", nil); err != nil {
          log.Fatalf("Metrics Server failed: %v", err)
       }
    }()

	// 6. Start Proxy
	log.Println("☁️ SemaMesh AI Proxy active on :8080")
	if *devMode {
		log.Println("🛡️ Identity Awareness: LOCAL (Mock)")
	} else {
		log.Println("🛡️ Identity Awareness: KUBERNETES (Real)")
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: semaHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Proxy Server failed: %v", err)
	}
}

# End of file