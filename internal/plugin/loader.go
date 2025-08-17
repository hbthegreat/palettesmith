// Package plugin discovers and loads plugins from manifests and specs
package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Spec struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	TemplateFile string  `json:"template_file"`
	Fields       []Field `json:"fields"`
}

type Field struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Type    string   `json:"type"` // "color"|"text"|"number"|"select"
	Default string   `json:"default,omitempty"`
	Help    string   `json:"help,omitempty"`
	Min     *float64 `json:"min,omitempty"`
	Max     *float64 `json:"max,omitempty"`
	Enum    []string `json:"enum,omitempty"`
}

type Manifest struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	SpecRelPath string   `json:"spec"` // relative to manifest dir (e.g., "spec.json")
	UserPaths   []string `json:"user_paths,omitempty"`
	SystemPaths []string `json:"system_paths,omitempty"`
	Reload      []string `json:"reload,omitempty"`

	Dir string `json:"-"` // absolute dir of the plugin (filled at load)
}

type Plugin struct {
	Manifest Manifest
	Spec     Spec
}

type Store struct {
	byID   map[string]Plugin
	list   []Plugin
	errors map[string][]error
}

type LoadResult struct {
	Store  *Store
	Errors map[string][]error
}

func Discover() *LoadResult {
	root := filepath.Join(findProjectRoot(), "plugins")
	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &LoadResult{
				Store:  &Store{byID: map[string]Plugin{}, list: nil, errors: map[string][]error{}},
				Errors: map[string][]error{},
			}
		}
		return &LoadResult{
			Store: &Store{byID: map[string]Plugin{}, list: nil, errors: map[string][]error{}},
			Errors: map[string][]error{
				"system": {fmt.Errorf("failed to read plugins directory: %w", err)},
			},
		}
	}

	seen := map[string]bool{}
	var plugs []Plugin
	allErrors := map[string][]error{}
	pluginErrors := map[string][]error{}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		
		pluginID := e.Name()
		dir := filepath.Join(root, pluginID)
		mf := filepath.Join(dir, "plugin.json")
		
		if _, err := os.Stat(mf); err != nil {
			allErrors[pluginID] = append(allErrors[pluginID], 
				fmt.Errorf("plugin.json not found: %w", err))
			continue
		}
		
		p, loadErrs := loadOne(mf)
		if len(loadErrs) > 0 {
			allErrors[pluginID] = append(allErrors[pluginID], loadErrs...)
			pluginErrors[pluginID] = loadErrs
		}
		
		if p.Manifest.ID != "" {
			if seen[p.Manifest.ID] {
				allErrors[pluginID] = append(allErrors[pluginID], 
					fmt.Errorf("duplicate plugin ID: %s", p.Manifest.ID))
				continue
			}
			seen[p.Manifest.ID] = true
			plugs = append(plugs, p)
		}
	}

	by := make(map[string]Plugin, len(plugs))
	for _, p := range plugs {
		by[p.Manifest.ID] = p
	}
	
	return &LoadResult{
		Store: &Store{
			byID:   by, 
			list:   plugs, 
			errors: pluginErrors,
		},
		Errors: allErrors,
	}
}

func loadOne(manifestPath string) (Plugin, []error) {
	var errs []error
	
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return Plugin{}, []error{fmt.Errorf("failed to read manifest: %w", err)}
	}
	
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return Plugin{}, []error{fmt.Errorf("invalid manifest JSON: %w", err)}
	}
	
	m.Dir = filepath.Dir(manifestPath)
	
	if m.ID == "" {
		errs = append(errs, errors.New("manifest missing required field: id"))
	}
	if m.SpecRelPath == "" {
		errs = append(errs, errors.New("manifest missing required field: spec"))
	}
	
	if len(errs) > 0 {
		return Plugin{Manifest: m}, errs
	}
	
	specPath := filepath.Join(m.Dir, filepath.FromSlash(m.SpecRelPath))
	sb, err := os.ReadFile(specPath)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to read spec file: %w", err))
		return Plugin{Manifest: m}, errs
	}
	
	var s Spec
	if err := json.Unmarshal(sb, &s); err != nil {
		errs = append(errs, fmt.Errorf("invalid spec JSON: %w", err))
		return Plugin{Manifest: m}, errs
	}

	m.ID = strings.ToLower(m.ID)
	if s.ID == "" {
		s.ID = m.ID
	}
	
	plugin := Plugin{Manifest: m, Spec: s}
	
	// TODO: Add plugin validation here when we integrate with validation package
	// validator := validation.NewPluginValidator()
	// validationErrs := validator.ValidatePlugin(plugin)
	// errs = append(errs, validationErrs...)
	
	return plugin, errs
}

func (st *Store) List() []Plugin               { return st.list }
func (st *Store) Get(id string) (Plugin, bool) { p, ok := st.byID[strings.ToLower(id)]; return p, ok }
func (st *Store) Errors() map[string][]error   { return st.errors }
func (st *Store) HasErrors(id string) bool     { _, ok := st.errors[id]; return ok }

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

// isProjectRoot checks if the given directory contains project root indicators
func isProjectRoot(dir string) bool {
	// Check for plugins directory
	if _, err := os.Stat(filepath.Join(dir, "plugins")); err == nil {
		return true
	}
	// Check for go.mod file
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return true
	}
	return false
}

func findProjectRoot() string {
	wd := mustGetwd()
	
	// Look for project root by finding go.mod or plugins directory
	current := wd
	for {
		if isProjectRoot(current) {
			return current
		}
		
		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root
			break
		}
		current = parent
	}
	
	// Fallback to current directory
	return wd
}
