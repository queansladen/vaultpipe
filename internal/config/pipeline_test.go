package config

import (
	"testing"
)

func minimalConfig() *Config {
	return &Config{
		VaultAddress: "http://127.0.0.1:8200",
		Command:      []string{"env"},
	}
}

func TestBuildPipeline_NilSubconfigs_ReturnsEmptyPipeline(t *testing.T) {
	cfg := minimalConfig()
	p, err := BuildPipeline(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages for empty config, got %d", p.Len())
	}
}

func TestBuildPipeline_WithCaser_AddsStage(t *testing.T) {
	cfg := minimalConfig()
	cfg.Env.Caser = &CaserConfig{Mode: "upper", Target: "key"}
	p, err := BuildPipeline(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestBuildPipeline_WithRename_AddsStage(t *testing.T) {
	cfg := minimalConfig()
	cfg.Env.Rename = []RenameRule{{From: "OLD", To: "NEW"}}
	p, err := BuildPipeline(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestBuildPipeline_MultipleOptions_CorrectStageCount(t *testing.T) {
	cfg := minimalConfig()
	cfg.Env.Caser = &CaserConfig{Mode: "upper", Target: "key"}
	cfg.Env.Rename = []RenameRule{{From: "A", To: "B"}}
	cfg.Env.Tag = &TagConfig{Prefix: "VP_"}
	p, err := BuildPipeline(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 3 {
		t.Fatalf("expected 3 stages, got %d", p.Len())
	}
}

func TestBuildPipeline_RunsStagesInOrder(t *testing.T) {
	cfg := minimalConfig()
	// Upper-case keys first, then rename the upper-cased key.
	cfg.Env.Caser = &CaserConfig{Mode: "upper", Target: "key"}
	cfg.Env.Rename = []RenameRule{{From: "MYKEY", To: "RENAMED"}}

	p, err := BuildPipeline(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := p.Run([]string{"mykey=hello"})
	if err != nil {
		t.Fatalf("pipeline run error: %v", err)
	}
	if len(out) != 1 || out[0] != "RENAMED=hello" {
		t.Fatalf("expected RENAMED=hello, got %v", out)
	}
}
