package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
}

func Load() (*Config, error) {
	// Env vars take precedence
	viper.SetEnvPrefix("PORKBUN")
	_ = viper.BindEnv("api_key")
	_ = viper.BindEnv("secret_key")

	// XDG config
	configDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(configDir, "overpork"))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read config file (ignore if not found)
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key not set (use PORKBUN_API_KEY env var or config file)")
	}
	if c.SecretKey == "" {
		return fmt.Errorf("secret key not set (use PORKBUN_SECRET_KEY env var or config file)")
	}
	return nil
}

func ConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "overpork"), nil
}
