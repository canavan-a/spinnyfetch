// Throwaway visual check for the spin transform: run with `go run ./cmd/spinpreview`.
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/canavan-a/spinnyfetch/internal/logo"
)

func colorize(ch rune, color string) string {
	if color == "" {
		return string(ch)
	}
	r, g, b := hexToRGB(color)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, ch)
}

func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func main() {
	grid := logo.Parse(logo.Source)
	theta := 0.0
	speed := 0.12 // radians per tick

	fmt.Print("\x1b[?25l")
	defer fmt.Print("\x1b[?25h")

	for i := 0; i < 400; i++ {
		cells := logo.FrameCells(grid, theta)
		var out string
		for _, row := range cells {
			for _, c := range row {
				out += colorize(c.Ch, c.Color)
			}
			out += "\n"
		}
		fmt.Print(out)
		time.Sleep(33 * time.Millisecond)
		fmt.Printf("\x1b[%dA", grid.Height)
		theta += speed
		if theta > 2*math.Pi {
			theta -= 2 * math.Pi
		}
	}
	fmt.Printf("\x1b[%dB", grid.Height)
}
