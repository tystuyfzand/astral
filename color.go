package astral

import (
	"fmt"
	"strconv"
	"strings"
)

type Color int32

// NewColor creates a Color from RGB values.
func NewColor(r, g, b uint8) Color {
	return Color(int32(r)<<16 | int32(g)<<8 | int32(b))
}

// R returns the red component.
func (c Color) R() uint8 {
	return uint8((c >> 16) & 0xFF)
}

// G returns the green component.
func (c Color) G() uint8 {
	return uint8((c >> 8) & 0xFF)
}

// B returns the blue component.
func (c Color) B() uint8 {
	return uint8(c & 0xFF)
}

// RGB returns all components.
func (c Color) RGB() (uint8, uint8, uint8) {
	return c.R(), c.G(), c.B()
}

// Hex returns the color as a hex string (#RRGGBB).
func (c Color) Hex() string {
	return fmt.Sprintf("#%06X", int32(c)&0xFFFFFF)
}

// String implements fmt.Stringer and returns the hex representation.
func (c Color) String() string {
	return c.Hex()
}

// ParseHex parses a hex color string.
// Supports: "#RRGGBB", "RRGGBB", "#RGB", "RGB".
func ParseHex(s string) (Color, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")

	switch len(s) {
	case 3: // RGB shorthand
		r, err := strconv.ParseUint(strings.Repeat(string(s[0]), 2), 16, 8)
		if err != nil {
			return 0, err
		}
		g, err := strconv.ParseUint(strings.Repeat(string(s[1]), 2), 16, 8)
		if err != nil {
			return 0, err
		}
		b, err := strconv.ParseUint(strings.Repeat(string(s[2]), 2), 16, 8)
		if err != nil {
			return 0, err
		}
		return NewColor(uint8(r), uint8(g), uint8(b)), nil

	case 6: // RRGGBB
		v, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return 0, err
		}
		return Color(v), nil
	}

	return 0, fmt.Errorf("invalid hex color: %q", s)
}
