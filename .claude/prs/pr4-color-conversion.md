# PR4: Color Format Conversion

## Overview
Implement comprehensive color format conversion system to handle different output formats required by various applications.

## Acceptance Criteria

### Color Format Support
- [ ] Support hex formats: #RGB, #RRGGBB, #RRGGBBAA
- [ ] Support rgb/rgba formats: rgb(r,g,b), rgba(r,g,b,a)
- [ ] Support HSL/HSLA formats: hsl(h,s,l), hsla(h,s,l,a)
- [ ] Support application-specific formats (0xRRGGBB, etc.)
- [ ] Handle alpha channel preservation and removal

### Internal Color Representation
- [ ] Define standard internal format (RGBA with float values)
- [ ] Create color parsing from all input formats
- [ ] Implement conversion to all output formats
- [ ] Preserve color accuracy during conversions
- [ ] Handle edge cases (transparent, out-of-range values)

### Validation Enhancement
- [ ] Validate colors can be represented in target format
- [ ] Show warnings for precision loss
- [ ] Prevent invalid color assignments
- [ ] Provide format-specific error messages

## Files to Create

### `internal/color/color.go`
```go
package color

type Color struct {
    R, G, B, A float64 // Values 0.0-1.0
}

type Format string

const (
    FormatHex3     Format = "hex3"      // #RGB
    FormatHex6     Format = "hex6"      // #RRGGBB
    FormatHex8     Format = "hex8"      // #RRGGBBAA
    FormatRGB      Format = "rgb"       // rgb(r,g,b)
    FormatRGBA     Format = "rgba"      // rgba(r,g,b,a)
    FormatHSL      Format = "hsl"       // hsl(h,s,l)
    FormatHSLA     Format = "hsla"      // hsla(h,s,l,a)
    FormatHex0x    Format = "hex0x"     // 0xRRGGBB
    FormatRGBNoPrefix Format = "rgb_no_prefix" // RRGGBB
)

func ParseColor(input string) (Color, error) {
    // Parse any supported color format
}

func (c Color) ToFormat(format Format) (string, error) {
    // Convert to specified output format
}

func (c Color) CanRepresentIn(format Format) bool {
    // Check if color can be accurately represented
}
```

### `internal/color/converter.go`
```go
package color

type Converter struct {
    formatSpecs map[string]FormatSpec
}

type FormatSpec struct {
    Format          Format   `json:"format"`
    SupportsAlpha   bool     `json:"supports_alpha"`
    BitDepth        int      `json:"bit_depth"`
    Template        string   `json:"template"`
    ValidationRegex string   `json:"validation_regex"`
}

func NewConverter() *Converter {
    // Initialize with all supported formats
}

func (c *Converter) Convert(input string, targetFormat Format) (string, error) {
    // Parse input and convert to target format
}

func (c *Converter) ValidateForFormat(input string, format Format) error {
    // Validate color is valid for specific format
}
```

### `internal/color/formats.go`
```go
package color

// Define all supported color formats and their specifications

var SupportedFormats = map[Format]FormatSpec{
    FormatHex6: {
        Format:          FormatHex6,
        SupportsAlpha:   false,
        BitDepth:        8,
        Template:        "#%02X%02X%02X",
        ValidationRegex: "^#[0-9a-fA-F]{6}$",
    },
    FormatRGBNoPrefix: {
        Format:          FormatRGBNoPrefix,
        SupportsAlpha:   false,
        BitDepth:        8,
        Template:        "%02x%02x%02x",
        ValidationRegex: "^[0-9a-fA-F]{6}$",
    },
    // ... other formats
}
```

## Files to Modify

### `internal/plugin/spec.go`
```go
// Add color format specification to field spec
type ColorSpec struct {
    Format          Format `json:"format"`
    SupportsAlpha   bool   `json:"supports_alpha,omitempty"`
    OutputFormat    Format `json:"output_format,omitempty"`
}

type Field struct {
    // ... existing fields
    Color *ColorSpec `json:"color,omitempty"`
}
```

### `internal/validation/service.go`
- Add color format validation
- Check target format compatibility
- Provide specific error messages for color issues

### `ui/tui/form.go`
- Use color conversion service
- Show format-specific validation errors
- Display color preview with target format

## Enhanced Plugin Specifications

### Updated `spec.json` Format
```json
{
  "id": "hyprland",
  "fields": [
    {
      "key": "accent",
      "type": "color",
      "color": {
        "format": "hex6",
        "output_format": "rgb_no_prefix",
        "supports_alpha": false
      },
      "default": "#89b4fa"
    },
    {
      "key": "opacity",
      "type": "color", 
      "color": {
        "format": "hex8",
        "output_format": "rgba",
        "supports_alpha": true
      },
      "default": "#89b4fa80"
    }
  ]
}
```

## Testing Requirements

### Unit Tests
- [ ] Test color parsing from all input formats
- [ ] Test conversion to all output formats
- [ ] Test alpha channel handling
- [ ] Test precision preservation
- [ ] Test invalid color handling
- [ ] Test edge cases (black, white, transparent)

### Conversion Tests
Create test cases for:
- Hex6 ↔ RGB ↔ HSL conversions
- Alpha channel preservation/removal
- Precision loss detection
- Format compatibility checking

### Integration Tests
- [ ] Test with real plugin specifications
- [ ] Verify template output matches expected format
- [ ] Test validation integration with form system

## Manual Testing Steps

1. **Format Conversion Test**
   - Enter colors in different formats
   - Verify conversion to target formats
   - Check precision is maintained
   - Test alpha channel handling

2. **Plugin Integration Test**
   - Configure Hyprland plugin with rgb_no_prefix
   - Configure Alacritty plugin with hex6
   - Verify outputs match expected formats
   - Check template substitution works

3. **Edge Case Test**
   - Test with pure colors (red, green, blue)
   - Test with transparent colors
   - Test with out-of-range values
   - Verify error messages are helpful

4. **Validation Test**
   - Enter invalid color formats
   - Test incompatible format combinations
   - Verify warnings for precision loss
   - Check error message clarity

## Error Handling

### Required Error Messages
- "Invalid color format: expected {format}"
- "Color {color} cannot be represented in {format}"
- "Alpha channel not supported in {format}"
- "Precision loss converting to {format}"

### Validation Rules
- Colors must be valid in input format
- Target format must support color representation
- Alpha channels handled appropriately
- Range validation for RGB/HSL values

## Dependencies
- PR1 (Project Restructure) must be completed
- PR2 (Enhanced Plugins) recommended for integration

## Definition of Done
- [ ] All color format conversions work correctly
- [ ] No precision loss in supported conversions
- [ ] Plugin specifications support color formats
- [ ] Validation prevents incompatible assignments
- [ ] All tests pass with good coverage
- [ ] Manual testing completed
- [ ] Performance acceptable for real-time conversion
- [ ] Code review completed