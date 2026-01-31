package cmd

import (
	"strings"

	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var glueCmd = &cobra.Command{
	Use:   "glue",
	Short: "Manage glue records (custom nameserver IPs)",
}

var glueListCmd = &cobra.Command{
	Use:   "list <domain>",
	Short: "List glue records for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		records, err := apiClient.GlueList(args[0])
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(records)
			return nil
		}

		if len(records) == 0 {
			output.Print("No glue records found")
			return nil
		}

		headers := []string{"SUBDOMAIN", "IPS"}
		rows := make([][]string, len(records))
		for i, r := range records {
			rows[i] = []string{r.Subdomain, strings.Join(r.IPs, ", ")}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

var glueCreateCmd = &cobra.Command{
	Use:   "create <domain> <subdomain> <ip> [ip...]",
	Short: "Create a glue record",
	Long: `Create a glue record for a custom nameserver.

Examples:
  overpork glue create example.com ns1 192.168.1.1
  overpork glue create example.com ns1 192.168.1.1 2001:db8::1`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		subdomain := args[1]
		ips := args[2:]

		if err := apiClient.GlueCreate(domain, subdomain, ips); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "created"})
		} else {
			output.Success("Created glue record for %s.%s", subdomain, domain)
		}
		return nil
	},
}

var glueUpdateCmd = &cobra.Command{
	Use:   "update <domain> <subdomain> <ip> [ip...]",
	Short: "Update a glue record",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		subdomain := args[1]
		ips := args[2:]

		if err := apiClient.GlueUpdate(domain, subdomain, ips); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "updated"})
		} else {
			output.Success("Updated glue record for %s.%s", subdomain, domain)
		}
		return nil
	},
}

var glueDeleteCmd = &cobra.Command{
	Use:   "delete <domain> <subdomain>",
	Short: "Delete a glue record",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.GlueDelete(args[0], args[1]); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "deleted"})
		} else {
			output.Success("Deleted glue record for %s.%s", args[1], args[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(glueCmd)
	glueCmd.AddCommand(glueListCmd)
	glueCmd.AddCommand(glueCreateCmd)
	glueCmd.AddCommand(glueUpdateCmd)
	glueCmd.AddCommand(glueDeleteCmd)
}
