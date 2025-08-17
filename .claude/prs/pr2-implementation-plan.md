# PR2: Enhanced Plugin System - Complete Implementation Plan

## Overview
This document captures the complete implementation plan for PR2 after extensive analysis and planning discussion.

## Architecture Decisions Made

### 1. JSON File Structure (Final)
- **plugin.json**: Plugin metadata, paths, reload commands (unchanged)
- **spec.json**: Enhanced with templates, validation, output formats (no separate templates.json)

### 2. Color Conversion Strategy  
**Decision**: Central color conversion service (Option A)
- Single converter that handles all input formats
- Output format specified per field in spec.json
- Cleaner architecture than per-template conversion

### 3. Error Handling Strategy
- **Plugin Loading**: Continue loading other plugins if one fails, collect errors
- **Display**: Show errors in TUI sidebar with red indicators  
- **CLI Flag**: Add `--validate-plugins` for plugin authors
- **Template Failures**: Skip individual fields, keep plugin if possible
- **Runtime**: Always revert to previous config on failure - never leave users broken

### 4. Template Strategy
**Decision**: Custom regex-based approach (not Go's text/template)
- Simple `{value}` placeholder replacement
- Regex patterns for extraction from config files
- Go's `text/template` is overkill for basic value substitution
- We need regex pattern matching for `FindExisting()` functionality anyway
- Templates stored in spec.json (not separate file)
- Pure string transformation (no file I/O in PR2)

### 5. Validation Architecture
**Decision**: Specific validators for different concerns
- `PluginValidator`: Validates plugin manifests and field definitions (PR2)
- `ConfigValidator`: Validates actual config files (Future PR)
- `TemplateValidator`: Validates template patterns work correctly (Future PR)
- `ThemeValidator`: Validates complete theme packages (Future PR)

**Naming Convention**: Avoid generic names like "service.go", use specific descriptive names

## Implementation Sequence

### Phase 1: Core Infrastructure
1. **Color Converter** (`internal/color/`)
   - Handle all format conversions
   - Comprehensive test coverage
   - Support formats: hex6, hex6_no_prefix, hex_0x, rgb_no_prefix

2. **Template Engine** (`internal/plugin/template.go`)
   - Pattern matching and replacement (regex-based)
   - Pure string transformation (no file I/O)
   - Error handling for invalid patterns

3. **Plugin Validator** (`internal/validation/plugin.go`)
   - Plugin manifest and field validation
   - Detailed error reporting with codes
   - Integration with color converter

### Phase 2: Enhanced Plugin System
1. **Enhanced Field Structure**
   ```json
   {
     "key": "accent",
     "label": "Accent Color",
     "type": "color",
     "default": "#89b4fa",
     "output_format": "rgb_no_prefix",
     "template": {
       "pattern": "col\\.active_border\\s*=\\s*rgb\\(([a-fA-F0-9]{6})\\)",
       "replacement": "col.active_border = rgb({value})"
     },
     "validation": {
       "required": true,
       "custom_regex": "^#[0-9a-fA-F]{6}$"
     },
     "help": "Highlight and focus color"
   }
   ```

2. **Plugin Loading Enhancement**
   - Load and validate templates
   - Collect errors instead of failing
   - Enhanced error reporting

### Phase 3: TUI Integration
1. **Form Enhancement**
   - Use validation service instead of inline validation
   - Real-time error display
   - Integration with template system

2. **Sidebar Enhancement**
   - Show plugin errors with red indicators
   - Plugin icons using Nerd Fonts

3. **Error Display System**
   - Context-aware error messages
   - In-TUI error reporting where possible

### Phase 4: Plugin Creation
Create plugins for all omarchy apps:

1. **Hyprland** (󰖟) - `.conf` format, `rgb()` values
2. **Alacritty** (󰆍) - `.toml` format, `"#hex"` and `"0xhex"`
3. **Waybar** (󰕮) - `.css` format, `#hex` values
4. **Mako** (󰵅) - `.ini` format, `#hex` values  
5. **BTTop** (󰍉) - `.theme` format, `theme[key]="#hex"`
6. **SwayOSD** (󰕾) - `.css` format, `@define-color`
7. **Walker** (󰍉) - `.css` format
8. **Hyprlock** (󰌾) - `.conf` format (like Hyprland)
9. **Neovim** (󰕷) - `.lua` format, colorscheme names

## File Structure Plan

```
internal/
├── color/
│   ├── converter.go           # NEW: Color format conversion
│   └── converter_test.go      # NEW: Comprehensive color tests
├── plugin/
│   ├── loader.go              # MODIFY: Enhanced plugin loading
│   ├── template.go            # NEW: Template pattern system
│   └── template_test.go       # NEW: Template tests
└── validation/
    ├── plugin.go              # NEW: Plugin manifest/field validator
    ├── plugin_test.go         # NEW: Plugin validation tests
    ├── config.go              # FUTURE: Config file validator
    ├── template.go            # FUTURE: Template pattern validator
    └── theme.go               # FUTURE: Complete theme validator

plugins/                       # EXPAND: Create all app plugins
├── hyprland/                  # EXISTS: Enhance with templates
├── alacritty/                 # NEW
├── waybar/                    # NEW  
├── mako/                      # NEW
├── btop/                      # NEW
├── swayosd/                   # NEW
├── walker/                    # NEW
├── hyprlock/                  # NEW
└── neovim/                    # NEW

ui/tui/
├── form.go                    # MODIFY: Use plugin validator
└── sidebar.go                 # MODIFY: Show error indicators

cmd/palettesmith/
└── main.go                    # MODIFY: Add --validate-plugins flag

tests/fixtures/                # EXPAND: Comprehensive test data
├── configs/                   # Sample files for each app
├── plugins/                   # Valid/invalid plugin examples
└── colors/                    # Color conversion test cases
```

## Key Technical Specifications

### Color Converter Interface
```go
type Converter struct {
    value string // normalized hex6 internally
}

func NewConverter(input string) (*Converter, error)
func (c *Converter) ToHex6() string           // #ffffff
func (c *Converter) ToHex6NoPrefix() string   // ffffff
func (c *Converter) ToHex0x() string          // 0xffffff  
func (c *Converter) ToRGBNoPrefix() string    // ffffff (for Hyprland)
func (c *Converter) ToFormat(format string) string
```

### Template Engine Interface
```go
type Template struct {
    Pattern     string `json:"pattern"`
    Replacement string `json:"replacement"`
}

type TemplateEngine struct{}

func (te *TemplateEngine) Apply(template Template, value string) (string, error)
func (te *TemplateEngine) FindExisting(content string, template Template) (string, error)
// Note: DetectFileType removed - file I/O belongs in PR3/PR7
```

### Plugin Validation Interface
```go
// Plugin manifest and field validation (PR2)
type PluginValidator struct {
    colorConverter *color.Converter
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

func (pv *PluginValidator) ValidateField(field Field, value string) []FieldError
func (pv *PluginValidator) ValidatePlugin(plugin Plugin) []error

// Future validators (separate PRs)
type ConfigValidator struct{}    // Validates actual config files
type TemplateValidator struct{}  // Validates template patterns work
type ThemeValidator struct{}     // Validates complete theme packages
```

### Enhanced Plugin Types
```go
type Field struct {
    Key          string     `json:"key"`
    Label        string     `json:"label"`
    Type         string     `json:"type"`
    Default      string     `json:"default,omitempty"`
    OutputFormat string     `json:"output_format,omitempty"`  // NEW
    Template     *Template  `json:"template,omitempty"`       // NEW
    Validation   *Validation `json:"validation,omitempty"`    // NEW
    Help         string     `json:"help,omitempty"`
    Min          *float64   `json:"min,omitempty"`
    Max          *float64   `json:"max,omitempty"`
    Enum         []string   `json:"enum,omitempty"`
}

type Validation struct {
    Required    bool   `json:"required,omitempty"`
    CustomRegex string `json:"custom_regex,omitempty"`
}

type PluginError struct {
    PluginID string
    Message  string
    Type     string // "parse_error", "missing_template", "validation_failed"
}

type Store struct {
    byID   map[string]Plugin
    list   []Plugin
    errors []PluginError  // NEW: Track loading errors
}
```

## Testing Strategy

### Unit Test Priority
1. **Color Converter** - All formats, edge cases, invalid inputs
2. **Template Engine** - Pattern matching, substitution, file detection
3. **Validation Service** - All field types, custom rules, error messages
4. **Plugin Loading** - Error collection, malformed JSON handling

### Integration Tests
1. **Plugin Discovery** - Load valid/invalid plugin directories
2. **Template Application** - Apply to real config files in temp dirs
3. **End-to-End** - Form → validation → template → file modification

### Test Fixtures Structure
```
tests/fixtures/
├── configs/
│   ├── hyprland.conf          # Various format examples
│   ├── alacritty.toml
│   ├── waybar.css
│   ├── mako.ini
│   ├── btop.theme
│   └── swayosd.css
├── plugins/
│   ├── valid-plugin/
│   ├── malformed-json/
│   ├── invalid-template/
│   ├── missing-fields/
│   └── regex-error/
└── colors/
    ├── valid-formats.json
    └── invalid-formats.json
```

## File Format Patterns

### Hyprland (.conf)
```
Pattern: col\.active_border\s*=\s*rgb\(([a-fA-F0-9]{6})\)
Replacement: col.active_border = rgb({value})
Output: rgb_no_prefix
```

### Alacritty (.toml)
```
Pattern: background\s*=\s*"#([a-fA-F0-9]{6})"
Replacement: background = "{value}"
Output: hex6
```

### CSS (.css)
```
Pattern: @define-color\s+foreground\s+#([a-fA-F0-9]{6});
Replacement: @define-color foreground {value};
Output: hex6
```

### BTTop (.theme)
```
Pattern: theme\[main_bg\]="#([a-fA-F0-9]{6})"
Replacement: theme[main_bg]="{value}"
Output: hex6
```

## Performance Considerations
- **Regex Caching**: Compile patterns once during plugin loading
- **Color Conversion**: Cache converted values per field update
- **File I/O**: Batch config file operations where possible
- **Error Collection**: Efficient storage and display of plugin errors

## Success Criteria
- [ ] All 9 app plugins load successfully
- [ ] Color conversion works for all formats
- [ ] Template application modifies config files correctly
- [ ] Error handling never leaves users in broken state
- [ ] TUI shows plugin errors clearly
- [ ] CLI validation flag works for plugin authors
- [ ] Comprehensive test coverage (>90%)
- [ ] Performance acceptable on reasonable plugin counts

## Implementation Notes
- No backward compatibility required (new project)
- Users have Nerd Fonts installed (can use proper icons)
- Business logic must stay separate from UI logic
- Error messages should be contextual and helpful
- Always provide rollback capability
- Follow existing Go project patterns and test structure

---

**Ready to begin implementation with Phase 1: Color Converter!**