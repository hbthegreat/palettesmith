package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"palettesmith/internal/config"
)

func TestNewManagerFirstRun(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// On first run, should return true for IsFirstRun
	assert.True(t, manager.IsFirstRun())

	cfg := manager.GetConfig()
	assert.Equal(t, "generic", cfg.Preset)

	// Verify config directory was created
	configDir := filepath.Join(tempHome, ".config/palettesmith")
	assert.DirExists(t, configDir)

	// Verify config file does NOT exist yet (until SaveConfig is called)
	configPath := filepath.Join(tempHome, ".config/palettesmith/config.json")
	assert.NoFileExists(t, configPath)

	// After calling SaveConfig, file should exist
	err = manager.SaveConfig()
	require.NoError(t, err)
	assert.FileExists(t, configPath)
}

// TestConfigFilePersistence tests that configs are properly saved and loaded
func TestConfigFilePersistence(t *testing.T) {
	// Create temporary home directory
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Create first manager and set preset
	manager1, err := config.NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if err := manager1.SetPreset("omarchy"); err != nil {
		t.Fatalf("SetPreset failed: %v", err)
	}

	if err := manager1.SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Create second manager - should load the saved config
	manager2, err := config.NewManager()
	if err != nil {
		t.Fatalf("Second NewManager failed: %v", err)
	}

	cfg := manager2.GetConfig()
	if cfg.Preset != "omarchy" {
		t.Errorf("Expected loaded preset 'omarchy', got '%s'", cfg.Preset)
	}

	// Verify the paths are correct for omarchy preset
	if !strings.Contains(cfg.TargetThemeDir, "omarchy") {
		t.Errorf("Expected omarchy theme dir, got: %s", cfg.TargetThemeDir)
	}
}

// TestConfigFileCorruption tests handling of corrupted config files
func TestConfigFileCorruption(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Create config directory
	configDir := filepath.Join(tempHome, ".config/palettesmith")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write corrupted JSON file
	configFile := filepath.Join(configDir, "config.json")
	corruptedJSON := `{"preset": "omarchy", "invalid": json}`
	if err := os.WriteFile(configFile, []byte(corruptedJSON), 0o644); err != nil {
		t.Fatalf("Failed to write corrupted config: %v", err)
	}

	// NewManager should fail with corrupted config
	_, err := config.NewManager()
	if err == nil {
		t.Error("Expected error with corrupted config file")
	}
	if !strings.Contains(err.Error(), "failed to load config") {
		t.Errorf("Expected config loading error, got: %v", err)
	}
}

// TestPresetSwitching tests switching between different presets
func TestPresetSwitching(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	manager, err := config.NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Test switching from generic to omarchy
	if err := manager.SetPreset("omarchy"); err != nil {
		t.Fatalf("SetPreset(omarchy) failed: %v", err)
	}

	cfg := manager.GetConfig()
	if cfg.Preset != "omarchy" {
		t.Errorf("Expected preset 'omarchy', got '%s'", cfg.Preset)
	}
	if !strings.Contains(cfg.TargetThemeDir, "omarchy") {
		t.Errorf("Expected omarchy paths, got: %s", cfg.TargetThemeDir)
	}

	// Test switching back to generic
	if err := manager.SetPreset("generic"); err != nil {
		t.Fatalf("SetPreset(generic) failed: %v", err)
	}

	cfg = manager.GetConfig()
	if cfg.Preset != "generic" {
		t.Errorf("Expected preset 'generic', got '%s'", cfg.Preset)
	}
	if !strings.Contains(cfg.TargetThemeDir, "palettesmith") {
		t.Errorf("Expected palettesmith paths, got: %s", cfg.TargetThemeDir)
	}
}

// TestConfigJSONFormat tests that saved config has correct JSON format
func TestConfigJSONFormat(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	manager, err := config.NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if err := manager.SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Read the saved config file
	configFile := filepath.Join(tempHome, ".config/palettesmith/config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Verify it's valid JSON
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Config file is not valid JSON: %v", err)
	}

	// Verify required fields exist
	requiredFields := []string{"target_theme_dir", "current_theme_link", "preset", "staging_dir"}
	for _, field := range requiredFields {
		if _, exists := cfg[field]; !exists {
			t.Errorf("Config missing required field: %s", field)
		}
	}

	// Verify JSON is formatted (indented)
	if !strings.Contains(string(data), "\n") {
		t.Error("Config JSON should be indented/formatted")
	}
}

// TestConcurrentAccess tests that multiple managers can work with the same config
func TestConcurrentAccess(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Create first manager and save config
	manager1, err := config.NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if err := manager1.SetPreset("omarchy"); err != nil {
		t.Fatalf("SetPreset failed: %v", err)
	}

	if err := manager1.SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Create second manager - should load the existing config
	manager2, err := config.NewManager()
	if err != nil {
		t.Fatalf("Second NewManager failed: %v", err)
	}

	// Both managers should have the same config
	cfg1 := manager1.GetConfig()
	cfg2 := manager2.GetConfig()

	if cfg1.Preset != cfg2.Preset {
		t.Errorf("Managers have different presets: %s vs %s", cfg1.Preset, cfg2.Preset)
	}
}
