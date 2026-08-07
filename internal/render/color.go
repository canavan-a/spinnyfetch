package render

import "fmt"

// ColorMode selects how hex colors are emitted as ANSI escapes.
type ColorMode int

const (
	Truecolor ColorMode = iota
	Color256
)

func ParseColorMode(s string) ColorMode {
	if s == "256" {
		return Color256
	}
	return Truecolor
}

// ansiColor wraps ch in the ANSI escape for hex under the given mode.
func ansiColor(mode ColorMode, hex string, ch rune) string {
	r, g, b := hexToRGB(hex)
	if mode == Color256 {
		code := rgbTo256(r, g, b)
		return fmt.Sprintf("\x1b[38;5;%dm%c\x1b[0m", code, ch)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, ch)
}

func ansiColorString(mode ColorMode, hex string, s string) string {
	r, g, b := hexToRGB(hex)
	if mode == Color256 {
		code := rgbTo256(r, g, b)
		return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", code, s)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, s)
}

// rgbTo256 maps 24-bit RGB to the xterm 256-color 6x6x6 cube (codes 16-231).
func rgbTo256(r, g, b int) int {
	toIdx := func(v int) int {
		if v < 48 {
			return 0
		}
		if v < 115 {
			return 1
		}
		return (v - 35) / 40
	}
	ri, gi, bi := toIdx(r), toIdx(g), toIdx(b)
	return 16 + 36*ri + 6*gi + bi
}
