package tui

import (
	"palettesmith/internal/plugin"
	"palettesmith/internal/theme"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormValidation(t *testing.T) {
	t.Run("should_use_validation_service_for_color_validation", func(t *testing.T) {
		spec := plugin.Spec{
			ID:    "test",
			Title: "Test Plugin",
			Fields: []plugin.Field{
				{
					Key:     "accent",
					Label:   "Accent Color",
					Type:    "color",
					Default: "#89b4fa",
					Help:    "Primary accent color",
				},
			},
		}

		th := theme.NewStore(theme.ThemeConfig{
			ThemeDefaults: map[string]string{},
		})

		form := newFormFromSpec(spec, "test", th)

		// Test that validator is initialized
		assert.NotNil(t, form.validator, "Form should have validation service")
		assert.Equal(t, 1, len(form.fields), "Should have one form field")

		// Simulate user input with valid color
		form.fields[0].input.SetValue("#ff0080")
		
		// Manually trigger validation (normally done in Update)
		errors := form.validator.ValidateField(form.fields[0].spec, form.fields[0].input.Value())
		
		assert.Empty(t, errors, "Valid color should pass validation")
	})

	t.Run("should_show_error_for_invalid_color", func(t *testing.T) {
		spec := plugin.Spec{
			ID:    "test",
			Title: "Test Plugin", 
			Fields: []plugin.Field{
				{
					Key:     "accent",
					Label:   "Accent Color",
					Type:    "color",
					Default: "#89b4fa",
				},
			},
		}

		form := newFormFromSpec(spec, "test", nil)

		// Test invalid color
		errors := form.validator.ValidateField(form.fields[0].spec, "not-a-color")
		
		assert.NotEmpty(t, errors, "Invalid color should fail validation")
		assert.Contains(t, errors[0].Message, "Invalid color format", "Should provide helpful error message")
	})

	t.Run("should_validate_number_fields", func(t *testing.T) {
		min := 0.0
		max := 100.0
		
		spec := plugin.Spec{
			ID:    "test",
			Title: "Test Plugin",
			Fields: []plugin.Field{
				{
					Key:     "opacity",
					Label:   "Opacity",
					Type:    "number", 
					Default: "50",
					Min:     &min,
					Max:     &max,
				},
			},
		}

		form := newFormFromSpec(spec, "test", nil)

		// Test valid number
		errors := form.validator.ValidateField(form.fields[0].spec, "75")
		assert.Empty(t, errors, "Valid number should pass validation")

		// Test invalid number
		errors = form.validator.ValidateField(form.fields[0].spec, "not-a-number")
		assert.NotEmpty(t, errors, "Invalid number should fail validation")

		// Test number below minimum
		errors = form.validator.ValidateField(form.fields[0].spec, "-10")
		assert.NotEmpty(t, errors, "Number below minimum should fail validation")

		// Test number above maximum
		errors = form.validator.ValidateField(form.fields[0].spec, "150")
		assert.NotEmpty(t, errors, "Number above maximum should fail validation")
	})
}

func TestIsValidColor(t *testing.T) {
	t.Run("should_validate_color_formats", func(t *testing.T) {
		validColors := []string{
			"#ff0080",
			"#FFFFFF",
			"#000000",
			"#123456",
			"ff0080",    // hex without #
			"0xff0080",  // hex with 0x prefix
			"#fff",      // hex3 format
		}

		for _, color := range validColors {
			assert.True(t, isValidColor(color), "Color %s should be valid", color)
		}

		invalidColors := []string{
			"not-a-color",
			"#gg0080",
			"#ff008",
			"",
		}

		for _, color := range invalidColors {
			assert.False(t, isValidColor(color), "Color %s should be invalid", color)
		}
	})
}

func TestFormFromSpec(t *testing.T) {
	t.Run("should_create_form_with_validation_service", func(t *testing.T) {
		spec := plugin.Spec{
			ID:    "test",
			Title: "Test Plugin",
			Fields: []plugin.Field{
				{
					Key:     "bg",
					Label:   "Background",
					Type:    "color",
					Default: "#1e1e2e",
				},
				{
					Key:     "fg", 
					Label:   "Foreground",
					Type:    "color",
					Default: "#cdd6f4",
				},
			},
		}

		form := newFormFromSpec(spec, "test", nil)

		require.NotNil(t, form.validator, "Form should have validator")
		assert.Equal(t, 2, len(form.fields), "Should have two fields")
		assert.Equal(t, "test", form.pluginID, "Should set plugin ID")
		assert.Equal(t, 10, form.labelW, "Should calculate label width correctly")
	})
}