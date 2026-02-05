package sniffer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/semamesh/SemaMesh/pkg/audit"
	"github.com/semamesh/SemaMesh/pkg/metrics"
)

// --- Structures ---
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type PartialResponse struct {
	Model   string       `json:"model"`
	Usage   *OpenAIUsage `json:"usage"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type RequestPayload struct {
	Messages []struct {
		Content string `json:"content"`
	} `json:"messages"`
}

// --- Logic ---

func ProxyAndSniff(w http.ResponseWriter, upstreamResp *http.Response, namespace string, reqBody []byte) error {
    // 1. Record the Request immediately 🚦
    statusStr := strconv.Itoa(upstreamResp.StatusCode)
    metrics.RequestsTotal.WithLabelValues(namespace, statusStr).Inc()

    // 2. Copy Headers
    for k, v := range upstreamResp.Header {
       w.Header()[k] = v
    }
    w.WriteHeader(upstreamResp.StatusCode)

    tapBuffer := bytes.NewBuffer(make([]byte, 0, 4096))
    splitStream := io.TeeReader(upstreamResp.Body, tapBuffer)

    _, err := io.Copy(w, splitStream)
    if err != nil {
       return err
    }

    go analyze(tapBuffer.Bytes(), reqBody, namespace)

    return nil
}

// ... imports ...

func analyze(respData []byte, reqBody []byte, namespace string) {
    if len(respData) == 0 {
       return
    }

    // Parse Response (Partial)
    var resp PartialResponse
    json.Unmarshal(respData, &resp)

    // Extract Prompt (Best Effort)
    var req RequestPayload
    promptText := "unknown"
    if err := json.Unmarshal(reqBody, &req); err == nil && len(req.Messages) > 0 {
       promptText = req.Messages[len(req.Messages)-1].Content
    }

    // Default Values
    model := resp.Model
    if model == "" {
       model = "error-response"
    }
    completionText := ""
    tokens := 0
    cost := 0.0

    // CASE 1: Success (Usage Data Exists)
    if resp.Usage != nil {
       tokens = resp.Usage.TotalTokens
       cost = estimateCost(resp.Model, resp.Usage.PromptTokens, resp.Usage.CompletionTokens)

       if len(resp.Choices) > 0 {
          completionText = resp.Choices[0].Message.Content
       }

       // Update Metrics
       metrics.TokenCounter.WithLabelValues("prompt", resp.Model, namespace).Add(float64(resp.Usage.PromptTokens))
       metrics.TokenCounter.WithLabelValues("completion", resp.Model, namespace).Add(float64(resp.Usage.CompletionTokens))
       metrics.CostCounter.WithLabelValues(resp.Model, namespace).Add(cost)
    } else {
       // CASE 2: Error / No Usage Data 🚨
       // We still want to log this!
       completionText = "Request Failed / No Token Usage"
    }

    // Always Submit to Audit Log
    audit.Submit(audit.LogEntry{
       Timestamp:      time.Now(),
       Namespace:      namespace,
       Model:          model,
       PromptText:     promptText,
       CompletionText: completionText,
       TotalTokens:    tokens,
       CostEst:        cost,
    })
}

func estimateCost(model string, promptTokens, completionTokens int) float64 {
	var promptPrice, completionPrice float64

	switch {
	case strings.Contains(model, "gpt-4"):
		promptPrice = 30.0 / 1000000.0
		completionPrice = 60.0 / 1000000.0
	case strings.Contains(model, "gpt-3.5") || strings.Contains(model, "gpt-4o"):
		promptPrice = 0.50 / 1000000.0
		completionPrice = 1.50 / 1000000.0
	default:
		return 0.0
	}
	return (float64(promptTokens) * promptPrice) + (float64(completionTokens) * completionPrice)
}