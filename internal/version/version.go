package version

import "fmt"

const (
	// Major version component.
	Major = 0
	// Minor version component.
	Minor = 1
	// Patch version component.
	Patch = 0
	// Pre-release label; empty for stable releases.
	PreRelease = "alpha"
)

// Version returns the full semantic version string.
func Version() string {
	if PreRelease != "" {
		return fmt.Sprintf("%d.%d.%d-%s", Major, Minor, Patch, PreRelease)
	}
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}

// BuildInfo holds optional build-time metadata injected via ldflags.
var BuildInfo = struct {
	Commit    string
	BuildDate string
	BuiltBy   string
}{
	Commit:    "unknown",
	BuildDate: "unknown",
	BuiltBy:   "unknown",
}

// FullString returns a human-readable version string including build metadata.
func FullString() string {
	return fmt.Sprintf(
		"vaultpipe %s (commit=%s built=%s by=%s)",
		Version(),
		BuildInfo.Commit,
		BuildInfo.BuildDate,
		BuildInfo.BuiltBy,
	)
}
