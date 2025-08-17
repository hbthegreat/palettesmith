package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateEngine(t *testing.T) {
	t.Run("should_create_template_engine_from_file", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "test.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)
		
		err := os.WriteFile(templatePath, []byte("Hello {{.name}}!"), 0644)
		require.NoError(t, err)
		
		engine, err := NewTemplateEngine(tmpDir, templateFile)
		
		require.NoError(t, err)
		assert.NotNil(t, engine)
	})

	t.Run("should_return_error_for_missing_template_file", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		engine, err := NewTemplateEngine(tmpDir, "nonexistent.tmpl")
		
		assert.Error(t, err)
		assert.Nil(t, engine)
		assert.Contains(t, err.Error(), "failed to read template file")
	})

	t.Run("should_return_error_for_invalid_template_syntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "invalid.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)
		
		err := os.WriteFile(templatePath, []byte("{{.invalid syntax}}"), 0644)
		require.NoError(t, err)
		
		engine, err := NewTemplateEngine(tmpDir, templateFile)
		
		assert.Error(t, err)
		assert.Nil(t, engine)
		assert.Contains(t, err.Error(), "failed to parse template")
	})
}

func TestTemplateEngine_Render(t *testing.T) {
	t.Run("should_render_simple_template", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "simple.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)
		
		err := os.WriteFile(templatePath, []byte("Hello {{.name}}!"), 0644)
		require.NoError(t, err)
		
		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)
		
		result, err := engine.Render(map[string]interface{}{
			"name": "World",
		})
		
		require.NoError(t, err)
		assert.Equal(t, "Hello World!", result)
	})

	t.Run("should_render_color_template", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "colors.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)
		
		template := `[colors]
background = "{{.bg}}"
foreground = "{{.fg}}"
accent = "{{.accent}}"`
		
		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)
		
		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)
		
		result, err := engine.Render(map[string]interface{}{
			"bg":     "#24273a",
			"fg":     "#cad3f5", 
			"accent": "#8aadf4",
		})
		
		require.NoError(t, err)
		expected := `[colors]
background = "#24273a"
foreground = "#cad3f5"
accent = "#8aadf4"`
		assert.Equal(t, expected, result)
	})

	t.Run("should_handle_missing_template_variables", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "missing.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)
		
		err := os.WriteFile(templatePath, []byte("{{.missing}}"), 0644)
		require.NoError(t, err)
		
		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)
		
		result, err := engine.Render(map[string]interface{}{})
		
		require.NoError(t, err)
		assert.Equal(t, "<no value>", result)
	})

	t.Run("should_return_error_for_nil_engine", func(t *testing.T) {
		var engine *TemplateEngine
		
		result, err := engine.Render(map[string]interface{}{})
		
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "template engine not initialized")
	})
}

func TestBuildFieldData(t *testing.T) {
	t.Run("should_build_data_from_fields_and_values", func(t *testing.T) {
		fields := []Field{
			{Key: "bg", Default: "#000000"},
			{Key: "fg", Default: "#ffffff"},
			{Key: "accent", Default: "#0088ff"},
		}
		
		values := map[string]string{
			"bg": "#24273a",
			"fg": "#cad3f5",
		}
		
		data := BuildFieldData(fields, values)
		
		assert.Equal(t, "#24273a", data["bg"])
		assert.Equal(t, "#cad3f5", data["fg"])
		assert.Equal(t, "#0088ff", data["accent"]) // should use default
	})

	t.Run("should_use_defaults_for_empty_values", func(t *testing.T) {
		fields := []Field{
			{Key: "bg", Default: "#000000"},
			{Key: "fg", Default: "#ffffff"},
		}
		
		values := map[string]string{
			"bg": "",
			"fg": "#custom",
		}
		
		data := BuildFieldData(fields, values)
		
		assert.Equal(t, "#000000", data["bg"]) // should use default
		assert.Equal(t, "#custom", data["fg"])
	})
}