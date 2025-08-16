// Package config provides configuration management for Palettesmith.
// It handles loading, saving, and managing application settings including
// theme directory paths and preset configurations for different target systems.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	PalettesmithConfigDir = ".config/palettesmith"
	OmarchyConfigDir      = ".config/omarchy"
)

type Config struct {
	TargetThemeDir   string `json:"target_theme_dir"`
	CurrentThemeLink string `json:"current_theme_link"`
	Preset           string `json:"preset"`
	StagingDir       string `json:"staging_dir"`
}

type Manager struct {
	cfg Config
}

func NewManager() (*Manager, error) {
	configDir, err := expandHome(PalettesmithConfigDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	configFile := filepath.Join(configDir, "config.json")

	var cfg Config

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// File doesn't exist, create default config
		cfg = Config{
			TargetThemeDir:   filepath.Join(configDir, "themes"),
			CurrentThemeLink: filepath.Join(configDir, "current/theme"),
			Preset:           "generic",
			StagingDir:       filepath.Join(configDir, "staging"),
		}

		// Save default config to file
		if err := saveConfigToFile(cfg, configFile); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	} else {
		cfg, err = loadConfigFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return &Manager{cfg: cfg}, nil
}

func (m *Manager) GetConfig() Config {
	return m.cfg
}

// SaveConfig persists the current configuration to disk
func (m *Manager) SaveConfig() error {
	configDir, err := expandHome(PalettesmithConfigDir)
	if err != nil {
		return fmt.Errorf("failed to save config: cannot determine home directory: %w", err)
	}
	
	configFile := filepath.Join(configDir, "config.json")
	if err := saveConfigToFile(m.cfg, configFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	return nil
}

func (m *Manager) SetPreset(preset string) error {
	palettesmithDir, err := expandHome(PalettesmithConfigDir)
	if err != nil {
		return fmt.Errorf("failed to set preset '%s': cannot determine home directory: %w", preset, err)
	}
	
	switch preset {
	case "generic":
		m.cfg = Config{
			TargetThemeDir:   filepath.Join(palettesmithDir, "themes"),
			CurrentThemeLink: filepath.Join(palettesmithDir, "current/theme"),
			Preset:           "generic",
			StagingDir:       filepath.Join(palettesmithDir, "staging"),
		}
	case "omarchy":
		omarchyDir, err := expandHome(OmarchyConfigDir)
		if err != nil {
			return fmt.Errorf("failed to set preset '%s': cannot determine home directory: %w", preset, err)
		}
		m.cfg = Config{
			TargetThemeDir:   filepath.Join(omarchyDir, "themes"),
			CurrentThemeLink: filepath.Join(omarchyDir, "current/theme"),
			Preset:           "omarchy",
			StagingDir:       filepath.Join(palettesmithDir, "staging"),
		}
	case "custom":
		return fmt.Errorf("preset '%s' is not yet supported", preset)
	default:
		return fmt.Errorf("unknown preset '%s': supported presets are 'generic', 'omarchy'", preset)
	}
	return nil
}

func expandHome(relativePath string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}
	return filepath.Join(home, relativePath), nil
}

// loadConfigFromFile reads and parses a JSON config file
func loadConfigFromFile(configFile string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return cfg, nil
}

// saveConfigToFile writes a config struct to a JSON file
func saveConfigToFile(cfg Config, configFile string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
