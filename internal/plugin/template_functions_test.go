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
}