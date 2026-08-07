package render

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/canavan-a/spinnyfetch/internal/logo"
	"github.com/canavan-a/spinnyfetch/internal/modules"
)

// SpinOptions controls the logo animation loop.
type SpinOptions struct {
	Enabled         bool
	SpeedDegPerTick float64
	FPS             int
}

const hideCursor = "\x1b[?25l"
const showCursor = "\x1b[?25h"

// Run renders grid + lines once (static frame) if spin is disabled, or
// animates the logo continuously until interrupted (SIGINT/SIGTERM).
// The terminal cursor is always restored on exit, however the loop ends.
func Run(grid logo.Grid, lines []modules.Line, opts Options) {
	if !opts.Spin.Enabled {
		cells := logo.FrameCells(grid, 0)
		fmt.Println(Compose(cells, lines, opts))
		return
	}

	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)

	// SIGWINCH (terminal resize) doesn't change our fixed art width, but if
	// the new terminal width is narrower than the art, lines may wrap and
	// throw off the cursor-up math. Rather than track wrapping, just drop
	// the "redraw in place" optimization for the next frame and repaint
	// fresh below wherever the cursor happens to be.
	resizeCh := make(chan os.Signal, 1)
	signal.Notify(resizeCh, syscall.SIGWINCH)

	fps := opts.Spin.FPS
	if fps <= 0 {
		fps = 30
	}
	interval := time.Second / time.Duration(fps)
	speed := opts.Spin.SpeedDegPerTick * math.Pi / 180
	if speed == 0 {
		speed = 6 * math.Pi / 180
	}

	theta := 0.0
	frameHeight := len(grid.Rows)
	if len(lines) > frameHeight {
		frameHeight = len(lines)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	first := true
	for {
		select {
		case <-exitCh:
			return
		case <-resizeCh:
			first = true
		case <-ticker.C:
			cells := logo.FrameCells(grid, theta)
			frame := Compose(cells, lines, opts)
			if !first {
				fmt.Printf("\x1b[%dA", frameHeight)
			}
			first = false
			fmt.Println(frame)
			theta += speed
			if theta > 2*math.Pi {
				theta -= 2 * math.Pi
			}
		}
	}
}

// StaticFrame renders a single non-animated frame as a string (used by
// --once and by tests).
func StaticFrame(grid logo.Grid, lines []modules.Line, opts Options) string {
	cells := logo.FrameCells(grid, 0)
	return Compose(cells, lines, opts)
}
