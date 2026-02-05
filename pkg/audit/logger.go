package audit

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Global channel to receive logs
var logChannel chan LogEntry

// LogEntry defines the structure of our audit JSON
type LogEntry struct {
	Timestamp      time.Time `json:"timestamp"`
	Namespace      string    `json:"namespace"`
	Model          string    `json:"model"`
	PromptText     string    `json:"prompt_text"`
	CompletionText string    `json:"completion_text"`
	TotalTokens    int       `json:"total_tokens"`
	CostEst        float64   `json:"cost_usd"`
}

// Init sets up the file writer in a background goroutine
func Init() {
	// 1. Define the absolute path
	filePath := "/var/log/semamesh/audit.log"

	// 2. Open the file (Append mode, Create if missing, Write only)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// If we can't write to /var/log, fall back to stdout or local dir
		log.Printf("❌ CRITICAL: Failed to open audit log at %s: %v", filePath, err)
		log.Println("⚠️ Falling back to './audit.log'")
		file, err = os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("❌ FATAL: Could not create audit log anywhere: %v", err)
		}
	} else {
		log.Printf("✅ Audit Logger active. Writing to: %s", filePath)
	}

	logChannel = make(chan LogEntry, 100)

	// 3. Start the worker
	go func() {
		defer file.Close()
		encoder := json.NewEncoder(file)
		for entry := range logChannel {
			if err := encoder.Encode(entry); err != nil {
				log.Printf("❌ Error writing audit entry: %v", err)
			}
		}
	}()
}

// Submit sends a log entry to the worker (Non-blocking)
func Submit(entry LogEntry) {
	select {
	case logChannel <- entry:
		// Success
	default:
		log.Println("⚠️ Audit Channel full! Dropping log entry.")
	}
}

// Helper for the sniffer to call
func LogAccess(namespace string, req []byte, resp []byte, tokens int) {
	// This is a legacy helper if you still use it in older code.
	// We prefer using 'Submit' directly with structured data.
}