package validation

import (
	"palettesmith/internal/plugin"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPluginValidator(t *testing.T) {
	t.Run("should_create_new_plugin_validator", func(t *testing.T) {
		validator := NewPluginValidator()
		
		assert.NotNil(t, validator)
	})
}

func TestPluginValidator_ValidateField(t *testing.T) {
	validator := NewPluginValidator()

	t.Run("should_validate_color_field_successfully", func(t *testing.T) {
		field := plugin.Field{
			Key:  "accent",
			Type: "color",
		}
		
		errors := validator.ValidateField(field, "#89b4fa")
		
		assert.Empty(t, errors)
	})

	t.Run("should_return_error_for_invalid_color", func(t *testing.T) {
		field := plugin.Field{
			Key:  "accent",
			Type: "color",
		}
		
		errors := validator.ValidateField(field, "invalid_color")
		
		require.Len(t, errors, 1)
		assert.Equal(t, "accent", errors[0].Field)
		assert.Equal(t, "invalid_color", errors[0].Code)
		assert.Contains(t, errors[0].Message, "Invalid color format")
	})

	t.Run("should_validate_number_field_successfully", func(t *testing.T) {
		field := plugin.Field{
			Key:  "opacity",
			Type: "number",
		}
		
		errors := validator.ValidateField(field, "0.8")
		
		assert.Empty(t, errors)
	})

	t.Run("should_return_error_for_invalid_number", func(t *testing.T) {
		field := plugin.Field{
			Key:  "opacity",
			Type: "number",
		}
		
		errors := validator.ValidateField(field, "not_a_number")
		
		require.Len(t, errors, 1)
		assert.Equal(t, "opacity", errors[0].Field)
		assert.Equal(t, "invalid_number", errors[0].Code)
		assert.Equal(t, "Must be a valid number", errors[0].Message)
	})

	t.Run("should_validate_number_with_min_max_constraints", func(t *testing.T) {
		min := 0.0
		max := 1.0
		field := plugin.Field{
			Key:  "opacity",
			Type: "number",
			Min:  &min,
			Max:  &max,
		}
		
		errors := validator.ValidateField(field, "1.5")
		
		require.Len(t, errors, 1)
		assert.Equal(t, "above_maximum", errors[0].Code)
		assert.Contains(t, errors[0].Message, "Must be at most 1")
	})

	t.Run("should_validate_select_field_successfully", func(t *testing.T) {
		field := plugin.Field{
			Key:  "theme",
			Type: "select",
			Enum: []string{"dark", "light", "auto"},
		}
		
		errors := validator.ValidateField(field, "dark")
		
		assert.Empty(t, errors)
	})

	t.Run("should_return_error_for_invalid_select_option", func(t *testing.T) {
		field := plugin.Field{
			Key:  "theme",
			Type: "select",
			Enum: []string{"dark", "light", "auto"},
		}
		
		errors := validator.ValidateField(field, "invalid_option")
		
		require.Len(t, errors, 1)
		assert.Equal(t, "theme", errors[0].Field)
		assert.Equal(t, "invalid_option", errors[0].Code)
		assert.Contains(t, errors[0].Message, "Must be one of: dark, light, auto")
	})

	t.Run("should_require_field_when_no_default", func(t *testing.T) {
		field := plugin.Field{
			Key:  "required_field",
			Type: "text",
		}
		
		errors := validator.ValidateField(field, "")
		
		require.Len(t, errors, 1)
		assert.Equal(t, "required_field", errors[0].Field)
		assert.Equal(t, "required", errors[0].Code)
		assert.Equal(t, "Field is required", errors[0].Message)
	})

	t.Run("should_allow_empty_field_when_default_exists", func(t *testing.T) {
		field := plugin.Field{
			Key:     "optional_field",
			Type:    "text",
			Default: "default_value",
		}
		
		errors := validator.ValidateField(field, "")
		
		assert.Empty(t, errors)
	})
}

func TestPluginValidator_ValidatePlugin(t *testing.T) {
	validator := NewPluginValidator()

	t.Run("should_validate_complete_plugin_successfully", func(t *testing.T) {
		p := plugin.Plugin{
			Manifest: plugin.Manifest{
				ID:          "test-app",
				Title:       "Test Application",
				SpecRelPath: "spec.json",
			},
			Spec: plugin.Spec{
				ID: "test-app",
				Fields: []plugin.Field{
					{
						Key:  "accent",
						Type: "color",
					},
				},
			},
		}
		
		errors := validator.ValidatePlugin(p)
		
		assert.Empty(t, errors)
	})

	t.Run("should_return_error_for_missing_manifest_id", func(t *testing.T) {
		p := plugin.Plugin{
			Manifest: plugin.Manifest{
				Title:       "Test Application",
				SpecRelPath: "spec.json",
			},
			Spec: plugin.Spec{
				ID: "test-app",
			},
		}
		
		errors := validator.ValidatePlugin(p)
		
		require.NotEmpty(t, errors)
		assert.Contains(t, errors[0].Error(), "plugin ID is required")
	})

	t.Run("should_return_error_for_invalid_field_type", func(t *testing.T) {
		p := plugin.Plugin{
			Manifest: plugin.Manifest{
				ID:          "test-app",
				Title:       "Test Application",
				SpecRelPath: "spec.json",
			},
			Spec: plugin.Spec{
				ID: "test-app",
				Fields: []plugin.Field{
					{
						Key:  "invalid_field",
						Type: "invalid_type",
					},
				},
			},
		}
		
		errors := validator.ValidatePlugin(p)
		
		require.NotEmpty(t, errors)
		assert.Contains(t, errors[0].Error(), "invalid type 'invalid_type'")
	})

	t.Run("should_return_error_for_missing_field_key", func(t *testing.T) {
		p := plugin.Plugin{
			Manifest: plugin.Manifest{
				ID:          "test-app",
				Title:       "Test Application",
				SpecRelPath: "spec.json",
			},
			Spec: plugin.Spec{
				ID: "test-app",
				Fields: []plugin.Field{
					{
						Type: "color",
					},
				},
			},
		}
		
		errors := validator.ValidatePlugin(p)
		
		require.NotEmpty(t, errors)
		assert.Contains(t, errors[0].Error(), "field key is required")
	})
}