package plugin

// PluginRegistry keeps plugins by ID and caches detection results.
type PluginRegistry struct {
    plugins map[string]Plugin
    detected map[string]bool
}

func NewRegistry() *PluginRegistry {
    return &PluginRegistry{plugins: map[string]Plugin{}, detected: map[string]bool{}}
}

func (r *PluginRegistry) Register(id string, p Plugin) {
    if r.plugins == nil { r.plugins = map[string]Plugin{} }
    r.plugins[id] = p
}

func (r *PluginRegistry) Get(id string) (Plugin, bool) {
    p, ok := r.plugins[id]
    return p, ok
}

func (r *PluginRegistry) List() map[string]Plugin { return r.plugins }

// DetectAll runs Detect for all plugins and caches the results.
func (r *PluginRegistry) DetectAll() map[string]bool {
    out := make(map[string]bool, len(r.plugins))
    for id, p := range r.plugins {
        ok := p.Detect()
        r.detected[id] = ok
        out[id] = ok
    }
    return out
}

