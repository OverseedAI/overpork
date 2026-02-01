package cmd

import (
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if output.JSONOutput {
			output.PrintJSON(map[string]string{"version": Version})
		} else {
			output.Print("opork " + Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
