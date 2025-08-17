package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover(t *testing.T) {
	t.Run("should_discover_plugins_and_collect_errors", func(t *testing.T) {
		result := Discover()
		
		assert.NotNil(t, result)
		assert.NotNil(t, result.Store)
		assert.NotNil(t, result.Errors)
	})

	t.Run("should_load_hyprland_plugin_successfully", func(t *testing.T) {
		result := Discover()
		
		store := result.Store
		plugins := store.List()
		
		t.Logf("Discovered %d plugins:", len(plugins))
		for _, p := range plugins {
			t.Logf("  - %s (%s)", p.Manifest.ID, p.Manifest.Title)
		}
		
		if len(result.Errors) > 0 {
			t.Logf("Discovery errors:")
			for id, errs := range result.Errors {
				for _, err := range errs {
					t.Logf("  %s: %v", id, err)
				}
			}
		}
		
		var hyprlandFound bool
		for _, p := range plugins {
			if p.Manifest.ID == "hyprland" {
				hyprlandFound = true
				assert.Equal(t, "Hyprland", p.Manifest.Title)
				assert.NotEmpty(t, p.Spec.Fields)
				
				// Verify spec has template file
				assert.NotEmpty(t, p.Spec.TemplateFile, "Spec should have template_file")
				break
			}
		}
		
		assert.True(t, hyprlandFound, "Hyprland plugin should be discovered")
		assert.NotEmpty(t, plugins, "Should discover at least one plugin")
	})
}

func TestStore_ErrorMethods(t *testing.T) {
	t.Run("should_track_plugin_errors", func(t *testing.T) {
		result := Discover()
		store := result.Store
		
		errors := store.Errors()
		assert.NotNil(t, errors)
		
		// Test HasErrors method
		for pluginID := range errors {
			assert.True(t, store.HasErrors(pluginID))
		}
		
		// Test with non-existent plugin
		assert.False(t, store.HasErrors("non_existent_plugin"))
	})
}

func TestLoadOne(t *testing.T) {
	t.Run("should_handle_missing_manifest_file", func(t *testing.T) {
		plugin, errs := loadOne("/nonexistent/plugin.json")
		
		assert.Empty(t, plugin.Manifest.ID)
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[0].Error(), "failed to read manifest")
	})
}

func TestFieldStruct(t *testing.T) {
	t.Run("should_support_field_properties", func(t *testing.T) {
		field := Field{
			Key:     "accent",
			Label:   "Accent Color",
			Type:    "color",
			Default: "#89b4fa",
			Help:    "Highlight and focus color",
		}
		
		assert.Equal(t, "accent", field.Key)
		assert.Equal(t, "color", field.Type)
		assert.Equal(t, "#89b4fa", field.Default)
		assert.Equal(t, "Highlight and focus color", field.Help)
	})
}