package core

import "testing"

func TestParseHexShort(t *testing.T) {
    c, err := ParseColor("#abc")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 0xaa || c.G != 0xbb || c.B != 0xcc || c.A != 0xff {
        t.Fatalf("got %#v", c)
    }
}

func TestParseHexLong(t *testing.T) {
    c, err := ParseColor("#112233")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 0x11 || c.G != 0x22 || c.B != 0x33 || c.A != 0xff {
        t.Fatalf("got %#v", c)
    }
}

func TestParseHex8(t *testing.T) {
    c, err := ParseColor("#11223344")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 0x11 || c.G != 0x22 || c.B != 0x33 || c.A != 0x44 {
        t.Fatalf("got %#v", c)
    }
}

func TestParseHyprRGB(t *testing.T) {
    c, err := ParseColor("rgb(112233)")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 0x11 || c.G != 0x22 || c.B != 0x33 || c.A != 0xff {
        t.Fatalf("got %#v", c)
    }
}

func TestParseHyprRGBA(t *testing.T) {
    c, err := ParseColor("rgba(11223344)")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 0x11 || c.G != 0x22 || c.B != 0x33 || c.A != 0x44 {
        t.Fatalf("got %#v", c)
    }
}

func TestParseNumericRGBA(t *testing.T) {
    c, err := ParseColor("rgba(255,0,127,0.5)")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 255 || c.G != 0 || c.B != 127 || c.A != 128 {
        t.Fatalf("got %#v", c)
    }
}

func TestParseHSLAndFormat(t *testing.T) {
    c, err := ParseColor("hsl(0, 0%, 100%)")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if c.R != 255 || c.G != 255 || c.B != 255 { t.Fatalf("got %#v", c) }
    if got := c.ToHex(); got != "#ffffff" { t.Fatalf("hex got %s", got) }
    if got, _ := c.ToFormat(FormatHyprRGB); got != "rgb(ffffff)" { t.Fatalf("hypr rgb got %s", got) }
}

func TestRoundTripHSL(t *testing.T) {
    // A color with mid values
    c := Color{R: 10, G: 200, B: 100, A: 255}
    hsl, _ := c.ToFormat(FormatHSL)
    // parse again through ToFormat->ParseColor path
    c2, err := ParseColor(hsl)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    // Colorspace conversions may differ slightly due to rounding
    const maxDelta = 2
    if diffAbs(int(c.R)-int(c2.R)) > maxDelta || diffAbs(int(c.G)-int(c2.G)) > maxDelta || diffAbs(int(c.B)-int(c2.B)) > maxDelta {
        t.Fatalf("round-trip mismatch: %#v vs %#v (via %s)", c, c2, hsl)
    }
}

func diffAbs(v int) int { if v < 0 { return -v }; return v }


