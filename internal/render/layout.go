package render

import (
	"fmt"
	"strings"

	"github.com/canavan-a/spinnyfetch/internal/logo"
	"github.com/canavan-a/spinnyfetch/internal/modules"
)

// Options controls static composition of a logo frame + info panel.
type Options struct {
	Gap       int
	KeyColor  string
	Separator string
	ColorMode ColorMode
	Spin      SpinOptions
}

// Compose renders logoCells beside the info lines, vertically centering the
// (usually shorter) info block against the logo's height. Returns the full
// multi-line frame as a single string, without a trailing newline.
func Compose(logoCells [][]logo.Cell, lines []modules.Line, opts Options) string {
	logoWidth := 0
	if len(logoCells) > 0 {
		logoWidth = len(logoCells[0])
	}
	logoHeight := len(logoCells)
	infoHeight := len(lines)

	startRow := 0
	if logoHeight > infoHeight {
		startRow = (logoHeight - infoHeight) / 2
	}

	var b strings.Builder
	totalRows := logoHeight
	if infoHeight > totalRows {
		totalRows = infoHeight
	}

	for y := 0; y < totalRows; y++ {
		if y < logoHeight {
			b.WriteString(renderLogoRow(logoCells[y], opts.ColorMode))
		} else {
			b.WriteString(strings.Repeat(" ", logoWidth))
		}
		b.WriteString(strings.Repeat(" ", opts.Gap))

		infoIdx := y - startRow
		if infoIdx >= 0 && infoIdx < infoHeight {
			b.WriteString(renderInfoLine(lines[infoIdx], opts))
		}
		if y < totalRows-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func renderLogoRow(row []logo.Cell, mode ColorMode) string {
	var b strings.Builder
	for _, c := range row {
		if c.Color == "" {
			b.WriteRune(c.Ch)
			continue
		}
		b.WriteString(ansiColor(mode, c.Color, c.Ch))
	}
	return b.String()
}

func renderInfoLine(l modules.Line, opts Options) string {
	sep := opts.Separator
	if sep == "" {
		sep = ":"
	}
	key := l.Key
	if opts.KeyColor != "" {
		key = ansiColorString(opts.ColorMode, opts.KeyColor, l.Key)
	}
	return fmt.Sprintf("%s%s %s", key, sep, l.Value)
}

func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
}
