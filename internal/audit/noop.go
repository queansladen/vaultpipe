package audit

// NoopLogger satisfies the same interface but discards all events.
// Useful when audit logging is disabled via config.
type NoopLogger struct{}

func (n *NoopLogger) Log(_, _, _ string, _ bool, _ string) error { return nil }
func (n *NoopLogger) SecretFetched(_, _ string) error            { return nil }
func (n *NoopLogger) SecretFetchFailed(_, _, _ string) error     { return nil }
func (n *NoopLogger) ProcessStarted(_ string) error              { return nil }

// Auditor is the common interface implemented by Logger and NoopLogger.
type Auditor interface {
	Log(action, path, key string, success bool, msg string) error
	SecretFetched(path, key string) error
	SecretFetchFailed(path, key, reason string) error
	ProcessStarted(command string) error
}

// NewAuditor returns a real Logger when enabled, otherwise a NoopLogger.
func NewAuditor(enabled bool, w interface{ Write([]byte) (int, error) }) Auditor {
	if !enabled {
		return &NoopLogger{}
	}
	return NewLogger(w)
}
