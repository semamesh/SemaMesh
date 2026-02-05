package proxy

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/semamesh/SemaMesh/pkg/identity"
	"github.com/semamesh/SemaMesh/pkg/sniffer"
)

type SemaHandler struct {
	target        *url.URL
	client        *http.Client
	identityMgr   *identity.Manager
}

func NewSemaHandler(targetURL string, idMgr *identity.Manager) (*SemaHandler, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	return &SemaHandler{
		target:      u,
		identityMgr: idMgr,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (h *SemaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. CAPTURE REQUEST BODY (The Prompt)
	// We read it into memory because we need it twice: once for OpenAI, once for Audit.
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	// Restore the body so the HTTP client can read it again
	r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))

	// 2. Prepare Upstream Request
	outReq, _ := http.NewRequestWithContext(r.Context(), r.Method, h.target.String()+r.URL.Path, bytes.NewBuffer(reqBodyBytes))
	for k, v := range r.Header {
		outReq.Header[k] = v
	}
	outReq.Host = h.target.Host

	// 3. Resolve Identity
	hostIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		hostIP = r.RemoteAddr
	}
	meta, found := h.identityMgr.GetIdentity(hostIP)
	namespace := "unknown-source"
	if found {
		namespace = meta.Namespace
	}

	// 4. Execute Request
	resp, err := h.client.Do(outReq)
	if err != nil {
		http.Error(w, "SemaMesh: Upstream LLM Unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 5. Sniff (Pass reqBodyBytes too!)
	err = sniffer.ProxyAndSniff(w, resp, namespace, reqBodyBytes)
	if err != nil {
		log.Printf("Error during proxy/sniff: %v", err)
	}
}