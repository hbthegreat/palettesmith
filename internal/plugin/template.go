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
			// Apply alpha to color (return as rgba)
			hex := strings.TrimPrefix(color, "#")
			if len(hex) != 6 {
				return color // fallback to original
			}
			
			r, _ := strconv.ParseInt(hex[0:2], 16, 64)
			g, _ := strconv.ParseInt(hex[2:4], 16, 64)
			b, _ := strconv.ParseInt(hex[4:6], 16, 64)
			
			return fmt.Sprintf("rgba(%d,%d,%d,%.1f)", r, g, b, opacity)
		},
		"brighten": func(color string, amount float64) string {
			// Brighten color by amount (0.0-1.0)
			hex := strings.TrimPrefix(color, "#")
			if len(hex) != 6 {
				return color // fallback to original
			}
			
			r, _ := strconv.ParseInt(hex[0:2], 16, 64)
			g, _ := strconv.ParseInt(hex[2:4], 16, 64)
			b, _ := strconv.ParseInt(hex[4:6], 16, 64)
			
			// Brighten by interpolating toward white
			r = int64(float64(r) + (255-float64(r))*amount)
			g = int64(float64(g) + (255-float64(g))*amount)
			b = int64(float64(b) + (255-float64(b))*amount)
			
			// Clamp to 0-255
			if r > 255 { r = 255 }
			if g > 255 { g = 255 }
			if b > 255 { b = 255 }
			
			return fmt.Sprintf("#%02x%02x%02x", r, g, b)
		},
		"mix": func(color1, color2 string, ratio float64) string {
			// Mix two colors by ratio (0.0 = color1, 1.0 = color2)
			hex1 := strings.TrimPrefix(color1, "#")
			hex2 := strings.TrimPrefix(color2, "#")
			if len(hex1) != 6 || len(hex2) != 6 {
				return color1 // fallback to first color
			}
			
			r1, _ := strconv.ParseInt(hex1[0:2], 16, 64)
			g1, _ := strconv.ParseInt(hex1[2:4], 16, 64)
			b1, _ := strconv.ParseInt(hex1[4:6], 16, 64)
			
			r2, _ := strconv.ParseInt(hex2[0:2], 16, 64)
			g2, _ := strconv.ParseInt(hex2[2:4], 16, 64)
			b2, _ := strconv.ParseInt(hex2[4:6], 16, 64)
			
			// Linear interpolation
			r := int64(float64(r1)*(1-ratio) + float64(r2)*ratio)
			g := int64(float64(g1)*(1-ratio) + float64(g2)*ratio)
			b := int64(float64(b1)*(1-ratio) + float64(b2)*ratio)
			
			return fmt.Sprintf("#%02x%02x%02x", r, g, b)
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