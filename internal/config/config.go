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
	cfg       Config
	isFirstRun bool
}

func NewManager() (*Manager, error) {
	configDir, err := expandHome(PalettesmithConfigDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Validate directory accessibility before proceeding
	if err := validateDirectoryAccess(configDir); err != nil {
		return nil, fmt.Errorf("config directory not accessible: %w", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	var cfg Config
	var isFirstRun bool

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// File doesn't exist, this is a first run
		isFirstRun = true
		cfg = Config{
			TargetThemeDir:   filepath.Join(configDir, "themes"),
			CurrentThemeLink: filepath.Join(configDir, "current/theme"),
			Preset:           "generic",
			StagingDir:       filepath.Join(configDir, "staging"),
		}

		// Don't save the config yet - let the setup flow do that
	} else if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	} else {
		// File exists, load it
		isFirstRun = false
		cfg, err = loadConfigFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return &Manager{cfg: cfg, isFirstRun: isFirstRun}, nil
}

func (m *Manager) GetConfig() Config {
	return m.cfg
}

// IsFirstRun returns whether this is the first time the application is running
func (m *Manager) IsFirstRun() bool {
	return m.isFirstRun
}

// MarkSetupComplete marks the first run as completed  
func (m *Manager) MarkSetupComplete() {
	m.isFirstRun = false
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

	var newCfg Config
	switch preset {
	case "generic":
		newCfg = Config{
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
		newCfg = Config{
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

	// Validate that all directories are accessible before committing the change
	if err := validatePresetDirectories(newCfg); err != nil {
		return fmt.Errorf("failed to set preset '%s': %w", preset, err)
	}

	// Only update the config if validation succeeded
	m.cfg = newCfg
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

// validateDirectoryAccess checks if a directory can be created and written to
func validateDirectoryAccess(dirPath string) error {
	// Try to create the directory
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("cannot create directory %s: %w", dirPath, err)
	}

	// Test write access by creating and removing a temporary file
	testFile := filepath.Join(dirPath, ".palettesmith_access_test")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		return fmt.Errorf("cannot write to directory %s: %w", dirPath, err)
	}
	
	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		// Log warning but don't fail - the main operation succeeded
		// In a real app, we might use a logger here
	}

	return nil
}

// validatePresetDirectories validates that the directories for a preset are accessible
func validatePresetDirectories(cfg Config) error {
	// Validate target theme directory
	if cfg.TargetThemeDir != "" {
		if err := validateDirectoryAccess(cfg.TargetThemeDir); err != nil {
			return fmt.Errorf("target theme directory not accessible: %w", err)
		}
	}

	// Validate staging directory
	if cfg.StagingDir != "" {
		if err := validateDirectoryAccess(cfg.StagingDir); err != nil {
			return fmt.Errorf("staging directory not accessible: %w", err)
		}
	}

	// For current theme link, just validate the parent directory exists
	if cfg.CurrentThemeLink != "" {
		parentDir := filepath.Dir(cfg.CurrentThemeLink)
		if err := validateDirectoryAccess(parentDir); err != nil {
			return fmt.Errorf("current theme link directory not accessible: %w", err)
		}
	}

	return nil
}
