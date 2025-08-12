# PaletteSmith Implementation Plan

## Executive Summary

PaletteSmith is a plugin-based TUI theme editor that discovers and themes applications on Linux systems. It provides both an interactive TUI and a composable CLI interface for integration with systems like Omarchy. Every themeable application is handled through plugins, making the system infinitely extensible.

## Technology Stack

- **Language**: Go 1.22+
- **TUI Framework**: Bubble Tea (github.com/charmbracelet/bubbletea)
- **Styling**: Lip Gloss (github.com/charmbracelet/lipgloss)
- **Components**: Bubbles (github.com/charmbracelet/bubbles)
- **Config Parsing**: 
  - CSS: github.com/tdewolff/parse/css
  - TOML: github.com/BurntSushi/toml
  - YAML: gopkg.in/yaml.v3
  - JSON: encoding/json (stdlib)
  - INI: gopkg.in/ini.v1
- **Color Manipulation**: github.com/lucasb-eyer/go-colorful
- **CLI Framework**: github.com/spf13/cobra
- **Testing**: testing (stdlib) + github.com/stretchr/testify

## Project Structure

```
palettesmith/
├── cmd/
│   └── palettesmith/
│       └── main.go                 # Entry point
├── internal/
│   ├── core/
│   │   ├── color.go                # Color type and manipulation
│   │   ├── palette.go              # Global palette management
│   │   ├── state.go                # State persistence
│   │   └── backup.go               # Backup/rollback system
│   ├── plugin/
│   │   ├── interface.go            # Plugin interface definition
│   │   ├── loader.go               # Plugin loading and discovery
│   │   ├── registry.go             # Plugin registry
│   │   └── manifest.go             # Plugin manifest parsing
│   ├── parser/
│   │   ├── interface.go            # Parser interface
│   │   ├── css.go                  # CSS parser implementation
│   │   ├── json.go                 # JSON parser implementation
│   │   ├── toml.go                 # TOML parser implementation
│   │   ├── yaml.go                 # YAML parser implementation
│   │   └── ini.go                  # INI parser implementation
│   ├── tui/
│   │   ├── app.go                  # Main TUI application
│   │   ├── navigation.go           # Navigation component
│   │   ├── colorpicker.go          # Color picker component
│   │   ├── preview.go              # Preview pane
│   │   ├── list.go                 # App/color list component
│   │   └── styles.go               # TUI styling
│   ├── cli/
│   │   ├── root.go                 # Root command
│   │   ├── list.go                 # List command
│   │   ├── get.go                  # Get command
│   │   ├── set.go                  # Set command
│   │   ├── apply.go                # Apply command
│   │   ├── rollback.go             # Rollback command
│   │   ├── export.go               # Export command
│   │   └── import.go               # Import command
│   └── config/
│       ├── paths.go                # Path resolution
│       └── settings.go             # User settings
├── plugins/
│   ├── waybar/
│   │   ├── plugin.go               # Waybar plugin implementation
│   │   └── manifest.yaml           # Waybar plugin manifest
│   ├── hyprland/
│   │   ├── plugin.go
│   │   └── manifest.yaml
│   ├── alacritty/
│   │   ├── plugin.go
│   │   └── manifest.yaml
│   ├── kitty/
│   │   ├── plugin.go
│   │   └── manifest.yaml
│   ├── btop/
│   │   ├── plugin.go
│   │   └── manifest.yaml
│   ├── mako/
│   │   ├── plugin.go
│   │   └── manifest.yaml
│   └── README.md                   # Plugin development guide
├── scripts/
│   ├── dev.sh                      # Development mode runner
│   ├── build.sh                    # Build script
│   ├── test.sh                     # Test runner
│   └── install.sh                  # Local installation
├── testdata/
│   ├── configs/                    # Sample config files for testing
│   └── themes/                     # Sample themes for testing
├── docs/
│   ├── PLUGIN_DEVELOPMENT.md       # Plugin API documentation
│   └── CLI_USAGE.md               # CLI documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Development Scripts

### scripts/dev.sh
```bash
#!/bin/bash
# Development mode - hot reload with air or simple rebuild
# Install: go install github.com/cosmtrek/air@latest
air || while true; do go run ./cmd/palettesmith "$@"; done
```

### scripts/build.sh
```bash
#!/bin/bash
# Build binary for current platform
CGO_ENABLED=0 go build -ldflags="-s -w" -o palettesmith ./cmd/palettesmith
```

### scripts/test.sh
```bash
#!/bin/bash
# Run all tests with coverage
go test -v -cover ./...
```

### scripts/install.sh
```bash
#!/bin/bash
# Install locally for testing
go build -o palettesmith ./cmd/palettesmith
sudo mv palettesmith /usr/local/bin/
```

### Makefile
```makefile
.PHONY: dev build test install clean

dev:
	@./scripts/dev.sh

build:
	@./scripts/build.sh

test:
	@./scripts/test.sh

install:
	@./scripts/install.sh

clean:
	rm -f palettesmith
	rm -rf dist/

run:
	go run ./cmd/palettesmith
```

## Implementation Phases

### Phase 0: Project Setup (Day 1)

1. **Initialize Go module**
   - Run: `go mod init github.com/yourusername/palettesmith`
   - Create directory structure as specified above
   - Create all empty .go files with package declarations
   - Start with this minimal main.go:

```go
// cmd/palettesmith/main.go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("PaletteSmith v0.0.1")
    if len(os.Args) > 1 && os.Args[1] == "--help" {
        fmt.Println("Usage: palettesmith [command] [flags]")
        fmt.Println("This is a placeholder. Implementation coming soon.")
    }
}
```

2. **Install dependencies**
   - Run: `go get github.com/charmbracelet/bubbletea`
   - Run: `go get github.com/charmbracelet/lipgloss`
   - Run: `go get github.com/charmbracelet/bubbles`
   - Run: `go get github.com/spf13/cobra`
   - Run: `go get github.com/lucasb-eyer/go-colorful`
   - Run: `go get github.com/tdewolff/parse/v2`
   - Run: `go get github.com/BurntSushi/toml`
   - Run: `go get gopkg.in/yaml.v3`
   - Run: `go get gopkg.in/ini.v1`
   - Run: `go get github.com/stretchr/testify`

3. **Create development scripts**
   - Copy all scripts from above into scripts/ directory
   - Make them executable: `chmod +x scripts/*.sh`
   - Test basic setup: `make run` (should fail but project should be set up)

4. **Setup version control**
   - Initialize git: `git init`
   - Create .gitignore with: `palettesmith`, `dist/`, `*.log`, `.air.toml`
   - Initial commit

### Phase 1: Core Types and Interfaces (Day 2-3)

1. **Implement core/color.go**
   - Define Color type that can handle hex, RGB, RGBA, HSL formats
   - Implement conversion methods between formats
   - Implement color parsing from strings
   - Implement color serialization to different formats
   - Add validation methods
   - Write unit tests

2. **Implement plugin/interface.go**
   - Define Plugin interface with methods:
     - GetManifest() returns plugin metadata
     - Detect() checks if app is installed
     - GetFiles() returns config file paths
     - ExtractColors() gets current colors from configs
     - ApplyColors() writes new colors to configs
     - Restart() restarts the application
     - Validate() checks if changes are valid
   - Define PluginManifest struct with all metadata fields
   - Define ColorDefinition struct for color mappings

3. **Implement plugin/manifest.go**
   - Create manifest parser that reads YAML manifests
   - Validate manifest structure
   - Handle path expansion (~/ to home directory)
   - Write tests with sample manifests

4. **Implement parser/interface.go**
   - Define Parser interface with methods:
     - Parse(content []byte) returns parsed document
     - FindColors(doc) returns map of color locations
     - ReplaceColor(doc, oldColor, newColor) modifies document
     - Serialize(doc) returns []byte
   - Create factory function to get parser by file type

5. **Implement core/state.go**
   - Define State struct containing:
     - Global palette (map of color names to Color)
     - Per-app overrides (map of app to color overrides)
     - Last applied timestamp
     - Backup references
   - Implement Load() and Save() methods using JSON
   - State file location: ~/.config/palettesmith/state.json
   - Write tests

### Phase 2: Parser Implementations (Day 4-5)

1. **Implement parser/css.go**
   - Use tdewolff/parse to parse CSS
   - Extract color values from:
     - Hex colors (#RGB, #RRGGBB, #RRGGBBAA)
     - RGB/RGBA functions
     - HSL/HSLA functions
     - CSS variables (@define-color for GTK CSS)
   - Implement color replacement preserving format
   - Handle comments and formatting
   - Write comprehensive tests

2. **Implement parser/json.go**
   - Use encoding/json with json.RawMessage
   - Walk JSON tree looking for color strings
   - Support JSONPath-style targeting from manifest
   - Preserve formatting as much as possible
   - Write tests with various JSON structures

3. **Implement parser/toml.go**
   - Use BurntSushi/toml
   - Extract colors from string values
   - Support dotted path notation from manifest
   - Preserve comments using raw parsing fallback
   - Write tests

4. **Implement parser/yaml.go**
   - Use gopkg.in/yaml.v3 with Node API
   - Preserve comments and formatting
   - Extract colors from string values
   - Support path notation from manifest
   - Write tests

5. **Implement parser/ini.go**
   - Use gopkg.in/ini.v1
   - Extract colors from values
   - Support section.key notation from manifest
   - Preserve comments
   - Write tests

### Phase 3: Plugin System (Day 6-7)

1. **Implement plugin/loader.go**
   - Create LoadBuiltinPlugins() function
   - Create LoadExternalPlugins() function that:
     - Scans ~/.config/palettesmith/plugins/
     - Loads and validates manifests
     - Creates plugin instances
   - Handle plugin conflicts (same app ID)
   - Write tests with mock plugins

2. **Implement plugin/registry.go**
   - Create PluginRegistry type
   - Implement Register(), Get(), List() methods
   - Implement detection logic that runs all Detect() methods
   - Cache detection results
   - Write tests

3. **Create first plugin: plugins/waybar/**
   - Write manifest.yaml with:
     - Detection rules (config file existence)
     - File paths and formats
     - Color definitions (CSS variables)
     - Restart command
   - Implement plugin.go:
     - Embed manifest
     - Implement all Plugin interface methods
     - Use CSS parser for style.css
     - Use JSON parser for config
   - Write tests with sample Waybar configs

4. **Create plugins/hyprland/**
   - Write manifest for hyprland.conf
   - Handle Hyprland's special color format: rgba() and rgb()
   - Parse color values from config
   - Implement hyprctl reload for restart
   - Write tests

5. **Create plugins/alacritty/**
   - Write manifest for alacritty.toml
   - Handle TOML color sections
   - Implement config reload mechanism
   - Write tests

### Phase 4: CLI Implementation (Day 8-9)

1. **Implement cli/root.go**
   - Setup Cobra root command
   - Add global flags:
     - --config (config directory)
     - --plugins (additional plugin directory)
     - --verbose (debug output)
     - --no-backup (skip backups)
   - Initialize plugin registry
   - Load state

2. **Implement cli/list.go**
   - Create 'list' command
   - Show detected applications
   - Show available colors per app
   - Format output as table or JSON (--json flag)
   - Include detection status

3. **Implement cli/get.go**
   - Create 'get' command with subcommands:
     - get global - show global palette
     - get app [name] - show app current colors
     - get overrides [app] - show app overrides
   - Support --json output
   - Show color values in multiple formats

4. **Implement cli/set.go**
   - Create 'set' command with subcommands:
     - set global [color=value...] 
     - set app [name] [color=value...]
   - Parse color values in any format
   - Update state file
   - Validate colors before saving
   - Support --from-file flag for JSON input

5. **Implement cli/apply.go**
   - Create 'apply' command
   - Load current state
   - For each detected app:
     - Read current configs
     - Apply global palette
     - Apply app overrides
     - Validate changes
     - Write configs
     - Restart app
   - Support --dry-run flag
   - Support --app flag to limit scope
   - Show progress and results

6. **Implement core/backup.go and cli/rollback.go**
   - Implement backup creation before apply
   - Store backups in ~/.config/palettesmith/backups/
   - Use timestamp directories
   - Create manifest.json with backup metadata
   - Implement rollback command to restore
   - Support --list to show available backups
   - Auto-cleanup old backups (keep last 10)

### Phase 5: TUI Foundation (Day 10-11)

1. **Implement tui/app.go**
   - Create main App model struct with:
     - Current view (navigation, color picker, etc.)
     - Plugin registry
     - Current state
     - Selected app
     - Selected color
   - Implement Init() for initial commands
   - Implement Update() for message handling
   - Implement View() for rendering
   - Setup keyboard handlers (tab, arrows, enter, esc, q)

2. **Implement tui/styles.go**
   - Define consistent styles using Lip Gloss
   - Create styles for:
     - Headers and titles
     - Navigation items (selected/unselected)
     - Color swatches
     - Input fields
     - Buttons
     - Status messages
   - Support both dark and light terminal themes

3. **Implement tui/navigation.go**
   - Create navigation sidebar component
   - List "Global" plus all detected apps
   - Handle selection with arrow keys
   - Show detection status icons
   - Implement scrolling for long lists
   - Send app selection messages

4. **Implement tui/list.go**
   - Create color list component
   - Show colors for selected app
   - Display color name, current value, override status
   - Handle selection with arrow keys
   - Support filtering/searching
   - Send color selection messages

5. **Implement tui/preview.go**
   - Create preview pane using Lip Gloss
   - Show color swatches in a grid
   - Display sample UI elements:
     - Text on background
     - Borders and accents
     - Status indicators
     - Button examples
   - Update in real-time with color changes

### Phase 6: TUI Color Picker (Day 12-13)

1. **Implement tui/colorpicker.go**
   - Create color picker modal component
   - Implement multiple input methods:
     - Hex input field
     - RGB sliders (0-255)
     - HSL sliders
     - Preset color grid
   - Show current and new color comparison
   - Handle OK/Cancel actions
   - Send color change messages

2. **Integrate color picker into main app**
   - Launch picker when selecting color from list
   - Update preview in real-time
   - Save changes to state on OK
   - Handle escape to cancel

3. **Add external picker support**
   - Add button to launch hyprpicker if available
   - Parse hyprpicker output
   - Update color picker with external value
   - Make it optional (check if installed)

4. **Implement copy/paste support**
   - Handle Ctrl+C to copy color value
   - Handle Ctrl+V to paste color value
   - Support multiple formats in clipboard

### Phase 7: TUI Integration (Day 14-15)

1. **Complete tui/app.go integration**
   - Wire all components together
   - Implement view switching
   - Add status bar with hints
   - Show unsaved changes indicator
   - Add confirmation dialogs

2. **Implement apply from TUI**
   - Add Apply button/command (Ctrl+S)
   - Show progress during apply
   - Display results (success/failure per app)
   - Handle errors gracefully
   - Update state after successful apply

3. **Implement TUI settings**
   - Add settings view for:
     - Backup preferences
     - Plugin directories
     - Color format preferences
     - Restart behavior
   - Save settings to config file

4. **Polish TUI experience**
   - Add help screen (F1 or ?)
   - Improve responsive layout
   - Add animations for transitions
   - Test on different terminal sizes
   - Ensure mouse support works

### Phase 8: Additional Plugins (Day 16-17)

1. **Create plugins/kitty/**
   - Parse kitty.conf format
   - Handle color definitions
   - Implement remote control reload
   - Test with real Kitty config

2. **Create plugins/btop/**
   - Parse btop theme format
   - Handle color definitions
   - Implement btop reload signal
   - Test with btop themes

3. **Create plugins/mako/**
   - Parse mako config format
   - Handle color definitions
   - Implement makoctl reload
   - Test with mako config

4. **Create plugins/wezterm/**
   - Parse Lua config (basic)
   - Extract color scheme
   - Handle reload
   - Test with WezTerm

5. **Document plugin development**
   - Write PLUGIN_DEVELOPMENT.md
   - Create plugin template
   - Document manifest format
   - Provide examples
   - Explain testing approach

### Phase 9: Testing and Validation (Day 18-19)

1. **Unit tests**
   - Test all parsers with various inputs
   - Test color conversions
   - Test state management
   - Test plugin loading
   - Test backup/restore
   - Achieve 80% code coverage

2. **Integration tests**
   - Test full apply cycle
   - Test with real config files
   - Test plugin detection
   - Test rollback scenarios
   - Test concurrent operations

3. **TUI testing**
   - Manual testing checklist
   - Test all keyboard shortcuts
   - Test on different terminals
   - Test color picker modes
   - Test error conditions

4. **Create test data**
   - Add sample configs for all supported apps
   - Create test themes
   - Add broken configs for error testing
   - Create minimal and maximal test cases

### Phase 10: Binary Building and Distribution (Day 20)

1. **Create build system**
   - Write build script for multiple platforms:
     - linux-amd64
     - linux-arm64
     - linux-arm
   - Use CGO_ENABLED=0 for static binaries
   - Strip binaries with -ldflags="-s -w"
   - Test binaries on different systems

2. **Create release process**
   - Setup GitHub Actions for CI/CD
   - Build on tag push
   - Generate checksums
   - Create release notes
   - Upload artifacts

3. **Create installation methods**
   - Write install.sh script
   - Create systemd user service (optional)
   - Document manual installation
   - Create uninstall instructions

4. **Package for distributions**
   - Create .deb package structure
   - Create .rpm package structure  
   - Create AUR PKGBUILD (later)
   - Document package installation

## Configuration Files

### User State (~/.config/palettesmith/state.json)
```json
{
  "global_palette": {
    "background": "#1e1e2e",
    "foreground": "#cdd6f4",
    "accent": "#89b4fa",
    "accent2": "#f38ba8",
    "muted": "#45475a",
    "border": "#313244"
  },
  "app_overrides": {
    "waybar": {
      "background": "#11111b"
    }
  },
  "last_applied": "2024-01-15T14:30:00Z",
  "last_backup": "2024-01-15-143000"
}
```

### Complete Waybar Plugin Manifest Example (plugins/waybar/manifest.yaml)
```yaml
metadata:
  id: waybar
  name: Waybar
  version: 1.0.0
  description: Waybar status bar theming
  author: palettesmith

detection:
  config_exists:
    - ~/.config/waybar/config
    - ~/.config/waybar/style.css
  binary_exists: waybar
  process_running: waybar

files:
  - path: ~/.config/waybar/style.css
    format: css
    parser: css
    backup: true
    optional: false

colors:
  - id: background
    label: Background
    description: Main background color
    default: "#1e1e2e"
    css_variables:
      - "@define-color background"
    
  - id: foreground
    label: Text
    description: Main text color
    default: "#cdd6f4"
    css_variables:
      - "@define-color foreground"
      - "@define-color text"
    
  - id: accent
    label: Accent
    description: Primary accent color
    default: "#89b4fa"
    css_variables:
      - "@define-color accent"
      - "@define-color blue"
    
  - id: warning
    label: Warning
    description: Warning state color
    default: "#f9e2af"
    css_variables:
      - "@define-color warning"
      - "@define-color yellow"
    
  - id: critical
    label: Critical
    description: Critical/error state color
    default: "#f38ba8"
    css_variables:
      - "@define-color critical"
      - "@define-color red"

restart:
  method: signal
  signal: SIGUSR2
  process: waybar
  fallback: killall waybar && waybar &
```

### Complete Hyprland Plugin Manifest Example (plugins/hyprland/manifest.yaml)
```yaml
metadata:
  id: hyprland
  name: Hyprland
  version: 1.0.0
  description: Hyprland window manager theming
  author: palettesmith

detection:
  config_exists:
    - ~/.config/hypr/hyprland.conf
  binary_exists: hyprctl
  process_running: Hyprland

files:
  - path: ~/.config/hypr/hyprland.conf
    format: conf
    parser: hyprland
    backup: true
    optional: false

colors:
  - id: active_border
    label: Active Border
    description: Active window border color
    default: "89b4fa"
    hypr_variables:
      - "$col.active_border"
    
  - id: inactive_border
    label: Inactive Border
    description: Inactive window border color
    default: "45475a"
    hypr_variables:
      - "$col.inactive_border"
    
  - id: group_border_active
    label: Group Border Active
    description: Active group border color
    default: "f9e2af"
    hypr_variables:
      - "$col.group_border_active"
    
  - id: shadow
    label: Shadow
    description: Window shadow color
    default: "1e1e2e"
    hypr_variables:
      - "$col.shadow"

restart:
  method: command
  command: hyprctl reload
  fallback: killall -SIGUSR2 Hyprland
```

### Complete Alacritty Plugin Manifest Example (plugins/alacritty/manifest.yaml)
```yaml
metadata:
  id: alacritty
  name: Alacritty
  version: 1.0.0
  description: Alacritty terminal emulator theming
  author: palettesmith

detection:
  config_exists:
    - ~/.config/alacritty/alacritty.toml
    - ~/.config/alacritty/alacritty.yml
  binary_exists: alacritty

files:
  - path: ~/.config/alacritty/alacritty.toml
    format: toml
    parser: toml
    backup: true
    optional: false

colors:
  - id: background
    label: Background
    description: Terminal background color
    default: "#1e1e2e"
    toml_path: colors.primary.background
    
  - id: foreground
    label: Foreground
    description: Terminal text color
    default: "#cdd6f4"
    toml_path: colors.primary.foreground
    
  - id: black
    label: Black
    description: ANSI black color
    default: "#45475a"
    toml_path: colors.normal.black
    
  - id: red
    label: Red
    description: ANSI red color
    default: "#f38ba8"
    toml_path: colors.normal.red
    
  - id: green
    label: Green
    description: ANSI green color
    default: "#a6e3a1"
    toml_path: colors.normal.green
    
  - id: yellow
    label: Yellow
    description: ANSI yellow color
    default: "#f9e2af"
    toml_path: colors.normal.yellow
    
  - id: blue
    label: Blue
    description: ANSI blue color
    default: "#89b4fa"
    toml_path: colors.normal.blue
    
  - id: magenta
    label: Magenta
    description: ANSI magenta color
    default: "#f5c2e7"
    toml_path: colors.normal.magenta
    
  - id: cyan
    label: Cyan
    description: ANSI cyan color
    default: "#94e2d5"
    toml_path: colors.normal.cyan
    
  - id: white
    label: White
    description: ANSI white color
    default: "#bac2de"
    toml_path: colors.normal.white

restart:
  method: none  # Alacritty auto-reloads on config change
```

## Error Handling Strategy

1. **Parser Errors**
   - Catch and report parsing errors
   - Fall back to regex-based replacement if parser fails
   - Always backup before modifications

2. **Plugin Errors**
   - Isolate plugin failures
   - Continue with other plugins if one fails
   - Report all errors at the end

3. **File System Errors**
   - Check permissions before operations
   - Handle missing directories
   - Provide clear error messages

4. **Color Validation**
   - Validate color formats
   - Check color contrast (optional)
   - Warn about potential issues

## Performance Considerations

1. **Plugin Loading**
   - Load plugins lazily
   - Cache detection results
   - Only load manifests initially

2. **File Operations**
   - Buffer file reads/writes
   - Use atomic operations
   - Minimize file system calls

3. **TUI Rendering**
   - Only redraw changed components
   - Use virtual scrolling for long lists
   - Debounce rapid updates

## Security Considerations

1. **File Safety**
   - Never modify files outside config directories
   - Validate all paths
   - Use safe path joining

2. **Plugin Security**
   - Validate plugin manifests
   - Sandbox plugin operations
   - Limit plugin permissions

3. **Backup Integrity**
   - Checksum backup files
   - Verify restore operations
   - Limit backup retention

## Success Criteria

1. **Prototype Complete When:**
   - Can detect at least 3 applications
   - Can extract and display colors
   - Can modify colors through TUI
   - Can apply changes and restart apps
   - Can rollback changes

2. **Production Ready When:**
   - All core plugins implemented
   - Full test coverage
   - Binary builds under 10MB
   - Documentation complete
   - No critical bugs for 1 week

## Development Tips for Junior Developers

1. **Start Small**
   - Implement one feature at a time
   - Test as you go
   - Commit frequently

2. **Use the Debugger**
   - Learn to use dlv (Delve) debugger
   - Add log statements liberally
   - Use fmt.Printf for quick debugging

3. **Read the Bubble Tea Tutorial**
   - Go through the official tutorial first
   - Study the example apps
   - Start with simple components

4. **Test Your Parsers**
   - Write tests before implementation
   - Use table-driven tests
   - Test edge cases

5. **Ask for Help**
   - Document what you tried
   - Share error messages
   - Create minimal reproductions

## Common Pitfalls to Avoid

1. **Don't Modify Original Configs Without Backup**
2. **Don't Assume File Formats** - Always validate
3. **Don't Ignore Errors** - Handle or propagate them
4. **Don't Parse Colors with Simple Regex** - Use proper parsers
5. **Don't Hardcode Paths** - Use path expansion
6. **Don't Skip Tests** - They catch bugs early
7. **Don't Overcomplicate** - Start simple, iterate

## Getting Started Checklist

- [ ] Install Go 1.22+ from https://go.dev/dl/
- [ ] Setup your editor (VSCode with Go extension recommended)
- [ ] Create GitHub repository (or local git repo)
- [ ] Clone repo and cd into it
- [ ] Follow Phase 0 completely
- [ ] Run `make dev` to verify setup (should print "PaletteSmith v0.0.1")
- [ ] Read Bubble Tea documentation at https://github.com/charmbracelet/bubbletea
- [ ] Study one existing Bubble Tea app (e.g., https://github.com/charmbracelet/glow)
- [ ] Begin Phase 1

## File Package Declarations Reference

Every Go file needs the correct package declaration at the top. Here's what each should have:

- `cmd/palettesmith/main.go` → `package main`
- `internal/core/*.go` → `package core`
- `internal/plugin/*.go` → `package plugin`
- `internal/parser/*.go` → `package parser`
- `internal/tui/*.go` → `package tui`
- `internal/cli/*.go` → `package cli`
- `internal/config/*.go` → `package config`
- `plugins/waybar/plugin.go` → `package waybar`
- `plugins/hyprland/plugin.go` → `package hyprland`
- `plugins/alacritty/plugin.go` → `package alacritty`

## Initial Directory Creation Commands

Run these commands from your project root to create the structure:

```bash
# Create all directories
mkdir -p cmd/palettesmith
mkdir -p internal/{core,plugin,parser,tui,cli,config}
mkdir -p plugins/{waybar,hyprland,alacritty,kitty,btop,mako}
mkdir -p scripts testdata/{configs,themes} docs

# Create empty Go files with correct package declarations
echo 'package main' > cmd/palettesmith/main.go
echo 'package core' > internal/core/color.go
echo 'package core' > internal/core/palette.go
echo 'package core' > internal/core/state.go
echo 'package core' > internal/core/backup.go
echo 'package plugin' > internal/plugin/interface.go
echo 'package plugin' > internal/plugin/loader.go
echo 'package plugin' > internal/plugin/registry.go
echo 'package plugin' > internal/plugin/manifest.go
echo 'package parser' > internal/parser/interface.go
echo 'package parser' > internal/parser/css.go
echo 'package parser' > internal/parser/json.go
echo 'package parser' > internal/parser/toml.go
echo 'package parser' > internal/parser/yaml.go
echo 'package parser' > internal/parser/ini.go
echo 'package parser' > internal/parser/hyprland.go
echo 'package tui' > internal/tui/app.go
echo 'package tui' > internal/tui/navigation.go
echo 'package tui' > internal/tui/colorpicker.go
echo 'package tui' > internal/tui/preview.go
echo 'package tui' > internal/tui/list.go
echo 'package tui' > internal/tui/styles.go
echo 'package cli' > internal/cli/root.go
echo 'package cli' > internal/cli/list.go
echo 'package cli' > internal/cli/get.go
echo 'package cli' > internal/cli/set.go
echo 'package cli' > internal/cli/apply.go
echo 'package cli' > internal/cli/rollback.go
echo 'package cli' > internal/cli/export.go
echo 'package cli' > internal/cli/import.go
echo 'package config' > internal/config/paths.go
echo 'package config' > internal/config/settings.go
echo 'package waybar' > plugins/waybar/plugin.go
echo 'package hyprland' > plugins/hyprland/plugin.go
echo 'package alacritty' > plugins/alacritty/plugin.go

# Create scripts (copy content from plan)
touch scripts/{dev.sh,build.sh,test.sh,install.sh}
chmod +x scripts/*.sh

# Create Makefile
touch Makefile

# Initialize git
git init
echo -e "palettesmith\ndist/\n*.log\n.air.toml" > .gitignore
```

## Questions to Answer Before Starting

1. What should happen if a plugin fails during apply?
2. Should we support custom color names or only standard ones?
3. How many backups should we keep?
4. Should the TUI auto-refresh when files change externally?
5. What's the minimum terminal size we support?
6. Should we support remote plugin repositories?
7. How do we handle config files with includes/imports?
8. Should we support color scheme import from other tools?

## Code Implementation Details

### Color Type Implementation (core/color.go)

The Color type must support these exact formats:
- Hex: `#RGB`, `#RRGGBB`, `#RRGGBBAA`
- Hyprland: `rgb(RRGGBB)`, `rgba(RRGGBBAA)` - note: no # symbol
- Numeric RGB: `rgb(255, 255, 255)`, `rgba(255, 255, 255, 1.0)`
- HSL: `hsl(360, 100%, 50%)`, `hsla(360, 100%, 50%, 1.0)`

Key methods to implement:
- `ParseColor(string) (*Color, error)` - detect format and parse
- `ToHex() string` - always returns #RRGGBB format
- `ToFormat(format string) string` - convert to specific format
- `IsValid() bool` - validate color values

### Plugin Interface Implementation (plugin/interface.go)

```go
type Plugin interface {
    GetManifest() *PluginManifest
    Detect() bool
    GetFiles() []string
    ExtractColors(state *State) (map[string]*Color, error)
    ApplyColors(state *State, colors map[string]*Color) error
    Restart() error
    Validate(colors map[string]*Color) error
}

type PluginManifest struct {
    Metadata   Metadata
    Detection  Detection
    Files      []FileConfig
    Colors     []ColorDefinition
    Restart    RestartConfig
}
```

### Plugin Loading Mechanism (plugin/loader.go)

1. **Built-in plugins**: Compile directly into binary
   - Each plugin in `plugins/*/plugin.go` implements the Plugin interface
   - Register in a map: `var builtinPlugins = map[string]Plugin{}`
   - Add each plugin to the map in init() functions

2. **External plugins**: Load from YAML manifests
   - Create a GenericPlugin type that reads YAML and implements Plugin interface
   - Use reflection or code generation to handle different formats
   - External plugins can only use parsers, not custom Go code

3. **Loading order**:
   - Load built-ins first
   - Then scan external directories
   - If duplicate IDs, external overrides built-in

### Parser Implementation Details

#### CSS Parser (parser/css.go)
- Use `github.com/tdewolff/parse/v2/css` to tokenize
- Walk tokens looking for color values
- Handle @define-color specially for GTK CSS
- Preserve comments by keeping original text and replacing only color values
- For replacement: use byte offsets to preserve formatting

#### JSON Parser (parser/json.go)
- Use `json.Decoder` with `UseNumber()` to preserve numeric formats
- Walk the decoded interface{} recursively
- For JSONPath support, implement simple dot notation: `colors.primary.background`
- When replacing, marshal back with `json.MarshalIndent` for formatting

#### TOML Parser (parser/toml.go)
- Use `github.com/BurntSushi/toml` with `MetaData` to get positions
- For comments, read file as text and merge after parsing
- Support paths like `colors.primary.background` using reflection
- Preserve order using `toml.Marshal` with original structure

#### Hyprland Parser (parser/hyprland.go)
- Custom parser for Hyprland's format
- Look for patterns: `$variable = rgba(RRGGBBAA)` or `$variable = rgb(RRGGBB)`
- Handle both with and without `0x` prefix
- Preserve exact format when replacing (if was rgba, keep as rgba)

### TUI Component Behaviors

#### Main App (tui/app.go)
States the app can be in:
1. `NavigationFocus` - sidebar has focus, arrow keys move between apps
2. `ListFocus` - color list has focus, arrow keys move between colors
3. `ColorPickerOpen` - modal is open, tab moves between inputs
4. `ApplyInProgress` - applying changes, show progress
5. `HelpOpen` - help screen visible

Key bindings:
- `Tab/Shift+Tab` - cycle between sidebar and list
- `Enter` - open color picker when color selected
- `Ctrl+S` - apply changes
- `Ctrl+Z` - undo last change (in memory)
- `q/Ctrl+C` - quit (with confirmation if unsaved)
- `?/F1` - show help

#### Color Picker (tui/colorpicker.go)
Input modes:
1. Hex input - text field, validate on type
2. RGB sliders - 3 sliders, 0-255 range
3. HSL sliders - H: 0-360, S/L: 0-100
4. Preset grid - 20 common colors in 4x5 grid

Behavior:
- All inputs update simultaneously when one changes
- Preview shows old and new color side by side
- Enter or button confirms, Escape cancels
- If hyprpicker available, show "Pick from Screen" button

#### Navigation (tui/navigation.go)
Display format:
```
Global           ✓
━━━━━━━━━━━━━━━━━
Waybar          ✓
Hyprland        ✓
Alacritty       ✓
Kitty           ✗
Btop            ✓
```
- ✓ = detected and available
- ✗ = not detected
- Highlight current selection with background color
- Show count of overrides if any: "Waybar (3)"

### CLI Command Examples

```bash
# List all detected apps and their colors
palettesmith list
palettesmith list --json

# Get current colors
palettesmith get global
palettesmith get app waybar
palettesmith get app waybar --format=hex

# Set colors
palettesmith set global background=#1e1e2e foreground=#cdd6f4
palettesmith set app waybar background=#11111b
palettesmith set --from-file mytheme.json

# Apply changes
palettesmith apply --dry-run
palettesmith apply --app waybar,hyprland
palettesmith apply --no-restart

# Rollback
palettesmith rollback
palettesmith rollback --to 2024-01-15-143000

# Export/Import
palettesmith export > mytheme.json
palettesmith import mytheme.json --merge
```

## Troubleshooting Guide

### Common Build Issues

1. **"cannot find package" errors**
   - Run `go mod tidy` to download dependencies
   - Check you're in the project root directory
   - Verify Go version with `go version` (needs 1.22+)

2. **"undefined: PluginInterface" errors**
   - Make sure all files have correct package declarations
   - Files in `internal/plugin/` should have `package plugin`
   - Run `go build ./...` to check all packages compile

3. **Binary too large**
   - Ensure using `CGO_ENABLED=0` when building
   - Use `-ldflags="-s -w"` to strip debug info
   - Check not embedding unnecessary assets

### Common Runtime Issues

1. **"permission denied" when applying**
   - Check config files are writable
   - Ensure running as correct user (not root)
   - Verify backup directory is writable

2. **Colors not changing**
   - Check the app's config format matches manifest
   - Verify restart command is working
   - Look for syntax errors in modified configs
   - Some apps cache configs - may need full restart

3. **Plugin not detected**
   - Check manifest.yaml syntax is valid
   - Verify detection conditions (file exists, etc.)
   - Run with `--verbose` to see detection attempts
   - Check plugin ID is unique

4. **TUI rendering issues**
   - Verify terminal supports 256 colors: `echo $TERM`
   - Try different terminal (Alacritty, Kitty recommended)
   - Check locale is UTF-8: `locale`
   - Resize terminal if layout broken

### Debug Mode

Add debug logging throughout:
```go
if verbose {
    log.Printf("Loading plugin: %s", manifest.Metadata.ID)
}
```

Run with verbose flag:
```bash
palettesmith --verbose
PALETTESMITH_DEBUG=1 palettesmith
```

### Testing Individual Components

```bash
# Test parser alone
go test ./internal/parser -v -run TestCSSParser

# Test specific plugin
go test ./plugins/waybar -v

# Test with coverage
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Decisions to Make

These questions from the original plan need answers:

1. **What should happen if a plugin fails during apply?**
   → Continue with other plugins, report all errors at end, mark state as "partially applied"

2. **Should we support custom color names or only standard ones?**
   → Support any names in global palette, but map to standard names for apps

3. **How many backups should we keep?**
   → Keep last 10 by default, configurable in settings

4. **Should the TUI auto-refresh when files change externally?**
   → No, too complex for v1. Add "Refresh" command instead

5. **What's the minimum terminal size we support?**
   → 80x24 minimum, degrade gracefully below that

6. **Should we support remote plugin repositories?**
   → Not in v1, just local directories

7. **How do we handle config files with includes/imports?**
   → Parse main file only in v1, warn about includes

8. **Should we support color scheme import from other tools?**
   → Yes, support base16 JSON format for import/export

## Next Steps After Prototype

1. Community feedback and testing
2. Plugin repository setup
3. Distribution packages
4. Integration with Omarchy
5. GUI version (optional)
6. Cloud sync (optional)