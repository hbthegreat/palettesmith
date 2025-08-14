// Package plugin discovers and loads plugins from manifests and specs
package plugin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Spec struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
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
	byID map[string]Plugin
	list []Plugin
}

func Discover() (*Store, error) {
	root := filepath.Join(mustGetwd(), "plugins")
	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Store{byID: map[string]Plugin{}, list: nil}, nil
		}
		return nil, err
	}

	seen := map[string]bool{}
	var plugs []Plugin

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		mf := filepath.Join(dir, "plugin.json")
		if _, err := os.Stat(mf); err != nil {
			continue
		}
		p, err := loadOne(mf)
		if err != nil {
			// ignore bad plugins
			continue
		}
		if seen[p.Manifest.ID] {
			continue
		}
		seen[p.Manifest.ID] = true
		plugs = append(plugs, p)
	}

	by := make(map[string]Plugin, len(plugs))
	for _, p := range plugs {
		by[p.Manifest.ID] = p
	}
	return &Store{byID: by, list: plugs}, nil
}

func loadOne(manifestPath string) (Plugin, error) {
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return Plugin{}, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return Plugin{}, err
	}
	m.Dir = filepath.Dir(manifestPath)
	if m.ID == "" || m.SpecRelPath == "" {
		return Plugin{}, errors.New("invalid plugin manifest (missing id/spec)")
	}
	specPath := filepath.Join(m.Dir, filepath.FromSlash(m.SpecRelPath))
	sb, err := os.ReadFile(specPath)
	if err != nil {
		return Plugin{}, err
	}
	var s Spec
	if err := json.Unmarshal(sb, &s); err != nil {
		return Plugin{}, err
	}

	m.ID = strings.ToLower(m.ID)
	if s.ID == "" {
		s.ID = m.ID
	}
	return Plugin{m, s}, nil
}

func (st *Store) List() []Plugin               { return st.list }
func (st *Store) Get(id string) (Plugin, bool) { p, ok := st.byID[strings.ToLower(id)]; return p, ok }

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
