package cmd

import (
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
)

var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Manage SSL certificates",
}

var sslGetCmd = &cobra.Command{
	Use:   "get <domain>",
	Short: "Retrieve SSL certificate bundle",
	Long: `Retrieve the SSL certificate bundle for a domain.
Outputs certificate, intermediate cert, and private key.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, err := apiClient.SSLRetrieve(args[0])
		if err != nil {
			return err
		}

		part, _ := cmd.Flags().GetString("part")

		if output.JSONOutput {
			if part != "" {
				switch part {
				case "cert":
					output.PrintJSON(map[string]string{"certificatechain": bundle.CertificateChain})
				case "key":
					output.PrintJSON(map[string]string{"privatekey": bundle.PrivateKey})
				case "intermediate":
					output.PrintJSON(map[string]string{"intermediatecertificate": bundle.IntermediateCertificate})
				default:
					output.PrintJSON(bundle)
				}
			} else {
				output.PrintJSON(bundle)
			}
			return nil
		}

		switch part {
		case "cert":
			output.Print(bundle.CertificateChain)
		case "key":
			output.Print(bundle.PrivateKey)
		case "intermediate":
			output.Print(bundle.IntermediateCertificate)
		case "public":
			output.Print(bundle.PublicKey)
		default:
			output.Print("=== Certificate Chain ===")
			output.Print(bundle.CertificateChain)
			output.Print("\n=== Private Key ===")
			output.Print(bundle.PrivateKey)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sslCmd)
	sslCmd.AddCommand(sslGetCmd)
	sslGetCmd.Flags().StringP("part", "p", "", "Output specific part: cert, key, intermediate, public")
}
