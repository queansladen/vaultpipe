package config

import "testing"

func TestResolveTruncate_NilConfig(t *testing.T) {
	tc := ResolveTruncate(nil)
	if tc.MaxBytes != 0 {
		t.Fatalf("expected 0, got %d", tc.MaxBytes)
	}
	if tc.Suffix != "" {
		t.Fatalf("expected empty suffix, got %q", tc.Suffix)
	}
}

func TestResolveTruncate_Defaults(t *testing.T) {
	cfg := &Config{}
	tc := ResolveTruncate(cfg)
	if tc.MaxBytes != 0 {
		t.Fatalf("expected 0 (unlimited), got %d", tc.MaxBytes)
	}
}

func TestResolveTruncate_CustomValues(t *testing.T) {
	cfg := &Config{
		Truncate: TruncateConfig{
			MaxBytes: 256,
			Suffix:   "...",
		},
	}
	tc := ResolveTruncate(cfg)
	if tc.MaxBytes != 256 {
		t.Fatalf("expected 256, got %d", tc.MaxBytes)
	}
	if tc.Suffix != "..." {
		t.Fatalf("expected '...', got %q", tc.Suffix)
	}
}
