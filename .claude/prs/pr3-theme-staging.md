# PR3: Theme Staging System

## Overview
Implement backup/staging system for safe theme editing with rollback capabilities and real-time diff display.

## Acceptance Criteria

### Staging Directory Management
- [ ] Create staging directories (`backup/` and `next/`)
- [ ] Import existing theme files to staging area
- [ ] Track which files came from where (for rollback)
- [ ] Handle symlinks and directory structures properly

### Backup System
- [ ] Backup original files before any modification
- [ ] Preserve file permissions and ownership
- [ ] Track backup metadata (source, timestamp, theme name)
- [ ] Support incremental backups for efficiency

### Diff System
- [ ] Show current vs next file differences
- [ ] Highlight changed lines with context
- [ ] Support different file formats (conf, toml, ini, etc.)
- [ ] Show diffs per plugin and per field

### Rollback Capabilities
- [ ] Full theme rollback (restore all original files)
- [ ] Field-level rollback (restore individual values)
- [ ] Plugin-level rollback (restore specific app configs)
- [ ] Confirm destructive rollback operations

## Files to Create

### `internal/staging/manager.go`
```go
package staging

type Manager struct {
    backupDir string
    nextDir   string
    config    *config.Config
}

type FileMetadata struct {
    OriginalPath string    `json:"original_path"`
    BackupPath   string    `json:"backup_path"`
    SourceTheme  string    `json:"source_theme"`
    Timestamp    time.Time `json:"timestamp"`
    Checksum     string    `json:"checksum"`
}

func NewManager(cfg *config.Config) (*Manager, error) {
    // Initialize staging directories
}

func (m *Manager) ImportTheme(themePath string) error {
    // Copy theme files to staging area
}

func (m *Manager) CreateBackup(files []string) error {
    // Backup original files before modification
}

func (m *Manager) GetDiff(pluginID, fieldKey string) (Diff, error) {
    // Generate diff for specific change
}
```

### `internal/staging/diff.go`
```go
package staging

type Diff struct {
    PluginID    string      `json:"plugin_id"`
    FieldKey    string      `json:"field_key"`
    FilePath    string      `json:"file_path"`
    Changes     []LineChange `json:"changes"`
    HasChanges  bool        `json:"has_changes"`
}

type LineChange struct {
    LineNumber int    `json:"line_number"`
    Type       string `json:"type"` // "added", "removed", "modified"
    OldContent string `json:"old_content"`
    NewContent string `json:"new_content"`
    Context    string `json:"context"`
}

func GenerateDiff(originalFile, modifiedFile string) ([]LineChange, error) {
    // Generate line-by-line diff
}
```

### `internal/staging/rollback.go`
```go
package staging

type RollbackService struct {
    manager *Manager
}

func NewRollbackService(m *Manager) *RollbackService {
    return &RollbackService{manager: m}
}

func (r *RollbackService) RollbackAll() error {
    // Restore all files from backup
}

func (r *RollbackService) RollbackPlugin(pluginID string) error {
    // Restore files for specific plugin
}

func (r *RollbackService) RollbackField(pluginID, fieldKey string) error {
    // Restore specific field value
}
```

## Files to Modify

### `internal/theme/store.go`
- Add staging integration
- Track changes in staging area
- Provide rollback for field-level changes
- Integrate with diff system

### `ui/tui/app.go`
- Add staging manager dependency
- Show diff information in UI
- Add rollback commands (U for undo)
- Display staging status

### `ui/tui/form.go`
- Update to work with staging system
- Show field-level rollback options
- Display diff preview for changes
- Integrate with rollback service

## Directory Structure
```
~/.config/palettesmith/staging/
├── metadata.json              # Staging session info
├── backup/                    # Original files
│   ├── hypr/
│   │   └── hyprland.conf
│   └── alacritty/
│       └── alacritty.toml
└── next/                      # Modified files  
    ├── hypr/
    │   └── hyprland.conf      # With applied changes
    └── alacritty/
        └── alacritty.toml     # With applied changes
```

## Testing Requirements

### Unit Tests
- [ ] Test staging directory creation
- [ ] Test file backup and restoration
- [ ] Test diff generation accuracy
- [ ] Test rollback operations
- [ ] Test metadata tracking

### Integration Tests
- [ ] Test full theme import workflow
- [ ] Test multi-plugin changes and rollbacks
- [ ] Test file permission preservation
- [ ] Test symlink handling
- [ ] Verify no data loss during operations

### Test Fixtures
Create in `tests/fixtures/staging/`:
- `sample-theme/` - Complete theme directory
- `modified-configs/` - Files with known changes
- `symlink-theme/` - Theme with symlinked files

## Manual Testing Steps

1. **Import Theme Test**
   - Start with clean staging area
   - Import existing omarchy theme
   - Verify all files copied correctly
   - Check metadata is accurate

2. **Diff Generation Test**
   - Make changes to theme values
   - Verify diffs show expected changes
   - Test different file formats
   - Check context lines are helpful

3. **Rollback Test**
   - Make several changes across plugins
   - Test field-level rollback
   - Test plugin-level rollback  
   - Test full theme rollback
   - Verify original state restored

4. **Safety Test**
   - Test with read-only files
   - Test with missing source files
   - Test with corrupted backups
   - Verify graceful error handling

## Error Handling

### Required Error Messages
- "Cannot create staging directory: {path}"
- "Backup failed for {file}: {error}"
- "Cannot generate diff: files not found"
- "Rollback failed: backup corrupted"
- "Permission denied accessing {file}"

### Safety Measures
- Never modify original files until user confirms
- Always verify backups before proceeding
- Check disk space before operations
- Validate file integrity with checksums
- Confirm destructive operations

## Dependencies
- PR1 (Project Restructure) must be completed
- PR2 (Enhanced Plugins) recommended for template integration

## Definition of Done
- [ ] Staging system creates proper backups
- [ ] Diff generation works for all file types
- [ ] Rollback operations restore correct state
- [ ] All tests pass
- [ ] Manual testing completed successfully
- [ ] File safety verified (no data loss)
- [ ] Performance acceptable for large themes
- [ ] Code review completed