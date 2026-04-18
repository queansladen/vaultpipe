package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/exec"
	"github.com/yourusername/vaultpipe/internal/resolve"
	"github.com/yourusername/vaultpipe/internal/signal"
	"github.com/yourusername/vaultpipe/internal/vault"
	"github.com/yourusername/vaultpipe/internal/version"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "vaultpipe",
	Short:   "Stream Vault secrets into a process environment",
	Version: version.FullString(),
	RunE:    run,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (required)")
	_ = rootCmd.MarkPersistentFlagRequired("config")
}

func run(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	auditor := audit.NewAuditor(cfg.Audit.Enabled, os.Stderr)

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	mappings, err := config.ToMappings(cfg)
	if err != nil {
		return fmt.Errorf("parsing mappings: %w", err)
	}

	resolver := resolve.NewResolver(client, auditor)

	h := signal.NewHandler()
	ctx, stop := h.WithContext(cmd.Context())
	defer stop()

	secrets, err := resolver.Resolve(ctx, mappings)
	if err != nil {
		return fmt.Errorf("resolving secrets: %w", err)
	}

	runner := exec.NewRunner(cfg.Command, os.Environ(), secrets)
	if err := runner.Run(ctx); err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
