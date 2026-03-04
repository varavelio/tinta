package tinta

import (
	"strings"
	"sync"
	"testing"

	"github.com/varavelio/tinta/internal/assert"
)

// --- Canvas basics ---

func TestCanvasEmpty(t *testing.T) {
	t.Run("empty canvas returns empty string", func(t *testing.T) {
		got := Canvas().String()
		assert.Equal(t, "", got)
	})
}

func TestCanvasSingleLayer(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("single layer at origin", func(t *testing.T) {
		box := Box().String("hi")
		got := Canvas().Add(box, 0, 0).String()
		// The box renders as:
		// ┌──┐
		// │hi│
		// └──┘
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "│hi│", lines[1])
		assert.Equal(t, "└──┘", lines[2])
		assert.Equal(t, 3, len(lines))
	})

	t.Run("single layer with offset", func(t *testing.T) {
		box := Box().String("x")
		got := Canvas().Add(box, 2, 1).String()
		lines := strings.Split(got, "\n")
		// Row 0 is blank (only spaces up to the box extent)
		// Row 1: 2 spaces + box top border
		assert.Equal(t, "  ┌─┐", lines[1])
		assert.Equal(t, "  │x│", lines[2])
		assert.Equal(t, "  └─┘", lines[3])
	})

	t.Run("plain text layer", func(t *testing.T) {
		got := Canvas().Add("hello", 0, 0).String()
		assert.Equal(t, "hello", got)
	})

	t.Run("multiline text layer", func(t *testing.T) {
		got := Canvas().Add("ab\ncd", 0, 0).String()
		lines := strings.Split(got, "\n")
		assert.Equal(t, "ab", lines[0])
		assert.Equal(t, "cd", lines[1])
	})
}

// --- Multi-layer compositing ---

func TestCanvasMultiLayer(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("two layers no overlap", func(t *testing.T) {
		got := Canvas().
			Add("AB", 0, 0).
			Add("CD", 5, 0).
			String()
		assert.Equal(t, "AB   CD", got)
	})

	t.Run("two layers overlapping, later wins", func(t *testing.T) {
		got := Canvas().
			Add("AAAA", 0, 0).
			Add("BB", 1, 0).
			String()
		// "AAAA" at x=0, then "BB" at x=1 overwrites positions 1 and 2
		assert.Equal(t, "ABBA", got)
	})

	t.Run("lower z is painted first higher z overwrites", func(t *testing.T) {
		got := Canvas().
			AddZ("XXXX", 0, 0, 10).
			AddZ("YY", 1, 0, 5).
			String()
		// z=5 painted first ("YY" at x=1), then z=10 paints "XXXX" at x=0
		// z=10 overwrites everything
		assert.Equal(t, "XXXX", got)
	})

	t.Run("same z uses insertion order", func(t *testing.T) {
		got := Canvas().
			AddZ("AAAA", 0, 0, 0).
			AddZ("BB", 1, 0, 0).
			String()
		// Same z=0, second insertion overwrites first at overlapping positions
		assert.Equal(t, "ABBA", got)
	})
}

// --- Z-ordering ---

func TestCanvasZOrder(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("high z covers low z", func(t *testing.T) {
		bg := strings.Repeat(".", 10)
		fg := "***"
		got := Canvas().
			AddZ(bg, 0, 0, 0).
			AddZ(fg, 3, 0, 1).
			String()
		assert.Equal(t, "...***....", got)
	})

	t.Run("low z does not cover high z", func(t *testing.T) {
		// Add high z first, then low z — low z cannot overwrite high z
		got := Canvas().
			AddZ("***", 3, 0, 1).
			AddZ(strings.Repeat(".", 10), 0, 0, 0).
			String()
		assert.Equal(t, "...***....", got)
	})
}

// --- Clipping ---

func TestCanvasClipping(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("negative x expands canvas left", func(t *testing.T) {
		got := Canvas().Add("ABCDE", -2, 0).String()
		// Content is NOT clipped — canvas expands to fit it.
		// The layer starts at x=-2 so the canvas shifts right by 2.
		// All 5 chars are visible at their shifted positions.
		assert.Equal(t, "ABCDE", got)
	})

	t.Run("negative y expands canvas up", func(t *testing.T) {
		got := Canvas().Add("row0\nrow1\nrow2", 0, -1).String()
		// Content is NOT clipped — canvas expands to fit all rows.
		lines := strings.Split(got, "\n")
		assert.Equal(t, 3, len(lines))
		assert.Equal(t, "row0", lines[0])
		assert.Equal(t, "row1", lines[1])
		assert.Equal(t, "row2", lines[2])
	})

	t.Run("completely negative still shows everything", func(t *testing.T) {
		got := Canvas().Add("hello", -10, -10).String()
		// Canvas expands to fit — content is visible.
		assert.Equal(t, "hello", got)
	})

	t.Run("fixed size clips right and bottom", func(t *testing.T) {
		got := Canvas().
			Width(3).Height(1).
			Add("ABCDE", 0, 0).
			String()
		assert.Equal(t, "ABC", got)
	})

	t.Run("negative x with fixed width clips", func(t *testing.T) {
		// Layer at x=-2 with 5 chars, canvas width=3.
		// After shift: chars at cols 0-4, but width=3 crops to cols 0-2.
		got := Canvas().Width(3).Add("ABCDE", -2, 0).String()
		assert.Equal(t, "ABC", got)
	})

	t.Run("negative y with fixed height clips", func(t *testing.T) {
		// Layer at y=-1 with 3 rows, canvas height=2.
		// After shift: rows at 0-2, but height=2 crops to rows 0-1.
		got := Canvas().Height(2).Add("row0\nrow1\nrow2", 0, -1).String()
		lines := strings.Split(got, "\n")
		assert.Equal(t, 2, len(lines))
		assert.Equal(t, "row0", lines[0])
		assert.Equal(t, "row1", lines[1])
	})

	t.Run("two layers one negative expand canvas to fit both", func(t *testing.T) {
		got := Canvas().
			Add("AA", 0, 0).
			Add("BB", -3, 0).
			String()
		// BB at x=-3 shifts everything by 3. BB at cols 0-1, AA at cols 3-4.
		assert.Equal(t, "BB AA", got)
	})
}

// --- Fixed dimensions ---

func TestCanvasFixedDimensions(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("width clips content", func(t *testing.T) {
		got := Canvas().Width(3).Add("ABCDE", 0, 0).String()
		assert.Equal(t, "ABC", got)
	})

	t.Run("height adds empty rows", func(t *testing.T) {
		got := Canvas().Height(3).Add("x", 0, 0).String()
		lines := strings.Split(got, "\n")
		assert.Equal(t, 3, len(lines))
		assert.Equal(t, "x", lines[0])
		// Empty rows are trimmed of trailing spaces, resulting in ""
		assert.Equal(t, "", lines[1])
		assert.Equal(t, "", lines[2])
	})
}

// --- Immutability ---

func TestCanvasImmutability(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("Add does not modify original", func(t *testing.T) {
		base := Canvas().Add("A", 0, 0)
		derived := base.Add("B", 1, 0)

		baseGot := base.String()
		derivedGot := derived.String()

		assert.Equal(t, "A", baseGot)
		assert.Equal(t, "AB", derivedGot)
	})

	t.Run("Width does not modify original", func(t *testing.T) {
		base := Canvas().Add("x", 0, 0)
		wide := base.Width(5)

		baseGot := base.String()
		wideGot := wide.String()

		// Trailing spaces are trimmed, so width only affects clipping
		assert.Equal(t, "x", baseGot)
		assert.Equal(t, "x", wideGot)
	})

	t.Run("Height does not modify original", func(t *testing.T) {
		base := Canvas().Add("x", 0, 0)
		tall := base.Height(3)

		baseLines := strings.Split(base.String(), "\n")
		tallLines := strings.Split(tall.String(), "\n")

		assert.Equal(t, 1, len(baseLines))
		assert.Equal(t, 3, len(tallLines))
	})
}

// --- Concurrent safety ---

func TestCanvasConcurrent(t *testing.T) {
	t.Run("concurrent reads from shared canvas", func(t *testing.T) {
		c := Canvas().
			Add("hello", 0, 0).
			Add("world", 0, 1)

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				got := c.String()
				if !strings.Contains(got, "hello") {
					t.Errorf("expected 'hello', got %q", got)
				}
			}()
		}
		wg.Wait()
	})
}

// --- 3D border acceptance test ---

func TestCanvas3DBorderEffect(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("nested box with offsets creates 3D effect", func(t *testing.T) {
		text := "Lorem ipsum"
		// Both boxes render independently with identical dimensions.
		front := Box().
			Border(BorderHeavy).
			PaddingX(5).
			PaddingY(1).
			String(text)
		shadow := Box().
			Border(BorderRounded).
			PaddingX(5).
			PaddingY(1).
			String(text)

		got := Canvas().
			Add(shadow, 1, 1). // shadow behind, offset right+down
			Add(front, 0, 0).  // front face on top
			String()

		expected := strings.TrimSpace(`
┏━━━━━━━━━━━━━━━━━━━━━┓
┃                     ┃╮
┃     Lorem ipsum     ┃│
┃                     ┃│
┗━━━━━━━━━━━━━━━━━━━━━┛│
 ╰─────────────────────╯
`)

		assert.Equal(t, expected, got)
	})
}

// --- ANSI compositing ---

func TestCanvasANSI(t *testing.T) {
	ForceColors(true)

	t.Run("styled layer preserves ANSI codes", func(t *testing.T) {
		styled := Text().Red().String("hi")
		got := Canvas().Add(styled, 0, 0).String()
		// Should contain the red ANSI code
		assert.Equal(t, true, strings.Contains(got, "\x1b[31m"))
		assert.Equal(t, true, strings.Contains(got, "hi"))
	})

	t.Run("styled box composites with ANSI", func(t *testing.T) {
		box := Box().Red().String("x")
		got := Canvas().Add(box, 0, 0).String()
		// Border should have red ANSI
		assert.Equal(t, true, strings.Contains(got, "\x1b[31m"))
	})

	t.Run("two styled layers compose correctly", func(t *testing.T) {
		red := Text().Red().String("RR")
		blue := Text().Blue().String("BB")
		got := Canvas().
			Add(red, 0, 0).
			Add(blue, 4, 0).
			String()
		// Both styles should be present
		assert.Equal(t, true, strings.Contains(got, "\x1b[31m"))
		assert.Equal(t, true, strings.Contains(got, "\x1b[34m"))
	})
}

// --- Auto z-increment ---

func TestCanvasAutoZ(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("auto z increments so later Add is on top", func(t *testing.T) {
		got := Canvas().
			Add("AAAA", 0, 0).
			Add("BB", 1, 0).
			String()
		// Second Add has higher auto-z, overwrites positions 1-2
		assert.Equal(t, "ABBA", got)
	})

	t.Run("AddZ does not break auto-z", func(t *testing.T) {
		got := Canvas().
			AddZ("AAAA", 0, 0, 0).
			Add("BB", 1, 0). // auto z should be 1
			String()
		assert.Equal(t, "ABBA", got)
	})

	t.Run("AddZ with high z bumps auto counter", func(t *testing.T) {
		got := Canvas().
			AddZ("AAAA", 0, 0, 100).
			Add("BB", 1, 0). // auto z should be 101
			String()
		assert.Equal(t, "ABBA", got)
	})
}

// --- Edge cases ---

func TestCanvasEdgeCases(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("empty string layer", func(t *testing.T) {
		got := Canvas().Add("", 0, 0).String()
		assert.Equal(t, "", got)
	})

	t.Run("multiple layers stacked vertically", func(t *testing.T) {
		got := Canvas().
			Add("top", 0, 0).
			Add("bot", 0, 2).
			String()
		lines := strings.Split(got, "\n")
		assert.Equal(t, "top", lines[0])
		assert.Equal(t, "", lines[1]) // gap row (trailing spaces trimmed)
		assert.Equal(t, "bot", lines[2])
	})

	t.Run("layer with trailing newline handled correctly", func(t *testing.T) {
		got := Canvas().Add("abc\n", 0, 0).String()
		// Trailing newline is stripped during parseGrid, so no extra row
		assert.Equal(t, "abc", got)
	})
}

// --- Real-world: multiple 3D depths ---

func TestCanvasMultipleDepths(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("three layers create cascading depth", func(t *testing.T) {
		box := Box().Border(BorderSimple).String("AB")
		// Three layers at different offsets, auto-z makes each successive one on top
		got := Canvas().
			Add(box, 4, 4). // deepest shadow
			Add(box, 2, 2). // middle shadow
			Add(box, 0, 0). // front
			String()

		lines := strings.Split(got, "\n")
		// Front box at (0,0)
		assert.Equal(t, true, strings.HasPrefix(lines[0], "┌──┐"))
		assert.Equal(t, true, strings.HasPrefix(lines[1], "│AB│"))
		assert.Equal(t, true, strings.HasPrefix(lines[2], "└──┘"))

		// Middle box peeks out at offset (2,2) — but its top-left area
		// is covered by the front box. Only the bottom-right edges show.
		// Deepest box peeks at (4,4).
		// Total height should be 4+3=7 rows
		assert.Equal(t, 7, len(lines))
	})
}
