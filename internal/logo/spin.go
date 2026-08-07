package logo

import "math"

// Cell is one rendered character of a logo frame, with its source color tag.
type Cell struct {
	Ch    rune
	Color string // ColorA, ColorB, or "" for empty
}

// Grid is the parsed source art: one []Cell per row, all rows padded to Width.
type Grid struct {
	Rows   [][]Cell
	Width  int
	Height int
}

// Parse converts Source into a Grid, tagging each non-space glyph with a
// color based on which horizontal half of the art it falls in. Uses the
// default ColorA/ColorB two-tone; see ParseWithColors to override them.
func Parse(lines []string) Grid {
	return ParseWithColors(lines, ColorA, ColorB)
}

// ParseWithColors is Parse with the two-tone colors overridden (e.g. from
// user config).
func ParseWithColors(lines []string, colorA, colorB string) Grid {
	width := 0
	runeRows := make([][]rune, len(lines))
	for i, l := range lines {
		r := []rune(l)
		runeRows[i] = r
		if len(r) > width {
			width = len(r)
		}
	}

	rows := make([][]Cell, len(lines))
	mid := width / 2
	for i, r := range runeRows {
		row := make([]Cell, width)
		for x := 0; x < width; x++ {
			var ch rune = ' '
			if x < len(r) {
				ch = r[x]
			}
			color := ""
			if ch != ' ' {
				if x < mid {
					color = colorA
				} else {
					color = colorB
				}
			}
			row[x] = Cell{Ch: ch, Color: color}
		}
		rows[i] = row
	}

	return Grid{Rows: rows, Width: width, Height: len(lines)}
}

// CellAspect is the assumed height/width ratio of a terminal character
// cell (typical monospace glyphs render roughly twice as tall as wide).
// Rotation coordinates are scaled by this so a "circle" of art actually
// looks circular instead of elliptical when spun.
const CellAspect = 2.0

// brailleBase is the start of the Unicode Braille Patterns block. Each of
// the 256 code points in the block is a bitmask over 8 dots arranged in a
// 2-column x 4-row grid within the cell, so adding a dot-coverage mask
// directly yields the matching glyph — no lookup table needed.
const brailleBase = 0x2800

// dotOffsets gives each Braille dot's (x, y) position within a cell, in
// cell units, and its bit within the standard Braille dot-numbering
// (1 2 3 7 down the left column, 4 5 6 8 down the right).
var dotOffsets = [8]struct {
	dx, dy float64
	bit    uint8
}{
	{-0.25, -0.375, 0}, // dot 1
	{-0.25, -0.125, 1}, // dot 2
	{-0.25, 0.125, 2},  // dot 3
	{0.25, -0.375, 3},  // dot 4
	{0.25, -0.125, 4},  // dot 5
	{0.25, 0.125, 5},   // dot 6
	{-0.25, 0.375, 6},  // dot 7
	{0.25, 0.375, 7},   // dot 8
}

// FrameCells renders the grid rotated in-plane by theta (radians) about its
// own center, like a flat gear/wheel spinning face-on — not a 3D flip.
// Each destination character cell is split into an 8-dot Braille sub-grid
// (2 columns x 4 rows); each dot is inverse-sampled independently from the
// source art (correcting for the non-square aspect ratio of terminal cells
// so the rotation reads as circular) and the resulting dot pattern is
// rendered with the matching Braille glyph. This gives much finer
// partial-coverage edges than a plain per-cell or even quadrant-block
// sample, since Braille packs 8 independent sub-pixels per character
// instead of 1 or 4. Each cell carries its render color so callers can
// emit ANSI-colored output directly.
func FrameCells(g Grid, theta float64) [][]Cell {
	cx := float64(g.Width-1) / 2
	cy := float64(g.Height-1) / 2

	// Inverse rotation: to fill destination (dx,dy), sample from the source
	// position that rotates *forward* by theta into (dx,dy).
	cos, sin := math.Cos(theta), math.Sin(theta)

	out := make([][]Cell, g.Height)
	for dy := 0; dy < g.Height; dy++ {
		line := make([]Cell, g.Width)
		for dx := 0; dx < g.Width; dx++ {
			line[dx] = Cell{Ch: ' '}

			var mask uint8
			counts := make(map[Cell]int, 8)
			for _, d := range dotOffsets {
				px := float64(dx) + d.dx - cx
				py := (float64(dy) + d.dy - cy) * CellAspect
				sx := px*cos + py*sin
				sy := -px*sin + py*cos

				srcX := int(math.Round(sx + cx))
				srcY := int(math.Round(sy/CellAspect + cy))
				if srcX < 0 || srcX >= g.Width || srcY < 0 || srcY >= g.Height {
					continue
				}
				cell := g.Rows[srcY][srcX]
				if cell.Ch == ' ' {
					continue
				}
				mask |= 1 << d.bit
				counts[cell]++
			}

			if mask != 0 {
				best := Cell{}
				bestN := 0
				for c, n := range counts {
					if n > bestN {
						best, bestN = c, n
					}
				}
				// Always render as Braille, even in the interior where all
				// dots hit the same source cell: mixing chunky source
				// glyphs with Braille edges made the frame look inconsistent
				// and jittery as it rotated, since two different glyph
				// densities were fighting for the same silhouette.
				line[dx] = Cell{Ch: rune(brailleBase + int(mask)), Color: best.Color}
			}
		}
		out[dy] = line
	}
	return out
}

// PlainText strips color from a FrameCells result, for tests and
// non-color terminal output.
func PlainText(cells [][]Cell) []string {
	out := make([]string, len(cells))
	for y, row := range cells {
		r := make([]rune, len(row))
		for x, c := range row {
			r[x] = c.Ch
		}
		out[y] = string(r)
	}
	return out
}
