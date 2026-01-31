package cmd

import (
	"github.com/OverseedAI/overpork/internal/api"
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains",
}

var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains in account",
	RunE: func(cmd *cobra.Command, args []string) error {
		start, _ := cmd.Flags().GetInt("start")
		domains, err := apiClient.DomainList(start)
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(domains)
			return nil
		}

		if len(domains) == 0 {
			output.Print("No domains found")
			return nil
		}

		headers := []string{"DOMAIN", "STATUS", "EXPIRES", "AUTO-RENEW", "PRIVACY"}
		rows := make([][]string, len(domains))
		for i, d := range domains {
			rows[i] = []string{d.Domain, d.Status, d.ExpireDate, d.AutoRenew, d.WhoisPrivacy}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

var domainGetCmd = &cobra.Command{
	Use:   "get <domain>",
	Short: "Get domain details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, err := apiClient.DomainGet(args[0])
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(domain)
			return nil
		}

		output.PrintTable([]string{"FIELD", "VALUE"}, [][]string{
			{"Domain", domain.Domain},
			{"Status", domain.Status},
			{"TLD", domain.TLD},
			{"Created", domain.CreateDate},
			{"Expires", domain.ExpireDate},
			{"Security Lock", domain.SecurityLock},
			{"WHOIS Privacy", domain.WhoisPrivacy},
			{"Auto-Renew", domain.AutoRenew},
		})
		return nil
	},
}

var domainNsGetCmd = &cobra.Command{
	Use:   "ns-get <domain>",
	Short: "Get nameservers for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns, err := apiClient.DomainGetNameservers(args[0])
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(ns)
			return nil
		}

		for _, n := range ns {
			output.Print(n)
		}
		return nil
	},
}

var domainNsSetCmd = &cobra.Command{
	Use:   "ns-set <domain> <ns1> [ns2] [ns3]...",
	Short: "Set nameservers for a domain",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		nameservers := args[1:]

		if err := apiClient.DomainUpdateNameservers(domain, nameservers); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "updated"})
		} else {
			output.Success("Nameservers updated")
		}
		return nil
	},
}

var domainForwardListCmd = &cobra.Command{
	Use:   "forward-list <domain>",
	Short: "List URL forwards for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		forwards, err := apiClient.DomainGetForwards(args[0])
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(forwards)
			return nil
		}

		if len(forwards) == 0 {
			output.Print("No forwards found")
			return nil
		}

		headers := []string{"ID", "SUBDOMAIN", "LOCATION", "TYPE", "WILDCARD"}
		rows := make([][]string, len(forwards))
		for i, f := range forwards {
			sub := f.Subdomain
			if sub == "" {
				sub = "@"
			}
			rows[i] = []string{f.ID, sub, f.Location, f.Type, f.Wildcard}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

var domainForwardAddCmd = &cobra.Command{
	Use:   "forward-add <domain> <location>",
	Short: "Add a URL forward",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		location := args[1]

		fwdType, _ := cmd.Flags().GetString("type")
		includePath, _ := cmd.Flags().GetBool("include-path")
		wildcard, _ := cmd.Flags().GetBool("wildcard")
		subdomain, _ := cmd.Flags().GetString("subdomain")

		opts := api.ForwardOpts{
			Type:        fwdType,
			IncludePath: includePath,
			Wildcard:    wildcard,
			Subdomain:   subdomain,
		}

		if err := apiClient.DomainAddForward(domain, location, opts); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "created"})
		} else {
			output.Success("Forward created")
		}
		return nil
	},
}

var domainForwardDeleteCmd = &cobra.Command{
	Use:   "forward-delete <domain> <id>",
	Short: "Delete a URL forward",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DomainDeleteForward(args[0], args[1]); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "deleted"})
		} else {
			output.Success("Forward deleted")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	domainCmd.AddCommand(domainListCmd)
	domainListCmd.Flags().Int("start", 0, "Pagination start index")

	domainCmd.AddCommand(domainGetCmd)
	domainCmd.AddCommand(domainNsGetCmd)
	domainCmd.AddCommand(domainNsSetCmd)

	domainCmd.AddCommand(domainForwardListCmd)
	domainCmd.AddCommand(domainForwardAddCmd)
	domainForwardAddCmd.Flags().String("type", "temporary", "Forward type: temporary or permanent")
	domainForwardAddCmd.Flags().Bool("include-path", false, "Include path in redirect")
	domainForwardAddCmd.Flags().Bool("wildcard", false, "Apply to all subdomains")
	domainForwardAddCmd.Flags().StringP("subdomain", "s", "", "Subdomain to forward")

	domainCmd.AddCommand(domainForwardDeleteCmd)
}
