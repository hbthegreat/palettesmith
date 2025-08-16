# PR7: Real-time Preview + Hot Reload

## Overview
Implement real-time theme preview with hot reload integration, command confirmation, and seamless integration with system theme mechanisms.

## Acceptance Criteria

### Real-time Preview System
- [ ] Apply changes to staging area immediately
- [ ] Update config files in staging as user types
- [ ] Show live preview without affecting original files
- [ ] Support rapid changes without performance issues
- [ ] Maintain staging/backup separation

### Hot Reload Integration
- [ ] Execute reload commands defined in plugin manifests
- [ ] Show commands before execution with confirmation
- [ ] Support different reload mechanisms per application
- [ ] Handle reload failures gracefully
- [ ] Integrate with omarchy system reload if available

### Command Execution Safety
- [ ] Display exact commands before running
- [ ] Require user confirmation for all commands
- [ ] Run commands with user permissions only
- [ ] Validate commands against allowlist
- [ ] Provide option to run commands manually

### Omarchy Integration
- [ ] Detect omarchy system and use native reload
- [ ] Fall back to individual app reloads if omarchy unavailable
- [ ] Respect omarchy theme switching mechanisms
- [ ] Support omarchy's systemd reload services

## Files to Create

### `internal/preview/manager.go`
```go
package preview

type Manager struct {
    stagingMgr *staging.Manager
    reloadSvc  *ReloadService
    config     *config.Config
    active     bool
}

func NewManager(stagingMgr *staging.Manager, config *config.Config) *Manager {
    // Initialize preview manager
}

func (m *Manager) StartPreview() error {
    // Begin real-time preview mode
}

func (m *Manager) StopPreview() error {
    // End preview mode, keep staging
}

func (m *Manager) ApplyChange(pluginID, fieldKey, value string) error {
    // Apply single field change to staging
}

func (m *Manager) CommitPreview() error {
    // Move staged changes to actual configs
}
```

### `internal/reload/service.go`
```go
package reload

type Service struct {
    omarchyDetected bool
    allowedCommands map[string]bool
}

type ReloadCommand struct {
    PluginID    string   `json:"plugin_id"`
    Command     []string `json:"command"`
    Description string   `json:"description"`
    Required    bool     `json:"required"`
}

func NewService() *Service {
    // Initialize reload service, detect omarchy
}

func (s *Service) GetReloadCommands(plugins []string) []ReloadCommand {
    // Get all reload commands for given plugins
}

func (s *Service) ExecuteReload(cmd ReloadCommand) error {
    // Execute single reload command safely
}

func (s *Service) ReloadOmarchy() error {
    // Use omarchy's system-wide reload if available
}

func (s *Service) ValidateCommand(cmd []string) error {
    // Validate command against security allowlist
}
```

### `internal/reload/omarchy.go`
```go
package reload

type OmarchyIntegration struct {
    available     bool
    reloadService string // systemd service name
    reloadCommand []string
}

func DetectOmarchy() *OmarchyIntegration {
    // Detect if omarchy is installed and configured
}

func (o *OmarchyIntegration) IsAvailable() bool {
    return o.available
}

func (o *OmarchyIntegration) ReloadTheme() error {
    // Use omarchy's theme reload mechanism
}

func (o *OmarchyIntegration) GetReloadCommand() []string {
    // Return the command that would be executed
}
```

### `ui/tui/preview.go`
```go
package tui

type PreviewState struct {
    active         bool
    pendingReloads []reload.ReloadCommand
    stagingDiff    map[string]staging.Diff
}

func (m Model) StartPreview() (Model, tea.Cmd) {
    // Enter preview mode
}

func (m Model) StopPreview() (Model, tea.Cmd) {
    // Exit preview mode
}

func (m Model) ShowReloadConfirmation() string {
    // Display reload commands for user confirmation
}
```

## Files to Modify

### `ui/tui/app.go`
- Add preview mode toggle (P key)
- Show preview status in status bar
- Add reload confirmation dialog
- Integrate with preview manager

### `ui/tui/form.go`
- Apply changes to preview in real-time
- Show staging status for each field
- Update preview as user types (debounced)
- Display reload status

### `internal/staging/manager.go`
- Add real-time staging updates
- Support incremental changes
- Optimize for frequent updates
- Maintain diff accuracy

## Enhanced Plugin Configuration

### `plugin.json` (Reload Enhancement)
```json
{
  "id": "hyprland",
  "title": "Hyprland",
  "reload": {
    "commands": [
      {
        "cmd": ["hyprctl", "reload"],
        "description": "Reload Hyprland configuration",
        "required": true
      },
      {
        "cmd": ["killall", "-SIGUSR2", "waybar"],
        "description": "Reload Waybar (if running)",
        "required": false
      }
    ],
    "omarchy_compatible": true,
    "delay_ms": 100
  },
  "spec": "spec.json",
  "user_paths": ["~/.config/hypr/hyprland.conf"],
  "system_paths": ["/etc/xdg/hypr/hyprland.conf"]
}
```

## User Interface Flow

### Preview Mode Workflow
1. User presses `P` to enter preview mode
2. Status bar shows "PREVIEW MODE - Changes staged"
3. User makes color/field changes
4. Changes apply to staging area immediately
5. User presses `A` to apply (shows reload confirmation)
6. User confirms or cancels reload commands
7. Changes applied to actual configs, apps reloaded

### Reload Confirmation Dialog
```
┌─ Apply Changes ─────────────────────────────┐
│ The following commands will be executed:    │
│                                             │
│ • hyprctl reload                           │
│   Reload Hyprland configuration            │
│                                             │
│ • killall -SIGUSR2 waybar                  │
│   Reload Waybar (if running)               │
│                                             │
│ [A]pply  [C]ancel  [M]anual                │
└─────────────────────────────────────────────┘
```

## Testing Requirements

### Unit Tests
- [ ] Test preview staging updates
- [ ] Test reload command validation
- [ ] Test omarchy detection logic
- [ ] Test command execution safety
- [ ] Test preview state management

### Integration Tests
- [ ] Test full preview workflow
- [ ] Test reload command execution
- [ ] Test omarchy integration (if available)
- [ ] Test error handling during reload
- [ ] Verify staging/backup integrity

### Security Tests
- [ ] Test command validation prevents dangerous commands
- [ ] Test privilege escalation prevention
- [ ] Test command injection prevention
- [ ] Verify commands run with user permissions only

## Manual Testing Steps

1. **Preview Mode Test**
   - Enter preview mode
   - Make several field changes
   - Verify changes appear in staging
   - Check real-time responsiveness
   - Exit preview mode without applying

2. **Reload Integration Test**
   - Apply changes with reload
   - Verify commands shown correctly
   - Test command confirmation flow
   - Check applications reload properly
   - Test manual command execution option

3. **Omarchy Integration Test** (if available)
   - Test on omarchy system
   - Verify omarchy reload is used
   - Test fallback to individual reloads
   - Check theme switching integration

4. **Error Handling Test**
   - Test with failed reload commands
   - Test with missing applications
   - Test with permission denied scenarios
   - Verify graceful degradation

## Security Considerations

### Command Allowlist
```go
var AllowedCommands = map[string]bool{
    "hyprctl":     true,
    "swaymsg":     true,
    "i3-msg":      true,
    "killall":     true, // with restrictions
    "systemctl":   true, // user only
    "pkill":       true, // with restrictions
}
```

### Restrictions
- No `sudo` or privilege escalation
- No shell execution or piping
- No file system modification commands
- No network commands
- User session only, no system-wide changes

## Error Handling

### Required Error Messages
- "Preview mode failed: cannot access staging area"
- "Reload command failed: {command} returned {error}"
- "Command not allowed: {command}"
- "Omarchy reload failed, falling back to individual reloads"

### Recovery Strategies
- Continue with other reload commands if one fails
- Provide manual command execution option
- Show which commands succeeded/failed
- Allow retry of failed commands

## Dependencies
- PR3 (Theme Staging) must be completed
- PR6 (App Detection) recommended for reload integration

## Definition of Done
- [ ] Preview mode provides real-time feedback
- [ ] Reload commands execute safely with confirmation
- [ ] Omarchy integration works when available
- [ ] All security validations in place
- [ ] Error handling provides clear feedback
- [ ] All tests pass
- [ ] Manual testing completed on different systems
- [ ] Performance acceptable for real-time updates
- [ ] Code review completed