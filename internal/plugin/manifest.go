package plugin

import (
    "io/fs"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"

    "palettesmith/internal/config"
)

// LoadManifestFromFile loads a YAML plugin manifest from disk.
func LoadManifestFromFile(path string) (*PluginManifest, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var m PluginManifest
    if err := yaml.Unmarshal(b, &m); err != nil {
        return nil, err
    }
    // expand paths in Files and Detection rules
    for i := range m.Files {
        m.Files[i].Path = config.ExpandPath(m.Files[i].Path)
    }
    for i := range m.Detection.ConfigExists {
        m.Detection.ConfigExists[i] = config.ExpandPath(m.Detection.ConfigExists[i])
    }
    if err := m.Validate(); err != nil {
        return nil, err
    }
    return &m, nil
}

// DiscoverExternalManifests scans the default external plugin dir for *.yaml files
// and returns their parsed manifests.
func DiscoverExternalManifests() ([]*PluginManifest, error) {
    dir := config.ExternalPluginsDir()
    var out []*PluginManifest
    _ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
        if err != nil { return nil }
        if d.IsDir() { return nil }
        if filepath.Ext(p) != ".yaml" && filepath.Ext(p) != ".yml" {
            return nil
        }
        m, err := LoadManifestFromFile(p)
        if err == nil {
            out = append(out, m)
        }
        return nil
    })
    return out, nil
}

