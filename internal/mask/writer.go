package mask

import (
	"io"
	"strings"

	"github.com/yourusername/vaultpipe/internal/redact"
)

// Writer wraps an io.Writer and redacts secrets before writing.
type Writer struct {
	underlying io.Writer
	redactor   *redact.Redactor
}

// NewWriter returns a Writer that redacts known secrets from all writes.
func NewWriter(w io.Writer, r *redact.Redactor) *Writer {
	return &Writer{
		underlying: w,
		redactor:   r,
	}
}

// Write redacts secrets from p before forwarding to the underlying writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	original := string(p)
	redacted := w.redactor.Redact(original)

	// Report original length to callers so they don't see a short-write error.
	_, err = io.WriteString(w.underlying, redacted)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// WrapStreams returns redacting writers for stdout and stderr.
func WrapStreams(stdout, stderr io.Writer, secrets []string) (io.Writer, io.Writer) {
	r := redact.New("")
	for _, s := range secrets {
		if strings.TrimSpace(s) != "" {
			r.Add(s)
		}
	}
	return NewWriter(stdout, r), NewWriter(stderr, r)
}
