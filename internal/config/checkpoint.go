package config

// CheckpointConfig controls whether pipeline stage checkpoints are captured
// and emitted to the audit log for diff inspection.
type CheckpointConfig struct {
	Enabled bool     `yaml:"enabled"`
	Stages  []string `yaml:"stages"`
}

// defaultStages are the pipeline stage names captured when no explicit list
// is provided.
var defaultStages = []string{"inherit", "resolve", "transform", "final"}

// ResolveCheckpoint returns a normalised CheckpointConfig.
// If cfg is nil or Enabled is false, an empty config is returned.
func ResolveCheckpoint(cfg *CheckpointConfig) CheckpointConfig {
	if cfg == nil || !cfg.Enabled {
		return CheckpointConfig{}
	}
	if len(cfg.Stages) == 0 {
		return CheckpointConfig{Enabled: true, Stages: defaultStages}
	}
	stages := make([]string, len(cfg.Stages))
	copy(stages, cfg.Stages)
	return CheckpointConfig{Enabled: true, Stages: stages}
}
