package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateFunctions(t *testing.T) {
	t.Run("should_convert_hex_to_rgba", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "rgba.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)

		template := `{{hexToRGBA .color 1.0}}`

		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)

		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)

		result, err := engine.Render(map[string]interface{}{
			"color": "#ff0080",
		})

		require.NoError(t, err)
		assert.Equal(t, "rgba(255,0,128,1.0)", result)
	})

	t.Run("should_convert_hex_to_rgba_with_alpha", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "rgba_alpha.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)

		template := `{{hexToRGBA .color 0.8}}`

		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)

		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)

		result, err := engine.Render(map[string]interface{}{
			"color": "#cdd6f4",
		})

		require.NoError(t, err)
		assert.Equal(t, "rgba(205,214,244,0.8)", result)
	})

	t.Run("should_apply_alpha_function", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "alpha.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)

		template := `{{alpha .color 0.5}}`

		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)

		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)

		result, err := engine.Render(map[string]interface{}{
			"color": "#ff0080",
		})

		require.NoError(t, err)
		assert.Equal(t, "rgba(255,0,128,0.5)", result)
	})

	t.Run("should_brighten_color", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "brighten.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)

		template := `{{brighten .color 0.3}}`

		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)

		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)

		result, err := engine.Render(map[string]interface{}{
			"color": "#808080",
		})

		require.NoError(t, err)
		// #808080 (128,128,128) brightened by 0.3 towards white (255,255,255)
		// 128 + (255-128) * 0.3 = 128 + 38.1 = 166 (0xa6)
		assert.Equal(t, "#a6a6a6", result)
	})

	t.Run("should_mix_two_colors", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "mix.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)

		template := `{{mix .color1 .color2 0.5}}`

		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)

		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)

		result, err := engine.Render(map[string]interface{}{
			"color1": "#000000",
			"color2": "#ffffff",
		})

		require.NoError(t, err)
		// 50% mix of black (0,0,0) and white (255,255,255)
		// 0 * 0.5 + 255 * 0.5 = 127.5 ≈ 127 (0x7f)
		assert.Equal(t, "#7f7f7f", result)
	})

	t.Run("should_handle_string_functions", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := "strings.tmpl"
		templatePath := filepath.Join(tmpDir, templateFile)

		template := `{{trimPrefix .value "prefix-"}}`

		err := os.WriteFile(templatePath, []byte(template), 0644)
		require.NoError(t, err)

		engine, err := NewTemplateEngine(tmpDir, templateFile)
		require.NoError(t, err)

		result, err := engine.Render(map[string]interface{}{
			"value": "prefix-content",
		})

		require.NoError(t, err)
		assert.Equal(t, "content", result)
	})
}

