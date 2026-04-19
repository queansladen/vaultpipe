package env

import (
	"os"
	"testing"
)

func TestExpand_PlainString(t *testing.T) {
	e := NewExpander(map[string]string{}, false)
	if got := e.Expand("hello"); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestExpand_ResolvesFromMap(t *testing.T) {
	e := NewExpander(map[string]string{"FOO": "bar"}, false)
	if got := e.Expand("${FOO}"); got != "bar" {
		t.Fatalf("expected 'bar', got %q", got)
	}
}

func TestExpand_UnknownNoFallback(t *testing.T) {
	e := NewExpander(map[string]string{}, false)
	if got := e.Expand("${MISSING}"); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestExpand_FallbackToEnv(t *testing.T) {
	os.Setenv("VP_TEST_EXPAND", "fromenv")
	defer os.Unsetenv("VP_TEST_EXPAND")
	e := NewExpander(map[string]string{}, true)
	if got := e.Expand("${VP_TEST_EXPAND}"); got != "fromenv" {
		t.Fatalf("expected 'fromenv', got %q", got)
	}
}

func TestExpand_MapTakesPrecedenceOverEnv(t *testing.T) {
	os.Setenv("VP_TEST_PREC", "fromenv")
	defer os.Unsetenv("VP_TEST_PREC")
	e := NewExpander(map[string]string{"VP_TEST_PREC": "frommap"}, true)
	if got := e.Expand("${VP_TEST_PREC}"); got != "frommap" {
		t.Fatalf("expected 'frommap', got %q", got)
	}
}

func TestExpandAll_ExpandsValues(t *testing.T) {
	e := NewExpander(map[string]string{"HOST": "localhost"}, false)
	input := []string{"ADDR=${HOST}:8080", "PLAIN=value", "NOEQUALS"}
	got := e.ExpandAll(input)
	if got[0] != "ADDR=localhost:8080" {
		t.Errorf("unexpected: %q", got[0])
	}
	if got[1] != "PLAIN=value" {
		t.Errorf("unexpected: %q", got[1])
	}
	if got[2] != "NOEQUALS" {
		t.Errorf("unexpected: %q", got[2])
	}
}
