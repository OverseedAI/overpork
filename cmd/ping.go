package cmd

import (
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Test API connectivity and authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.Ping(); err != nil {
			return err
		}
		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "OK"})
		} else {
			output.Success("OK")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
