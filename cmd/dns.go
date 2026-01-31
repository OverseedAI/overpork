package cmd

import (
	"github.com/OverseedAI/overpork/internal/api"
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS records",
}

var dnsListCmd = &cobra.Command{
	Use:   "list <domain>",
	Short: "List DNS records for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType, _ := cmd.Flags().GetString("type")
		subdomain, _ := cmd.Flags().GetString("subdomain")

		var records []api.DNSRecord
		var err error

		if recordType != "" && subdomain != "" {
			records, err = apiClient.DNSListByTypeAndSubdomain(domain, recordType, subdomain)
		} else if recordType != "" {
			records, err = apiClient.DNSListByType(domain, recordType)
		} else {
			records, err = apiClient.DNSList(domain)
		}
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(records)
			return nil
		}

		if len(records) == 0 {
			output.Print("No records found")
			return nil
		}

		headers := []string{"ID", "TYPE", "NAME", "CONTENT", "TTL", "PRIO"}
		rows := make([][]string, len(records))
		for i, r := range records {
			rows[i] = []string{r.ID, r.Type, r.Name, r.Content, r.TTL, r.Prio}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

var dnsCreateCmd = &cobra.Command{
	Use:   "create <domain> <type> <content>",
	Short: "Create a DNS record",
	Long: `Create a DNS record for a domain.

Examples:
  overpork dns create example.com A 192.168.1.1
  overpork dns create example.com A 192.168.1.1 --name www
  overpork dns create example.com MX mail.example.com --prio 10
  overpork dns create example.com TXT "v=spf1 include:_spf.google.com ~all"`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType := args[1]
		content := args[2]

		name, _ := cmd.Flags().GetString("name")
		ttl, _ := cmd.Flags().GetString("ttl")
		prio, _ := cmd.Flags().GetString("prio")

		opts := api.DNSCreateOpts{
			Name: name,
			TTL:  ttl,
			Prio: prio,
		}

		id, err := apiClient.DNSCreate(domain, recordType, content, opts)
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]any{"id": id, "status": "created"})
		} else {
			output.Success("Created record %d", id)
		}
		return nil
	},
}

var dnsUpdateCmd = &cobra.Command{
	Use:   "update <domain> <id> <type> <content>",
	Short: "Update a DNS record by ID",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordID := args[1]
		recordType := args[2]
		content := args[3]

		name, _ := cmd.Flags().GetString("name")
		ttl, _ := cmd.Flags().GetString("ttl")
		prio, _ := cmd.Flags().GetString("prio")

		opts := api.DNSCreateOpts{
			Name: name,
			TTL:  ttl,
			Prio: prio,
		}

		if err := apiClient.DNSUpdate(domain, recordID, recordType, content, opts); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "updated"})
		} else {
			output.Success("Updated record %s", recordID)
		}
		return nil
	},
}

var dnsSetCmd = &cobra.Command{
	Use:   "set <domain> <type> <subdomain> <content>",
	Short: "Update a DNS record by type and subdomain",
	Long: `Update a DNS record identified by type and subdomain.
Use @ for the root domain.

Examples:
  overpork dns set example.com A www 192.168.1.1
  overpork dns set example.com A @ 192.168.1.1`,
	Args: cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType := args[1]
		subdomain := args[2]
		content := args[3]

		// Handle @ as empty string for root
		if subdomain == "@" {
			subdomain = ""
		}

		name, _ := cmd.Flags().GetString("name")
		ttl, _ := cmd.Flags().GetString("ttl")
		prio, _ := cmd.Flags().GetString("prio")

		opts := api.DNSCreateOpts{
			Name: name,
			TTL:  ttl,
			Prio: prio,
		}

		if err := apiClient.DNSUpdateByTypeAndSubdomain(domain, recordType, subdomain, content, opts); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "updated"})
		} else {
			output.Success("Updated %s record for %s", recordType, subdomain)
		}
		return nil
	},
}

var dnsDeleteCmd = &cobra.Command{
	Use:   "delete <domain> <id>",
	Short: "Delete a DNS record by ID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordID := args[1]

		if err := apiClient.DNSDelete(domain, recordID); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "deleted"})
		} else {
			output.Success("Deleted record %s", recordID)
		}
		return nil
	},
}

var dnsDeleteByNameCmd = &cobra.Command{
	Use:   "delete-by-name <domain> <type> <subdomain>",
	Short: "Delete a DNS record by type and subdomain",
	Long: `Delete a DNS record identified by type and subdomain.
Use @ for the root domain.

Examples:
  overpork dns delete-by-name example.com A www
  overpork dns delete-by-name example.com A @`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType := args[1]
		subdomain := args[2]

		if subdomain == "@" {
			subdomain = ""
		}

		if err := apiClient.DNSDeleteByTypeAndSubdomain(domain, recordType, subdomain); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "deleted"})
		} else {
			displayName := subdomain
			if displayName == "" {
				displayName = "@"
			}
			output.Success("Deleted %s record for %s", recordType, displayName)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)

	dnsCmd.AddCommand(dnsListCmd)
	dnsListCmd.Flags().StringP("type", "t", "", "Filter by record type (A, AAAA, MX, etc.)")
	dnsListCmd.Flags().StringP("subdomain", "s", "", "Filter by subdomain (requires --type)")

	dnsCmd.AddCommand(dnsCreateCmd)
	dnsCreateCmd.Flags().StringP("name", "n", "", "Subdomain (empty for root)")
	dnsCreateCmd.Flags().String("ttl", "", "TTL in seconds")
	dnsCreateCmd.Flags().String("prio", "", "Priority (for MX/SRV records)")

	dnsCmd.AddCommand(dnsUpdateCmd)
	dnsUpdateCmd.Flags().StringP("name", "n", "", "Subdomain (empty for root)")
	dnsUpdateCmd.Flags().String("ttl", "", "TTL in seconds")
	dnsUpdateCmd.Flags().String("prio", "", "Priority (for MX/SRV records)")

	dnsCmd.AddCommand(dnsSetCmd)
	dnsSetCmd.Flags().StringP("name", "n", "", "New subdomain name")
	dnsSetCmd.Flags().String("ttl", "", "TTL in seconds")
	dnsSetCmd.Flags().String("prio", "", "Priority (for MX/SRV records)")

	dnsCmd.AddCommand(dnsDeleteCmd)
	dnsCmd.AddCommand(dnsDeleteByNameCmd)
}
