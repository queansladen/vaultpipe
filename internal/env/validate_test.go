package env

import (
	"strings"
	"testing"
)

func TestValidate_ValidPairs(t *testing.T) {
	v := NewValidator(false)
	pairs := []string{"FOO=bar", "BAZ=qux"}
	if err := v.Validate(pairs); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingSeparator(t *testing.T) {
	v := NewValidator(false)
	err := v.Validate([]string{"NOSEPARATOR"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
	if !strings.Contains(err.Error(), "missing '=' separator") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	v := NewValidator(false)
	err := v.Validate([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if !strings.Contains(err.Error(), "empty key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_EmptyValue_NotAllowed(t *testing.T) {
	v := NewValidator(false)
	err := v.Validate([]string{"FOO="})
	if err == nil {
		t.Fatal("expected error for empty value when not allowed")
	}
	if !strings.Contains(err.Error(), "empty value") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_EmptyValue_Allowed(t *testing.T) {
	v := NewValidator(true)
	if err := v.Validate([]string{"FOO="}); err != nil {
		t.Fatalf("expected no error when empty values allowed, got: %v", err)
	}
}

func TestValidate_Duplicate(t *testing.T) {
	v := NewValidator(true)
	err := v.Validate([]string{"FOO=bar", "BAZ=qux", "FOO=dup"})
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestIsValidationError(t *testing.T) {
	v := NewValidator(false)
	err := v.Validate([]string{"=bad"})
	if !IsValidationError(err) {
		t.Errorf("expected IsValidationError to return true")
	}
}

func TestValidate_Empty(t *testing.T) {
	v := NewValidator(false)
	if err := v.Validate(nil); err != nil {
		t.Fatalf("expected no error for nil slice, got: %v", err)
	}
}
