package plugin

import (
    "os"
    "os/exec"
)

// SimplePlugin is a minimal implementation for detection and metadata only.
// Deeper color extraction/apply will be implemented in later phases.
type SimplePlugin struct {
    manifest *PluginManifest
}

func (p *SimplePlugin) GetManifest() *PluginManifest { return p.manifest }

func (p *SimplePlugin) Detect() bool {
    // Check config files
    for _, f := range p.manifest.Detection.ConfigExists {
        if _, err := os.Stat(f); err == nil {
            return true
        }
    }
    // Check binary
    if p.manifest.Detection.BinaryExists != "" {
        if _, err := exec.LookPath(p.manifest.Detection.BinaryExists); err == nil {
            return true
        }
    }
    // Process check is skipped in v1 list detection (non-portable, requires /proc scan)
    return false
}

func (p *SimplePlugin) GetFiles() []string {
    files := make([]string, 0, len(p.manifest.Files))
    for _, f := range p.manifest.Files {
        files = append(files, f.Path)
    }
    return files
}

// LoadBuiltinPlugins returns in-binary plugin implementations.
// For v1 we have none; later we may add Waybar/Hyprland/Alacritty baked-in.
func LoadBuiltinPlugins() map[string]Plugin {
    return map[string]Plugin{}
}

// LoadExternalPlugins loads all manifests from the external directory and materializes SimplePlugin instances.
func LoadExternalPlugins() (map[string]Plugin, error) {
    manifests, err := DiscoverExternalManifests()
    if err != nil {
        return nil, err
    }
    out := make(map[string]Plugin, len(manifests))
    for _, m := range manifests {
        out[m.Metadata.ID] = &SimplePlugin{manifest: m}
    }
    return out, nil
}

