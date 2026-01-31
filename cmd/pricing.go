package cmd

import (
	"fmt"
	"sort"

	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var pricingCmd = &cobra.Command{
	Use:   "pricing",
	Short: "View domain pricing",
}

var pricingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all TLD pricing",
	RunE: func(cmd *cobra.Command, args []string) error {
		pricing, err := apiClient.PricingList()
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(pricing)
			return nil
		}

		// Sort TLDs alphabetically
		tlds := make([]string, 0, len(pricing))
		for tld := range pricing {
			tlds = append(tlds, tld)
		}
		sort.Strings(tlds)

		headers := []string{"TLD", "REGISTER", "RENEW", "TRANSFER"}
		rows := make([][]string, len(tlds))
		for i, tld := range tlds {
			p := pricing[tld]
			rows[i] = []string{tld, "$" + p.Registration, "$" + p.Renewal, "$" + p.Transfer}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

var pricingCheckCmd = &cobra.Command{
	Use:   "check <domain>",
	Short: "Check domain availability and price",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		available, price, err := apiClient.DomainCheck(domain)
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]any{
				"domain":    domain,
				"available": available,
				"price":     price,
			})
			return nil
		}

		if available {
			output.Success("%s is available - $%.2f", domain, price)
		} else {
			output.Print(fmt.Sprintf("%s is not available", domain))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pricingCmd)
	pricingCmd.AddCommand(pricingListCmd)
	pricingCmd.AddCommand(pricingCheckCmd)
}
