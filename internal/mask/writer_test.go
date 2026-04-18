package mask_test

import (
	"bytes"
	"testing"

	"github.com/yourusername/vaultpipe/internal/mask"
	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestWrite_PassthroughWhenNoSecrets(t *testing.T) {
	var buf bytes.Buffer
	r := redact.New("")
	w := mask.NewWriter(&buf, r)

	_, err := w.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestWrite_RedactsSecret(t *testing.T) {
	var buf bytes.Buffer
	r := redact.New("")
	r.Add("s3cr3t")
	w := mask.NewWriter(&buf, r)

	_, err := w.Write([]byte("token=s3cr3t here"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got == "token=s3cr3t here" {
		t.Error("expected secret to be redacted")
	}
}

func TestWrite_ReturnsOriginalLength(t *testing.T) {
	var buf bytes.Buffer
	r := redact.New("")
	r.Add("secret")
	w := mask.NewWriter(&buf, r)

	input := []byte("secret")
	n, err := w.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("expected n=%d, got %d", len(input), n)
	}
}

func TestWrapStreams_BothWritersRedact(t *testing.T) {
	var out, errBuf bytes.Buffer
	secrets := []string{"topsecret"}

	stdout, stderr := mask.WrapStreams(&out, &errBuf, secrets)

	stdout.Write([]byte("value=topsecret"))
	stderr.Write([]byte("err=topsecret"))

	if out.String() == "value=topsecret" {
		t.Error("stdout: expected secret to be redacted")
	}
	if errBuf.String() == "err=topsecret" {
		t.Error("stderr: expected secret to be redacted")
	}
}

func TestWrapStreams_IgnoresEmptySecrets(t *testing.T) {
	var out, errBuf bytes.Buffer
	stdout, _ := mask.WrapStreams(&out, &errBuf, []string{"", "  "})

	stdout.Write([]byte("clean output"))
	if got := out.String(); got != "clean output" {
		t.Errorf("expected %q, got %q", "clean output", got)
	}
}
