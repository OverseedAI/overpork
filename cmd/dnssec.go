package cmd

import (
	"github.com/OverseedAI/overpork/internal/api"
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var dnssecCmd = &cobra.Command{
	Use:   "dnssec",
	Short: "Manage DNSSEC records",
}

var dnssecListCmd = &cobra.Command{
	Use:   "list <domain>",
	Short: "List DNSSEC records for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		records, err := apiClient.DNSSECList(args[0])
		if err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(records)
			return nil
		}

		if len(records) == 0 {
			output.Print("No DNSSEC records found")
			return nil
		}

		headers := []string{"KEYTAG", "ALGORITHM", "DIGEST_TYPE", "DIGEST"}
		rows := make([][]string, len(records))
		for i, r := range records {
			digest := r.Digest
			if len(digest) > 32 {
				digest = digest[:32] + "..."
			}
			rows[i] = []string{r.KeyTag, r.Algorithm, r.DigestType, digest}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

var dnssecCreateCmd = &cobra.Command{
	Use:   "create <domain>",
	Short: "Create a DNSSEC record",
	Long: `Create a DNSSEC DS record at the registry.

Required flags: --keytag, --algorithm, --digest-type, --digest`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		keyTag, _ := cmd.Flags().GetString("keytag")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		digestType, _ := cmd.Flags().GetString("digest-type")
		digest, _ := cmd.Flags().GetString("digest")
		publicKey, _ := cmd.Flags().GetString("public-key")
		flags, _ := cmd.Flags().GetString("flags")

		record := api.DNSSECRecord{
			KeyTag:     keyTag,
			Algorithm:  algorithm,
			DigestType: digestType,
			Digest:     digest,
			PublicKey:  publicKey,
			Flags:      flags,
		}

		if err := apiClient.DNSSECCreate(domain, record); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "created"})
		} else {
			output.Success("Created DNSSEC record for %s", domain)
		}
		return nil
	},
}

var dnssecDeleteCmd = &cobra.Command{
	Use:   "delete <domain> <keytag>",
	Short: "Delete a DNSSEC record",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DNSSECDelete(args[0], args[1]); err != nil {
			return err
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"status": "deleted"})
		} else {
			output.Success("Deleted DNSSEC record %s", args[1])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnssecCmd)

	dnssecCmd.AddCommand(dnssecListCmd)

	dnssecCmd.AddCommand(dnssecCreateCmd)
	dnssecCreateCmd.Flags().String("keytag", "", "Key tag (required)")
	dnssecCreateCmd.Flags().String("algorithm", "", "Algorithm number (required)")
	dnssecCreateCmd.Flags().String("digest-type", "", "Digest type (required)")
	dnssecCreateCmd.Flags().String("digest", "", "Digest value (required)")
	dnssecCreateCmd.Flags().String("public-key", "", "Public key (optional)")
	dnssecCreateCmd.Flags().String("flags", "", "Flags (optional)")
	dnssecCreateCmd.MarkFlagRequired("keytag")
	dnssecCreateCmd.MarkFlagRequired("algorithm")
	dnssecCreateCmd.MarkFlagRequired("digest-type")
	dnssecCreateCmd.MarkFlagRequired("digest")

	dnssecCmd.AddCommand(dnssecDeleteCmd)
}
