package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/vaultpipe/internal/audit"
	"github.com/vaultpipe/vaultpipe/internal/config"
	"github.com/vaultpipe/vaultpipe/internal/exec"
	"github.com/vaultpipe/vaultpipe/internal/resolve"
	"github.com/vaultpipe/vaultpipe/internal/vault"
	"github.com/vaultpipe/vaultpipe/internal/version"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Stream secrets from Vault into process environments",
	Long:  version.FullString(),
	RunE:  run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "path to config file (required)")
	_ = rootCmd.MarkPersistentFlagRequired("config")
}

func run(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	auditor := audit.NewAuditor(cfg.Audit.Enabled, cfg.Audit.LogPath)

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	mappings, err := config.ToMappings(cfg)
	if err != nil {
		return fmt.Errorf("parsing mappings: %w", err)
	}

	resolver := resolve.NewResolver(client, auditor)
	secrets, err := resolver.Resolve(cmd.Context(), mappings)
	if err != nil {
		return fmt.Errorf("resolving secrets: %w", err)
	}

	runner := exec.NewRunner(cfg.Command, secrets)
	return runner.Run(cmd.Context())
}
