package plugin

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type TemplateEngine struct {
	template *template.Template
}

func NewTemplateEngine(pluginDir, templateFile string) (*TemplateEngine, error) {
	templatePath := filepath.Join(pluginDir, templateFile)
	
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}
	
	funcMap := template.FuncMap{
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"alpha": func(color string, opacity float64) string {
			// For now, just return the color as-is
			// TODO: Implement actual alpha blending
			return color
		},
		"brighten": func(color string, amount float64) string {
			// For now, just return the color as-is  
			// TODO: Implement color brightening
			return color
		},
		"mix": func(color1, color2 string, ratio float64) string {
			// For now, just return the first color
			// TODO: Implement color mixing
			return color1
		},
		"hexToRGBA": func(hex string, alpha float64) string {
			// Convert #ffffff to rgba(255,255,255,1.0)
			hex = strings.TrimPrefix(hex, "#")
			if len(hex) != 6 {
				return "rgba(0,0,0,1.0)" // fallback
			}
			
			r, _ := strconv.ParseInt(hex[0:2], 16, 64)
			g, _ := strconv.ParseInt(hex[2:4], 16, 64)
			b, _ := strconv.ParseInt(hex[4:6], 16, 64)
			
			return fmt.Sprintf("rgba(%d,%d,%d,%.1f)", r, g, b, alpha)
		},
	}
	
	tmpl, err := template.New(templateFile).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	
	return &TemplateEngine{template: tmpl}, nil
}

func (te *TemplateEngine) Render(data map[string]interface{}) (string, error) {
	if te == nil || te.template == nil {
		return "", fmt.Errorf("template engine not initialized")
	}
	
	var buf bytes.Buffer
	if err := te.template.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	
	return buf.String(), nil
}

func BuildFieldData(fields []Field, values map[string]string) map[string]interface{} {
	data := make(map[string]interface{})
	
	for _, field := range fields {
		value := values[field.Key]
		if value == "" {
			value = field.Default
		}
		data[field.Key] = value
	}
	
	return data
}