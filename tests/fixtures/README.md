# Test Fixtures

This directory contains sample configurations and test data for validating palettesmith functionality.

## Directory Structure

```
fixtures/
├── configs/           # Sample application configs
│   ├── hyprland.conf
│   ├── alacritty.toml
│   └── waybar.css
├── plugins/           # Test plugin definitions
│   ├── test-app/
│   └── malformed-plugin/
├── themes/            # Sample theme directories
│   ├── catppuccin/
│   └── test-theme/
└── staging/           # Staging workflow test data
    ├── backup/
    └── next/
```

## Usage

These fixtures are used by integration tests in `tests/integration/` to verify:
- Plugin loading and validation
- Config file parsing and modification
- Theme import/export workflows
- Staging and rollback operations

Unit tests (colocated with source code) create minimal test data inline and should not use these fixtures unless necessary.

## Adding New Fixtures

When adding new test data:
1. Follow the existing directory structure
2. Include both valid and invalid examples
3. Document any special test requirements
4. Ensure fixtures work across different systems