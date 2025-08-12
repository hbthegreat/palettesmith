package plugin

import (
    "errors"
    "fmt"
)

// Basic metadata for a plugin as declared in its manifest.
type Metadata struct {
    ID          string `yaml:"id" json:"id"`
    Name        string `yaml:"name" json:"name"`
    Version     string `yaml:"version" json:"version"`
    Description string `yaml:"description" json:"description"`
    Author      string `yaml:"author" json:"author"`
}

type Detection struct {
    ConfigExists   []string `yaml:"config_exists" json:"config_exists"`
    BinaryExists   string   `yaml:"binary_exists" json:"binary_exists"`
    ProcessRunning string   `yaml:"process_running" json:"process_running"`
}

type FileConfig struct {
    Path     string `yaml:"path" json:"path"`
    Format   string `yaml:"format" json:"format"`
    Parser   string `yaml:"parser" json:"parser"`
    Backup   bool   `yaml:"backup" json:"backup"`
    Optional bool   `yaml:"optional" json:"optional"`
}

// ColorDefinition intentionally remains flexible across formats.
// For v1 we only validate presence of ID and at least one targeting hint.
type ColorDefinition struct {
    ID           string   `yaml:"id" json:"id"`
    Label        string   `yaml:"label" json:"label"`
    Description  string   `yaml:"description" json:"description"`
    Default      string   `yaml:"default" json:"default"`
    CSSVariables []string `yaml:"css_variables" json:"css_variables"`
    TOMLPath     string   `yaml:"toml_path" json:"toml_path"`
    YAMLPath     string   `yaml:"yaml_path" json:"yaml_path"`
    JSONPath     string   `yaml:"json_path" json:"json_path"`
    INIKeys      []string `yaml:"ini_keys" json:"ini_keys"`
    HyprVars     []string `yaml:"hypr_variables" json:"hypr_variables"`
}

type RestartConfig struct {
    Method   string `yaml:"method" json:"method"`   // none|signal|command
    Signal   string `yaml:"signal" json:"signal"`
    Process  string `yaml:"process" json:"process"`
    Command  string `yaml:"command" json:"command"`
    Fallback string `yaml:"fallback" json:"fallback"`
}

type PluginManifest struct {
    Metadata  Metadata          `yaml:"metadata" json:"metadata"`
    Detection Detection         `yaml:"detection" json:"detection"`
    Files     []FileConfig      `yaml:"files" json:"files"`
    Colors    []ColorDefinition `yaml:"colors" json:"colors"`
    Restart   RestartConfig     `yaml:"restart" json:"restart"`
}

// Validate performs minimal schema sanity checks.
func (m *PluginManifest) Validate() error {
    if m.Metadata.ID == "" {
        return errors.New("manifest.metadata.id is required")
    }
    if len(m.Files) == 0 {
        return errors.New("manifest.files must not be empty")
    }
    // Ensure each color has an ID and at least one target hint
    for _, c := range m.Colors {
        if c.ID == "" {
            return fmt.Errorf("manifest.colors entry missing id")
        }
        if c.TOMLPath == "" && c.YAMLPath == "" && c.JSONPath == "" && len(c.CSSVariables) == 0 && len(c.INIKeys) == 0 && len(c.HyprVars) == 0 {
            return fmt.Errorf("color %q has no targeting hints", c.ID)
        }
    }
    return nil
}

// Plugin is the runtime behavior for an application integration.
// The state and color types will be defined in internal/core; we avoid import cycles by using interfaces in higher layers later.
type Plugin interface {
    GetManifest() *PluginManifest
    Detect() bool
    GetFiles() []string
}

