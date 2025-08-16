# Architecture Overview

## Directory Structure

```
cmd/palettesmith/           # Application entry point
internal/                   # Business logic (unexported)
  ├── config/              # Configuration management
  ├── plugin/              # Plugin discovery and loading
  ├── theme/               # Theme management and resolution
  └── validation/          # Input validation service (future)
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

## Data Flow

1. **Plugin Discovery**: `internal/plugin` scans `./plugins/` directories
2. **Theme Loading**: `internal/theme` manages defaults + overrides 
3. **UI Consumption**: `ui/tui` consumes business logic without modification
4. **Config Application**: Plugins define how to read/write actual config files

## Plugin Architecture

### Manifest (plugin.json)
```json
{
  "id": "hyprland",
  "title": "Hyprland", 
  "spec": "spec.json",
  "user_paths": ["~/.config/hypr/hyprland.conf"],
  "system_paths": ["/etc/xdg/hypr/hyprland.conf"],
  "reload": ["hyprctl", "reload"],
  "templates": "templates.json"
}
```

### Specification (spec.json)
```json
{
  "id": "hyprland",
  "fields": [
    {
      "key": "accent",
      "label": "Accent Color",
      "type": "color",
      "default": "#89b4fa",
      "color": {"format": "hex6"},
      "help": "Highlight and focus color"
    }
  ]
}
```

## Theme Resolution Hierarchy

1. **User Override** (per-target customization)
2. **Theme Default** (global theme values)  
3. **Field Default** (plugin-defined fallback)

## Configuration System

### Palettesmith Config (~/.config/palettesmith/config.json)
```json
{
  "target_theme_dir": "~/.config/omarchy/themes",
  "current_theme_link": "~/.config/omarchy/current/theme", 
  "preset": "omarchy"
}
```

### Directory Layout
```
~/.config/palettesmith/
├── config.json           # System integration config
├── themes/              # User-created themes
│   └── my-theme/        # Custom theme directory
└── staging/             # Temporary changes during editing
    ├── backup/          # Original files backup
    └── next/            # Modified files preview
```

## Safety Principles

- **Always backup** original configs before modification
- **Show diffs** before applying changes
- **Confirm destructive operations** (reload commands, file overwrites)
- **Graceful degradation** when plugins fail to load
- **Rollback capability** at field and theme level