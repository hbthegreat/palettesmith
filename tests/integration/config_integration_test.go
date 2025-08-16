package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	require.NoError(t, err)

	err = manager1.SetPreset("omarchy")
	require.NoError(t, err)

	err = manager1.SaveConfig()
	require.NoError(t, err)

	// Create second manager - should load the saved config
	manager2, err := config.NewManager()
	require.NoError(t, err)

	cfg := manager2.GetConfig()
	assert.Equal(t, "omarchy", cfg.Preset)

	// Verify the paths are correct for omarchy preset
	assert.Contains(t, cfg.TargetThemeDir, "omarchy")
}

// TestConfigFileCorruption tests handling of corrupted config files
func TestConfigFileCorruption(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Create config directory
	configDir := filepath.Join(tempHome, ".config/palettesmith")
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	// Write corrupted JSON file
	configFile := filepath.Join(configDir, "config.json")
	corruptedJSON := `{"preset": "omarchy", "invalid": json}`
	err = os.WriteFile(configFile, []byte(corruptedJSON), 0o644)
	require.NoError(t, err)

	// NewManager should fail with corrupted config
	_, err = config.NewManager()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

// TestPresetSwitching tests switching between different presets
func TestPresetSwitching(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Test switching from generic to omarchy
	err = manager.SetPreset("omarchy")
	require.NoError(t, err)

	cfg := manager.GetConfig()
	assert.Equal(t, "omarchy", cfg.Preset)
	assert.Contains(t, cfg.TargetThemeDir, "omarchy")

	// Test switching back to generic
	err = manager.SetPreset("generic")
	require.NoError(t, err)

	cfg = manager.GetConfig()
	assert.Equal(t, "generic", cfg.Preset)
	assert.Contains(t, cfg.TargetThemeDir, "palettesmith")
}

// TestConfigJSONFormat tests that saved config has correct JSON format
func TestConfigJSONFormat(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	manager, err := config.NewManager()
	require.NoError(t, err)

	err = manager.SaveConfig()
	require.NoError(t, err)

	// Read the saved config file
	configFile := filepath.Join(tempHome, ".config/palettesmith/config.json")
	data, err := os.ReadFile(configFile)
	require.NoError(t, err)

	// Verify it's valid JSON
	var cfg map[string]interface{}
	err = json.Unmarshal(data, &cfg)
	require.NoError(t, err)

	// Verify required fields exist
	requiredFields := []string{"target_theme_dir", "current_theme_link", "preset", "staging_dir"}
	for _, field := range requiredFields {
		assert.Contains(t, cfg, field, "Config missing required field: %s", field)
	}

	// Verify JSON is formatted (indented)
	assert.Contains(t, string(data), "\n", "Config JSON should be indented/formatted")
}

// TestConcurrentAccess tests that multiple managers can work with the same config
func TestConcurrentAccess(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Create first manager and save config
	manager1, err := config.NewManager()
	require.NoError(t, err)

	err = manager1.SetPreset("omarchy")
	require.NoError(t, err)

	err = manager1.SaveConfig()
	require.NoError(t, err)

	// Create second manager - should load the existing config
	manager2, err := config.NewManager()
	require.NoError(t, err)

	// Both managers should have the same config
	cfg1 := manager1.GetConfig()
	cfg2 := manager2.GetConfig()

	assert.Equal(t, cfg1.Preset, cfg2.Preset)
	assert.Equal(t, cfg1.TargetThemeDir, cfg2.TargetThemeDir)
	assert.Equal(t, cfg1.CurrentThemeLink, cfg2.CurrentThemeLink)
	assert.Equal(t, cfg1.StagingDir, cfg2.StagingDir)
}
