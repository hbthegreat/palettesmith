# Palettesmith

> **⚠️ DEVELOPMENT STATUS**: This project is currently under active development and is not yet ready for production use. Features and APIs may change significantly. Please do not use this in production environments.

A powerful TUI application for theme management with a plugin-based architecture. Create, edit, and apply color themes across multiple applications with sophisticated template-based configuration generation.

## Features

- 🎨 **Visual Theme Editor** - Interactive TUI for creating and editing color themes
- 🔌 **Plugin System** - JSON-based plugin architecture supporting any application
- 📝 **Advanced Templates** - Go text/template engine with color manipulation functions
- ✅ **Validation System** - Comprehensive field and plugin validation with error reporting
- 🛠️ **Developer Tools** - CLI validation tool for plugin development and CI/CD
- 🔄 **Multiple Presets** - Support for different workspace configurations (generic, omarchy)
- 🎯 **Safe Application** - Backup and preview changes before applying to your system

## Quick Start

### Installation

```bash
git clone https://github.com/yourusername/palettesmith
cd palettesmith
go build -o palettesmith cmd/palettesmith/main.go
```

### Basic Usage

```bash
# Start the interactive TUI
./palettesmith

# Validate all plugins (useful for development)
./palettesmith --validate-plugins

# Get help
./palettesmith --help
```

### First Run Setup

When you first run Palettesmith, you'll be guided through a setup process to choose your workspace preset:

- **Generic**: Standard configuration suitable for most setups
- **Omarchy**: Optimized for the Omarchy workspace environment

## Supported Applications

Palettesmith comes with built-in support for popular Linux applications:

| Application | Description | Template Features |
|-------------|-------------|-------------------|
| **Alacritty** | Terminal emulator | Color schemes, transparency |
| **BTTop** | System monitor | Color themes, interface styling |
| **Hyprland** | Wayland compositor | Window decorations, workspace themes |
| **Hyprlock** | Screen locker | Background colors, lock screen styling |
| **Mako** | Notification daemon | Alert colors, notification styling |
| **SwayOSD** | On-screen display | OSD colors and styling |
| **Walker** | Application launcher | Interface themes, selection colors |
| **Waybar** | Status bar | Bar themes, module styling |

## Plugin Development

### Plugin Structure

Each plugin consists of two main files in `./plugins/<id>/`:

```
plugins/myapp/
├── plugin.json     # Plugin manifest
└── spec.json       # Field specifications
```

### Plugin Manifest (`plugin.json`)

```json
{
  "id": "myapp",
  "title": "My Application",
  "spec": "spec.json",
  "user_paths": ["~/.config/myapp/config.conf"],
  "system_paths": ["/etc/myapp/config.conf"],
  "reload": ["myapp", "--reload"],
  "templates": "templates.json"
}
```

### Field Specification (`spec.json`)

```json
{
  "id": "myapp",
  "fields": [
    {
      "key": "primary_color",
      "label": "Primary Color",
      "type": "color",
      "default": "#89b4fa",
      "color": {"format": "hex6"},
      "help": "Main application accent color"
    },
    {
      "key": "opacity",
      "label": "Window Opacity",
      "type": "number",
      "default": "0.9",
      "number": {"min": 0.0, "max": 1.0, "step": 0.1},
      "help": "Window transparency level"
    }
  ]
}
```

### Template Functions

Palettesmith provides powerful template functions for color manipulation:

```go
// Convert hex color to RGBA format
{{hexToRGBA .primary_color 1.0}}
// Output: rgba(137,180,250,1.0)

// Apply alpha transparency
{{alpha .primary_color 0.8}}
// Output: rgba(137,180,250,0.8)

// Brighten color by percentage
{{brighten .primary_color 0.3}}
// Output: #b5d4fb

// Mix two colors
{{mix .primary_color .secondary_color 0.5}}
// Output: blended color at 50% mix

// Standard string functions
{{trimPrefix .value "prefix-"}}
{{toUpper .text}}
{{toLower .text}}
```

### Plugin Validation

Validate your plugins during development:

```bash
# Validate all plugins
./palettesmith --validate-plugins

# Example output
Validating plugins...
Found 8 plugins:

▸ myapp (My Application)
  ✓ Plugin is valid

Validation complete: ✓ All plugins valid
```

## Development

### Prerequisites

- Go 1.19 or later
- Linux environment (Wayland/X11)

### Development Commands

```bash
# Run the application
go run cmd/palettesmith/main.go

# Build executable
go build -o palettesmith cmd/palettesmith/main.go

# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Static analysis
go vet ./...

# Format code
go fmt ./...
```

### Project Structure

```
cmd/palettesmith/           # Application entry point
internal/                   # Business logic (unexported)
  ├── config/              # Configuration management
  ├── plugin/              # Plugin discovery and loading
  ├── theme/               # Theme management and resolution
  └── validation/          # Input validation service
ui/                        # Interface implementations
  ├── tui/                # Terminal interface (Bubble Tea)
  └── native/             # Future native UI
  └── web/                # Future web UI
plugins/                   # Plugin definitions
  └── <id>/
    ├── plugin.json        # Plugin manifest
    └── spec.json          # Field specifications
tests/
  ├── fixtures/            # Sample configs for testing
  └── integration/         # Integration tests
```

### Architecture Principles

- **Separation of Concerns**: Business logic in `internal/`, UI implementations in `ui/`
- **Plugin System**: JSON-based plugin discovery from `./plugins/<id>/` directories
- **Theme Safety**: Never break existing user configurations - protection first, enhancement second
- **Modular Design**: Each component must be independently testable and UI-agnostic

## Configuration

### System Configuration

Palettesmith stores its configuration in `~/.config/palettesmith/config.json`:

```json
{
  "target_theme_dir": "~/.config/omarchy/themes",
  "current_theme_link": "~/.config/omarchy/current/theme",
  "preset": "omarchy",
  "first_run": false
}
```

### Theme Resolution Hierarchy

1. **User Override** (per-target customization)
2. **Theme Default** (global theme values)
3. **Field Default** (plugin-defined fallback)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Validate plugins (`./palettesmith --validate-plugins`)
7. Commit your changes (`git commit -am 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Plugin Contributions

When contributing new plugins:

1. Follow the plugin structure guidelines
2. Include comprehensive field specifications
3. Test with `--validate-plugins` flag
4. Add integration tests in `tests/integration/`
5. Document any special requirements

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/yourusername/palettesmith/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/yourusername/palettesmith/discussions)
- 📖 **Documentation**: Check the [.claude/](/.claude/) directory for detailed architecture docs

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- Uses [Bubbles](https://github.com/charmbracelet/bubbles) for UI components
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) for terminal styling
