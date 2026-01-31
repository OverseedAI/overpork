package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	// Set env vars
	os.Setenv("PORKBUN_API_KEY", "test-api-key")
	os.Setenv("PORKBUN_SECRET_KEY", "test-secret-key")
	defer func() {
		os.Unsetenv("PORKBUN_API_KEY")
		os.Unsetenv("PORKBUN_SECRET_KEY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "test-api-key")
	}
	if cfg.SecretKey != "test-secret-key" {
		t.Errorf("SecretKey = %q, want %q", cfg.SecretKey, "test-secret-key")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "valid",
			cfg:     Config{APIKey: "pk1_xxx", SecretKey: "sk1_xxx"},
			wantErr: false,
		},
		{
			name:    "missing api key",
			cfg:     Config{SecretKey: "sk1_xxx"},
			wantErr: true,
		},
		{
			name:    "missing secret key",
			cfg:     Config{APIKey: "pk1_xxx"},
			wantErr: true,
		},
		{
			name:    "both missing",
			cfg:     Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error = %v", err)
	}
	if dir == "" {
		t.Error("ConfigDir() returned empty string")
	}
}
