package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/OverseedAI/overpork/internal/config"
	"github.com/OverseedAI/overpork/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create config file with API credentials",
	Long:  "Create config file with API credentials. Interactive by default, or use flags for non-interactive mode.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		secretKey, _ := cmd.Flags().GetString("secret-key")

		reader := bufio.NewReader(os.Stdin)

		if apiKey == "" {
			fmt.Print("API Key: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read API key: %w", err)
			}
			apiKey = strings.TrimSpace(input)
		}

		if secretKey == "" {
			fmt.Print("Secret Key: ")
			if term.IsTerminal(int(os.Stdin.Fd())) {
				secret, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return fmt.Errorf("failed to read secret key: %w", err)
				}
				fmt.Println()
				secretKey = string(secret)
			} else {
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read secret key: %w", err)
				}
				secretKey = strings.TrimSpace(input)
			}
		}

		if apiKey == "" || secretKey == "" {
			return fmt.Errorf("both API key and secret key are required")
		}

		configDir, err := config.ConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}

		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		configPath := filepath.Join(configDir, "config.yaml")
		content := fmt.Sprintf("api_key: %s\nsecret_key: %s\n", apiKey, secretKey)

		if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"path": configPath, "status": "created"})
		} else {
			output.Success("Config saved to %s", configPath)
		}
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		configPath := filepath.Join(configDir, "config.yaml")

		if output.JSONOutput {
			output.PrintJSON(map[string]string{"path": configPath})
		} else {
			output.Print(configPath)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().String("api-key", "", "Porkbun API key")
	configInitCmd.Flags().String("secret-key", "", "Porkbun secret key")

	configCmd.AddCommand(configPathCmd)
}
