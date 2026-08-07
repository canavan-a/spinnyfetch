package logo

import (
	"math"
	"testing"
)

// TestFrameCellsAtZero checks that at theta=0 the rendered silhouette
// (which cells are filled vs. blank) matches the source art exactly. Glyphs
// themselves are expected to differ: FrameCells always renders through the
// Braille dot grid now, even for unrotated content, so the frame is a
// consistent glyph density throughout the spin instead of mixing chunky
// source glyphs with Braille only at rotated edges.
func TestFrameCellsAtZero(t *testing.T) {
	g := Parse(Source)
	text := PlainText(FrameCells(g, 0))
	for y := range text {
		srcRow := []rune(Source[y])
		outRow := []rune(text[y])
		for x, sc := range srcRow {
			wantFilled := sc != ' '
			gotFilled := x < len(outRow) && outRow[x] != ' '
			if wantFilled != gotFilled {
				t.Fatalf("row %d col %d at theta=0: filled=%v, want %v", y, x, gotFilled, wantFilled)
			}
		}
	}
}

// countNonSpace returns how many cells in a frame carry a non-space glyph.
func countNonSpace(cells [][]Cell) int {
	n := 0
	for _, row := range cells {
		for _, c := range row {
			if c.Ch != ' ' {
				n++
			}
		}
	}
	return n
}

// TestFrameCellsPreservesContentAcrossRotation checks that rotating (an
// in-plane spin, not a 3D flip) doesn't collapse or blow up the amount of
// visible art at any angle — a regression guard for the old edge-on/mirror
// behavior this replaced.
func TestFrameCellsPreservesContentAcrossRotation(t *testing.T) {
	g := Parse(Source)
	base := countNonSpace(FrameCells(g, 0))
	if base == 0 {
		t.Fatal("theta=0 frame has no content")
	}
	for _, theta := range []float64{0.3, math.Pi / 4, math.Pi / 2, 2, math.Pi, 3 * math.Pi / 2} {
		n := countNonSpace(FrameCells(g, theta))
		if n < base/3 || n > base*3 {
			t.Fatalf("theta=%v non-space count = %d, too far from baseline %d (rotation shouldn't collapse/inflate content like the old flip transform did)", theta, n, base)
		}
	}
}

// TestFrameCellsAt180DegreesIsPointReflection checks that a half turn maps
// each source cell to its point-reflection through the grid center, which
// is what an in-plane 180-degree rotation should do (as opposed to a
// left-right mirror, which is what the old 3D-flip transform produced).
func TestFrameCellsAt180DegreesIsPointReflection(t *testing.T) {
	g := Parse(Source)
	rotated := FrameCells(g, math.Pi)
	w, h := g.Width, g.Height
	matches, total := 0, 0
	for y, row := range g.Rows {
		for x, c := range row {
			if c.Ch == ' ' {
				continue
			}
			total++
			rx, ry := w-1-x, h-1-y
			if rotated[ry][rx].Ch != ' ' {
				matches++
			}
		}
	}
	if float64(matches)/float64(total) < 0.8 {
		t.Fatalf("theta=180deg: only %d/%d source cells landed at their point-reflection (nearest-neighbor rounding allows some slack, but this is too low)", matches, total)
	}
}

func TestFrameCellsSizeStaysConstant(t *testing.T) {
	g := Parse(Source)
	for _, theta := range []float64{0, 0.3, math.Pi / 4, math.Pi / 2, 2, math.Pi} {
		cells := FrameCells(g, theta)
		if len(cells) != g.Height {
			t.Fatalf("theta=%v height = %d, want %d", theta, len(cells), g.Height)
		}
		for y, row := range cells {
			if len(row) != g.Width {
				t.Fatalf("theta=%v row %d width = %d, want %d", theta, y, len(row), g.Width)
			}
		}
	}
}
