# PR8: Theme Save/Export

## Overview
Implement theme saving, naming, and export functionality with symlink integration for omarchy and other target systems.

## Acceptance Criteria

### Theme Saving System
- [ ] Save custom themes to palettesmith directory
- [ ] Prompt user for theme name when saving
- [ ] Generate complete theme package (all app configs)
- [ ] Preserve theme metadata and creation info
- [ ] Support theme versioning and updates

### Export Integration
- [ ] Create symlinks to target system theme directory
- [ ] Support omarchy theme directory integration
- [ ] Handle conflicts with existing themes
- [ ] Provide backup of overwritten themes
- [ ] Support different target system layouts

### Theme Metadata
- [ ] Store creation date, author, description
- [ ] Track which plugins/apps are included
- [ ] Save base theme information (if derived)
- [ ] Include preview images or color swatches
- [ ] Support theme tagging and categories

### Import/Export Formats
- [ ] Support palettesmith native format
- [ ] Export to omarchy format
- [ ] Support theme sharing (portable format)
- [ ] Import from other theme systems
- [ ] Validate imported themes

## Files to Create

### `internal/export/manager.go`
```go
package export

type Manager struct {
    config     *config.Config
    stagingMgr *staging.Manager
}

type ThemeExport struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Author      string    `json:"author"`
    Created     time.Time `json:"created"`
    BaseTheme   string    `json:"base_theme,omitempty"`
    Plugins     []string  `json:"plugins"`
    Files       map[string]string `json:"files"` // relative path -> content
}

func NewManager(config *config.Config, stagingMgr *staging.Manager) *Manager {
    // Initialize export manager
}

func (m *Manager) SaveTheme(name, description string) error {
    // Save current staging state as named theme
}

func (m *Manager) ExportToTarget(themeName string) error {
    // Export theme to target system (omarchy, etc.)
}

func (m *Manager) ListThemes() ([]ThemeInfo, error) {
    // List all saved themes
}
```

### `internal/export/symlink.go`
```go
package export

type SymlinkManager struct {
    targetDir string
    sourceDir string
}

func NewSymlinkManager(targetDir, sourceDir string) *SymlinkManager {
    // Initialize symlink manager
}

func (s *SymlinkManager) CreateThemeSymlink(themeName string) error {
    // Create symlink from target to palettesmith theme
}

func (s *SymlinkManager) RemoveThemeSymlink(themeName string) error {
    // Remove symlink safely
}

func (s *SymlinkManager) ValidateTarget() error {
    // Ensure target directory is suitable for symlinking
}

func (s *SymlinkManager) BackupExisting(themeName string) error {
    // Backup existing theme if it would be overwritten
}
```

### `internal/export/formats.go`
```go
package export

type ExportFormat interface {
    Name() string
    Export(theme ThemeExport, targetDir string) error
    Validate(themePath string) error
}

// Palettesmith native format
type PalettesmithFormat struct{}

func (p PalettesmithFormat) Export(theme ThemeExport, targetDir string) error {
    // Export in palettesmith native format
}

// Omarchy format
type OmarchyFormat struct{}

func (o OmarchyFormat) Export(theme ThemeExport, targetDir string) error {
    // Export in omarchy-compatible format
}

// Portable format (for sharing)
type PortableFormat struct{}

func (p PortableFormat) Export(theme ThemeExport, targetDir string) error {
    // Export as portable archive
}
```

### `ui/tui/save.go`
```go
package tui

type SaveModel struct {
    nameInput        textinput.Model
    descriptionInput textinput.Model
    authorInput      textinput.Model
    stage           SaveStage
    err             string
}

type SaveStage int

const (
    SaveStageName SaveStage = iota
    SaveStageDescription
    SaveStageAuthor
    SaveStageConfirm
)

func NewSaveModel() SaveModel {
    // Initialize save dialog
}

func (m SaveModel) Update(msg tea.Msg) (SaveModel, tea.Cmd) {
    // Handle save form input
}

func (m SaveModel) View() string {
    // Render save dialog
}
```

## Files to Modify

### `ui/tui/app.go`
- Add save command (S key)
- Show save dialog when requested
- Integrate with export manager
- Display save status and results

### `internal/config/config.go`
- Add theme directory configuration
- Add author information for saved themes
- Support different export formats

## Theme Directory Structure

### Palettesmith Themes
```
~/.config/palettesmith/themes/
├── my-custom-theme/
│   ├── theme.json                 # Theme metadata
│   ├── hypr/
│   │   └── hyprland.conf
│   ├── alacritty/
│   │   └── alacritty.toml
│   └── waybar/
│       └── style.css
├── catppuccin-modified/
│   ├── theme.json
│   └── ...
└── shared-themes/
    └── downloaded-theme.tar.gz
```

### Theme Metadata Format
```json
{
  "name": "My Custom Theme",
  "description": "Dark theme with blue accents",
  "author": "username",
  "created": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "base_theme": "catppuccin",
  "plugins": ["hyprland", "alacritty", "waybar"],
  "colors": {
    "bg": "#1e1e2e",
    "fg": "#cdd6f4", 
    "accent": "#89b4fa"
  },
  "tags": ["dark", "blue", "minimal"]
}
```

## Testing Requirements

### Unit Tests
- [ ] Test theme saving with metadata
- [ ] Test symlink creation and removal
- [ ] Test export format conversion
- [ ] Test theme validation
- [ ] Test conflict resolution

### Integration Tests
- [ ] Test full save workflow
- [ ] Test omarchy integration
- [ ] Test theme import/export
- [ ] Test symlink management
- [ ] Verify theme portability

### File System Tests
- [ ] Test with different file permissions
- [ ] Test with existing theme conflicts
- [ ] Test symlink creation on different filesystems
- [ ] Test cleanup after failures

## Manual Testing Steps

1. **Theme Saving Test**
   - Make changes to multiple plugins
   - Save theme with custom name
   - Verify all files saved correctly
   - Check metadata is accurate

2. **Export Integration Test**
   - Save theme to palettesmith directory
   - Export to omarchy themes directory
   - Verify symlink created correctly
   - Test theme appears in omarchy theme switcher

3. **Conflict Resolution Test**
   - Try saving theme with existing name
   - Test overwrite confirmation
   - Verify backup creation
   - Test rollback if export fails

4. **Theme Sharing Test**
   - Export theme as portable format
   - Import theme on different system
   - Verify theme works correctly
   - Test theme validation

## Save Dialog Interface

### Save Flow
```
┌─ Save Theme ────────────────────────────────┐
│ Name: My Custom Theme                       │
│ ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔  │
│                                             │
│ Description: Dark theme with blue accents   │
│ ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔  │
│                                             │
│ Author: username                            │
│ ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔  │
│                                             │
│ ☑ Export to omarchy themes                  │
│ ☑ Create symlink                            │
│                                             │
│ [S]ave  [C]ancel                            │
└─────────────────────────────────────────────┘
```

## Export Options

### Target Systems
- **Omarchy**: Export to `~/.config/omarchy/themes/`
- **Generic**: Keep in palettesmith directory only
- **Portable**: Create shareable archive
- **Custom**: User-specified directory

### Symlink Strategies
- **Replace**: Replace existing theme (with backup)
- **Merge**: Merge with existing theme
- **Skip**: Keep both themes separate
- **Rename**: Auto-rename to avoid conflicts

## Error Handling

### Required Error Messages
- "Theme name already exists, choose different name"
- "Cannot create symlink: permission denied"
- "Export failed: target directory not writable"
- "Theme validation failed: {specific errors}"

### Recovery Strategies
- Prompt for different theme name on conflict
- Offer manual symlink creation instructions
- Provide rollback on partial export failure
- Validate before starting export process

## Dependencies
- PR3 (Theme Staging) must be completed
- PR1 (Project Restructure) for config integration

## Definition of Done
- [ ] Themes save with complete metadata
- [ ] Symlink integration works with omarchy
- [ ] Export formats produce valid themes
- [ ] Save dialog provides good user experience
- [ ] Conflict resolution works reliably
- [ ] All tests pass
- [ ] Manual testing completed
- [ ] Theme portability verified
- [ ] Code review completed
- [ ] Documentation includes theme format specification