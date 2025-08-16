# PR5: Basic Grid Color Picker

## Overview
Implement a grid-based color picker component that can be used by any UI implementation, with keyboard navigation and visual feedback.

## Acceptance Criteria

### Color Grid Component
- [ ] Generate color palette grids (hue, saturation, lightness)
- [ ] Support different grid sizes and arrangements
- [ ] Provide keyboard navigation (arrow keys)
- [ ] Show current color selection clearly
- [ ] Support color preview and comparison

### UI Integration
- [ ] Integrate with existing form system
- [ ] Replace text-only color input with picker
- [ ] Maintain compatibility with direct hex input
- [ ] Show both picker and text input options

### Navigation System
- [ ] Arrow keys move through color grid
- [ ] Enter/Space selects current color
- [ ] Escape cancels color selection
- [ ] Tab switches between grid areas (hue/sat/light)
- [ ] Support vim-style navigation (optional)

### Visual Design
- [ ] Show color swatches clearly in terminal
- [ ] Highlight currently selected color
- [ ] Show current theme colors for reference
- [ ] Display color values in multiple formats

## Files to Create

### `internal/colorpicker/grid.go`
```go
package colorpicker

type Grid struct {
    Width     int
    Height    int
    Colors    [][]Color
    CurrentX  int
    CurrentY  int
    Type      GridType
}

type GridType string

const (
    GridTypeHue        GridType = "hue"
    GridTypeSaturation GridType = "saturation"
    GridTypeLightness  GridType = "lightness"
    GridTypeCustom     GridType = "custom"
)

func NewHueGrid(width, height int) *Grid {
    // Generate hue-based color grid
}

func NewSaturationGrid(baseHue float64, width, height int) *Grid {
    // Generate saturation/lightness grid for given hue
}

func (g *Grid) MoveUp() bool    { /* Move cursor up */ }
func (g *Grid) MoveDown() bool  { /* Move cursor down */ }
func (g *Grid) MoveLeft() bool  { /* Move cursor left */ }
func (g *Grid) MoveRight() bool { /* Move cursor right */ }

func (g *Grid) GetCurrentColor() Color {
    return g.Colors[g.CurrentY][g.CurrentX]
}
```

### `internal/colorpicker/picker.go`
```go
package colorpicker

type Picker struct {
    hueGrid    *Grid
    satGrid    *Grid
    activeGrid GridType
    baseColor  Color
    callback   func(Color)
}

func NewPicker(initialColor Color) *Picker {
    // Initialize picker with grids
}

func (p *Picker) SetBaseColor(c Color) {
    // Update saturation grid based on selected hue
}

func (p *Picker) HandleKeyPress(key string) bool {
    // Handle navigation and selection
}

func (p *Picker) GetSelectedColor() Color {
    // Return currently selected color
}

func (p *Picker) SetCallback(fn func(Color)) {
    // Set callback for color selection
}
```

### `ui/tui/colorpicker.go`
```go
package tui

type ColorPickerModel struct {
    picker     *colorpicker.Picker
    width      int
    height     int
    showPicker bool
    currentVal string
}

func NewColorPickerModel(initialColor string) ColorPickerModel {
    // Initialize TUI color picker component
}

func (m ColorPickerModel) Update(msg tea.Msg) (ColorPickerModel, tea.Cmd) {
    // Handle Bubble Tea messages
}

func (m ColorPickerModel) View() string {
    // Render color picker or text input
}

func (m ColorPickerModel) TogglePicker() ColorPickerModel {
    // Switch between picker and text input
}
```

## Files to Modify

### `ui/tui/form.go`
- Replace simple text input with color picker component
- Add toggle between picker and text input modes
- Integrate keyboard shortcuts for picker activation
- Show current theme colors as reference

### `internal/validation/service.go`
- Add validation for colors selected from picker
- Ensure picked colors are compatible with target format
- Validate color transitions and accessibility

## Color Grid Layout

### Hue Grid (Primary Selection)
```
[Red][Orange][Yellow][Green][Cyan][Blue][Purple][Pink]
```

### Saturation/Lightness Grid (Secondary Selection)
```
    Light ←→ Dark
  ┌─────────────────┐ ↑
S │ ░░░▓▓▓▓▓▓▓▓███ │ │ Saturated
a │ ░░░▓▓▓▓▓▓▓▓███ │ │
t │ ░░░▓▓▓▓▓▓▓▓███ │ │
  │ ░░░▓▓▓▓▓▓▓▓███ │ │
  │ ░░░▓▓▓▓▓▓▓▓███ │ ↓ Desaturated
  └─────────────────┘
```

### Reference Palette
Show current theme colors:
```
Current Theme: [bg] [fg] [accent] [border] [...]
```

## Testing Requirements

### Unit Tests
- [ ] Test grid generation algorithms
- [ ] Test navigation boundaries (edge cases)
- [ ] Test color calculation accuracy
- [ ] Test keyboard input handling
- [ ] Test color format integration

### Visual Tests
Create fixtures for:
- Different terminal color capabilities
- Various grid sizes and layouts
- Color accuracy verification
- Navigation flow testing

### Integration Tests
- [ ] Test integration with form system
- [ ] Test color picker with different field types
- [ ] Verify picked colors match expected formats
- [ ] Test theme color reference display

## Manual Testing Steps

1. **Grid Navigation Test**
   - Open color picker for accent field
   - Test arrow key navigation in all directions
   - Verify cursor wraps appropriately at edges
   - Test tab switching between grid areas

2. **Color Selection Test**
   - Select colors from different grid areas
   - Verify color values are accurate
   - Test color format conversion
   - Check visual feedback is clear

3. **Integration Test**
   - Use color picker in form context
   - Test switching between picker and text input
   - Verify selected colors update form properly
   - Check validation works with picked colors

4. **Reference Display Test**
   - Verify current theme colors show correctly
   - Test color comparison functionality
   - Check color accessibility in terminal

## Keyboard Shortcuts

### Navigation
- `Arrow Keys` - Move through color grid
- `Tab` - Switch between hue/saturation grids
- `Enter/Space` - Select current color
- `Escape` - Cancel picker, return to text input

### Optional (Future)
- `h/j/k/l` - Vim-style navigation
- `t` - Toggle picker/text input mode
- `r` - Reset to original color

## Error Handling

### Required Error Messages
- "Color picker not available in this terminal"
- "Invalid color selected"
- "Cannot convert color to required format"

### Graceful Degradation
- Fall back to text input if terminal lacks color support
- Provide clear visual feedback for selected colors
- Handle terminal resize gracefully

## Dependencies
- PR4 (Color Format Conversion) must be completed
- PR1 (Project Restructure) for UI separation

## Definition of Done
- [ ] Color picker displays correctly in terminal
- [ ] Keyboard navigation works smoothly
- [ ] Selected colors integrate with form system
- [ ] Color accuracy verified across formats
- [ ] All tests pass
- [ ] Manual testing completed across different terminals
- [ ] Performance acceptable for real-time navigation
- [ ] Code review completed
- [ ] Documentation includes usage examples