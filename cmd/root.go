package cmd

import (
	"os"

	"github.com/OverseedAI/overpork/internal/api"
	"github.com/OverseedAI/overpork/internal/config"
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	apiClient *api.Client
)

var rootCmd = &cobra.Command{
	Use:   "opork",
	Short: "CLI wrapper for Porkbun API",
	Long:  "opork is a CLI tool for managing domains, DNS records, and SSL certificates via the Porkbun API.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip auth for commands that don't need it
		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}
		if cmd.Parent() != nil && cmd.Parent().Name() == "overpork" && cmd.Name() == "config" {
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		apiClient = api.NewClient(cfg)
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&output.JSONOutput, "json", false, "Output in JSON format")
}
