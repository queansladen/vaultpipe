package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Path      string    `json:"path,omitempty"`
	Key       string    `json:"key,omitempty"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit events as JSON lines.
type Logger struct {
	writer  io.Writer
	encoder *json.Encoder
}

// NewLogger creates a Logger writing to w. Pass nil to use stderr.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Logger{writer: w, encoder: enc}
}

// Log emits an audit event.
func (l *Logger) Log(action, path, key string, success bool, msg string) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Path:      path,
		Key:       key,
		Success:   success,
		Message:   msg,
	}
	return l.encoder.Encode(e)
}

// SecretFetched is a convenience wrapper for a successful secret read.
func (l *Logger) SecretFetched(path, key string) error {
	return l.Log("secret_fetch", path, key, true, "")
}

// SecretFetchFailed logs a failed secret read.
func (l *Logger) SecretFetchFailed(path, key, reason string) error {
	return l.Log("secret_fetch", path, key, false, reason)
}

// ProcessStarted logs that the child process was launched.
func (l *Logger) ProcessStarted(command string) error {
	return l.Log("process_start", "", "", true, command)
}

// ProcessFailed logs that the child process could not be started or exited with an error.
func (l *Logger) ProcessFailed(command, reason string) error {
	return l.Log("process_start", "", "", false, command+": "+reason)
}
