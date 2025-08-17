package config_test

import (
	"strings"
	"testing"

	"palettesmith/internal/plugin"
	"palettesmith/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateIntegration(t *testing.T) {
	t.Run("should_render_alacritty_template", func(t *testing.T) {
		result := plugin.Discover()
		store := result.Store
		
		alacritty, found := store.Get("alacritty")
		require.True(t, found, "Alacritty plugin should be found")
		
		engine, err := plugin.NewTemplateEngine(alacritty.Manifest.Dir, alacritty.Spec.TemplateFile)
		require.NoError(t, err)
		
		values := map[string]string{
			"background": "#24273a",
			"foreground": "#cad3f5",
			"red":        "#ed8796",
			"green":      "#a6da95",
		}
		
		data := plugin.BuildFieldData(alacritty.Spec.Fields, values)
		output, err := engine.Render(data)
		
		require.NoError(t, err)
		assert.Contains(t, output, `background = "#24273a"`)
		assert.Contains(t, output, `foreground = "#cad3f5"`)
		assert.Contains(t, output, `red = "#ed8796"`)
		assert.Contains(t, output, `green = "#a6da95"`)
	})

	t.Run("should_render_hyprland_template", func(t *testing.T) {
		result := plugin.Discover()
		store := result.Store
		
		hyprland, found := store.Get("hyprland")
		require.True(t, found, "Hyprland plugin should be found")
		
		engine, err := plugin.NewTemplateEngine(hyprland.Manifest.Dir, hyprland.Spec.TemplateFile)
		require.NoError(t, err)
		
		values := map[string]string{
			"active_border": "#c6d0f5",
		}
		
		data := plugin.BuildFieldData(hyprland.Spec.Fields, values)
		output, err := engine.Render(data)
		
		require.NoError(t, err)
		assert.Contains(t, output, "general {")
		assert.Contains(t, output, "col.active_border = rgb(c6d0f5)")
	})
}

func TestAllPluginsEndToEnd(t *testing.T) {
	t.Run("should_validate_and_render_all_plugins", func(t *testing.T) {
		result := plugin.Discover()
		store := result.Store
		
		// Ensure we have the expected number of plugins
		plugins := store.List()
		assert.GreaterOrEqual(t, len(plugins), 8, "Should discover at least 8 plugins")
		
		// Test data with various color formats
		testValues := map[string]string{
			"background":      "#24273a",
			"foreground":      "#cad3f5", 
			"text":            "#c6d0f5",
			"accent":          "#8aadf4",
			"red":             "#ed8796",
			"green":           "#a6da95", 
			"blue":            "#8aadf4",
			"yellow":          "#eed49f",
			"cyan":            "#91d7e3",
			"magenta":         "#f5bde6",
			"black":           "#181926",
			"white":           "#f4dbd6",
			"active_border":   "#c6d0f5",
			"inactive_border": "#363a4f",
			"cursor":          "#f4dbd6",
			"selection":       "#5b6078",
			"background_color": "#24273a",
			"border_color":    "#c6d0f5",
			"text_color":      "#cad3f5",
			"selected_text":   "#8caaee",
			"base":            "#24273a",
			"border":          "#c6d0f5",
		}
		
		for _, p := range plugins {
			t.Run(p.Manifest.ID, func(t *testing.T) {
				// Test plugin validation
				validator := validation.NewPluginValidator()
				pluginErrors := validator.ValidatePlugin(p)
				assert.Empty(t, pluginErrors, "Plugin %s should be valid", p.Manifest.ID)
				
				// Test field validation
				for _, field := range p.Spec.Fields {
					if field.Type == "color" {
						// Test with valid color
						testColor := testValues[field.Key]
						if testColor == "" {
							testColor = field.Default
						}
						
						fieldErrors := validator.ValidateField(field, testColor)
						assert.Empty(t, fieldErrors, "Field %s should validate color %s", field.Key, testColor)
						
						// Test with invalid color
						invalidErrors := validator.ValidateField(field, "not-a-color")
						assert.NotEmpty(t, invalidErrors, "Field %s should reject invalid color", field.Key)
					}
				}
				
				// Test template rendering
				engine, err := plugin.NewTemplateEngine(p.Manifest.Dir, p.Spec.TemplateFile)
				require.NoError(t, err, "Should create template engine for %s", p.Manifest.ID)
				
				data := plugin.BuildFieldData(p.Spec.Fields, testValues)
				output, err := engine.Render(data)
				require.NoError(t, err, "Should render template for %s", p.Manifest.ID)
				assert.NotEmpty(t, output, "Template output should not be empty for %s", p.Manifest.ID)
				
				// Verify output contains expected color values in appropriate formats
				hasColorOutput := false
				for _, field := range p.Spec.Fields {
					if field.Type == "color" {
						testColor := testValues[field.Key]
						if testColor == "" {
							testColor = field.Default
						}
						if testColor != "" {
							// Check if output contains color in any expected format
							foundColorFormat := false
							
							// Check original hex format
							if strings.Contains(output, testColor) {
								foundColorFormat = true
							}
							
							// Check converted formats based on plugin
							switch p.Manifest.ID {
							case "hyprland":
								// Converts #c6d0f5 -> rgb(c6d0f5)
								hexNoPrefix := strings.TrimPrefix(testColor, "#")
								if strings.Contains(output, "rgb("+hexNoPrefix+")") {
									foundColorFormat = true
								}
							case "hyprlock":
								// Converts #c6d0f5 -> rgba(198,208,245,1.0)
								if strings.Contains(output, "rgba(") {
									foundColorFormat = true
								}
							}
							
							if foundColorFormat {
								hasColorOutput = true
							}
						}
					}
				}
				assert.True(t, hasColorOutput, "Plugin %s should have color output in some format", p.Manifest.ID)
			})
		}
	})
}

func TestTemplateFunctionIntegration(t *testing.T) {
	t.Run("should_use_template_functions_in_real_plugins", func(t *testing.T) {
		result := plugin.Discover()
		store := result.Store
		
		// Test hexToRGBA function with hyprlock (uses rgba format)
		hyprlock, found := store.Get("hyprlock")
		if found {
			engine, err := plugin.NewTemplateEngine(hyprlock.Manifest.Dir, hyprlock.Spec.TemplateFile)
			require.NoError(t, err)
			
			values := map[string]string{
				"bg_color": "#181824",
			}
			
			data := plugin.BuildFieldData(hyprlock.Spec.Fields, values)
			output, err := engine.Render(data)
			
			require.NoError(t, err)
			// Should contain rgba format from hexToRGBA template function
			assert.Contains(t, output, "rgba(24,24,36,", "Should convert hex to rgba format")
		}
	})
}