package core

import (
    "fmt"
    "math"
    "regexp"
    "strconv"
    "strings"

    colorful "github.com/lucasb-eyer/go-colorful"
)

// Color represents an sRGBA color in 8-bit channels.
// Alpha is optional in many formats; when absent we treat it as 255 (opaque).
type Color struct {
    R uint8
    G uint8
    B uint8
    A uint8
}

// Common output format identifiers used by ToFormat.
const (
    FormatHex      = "hex"       // #RRGGBB
    FormatHex8     = "hex8"      // #RRGGBBAA
    FormatRGB      = "rgb"       // rgb(r, g, b)
    FormatRGBA     = "rgba"      // rgba(r, g, b, aFloat)
    FormatHSL      = "hsl"       // hsl(h, s%, l%)
    FormatHSLA     = "hsla"      // hsla(h, s%, l%, aFloat)
    FormatHyprRGB  = "hypr_rgb"  // rgb(RRGGBB) (hex inside, no #)
    FormatHyprRGBA = "hypr_rgba" // rgba(RRGGBBAA) (hex inside, no #)
)

// IsValid returns true if the color channels are in range.
func (c Color) IsValid() bool {
    // uint8 fields are always 0..255, but keep for symmetry/future changes.
    return true
}

// ToHex returns the color in #RRGGBB form (alpha discarded).
func (c Color) ToHex() string {
    return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// ToHex8 returns the color in #RRGGBBAA form (alpha included).
func (c Color) ToHex8() string {
    return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

// ToFormat converts the color into one of the supported textual formats.
func (c Color) ToFormat(format string) (string, error) {
    switch strings.ToLower(format) {
    case FormatHex:
        return c.ToHex(), nil
    case FormatHex8:
        return c.ToHex8(), nil
    case FormatRGB:
        return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B), nil
    case FormatRGBA:
        a := float64(c.A) / 255.0
        return fmt.Sprintf("rgba(%d, %d, %d, %.3f)", c.R, c.G, c.B, trimFloat(a)), nil
    case FormatHSL:
        h, s, l := c.toHSL()
        return fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l), nil
    case FormatHSLA:
        h, s, l := c.toHSL()
        a := float64(c.A) / 255.0
        return fmt.Sprintf("hsla(%d, %d%%, %d%%, %.3f)", h, s, l, trimFloat(a)), nil
    case FormatHyprRGB:
        return fmt.Sprintf("rgb(%02x%02x%02x)", c.R, c.G, c.B), nil
    case FormatHyprRGBA:
        return fmt.Sprintf("rgba(%02x%02x%02x%02x)", c.R, c.G, c.B, c.A), nil
    default:
        return "", fmt.Errorf("unknown color format: %s", format)
    }
}

// ParseColor parses a color string in one of the supported formats.
// Supported:
// - #RGB, #RRGGBB, #RRGGBBAA
// - rgb(RRGGBB), rgba(RRGGBBAA)                  (Hyprland hex-in-function form)
// - rgb(r, g, b), rgba(r, g, b, aFloat)
// - hsl(h, s%, l%), hsla(h, s%, l%, aFloat)
func ParseColor(input string) (Color, error) {
    s := strings.TrimSpace(strings.ToLower(input))

    // Hex forms
    if strings.HasPrefix(s, "#") {
        return parseHex(s)
    }

    // hyprland-like hex-in-function: rgb(112233) or rgba(11223344) or with 0x prefix
    if m := hyprFnHexRegex.FindStringSubmatch(s); len(m) == 3 {
        fn := m[1]
        hex := strings.TrimPrefix(m[2], "0x")
        if len(hex) == 6 {
            col, err := parseHex("#" + hex)
            if err != nil {
                return Color{}, err
            }
            // If rgba() but only 6 hex, assume alpha=FF
            if fn == "rgba" {
                col.A = 255
            }
            return col, nil
        } else if len(hex) == 8 {
            // In hypr land, rgba is RRGGBBAA order
            col, err := parseHex("#" + hex)
            if err != nil {
                return Color{}, err
            }
            return col, nil
        }
        return Color{}, fmt.Errorf("invalid hypr hex length: %s", input)
    }

    // numeric rgb/rgba
    if m := rgbRegex.FindStringSubmatch(s); len(m) == 4 {
        r, _ := strconv.Atoi(m[1])
        g, _ := strconv.Atoi(m[2])
        b, _ := strconv.Atoi(m[3])
        return Color{R: clamp8(r), G: clamp8(g), B: clamp8(b), A: 255}, nil
    }
    if m := rgbaRegex.FindStringSubmatch(s); len(m) == 5 {
        r, _ := strconv.Atoi(m[1])
        g, _ := strconv.Atoi(m[2])
        b, _ := strconv.Atoi(m[3])
        aFloat, _ := strconv.ParseFloat(m[4], 64)
        return Color{R: clamp8(r), G: clamp8(g), B: clamp8(b), A: clampAlpha01(aFloat)}, nil
    }

    // hsl/hsla
    if m := hslRegex.FindStringSubmatch(s); len(m) == 4 {
        h, _ := strconv.ParseFloat(m[1], 64)
        sPerc, _ := strconv.ParseFloat(m[2], 64)
        lPerc, _ := strconv.ParseFloat(m[3], 64)
        return fromHSL(h, sPerc, lPerc, 1.0), nil
    }
    if m := hslaRegex.FindStringSubmatch(s); len(m) == 5 {
        h, _ := strconv.ParseFloat(m[1], 64)
        sPerc, _ := strconv.ParseFloat(m[2], 64)
        lPerc, _ := strconv.ParseFloat(m[3], 64)
        a, _ := strconv.ParseFloat(m[4], 64)
        return fromHSL(h, sPerc, lPerc, a), nil
    }

    return Color{}, fmt.Errorf("unsupported color format: %q", input)
}

// Helpers

var (
    // rgb(255, 255, 255)
    rgbRegex  = regexp.MustCompile(`^rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$`)
    // rgba(255, 255, 255, 0.5)
    rgbaRegex = regexp.MustCompile(`^rgba\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*([01]?(?:\.\d+)?)\s*\)$`)
    // hsl(180, 50%, 50%)
    hslRegex  = regexp.MustCompile(`^hsl\(\s*(-?\d+(?:\.\d+)?)\s*,\s*(\d{1,3})%\s*,\s*(\d{1,3})%\s*\)$`)
    // hsla(180, 50%, 50%, 0.5)
    hslaRegex = regexp.MustCompile(`^hsla\(\s*(-?\d+(?:\.\d+)?)\s*,\s*(\d{1,3})%\s*,\s*(\d{1,3})%\s*,\s*([01]?(?:\.\d+)?)\s*\)$`)
    // hypr: rgb(112233) or rgba(11223344), optionally 0x prefix
    hyprFnHexRegex = regexp.MustCompile(`^(rgb|rgba)\(\s*(?:0x)?([0-9a-f]{6,8})\s*\)$`)
)

func parseHex(s string) (Color, error) {
    hex := strings.TrimPrefix(s, "#")
    switch len(hex) {
    case 3:
        r := strings.Repeat(string(hex[0]), 2)
        g := strings.Repeat(string(hex[1]), 2)
        b := strings.Repeat(string(hex[2]), 2)
        return parseHex("#" + r + g + b)
    case 6:
        r, err := strconv.ParseUint(hex[0:2], 16, 8)
        if err != nil { return Color{}, err }
        g, err := strconv.ParseUint(hex[2:4], 16, 8)
        if err != nil { return Color{}, err }
        b, err := strconv.ParseUint(hex[4:6], 16, 8)
        if err != nil { return Color{}, err }
        return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
    case 8:
        r, err := strconv.ParseUint(hex[0:2], 16, 8)
        if err != nil { return Color{}, err }
        g, err := strconv.ParseUint(hex[2:4], 16, 8)
        if err != nil { return Color{}, err }
        b, err := strconv.ParseUint(hex[4:6], 16, 8)
        if err != nil { return Color{}, err }
        a, err := strconv.ParseUint(hex[6:8], 16, 8)
        if err != nil { return Color{}, err }
        return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
    default:
        return Color{}, fmt.Errorf("invalid hex length: %s", s)
    }
}

func (c Color) toHSL() (h int, s int, l int) {
    col := colorful.Color{R: float64(c.R) / 255.0, G: float64(c.G) / 255.0, B: float64(c.B) / 255.0}
    hh, ss, ll := col.Hsl()
    // Normalize ranges: h 0..360, s/l 0..100
    if hh < 0 {
        hh = math.Mod(hh, 360) + 360
    }
    return int(math.Round(hh)), int(math.Round(ss * 100)), int(math.Round(ll * 100))
}

func fromHSL(h float64, sPerc float64, lPerc float64, a float64) Color {
    // clamp inputs
    for h < 0 {
        h += 360
    }
    h = math.Mod(h, 360)
    if sPerc < 0 { sPerc = 0 } else if sPerc > 100 { sPerc = 100 }
    if lPerc < 0 { lPerc = 0 } else if lPerc > 100 { lPerc = 100 }
    if a < 0 { a = 0 } else if a > 1 { a = 1 }
    col := colorful.Hsl(h, sPerc/100.0, lPerc/100.0)
    r := uint8(math.Round(col.R * 255))
    g := uint8(math.Round(col.G * 255))
    b := uint8(math.Round(col.B * 255))
    return Color{R: r, G: g, B: b, A: uint8(math.Round(a * 255))}
}

func clamp8(v int) uint8 {
    if v < 0 { v = 0 }
    if v > 255 { v = 255 }
    return uint8(v)
}

func clampAlpha01(a float64) uint8 {
    if a < 0 { a = 0 }
    if a > 1 { a = 1 }
    return uint8(math.Round(a * 255))
}

func trimFloat(f float64) float64 {
    // Round to 3 decimals for stable output
    return math.Round(f*1000) / 1000
}

