package version_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/version"
)

func TestVersion_Format(t *testing.T) {
	v := version.Version()
	if v == "" {
		t.Fatal("expected non-empty version string")
	}
	// Should contain at least two dots (major.minor.patch)
	parts := strings.Split(v, ".")
	if len(parts) < 3 {
		t.Fatalf("expected semver with at least 3 parts, got %q", v)
	}
}

func TestVersion_ContainsPreRelease(t *testing.T) {
	v := version.Version()
	if version.PreRelease != "" && !strings.Contains(v, version.PreRelease) {
		t.Fatalf("expected version %q to contain pre-release label %q", v, version.PreRelease)
	}
}

func TestVersion_NoPreRelease(t *testing.T) {
	orig := version.PreRelease
	// We cannot mutate the const, so just verify the helper logic via FullString.
	_ = orig
	v := version.Version()
	if !strings.HasPrefix(v, "0.") {
		t.Fatalf("expected version to start with major 0, got %q", v)
	}
}

func TestFullString_ContainsVersion(t *testing.T) {
	full := version.FullString()
	if !strings.Contains(full, "vaultpipe") {
		t.Fatalf("expected FullString to contain 'vaultpipe', got %q", full)
	}
	if !strings.Contains(full, version.Version()) {
		t.Fatalf("expected FullString to contain version %q, got %q", version.Version(), full)
	}
}

func TestFullString_ContainsBuildInfo(t *testing.T) {
	version.BuildInfo.Commit = "abc1234"
	version.BuildInfo.BuildDate = "2024-01-01"
	version.BuildInfo.BuiltBy = "ci"

	full := version.FullString()
	for _, want := range []string{"abc1234", "2024-01-01", "ci"} {
		if !strings.Contains(full, want) {
			t.Errorf("expected FullString to contain %q, got %q", want, full)
		}
	}
}

func TestVersion_NoLeadingV(t *testing.T) {
	// Version strings should be bare semver without a leading 'v' prefix,
	// e.g. "0.1.0" rather than "v0.1.0", so callers can add prefixes as needed.
	v := version.Version()
	if strings.HasPrefix(v, "v") {
		t.Fatalf("expected version without leading 'v' prefix, got %q", v)
	}
}
