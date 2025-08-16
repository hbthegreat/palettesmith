# PR2: Enhanced Plugin System

## Overview
Add templating system, regex patterns, and improved validation to the plugin system for robust config file handling.

## Acceptance Criteria

### Plugin Schema Enhancement
- [ ] Add template definitions to plugin manifests
- [ ] Add regex patterns for finding existing values
- [ ] Create validation service for plugin schemas
- [ ] Support multiple config file formats per plugin
- [ ] Add error handling for malformed plugins

### Template System
- [ ] Support placeholder syntax like `{accent}`, `{bg}`, `{fg}`
- [ ] Handle different color format outputs per field
- [ ] Support conditional templates based on field types
- [ ] Validate templates compile correctly

### Validation Service
- [ ] Separate validation logic from UI code
- [ ] Support all field types (color, number, text, select)
- [ ] Add custom validation rules per plugin
- [ ] Provide detailed error messages

## Files to Create

### `internal/plugin/template.go`
```go
package plugin

type Template struct {
    Pattern     string `json:"pattern"`      // Regex to find existing value
    Replacement string `json:"replacement"`   // Template with {placeholders}
    FileType    string `json:"file_type"`    // "conf", "toml", "ini", etc.
}

type TemplateEngine struct {
    templates map[string]Template
}

func (te *TemplateEngine) Apply(fieldKey, value string) (string, error) {
    // Apply template with value substitution
}

func (te *TemplateEngine) FindExisting(content, fieldKey string) (string, error) {
    // Use regex to find current value in config file
}
```

### `internal/validation/service.go`
```go
package validation

type Service struct{}

type FieldError struct {
    Field   string
    Message string
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) ValidateField(spec plugin.Field, value string) []FieldError {
    // Validate based on field type and constraints
}

func (s *Service) ValidatePlugin(p plugin.Plugin) []error {
    // Validate plugin configuration is valid
}
```

## Enhanced Plugin Schema

### `plugin.json` (Updated)
```json
{
  "id": "hyprland",
  "title": "Hyprland",
  "spec": "spec.json",
  "templates": "templates.json",
  "user_paths": ["~/.config/hypr/hyprland.conf"],
  "system_paths": ["/etc/xdg/hypr/hyprland.conf"],
  "reload": ["hyprctl", "reload"],
  "config_format": "conf"
}
```

### `templates.json` (New)
```json
{
  "accent": {
    "pattern": "col\\.active_border\\s*=\\s*rgb\\(([a-fA-F0-9]{6})\\)",
    "replacement": "col.active_border = rgb({accent})",
    "file_type": "conf"
  },
  "bg": {
    "pattern": "col\\.inactive_border\\s*=\\s*rgb\\(([a-fA-F0-9]{6})\\)", 
    "replacement": "col.inactive_border = rgb({bg})",
    "file_type": "conf"
  }
}
```

### `spec.json` (Enhanced)
```json
{
  "id": "hyprland",
  "fields": [
    {
      "key": "accent",
      "label": "Accent Color",
      "type": "color",
      "default": "#89b4fa",
      "color": {
        "format": "hex6",
        "output_format": "rgb_no_prefix"
      },
      "validation": {
        "required": true,
        "custom_regex": "^#[0-9a-fA-F]{6}$"
      },
      "help": "Highlight and focus color"
    }
  ]
}
```

## Files to Modify

### `internal/plugin/loader.go`
- Add template loading functionality
- Enhanced error handling for plugin validation
- Support for multiple file formats
- Validate plugin schemas on load

### `ui/tui/form.go`
- Remove validation logic (move to validation service)
- Use validation service for field validation
- Improve error display with detailed messages

## Testing Requirements

### Unit Tests
- [ ] Test template pattern matching
- [ ] Test template value substitution
- [ ] Test validation service with all field types
- [ ] Test malformed plugin handling
- [ ] Test regex pattern compilation

### Integration Tests  
- [ ] Test template application to real config files
- [ ] Test finding existing values in config files
- [ ] Test plugin loading with templates
- [ ] Verify validation errors are user-friendly

### Test Fixtures
Create sample configs in `tests/fixtures/`:
- `hyprland.conf` - Sample Hyprland config
- `alacritty.toml` - Sample Alacritty config  
- `malformed-plugin/` - Invalid plugin for error testing

## Manual Testing Steps

1. **Template System Test**
   - Create test plugin with templates
   - Verify pattern matching finds existing values
   - Test template substitution works correctly
   - Check different color format outputs

2. **Validation Test**
   - Enter invalid color values
   - Test number field min/max validation
   - Verify custom regex validation works
   - Check error messages are helpful

3. **Plugin Loading Test**
   - Test with malformed JSON
   - Test with missing template files
   - Verify graceful degradation
   - Check error reporting

## Error Handling

### Required Error Messages
- "Plugin {id}: template file not found"
- "Plugin {id}: invalid regex pattern in {field}"
- "Plugin {id}: template compilation failed"
- "Validation failed for {field}: {specific error}"

### Graceful Degradation
- Skip plugins with invalid templates
- Continue loading other plugins if one fails
- Show which plugins failed to load
- Provide fallback validation for missing rules

## Dependencies
- PR1 (Project Restructure) must be completed

## Definition of Done
- [ ] Template system works with existing Hyprland plugin
- [ ] Validation service handles all current field types
- [ ] Plugin loading robust against malformed configs
- [ ] All tests pass
- [ ] Manual testing completed
- [ ] Code review completed
- [ ] Error handling provides useful feedback