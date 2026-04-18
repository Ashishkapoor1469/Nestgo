package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GlobalConfig represents the preferences stored in ~/.nestgo/config.json
type GlobalConfig struct {
	DefaultPackageManager string `json:"defaultPackageManager"`
	TelemetryEnabled      bool   `json:"telemetryEnabled"`
	TemplatesPath         string `json:"templatesPath"`
}

// LoadGlobalConfig reads the config from the user's home directory.
func LoadGlobalConfig() (*GlobalConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	configDir := filepath.Join(home, ".nestgo")
	configPath := filepath.Join(configDir, "config.json")

	// Create defaults if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &GlobalConfig{
			DefaultPackageManager: "go",
			TelemetryEnabled:      true,
			TemplatesPath:         "",
		}
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return defaultConfig, nil // Degrading gracefully
		}
		
		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		_ = os.WriteFile(configPath, data, 0644)
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg GlobalConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
