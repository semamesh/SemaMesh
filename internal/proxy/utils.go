package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// NotifyViolation sends a structured alert to a Slack Webhook
func NotifyViolation(agentName, reason string) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return // Silently skip if no webhook is configured
	}

	// Create a professional-looking Slack message
	payload := map[string]interface{}{
		"text": fmt.Sprintf("⚠️ *SemaMesh Security Alert*\n*Agent:* %s\n*Violation:* %s\n*Timestamp:* %s",
			agentName, reason, time.Now().Format(time.RFC822)),
	}

	jsonPayload, _ := json.Marshal(payload)

	// Send it asynchronously so we don't slow down the proxy
	go func() {
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Printf("[Error] Failed to send Slack alert: %v\n", err)
		}
	}()
}