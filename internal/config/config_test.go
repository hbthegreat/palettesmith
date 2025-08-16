package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandHome(t *testing.T) {
	t.Run("should_expand_relative_path_to_absolute_home_path", func(t *testing.T) {
		result, err := expandHome(".config/palettesmith")

		require.NoError(t, err)
		assert.Contains(t, result, ".config/palettesmith")
		assert.True(t, filepath.IsAbs(result))
	})

	t.Run("should_return_error_when_home_directory_unavailable", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		os.Unsetenv("HOME")

		_, err := expandHome(".config/palettesmith")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "home directory")
	})
}

func TestSetPreset(t *testing.T) {
	t.Run("should_switch_to_generic_preset_and_update_paths", func(t *testing.T) {
		manager := &Manager{cfg: Config{Preset: "omarchy"}}

		err := manager.SetPreset("generic")

		assert.NoError(t, err)
		assert.Equal(t, "generic", manager.cfg.Preset)
		assert.Contains(t, manager.cfg.TargetThemeDir, "palettesmith")
	})

	t.Run("should_switch_to_omarchy_preset_and_update_paths", func(t *testing.T) {
		manager := &Manager{cfg: Config{Preset: "generic"}}

		err := manager.SetPreset("omarchy")

		assert.NoError(t, err)
		assert.Equal(t, "omarchy", manager.cfg.Preset)
		assert.Contains(t, manager.cfg.TargetThemeDir, "omarchy")
	})

	t.Run("should_return_error_for_unknown_preset", func(t *testing.T) {
		manager := &Manager{cfg: Config{Preset: "generic"}}

		err := manager.SetPreset("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown preset")
		assert.Contains(t, err.Error(), "nonexistent")
	})

	t.Run("should_return_error_for_unsupported_custom_preset", func(t *testing.T) {
		manager := &Manager{cfg: Config{Preset: "generic"}}

		err := manager.SetPreset("custom")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet supported")
	})

	t.Run("should_return_error_when_home_directory_unavailable_for_generic", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		os.Unsetenv("HOME")

		manager := &Manager{cfg: Config{Preset: "omarchy"}}
		err := manager.SetPreset("generic")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set preset 'generic'")
		assert.Contains(t, err.Error(), "cannot determine home directory")
	})

	t.Run("should_return_error_when_home_directory_unavailable_for_omarchy", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		os.Unsetenv("HOME")

		manager := &Manager{cfg: Config{Preset: "generic"}}
		err := manager.SetPreset("omarchy")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set preset 'omarchy'")
		assert.Contains(t, err.Error(), "cannot determine home directory")
	})
}

func TestConfig_JSONSerialization(t *testing.T) {
	t.Run("should_marshal_and_unmarshal_config_without_data_loss", func(t *testing.T) {
		original := Config{
			TargetThemeDir:   "/home/user/.config/palettesmith/themes",
			CurrentThemeLink: "/home/user/.config/palettesmith/current/theme",
			Preset:           "generic",
			StagingDir:       "/home/user/.config/palettesmith/staging",
		}

		data, err := json.MarshalIndent(original, "", "  ")
		require.NoError(t, err)

		var restored Config
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Equal(t, original, restored)
	})
}

func TestNewManager(t *testing.T) {
	t.Run("should_create_manager_with_default_config_when_no_config_file_exists", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		manager, err := NewManager()

		require.NoError(t, err)
		assert.NotNil(t, manager)

		cfg := manager.GetConfig()
		assert.Equal(t, "generic", cfg.Preset)
		assert.Contains(t, cfg.TargetThemeDir, "palettesmith")
	})

	t.Run("should_load_existing_config_when_config_file_exists", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// Create config directory and file
		configDir := filepath.Join(tempHome, ".config/palettesmith")
		err := os.MkdirAll(configDir, 0o755)
		require.NoError(t, err)

		existingConfig := Config{
			TargetThemeDir:   "/existing/themes",
			CurrentThemeLink: "/existing/current",
			Preset:           "omarchy",
			StagingDir:       "/existing/staging",
		}

		configFile := filepath.Join(configDir, "config.json")
		err = saveConfigToFile(existingConfig, configFile)
		require.NoError(t, err)

		manager, err := NewManager()

		require.NoError(t, err)
		cfg := manager.GetConfig()
		assert.Equal(t, "omarchy", cfg.Preset)
		assert.Equal(t, "/existing/themes", cfg.TargetThemeDir)
	})

	t.Run("should_return_error_when_home_directory_unavailable", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		os.Unsetenv("HOME")

		_, err := NewManager()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get home directory")
	})

	t.Run("should_return_error_when_config_file_corrupted", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// Create config directory and corrupted file
		configDir := filepath.Join(tempHome, ".config/palettesmith")
		err := os.MkdirAll(configDir, 0o755)
		require.NoError(t, err)

		configFile := filepath.Join(configDir, "config.json")
		err = os.WriteFile(configFile, []byte("{invalid json"), 0o644)
		require.NoError(t, err)

		_, err = NewManager()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load config")
	})
}

func TestManager_SaveConfig(t *testing.T) {
	t.Run("should_persist_config_to_disk", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// Create the config directory first
		configDir := filepath.Join(tempHome, ".config/palettesmith")
		err := os.MkdirAll(configDir, 0o755)
		require.NoError(t, err)

		manager := &Manager{
			cfg: Config{
				TargetThemeDir:   "/test/themes",
				CurrentThemeLink: "/test/current",
				Preset:           "test",
				StagingDir:       "/test/staging",
			},
		}

		err = manager.SaveConfig()
		require.NoError(t, err)

		// Verify file was created
		configFile := filepath.Join(tempHome, ".config/palettesmith/config.json")
		assert.FileExists(t, configFile)

		// Verify content is correct
		loadedConfig, err := loadConfigFromFile(configFile)
		require.NoError(t, err)
		assert.Equal(t, manager.cfg, loadedConfig)
	})

	t.Run("should_return_error_when_home_directory_unavailable", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		os.Unsetenv("HOME")

		manager := &Manager{cfg: Config{Preset: "test"}}
		err := manager.SaveConfig()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save config")
		assert.Contains(t, err.Error(), "cannot determine home directory")
	})
}

func TestManager_GetConfig(t *testing.T) {
	t.Run("should_return_copy_of_internal_config", func(t *testing.T) {
		expectedConfig := Config{
			TargetThemeDir:   "/test/themes",
			CurrentThemeLink: "/test/current",
			Preset:           "test",
			StagingDir:       "/test/staging",
		}

		manager := &Manager{cfg: expectedConfig}
		result := manager.GetConfig()

		assert.Equal(t, expectedConfig, result)
	})
}

func TestLoadConfigFromFile_ErrorHandling(t *testing.T) {
	t.Run("should_return_error_when_config_file_does_not_exist", func(t *testing.T) {
		_, err := loadConfigFromFile("/nonexistent/path/config.json")
		assert.Error(t, err)
	})

	t.Run("should_return_error_when_config_file_contains_invalid_json", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "invalid.json")
		err := os.WriteFile(tempFile, []byte("{invalid json}"), 0o644)
		require.NoError(t, err)

		_, err = loadConfigFromFile(tempFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config JSON")
	})
}

func TestSaveConfigToFile_ErrorHandling(t *testing.T) {
	t.Run("should_return_error_when_cannot_write_to_restricted_directory", func(t *testing.T) {
		cfg := Config{Preset: "test"}

		err := saveConfigToFile(cfg, "/root/config.json")
		assert.Error(t, err)
	})
}

func TestConstants_Validation(t *testing.T) {
	t.Run("should_have_non_empty_config_directory_constants", func(t *testing.T) {
		assert.NotEmpty(t, PalettesmithConfigDir)
		assert.NotEmpty(t, OmarchyConfigDir)
	})

	t.Run("should_use_relative_paths_for_config_directories", func(t *testing.T) {
		assert.NotEqual(t, "/", string(PalettesmithConfigDir[0]))
		assert.NotEqual(t, "/", string(OmarchyConfigDir[0]))
	})
}

func TestManager_IsFirstRun(t *testing.T) {
	t.Run("should_return_true_for_first_run_when_no_config_exists", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		manager, err := NewManager()

		require.NoError(t, err)
		assert.True(t, manager.IsFirstRun())
	})

	t.Run("should_return_false_when_config_file_exists", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// Create config file first
		configDir := filepath.Join(tempHome, ".config/palettesmith")
		err := os.MkdirAll(configDir, 0o755)
		require.NoError(t, err)

		existingConfig := Config{
			TargetThemeDir:   "/existing/themes",
			CurrentThemeLink: "/existing/current",
			Preset:           "omarchy",
			StagingDir:       "/existing/staging",
		}

		configFile := filepath.Join(configDir, "config.json")
		err = saveConfigToFile(existingConfig, configFile)
		require.NoError(t, err)

		manager, err := NewManager()

		require.NoError(t, err)
		assert.False(t, manager.IsFirstRun())
	})
}

func TestManager_MarkSetupComplete(t *testing.T) {
	t.Run("should_mark_first_run_as_false", func(t *testing.T) {
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		manager, err := NewManager()
		require.NoError(t, err)

		// Should be first run initially
		assert.True(t, manager.IsFirstRun())

		// Mark setup complete
		manager.MarkSetupComplete()

		// Should no longer be first run
		assert.False(t, manager.IsFirstRun())
	})
}
