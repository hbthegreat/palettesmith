package color

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Converter handles color format conversion between different representations.
// It normalizes all inputs to a hex6 format internally (#ffffff) and provides
// methods to output in various formats required by different config files.
type Converter struct {
	value string // Always stored as lowercase hex6 with # prefix
}

// NewConverter creates a new color converter from various input formats.
// Supported formats:
//   - #ffffff, #fff (hex with hash)
//   - ffffff, fff (hex without hash)
//   - 0xffffff, 0Xffffff (hex with 0x prefix)
//
// All inputs are normalized to lowercase hex6 format internally.
func NewConverter(input string) (*Converter, error) {
	if input == "" {
		return nil, fmt.Errorf("color input cannot be empty")
	}

	// Normalize whitespace
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("color input cannot be empty")
	}

	// Convert to lowercase for consistent processing
	input = strings.ToLower(input)

	var hexValue string
	var err error

	switch {
	case strings.HasPrefix(input, "0x"):
		hexValue, err = parseHex0x(input)
	case strings.HasPrefix(input, "#"):
		hexValue, err = parseHexWithHash(input)
	default:
		hexValue, err = parseHexNoPrefix(input)
	}

	if err != nil {
		return nil, err
	}

	return &Converter{value: hexValue}, nil
}

// ToHex6 returns the color in #ffffff format
func (c *Converter) ToHex6() string {
	if c == nil {
		return ""
	}
	return c.value
}

// ToHex6NoPrefix returns the color in ffffff format (no # prefix)
func (c *Converter) ToHex6NoPrefix() string {
	if c == nil {
		return ""
	}
	return strings.TrimPrefix(c.value, "#")
}

// ToHex0x returns the color in 0xffffff format
func (c *Converter) ToHex0x() string {
	if c == nil {
		return ""
	}
	return "0x" + strings.TrimPrefix(c.value, "#")
}

// ToRGBNoPrefix returns the color in ffffff format (for Hyprland rgb() usage)
// This is the same as ToHex6NoPrefix but named for clarity of intent
func (c *Converter) ToRGBNoPrefix() string {
	if c == nil {
		return ""
	}
	return strings.TrimPrefix(c.value, "#")
}

// ToFormat converts to the specified format string.
// Supported formats: "hex6", "hex6_no_prefix", "hex_0x", "rgb_no_prefix"
// Returns hex6 format for unknown format strings.
func (c *Converter) ToFormat(format string) string {
	if c == nil {
		return ""
	}

	switch format {
	case "hex6":
		return c.ToHex6()
	case "hex6_no_prefix":
		return c.ToHex6NoPrefix()
	case "hex_0x":
		return c.ToHex0x()
	case "rgb_no_prefix":
		return c.ToRGBNoPrefix()
	default:
		// Return hex6 as safe default for unknown formats
		return c.ToHex6()
	}
}

// parseHex0x parses "0xffffff" or "0xfff" format
func parseHex0x(input string) (string, error) {
	if len(input) < 3 {
		return "", fmt.Errorf("invalid 0x color format: too short")
	}

	hexPart := input[2:] // Remove "0x" prefix
	return parseHexDigits(hexPart)
}

// parseHexWithHash parses "#ffffff" or "#fff" format  
func parseHexWithHash(input string) (string, error) {
	if len(input) < 2 {
		return "", fmt.Errorf("invalid hex color format: too short")
	}

	hexPart := input[1:] // Remove "#" prefix
	return parseHexDigits(hexPart)
}

// parseHexNoPrefix parses "ffffff" or "fff" format
func parseHexNoPrefix(input string) (string, error) {
	return parseHexDigits(input)
}

// parseHexDigits validates and normalizes hex digits to hex6 format
func parseHexDigits(hexPart string) (string, error) {
	if hexPart == "" {
		return "", fmt.Errorf("invalid hex color: no hex digits provided")
	}

	// Validate hex characters
	hexRegex := regexp.MustCompile(`^[0-9a-f]+$`)
	if !hexRegex.MatchString(hexPart) {
		return "", fmt.Errorf("invalid hex color: contains invalid hex characters")
	}

	switch len(hexPart) {
	case 3:
		// Expand hex3 to hex6: "abc" -> "aabbcc"
		return fmt.Sprintf("#%c%c%c%c%c%c", 
			hexPart[0], hexPart[0],
			hexPart[1], hexPart[1], 
			hexPart[2], hexPart[2]), nil
	case 6:
		// Already hex6 format
		return "#" + hexPart, nil
	default:
		return "", fmt.Errorf("invalid hex color: invalid length %d (expected 3 or 6)", len(hexPart))
	}
}

// ToRGBAString returns the color as rgba(r,g,b,a) format
func (c *Converter) ToRGBAString(alpha float64) string {
	if c == nil {
		return ""
	}
	
	r, g, b := c.getRGBValues()
	return fmt.Sprintf("rgba(%d,%d,%d,%.1f)", r, g, b, alpha)
}

// Brighten returns a brightened version of the color by interpolating toward white
func (c *Converter) Brighten(amount float64) string {
	if c == nil {
		return ""
	}
	
	r, g, b := c.getRGBValues()
	
	// Brighten by interpolating toward white (255,255,255)
	r = int64(float64(r) + (255-float64(r))*amount)
	g = int64(float64(g) + (255-float64(g))*amount)
	b = int64(float64(b) + (255-float64(b))*amount)
	
	// Clamp to 0-255
	if r > 255 { r = 255 }
	if g > 255 { g = 255 }
	if b > 255 { b = 255 }
	
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// Mix blends this color with another color by the specified ratio
// ratio 0.0 = this color, ratio 1.0 = other color
func (c *Converter) Mix(other *Converter, ratio float64) string {
	if c == nil || other == nil {
		return ""
	}
	
	r1, g1, b1 := c.getRGBValues()
	r2, g2, b2 := other.getRGBValues()
	
	// Linear interpolation
	r := int64(float64(r1)*(1-ratio) + float64(r2)*ratio)
	g := int64(float64(g1)*(1-ratio) + float64(g2)*ratio)
	b := int64(float64(b1)*(1-ratio) + float64(b2)*ratio)
	
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// getRGBValues extracts RGB values from the internal hex6 format
func (c *Converter) getRGBValues() (int64, int64, int64) {
	hex := strings.TrimPrefix(c.value, "#")
	
	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	
	return r, g, b
}