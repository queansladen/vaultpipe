package config

import (
	"fmt"

	"github.com/your-org/vaultpipe/internal/env"
)

// PassthroughConfig holds the raw YAML representation for env passthrough.
type PassthroughConfig struct {
	Mode   string   `yaml:"mode"`   // all | none | filtered
	Allow  []string `yaml:"allow"`  // prefix or exact rules (filtered mode)
}

// ResolvePassthroughMode converts the string mode from config into the
// typed PassthroughMode understood by the env package.
func ResolvePassthroughMode(cfg PassthroughConfig) (env.PassthroughMode, *env.Filter, error) {
	switch cfg.Mode {
	case "", "all":
		return env.PassthroughAll, nil, nil
	case "none":
		return env.PassthroughNone, nil, nil
	case "filtered":
		f := env.NewFilter(cfg.Allow)
		return env.PassthroughFiltered, f, nil
	default:
		return 0, nil, fmt.Errorf("unknown passthrough mode %q: must be all, none, or filtered", cfg.Mode)
	}
}
