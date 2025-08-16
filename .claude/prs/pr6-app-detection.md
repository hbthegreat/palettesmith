# PR6: App Auto-Detection + Plugin Loading

## Overview
Implement automatic detection of installed applications and intelligent plugin loading with installation status indicators.

## Acceptance Criteria

### Application Detection
- [ ] Detect installed applications by checking common paths
- [ ] Support multiple detection methods (binary, config files, package managers)
- [ ] Cache detection results for performance
- [ ] Update detection status dynamically
- [ ] Handle different Linux distributions

### Plugin Loading Logic
- [ ] Show only plugins for detected applications by default
- [ ] Provide "show all" mode with installation badges
- [ ] Allow manual plugin enable/disable
- [ ] Sort plugins by installation status
- [ ] Remember user preferences for plugin visibility

### Status Indicators
- [ ] "Installed" - App detected and ready
- [ ] "Not Detected" - App not found on system
- [ ] "Config Missing" - App installed but no config found
- [ ] "Disabled" - User manually disabled plugin
- [ ] "Error" - Plugin failed to load

## Files to Create

### `internal/detection/detector.go`
```go
package detection

type Detector struct {
    cache    map[string]DetectionResult
    methods  []DetectionMethod
}

type DetectionResult struct {
    AppID         string    `json:"app_id"`
    Installed     bool      `json:"installed"`
    Version       string    `json:"version,omitempty"`
    ConfigPaths   []string  `json:"config_paths"`
    BinaryPaths   []string  `json:"binary_paths"`
    DetectedBy    string    `json:"detected_by"`
    LastChecked   time.Time `json:"last_checked"`
}

type DetectionMethod interface {
    Name() string
    Detect(appID string, manifest plugin.Manifest) DetectionResult
}

func NewDetector() *Detector {
    // Initialize with detection methods
}

func (d *Detector) DetectApp(appID string, manifest plugin.Manifest) DetectionResult {
    // Run detection for specific app
}

func (d *Detector) DetectAll(plugins []plugin.Plugin) map[string]DetectionResult {
    // Detect all apps in parallel
}
```

### `internal/detection/methods.go`
```go
package detection

// Binary detection method
type BinaryDetection struct{}

func (b BinaryDetection) Name() string { return "binary" }

func (b BinaryDetection) Detect(appID string, manifest plugin.Manifest) DetectionResult {
    // Check if binary exists in PATH
}

// Config file detection method  
type ConfigDetection struct{}

func (c ConfigDetection) Name() string { return "config" }

func (c ConfigDetection) Detect(appID string, manifest plugin.Manifest) DetectionResult {
    // Check if config files exist
}

// Package manager detection method
type PackageDetection struct{}

func (p PackageDetection) Name() string { return "package" }

func (p PackageDetection) Detect(appID string, manifest plugin.Manifest) DetectionResult {
    // Check package manager databases (apt, pacman, etc.)
}
```

### `internal/plugin/manager.go`
```go
package plugin

type Manager struct {
    store         *Store
    detector      *detection.Detector
    preferences   *UserPreferences
    showAll       bool
}

type UserPreferences struct {
    EnabledPlugins  map[string]bool `json:"enabled_plugins"`
    ShowAll         bool            `json:"show_all"`
    AutoDetect      bool            `json:"auto_detect"`
}

func NewManager(store *Store, detector *detection.Detector) *Manager {
    // Initialize plugin manager
}

func (m *Manager) GetVisiblePlugins() []PluginStatus {
    // Return plugins that should be shown based on detection + preferences
}

func (m *Manager) SetPluginEnabled(pluginID string, enabled bool) error {
    // Enable/disable plugin manually
}

func (m *Manager) SetShowAll(showAll bool) {
    // Toggle between detected-only and all plugins
}
```

### `internal/plugin/status.go`
```go
package plugin

type PluginStatus struct {
    Plugin     Plugin                  `json:"plugin"`
    Detection  detection.DetectionResult `json:"detection"`
    Enabled    bool                    `json:"enabled"`
    Status     Status                  `json:"status"`
    Error      string                  `json:"error,omitempty"`
}

type Status string

const (
    StatusInstalled    Status = "installed"
    StatusNotDetected  Status = "not_detected"
    StatusConfigMissing Status = "config_missing"
    StatusDisabled     Status = "disabled"
    StatusError        Status = "error"
)

func (s Status) Badge() string {
    // Return colored badge for TUI display
}

func (s Status) Description() string {
    // Return human-readable description
}
```

## Files to Modify

### `ui/tui/sidebar.go`
- Add installation status badges to plugin list
- Show/hide plugins based on detection status
- Add toggle for "show all" mode
- Sort plugins by status (installed first)

### `ui/tui/app.go`
- Integrate plugin manager
- Add commands for toggling plugin visibility
- Show detection status in status bar
- Add refresh detection command

### Enhanced Plugin Manifest

### `plugin.json` (Detection Enhancement)
```json
{
  "id": "hyprland",
  "title": "Hyprland",
  "detection": {
    "binary_names": ["hyprland", "Hyprland"],
    "config_paths": ["~/.config/hypr/hyprland.conf"],
    "package_names": {
      "arch": ["hyprland"],
      "ubuntu": ["hyprland"],
      "fedora": ["hyprland"]
    },
    "version_command": ["hyprland", "--version"]
  },
  "spec": "spec.json",
  "user_paths": ["~/.config/hypr/hyprland.conf"],
  "system_paths": ["/etc/xdg/hypr/hyprland.conf"],
  "reload": ["hyprctl", "reload"]
}
```

## Testing Requirements

### Unit Tests
- [ ] Test each detection method independently
- [ ] Test plugin filtering logic
- [ ] Test user preference persistence
- [ ] Test status badge generation
- [ ] Test concurrent detection

### Integration Tests
- [ ] Test full detection workflow
- [ ] Test plugin manager with real plugins
- [ ] Test preference saving/loading
- [ ] Verify detection caching works

### Mock Tests
Create test environments:
- System with all apps installed
- System with no apps installed
- Mixed installation states
- Different package managers

## Manual Testing Steps

1. **Detection Accuracy Test**
   - Test on system with Hyprland installed
   - Test on system without Hyprland
   - Verify detection methods work correctly
   - Check version detection accuracy

2. **Plugin Visibility Test**
   - Start with default (detected only) view
   - Toggle "show all" mode
   - Manually disable/enable plugins
   - Verify preferences persist across restarts

3. **Status Badge Test**
   - Verify badges display correctly
   - Test badge colors in different terminals
   - Check status descriptions are helpful
   - Test sorting by installation status

4. **Performance Test**
   - Test detection speed with many plugins
   - Verify caching reduces detection time
   - Test concurrent detection doesn't block UI

## Cache Strategy

### Detection Cache
```json
{
  "last_updated": "2024-01-15T10:30:00Z",
  "detection_results": {
    "hyprland": {
      "installed": true,
      "version": "0.34.0",
      "detected_by": "binary",
      "last_checked": "2024-01-15T10:30:00Z"
    }
  }
}
```

### Cache Rules
- Cache detection results for 1 hour
- Invalidate cache when config changes detected
- Allow manual cache refresh
- Cache per-user (different users, different results)

## Error Handling

### Required Error Messages
- "Detection failed for {app}: {error}"
- "Plugin {id} disabled due to detection failure"
- "Cannot refresh detection: {error}"
- "Package manager query failed: {error}"

### Graceful Degradation
- Continue with other plugins if one detection fails
- Show error status but don't crash
- Provide manual override options
- Fall back to basic file existence checks

## Dependencies
- PR1 (Project Restructure) must be completed
- PR2 (Enhanced Plugins) for manifest enhancements

## Definition of Done
- [ ] Detection works for common Linux applications
- [ ] Plugin visibility toggles work correctly
- [ ] Status badges display helpfully
- [ ] Performance acceptable for many plugins
- [ ] User preferences persist correctly
- [ ] All tests pass
- [ ] Manual testing on different distributions
- [ ] Code review completed
- [ ] Caching prevents excessive system checks