# PR1: Project Restructure + Config System

## Overview
Restructure project to separate UI from business logic and add configuration system for different target environments.

## Acceptance Criteria

### Directory Restructure
- [ ] Move `internal/tui/` to `ui/tui/`
- [ ] Update all import paths in moved files
- [ ] Verify application still builds and runs
- [ ] Update go.mod if needed

### Configuration System
- [ ] Create `internal/config/` package
- [ ] Implement config loading from `~/.config/palettesmith/config.json`
- [ ] Add "omarchy" preset that auto-configures paths
- [ ] Create config on first run if it doesn't exist
- [ ] Add validation for config file format

### Config Schema
```json
{
  "target_theme_dir": "~/.config/omarchy/themes",
  "current_theme_link": "~/.config/omarchy/current/theme",
  "preset": "omarchy",
  "staging_dir": "~/.config/palettesmith/staging"
}
```

## Files to Create

### `internal/config/config.go`
```go
package config

type Config struct {
    TargetThemeDir    string `json:"target_theme_dir"`
    CurrentThemeLink  string `json:"current_theme_link"`
    Preset           string `json:"preset"`
    StagingDir       string `json:"staging_dir"`
}

type Manager struct {
    cfg Config
}

func NewManager() (*Manager, error) {
    // Load or create config
}

func (m *Manager) GetConfig() Config {
    return m.cfg
}

func (m *Manager) SetPreset(preset string) error {
    // Apply preset configurations
}
```

### Presets to Support
- **omarchy**: Uses `~/.config/omarchy/themes` and `~/.config/omarchy/current/theme`
- **generic**: Uses `~/.config/palettesmith/themes` only
- **custom**: User-defined paths

## Files to Modify

### `ui/tui/app.go` (moved from `internal/tui/app.go`)
- Update package declaration to `package tui`
- Update imports to use `internal/config`
- Add config manager dependency

### `cmd/palettesmith/main.go`
- Update import path from `palettesmith/internal/tui` to `palettesmith/ui/tui`
- Initialize config manager before TUI

## Testing Requirements

### Unit Tests
- [ ] Test config loading from file
- [ ] Test config creation when file doesn't exist
- [ ] Test preset application (omarchy, generic, custom)
- [ ] Test path expansion (`~` to home directory)
- [ ] Test invalid JSON handling

### Integration Tests
- [ ] Test full application startup with config
- [ ] Verify omarchy preset creates correct paths
- [ ] Test config persistence across restarts

## Manual Testing Steps

1. **Clean Environment Test**
   - Remove `~/.config/palettesmith/` if it exists
   - Run application
   - Verify config file created with omarchy preset
   - Check that directories are created correctly

2. **Existing Config Test**
   - Create config file manually
   - Run application  
   - Verify config is loaded correctly
   - Modify config and restart to confirm persistence

3. **Path Expansion Test**
   - Set config with `~` in paths
   - Verify paths are expanded to actual home directory
   - Test with different users if possible

## Error Handling

### Required Error Messages
- "Config file invalid JSON: {error details}"
- "Cannot create config directory: {path}"
- "Cannot write config file: {error}"
- "Unknown preset: {preset name}"

### Graceful Degradation
- If config loading fails, use sensible defaults
- If directory creation fails, show clear error and exit
- Never start with broken configuration

## Dependencies
- None (this is the foundation PR)

## Definition of Done
- [ ] All tests pass
- [ ] Application builds and runs with restructured code
- [ ] Config system works with omarchy preset
- [ ] Manual testing completed successfully
- [ ] Code review completed
- [ ] Documentation updated