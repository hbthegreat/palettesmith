package plugin

import (
	"bytes"
	"fmt"
	"os"
	"palettesmith/internal/color"
	"path/filepath"
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
		"alpha": func(colorStr string, opacity float64) string {
			// Clamp opacity to [0.0, 1.0]
			if opacity < 0.0 {
				opacity = 0.0
			} else if opacity > 1.0 {
				opacity = 1.0
			}
			
			// Use color converter for robust parsing
			converter, err := color.NewConverter(colorStr)
			if err != nil {
				return colorStr // fallback to original
			}
			
			return converter.ToRGBAString(opacity)
		},
		"brighten": func(colorStr string, amount float64) string {
			// Clamp amount to [0.0, 1.0]
			if amount < 0.0 {
				amount = 0.0
			} else if amount > 1.0 {
				amount = 1.0
			}
			
			// Use color converter for robust parsing
			converter, err := color.NewConverter(colorStr)
			if err != nil {
				return colorStr // fallback to original
			}
			
			return converter.Brighten(amount)
		},
		"mix": func(color1, color2 string, ratio float64) string {
			// Clamp ratio to [0.0, 1.0]
			if ratio < 0.0 {
				ratio = 0.0
			} else if ratio > 1.0 {
				ratio = 1.0
			}
			
			// Use color converter for robust parsing
			conv1, err1 := color.NewConverter(color1)
			conv2, err2 := color.NewConverter(color2)
			if err1 != nil || err2 != nil {
				return color1 // fallback to first color
			}
			
			return conv1.Mix(conv2, ratio)
		},
		"hexToRGBA": func(hexStr string, alpha float64) string {
			// Clamp alpha to [0.0, 1.0]
			if alpha < 0.0 {
				alpha = 0.0
			} else if alpha > 1.0 {
				alpha = 1.0
			}
			
			// Use color converter for robust parsing
			converter, err := color.NewConverter(hexStr)
			if err != nil {
				return "rgba(0,0,0,1.0)" // fallback
			}
			
			return converter.ToRGBAString(alpha)
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