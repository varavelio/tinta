package tinta

import (
	"strings"
	"sync"
	"testing"

	"github.com/varavelio/tinta/internal/assert"
)

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
		assert.Equal(t, "ABBA", got)
	})

	t.Run("lower z is painted first higher z overwrites", func(t *testing.T) {
		got := Canvas().
			AddZ("XXXX", 0, 0, 10).
			AddZ("YY", 1, 0, 5).
			String()
		assert.Equal(t, "XXXX", got)
	})

	t.Run("same z uses insertion order", func(t *testing.T) {
		got := Canvas().
			AddZ("AAAA", 0, 0, 0).
			AddZ("BB", 1, 0, 0).
			String()
		assert.Equal(t, "ABBA", got)
	})
}

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
		got := Canvas().
			AddZ("***", 3, 0, 1).
			AddZ(strings.Repeat(".", 10), 0, 0, 0).
			String()
		assert.Equal(t, "...***....", got)
	})
}

func TestCanvasClipping(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("negative x expands canvas left", func(t *testing.T) {
		got := Canvas().Add("ABCDE", -2, 0).String()
		assert.Equal(t, "ABCDE", got)
	})

	t.Run("negative y expands canvas up", func(t *testing.T) {
		got := Canvas().Add("row0\nrow1\nrow2", 0, -1).String()
		lines := strings.Split(got, "\n")
		assert.Equal(t, 3, len(lines))
		assert.Equal(t, "row0", lines[0])
		assert.Equal(t, "row1", lines[1])
		assert.Equal(t, "row2", lines[2])
	})

	t.Run("completely negative still shows everything", func(t *testing.T) {
		got := Canvas().Add("hello", -10, -10).String()
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
		got := Canvas().Width(3).Add("ABCDE", -2, 0).String()
		assert.Equal(t, "ABC", got)
	})

	t.Run("negative y with fixed height clips", func(t *testing.T) {
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
		assert.Equal(t, "BB AA", got)
	})
}

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
		assert.Equal(t, "", lines[1])
		assert.Equal(t, "", lines[2])
	})
}

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

func TestCanvas3DBorderEffect(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("nested box with offsets creates 3D effect", func(t *testing.T) {
		text := "Lorem ipsum"
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
			Add(shadow, 1, 1).
			Add(front, 0, 0).
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

func TestCanvasANSI(t *testing.T) {
	ForceColors(true)

	t.Run("styled layer preserves ANSI codes", func(t *testing.T) {
		styled := Text().Red().String("hi")
		got := Canvas().Add(styled, 0, 0).String()
		assert.Equal(t, true, strings.Contains(got, "\x1b[31m"))
		assert.Equal(t, true, strings.Contains(got, "hi"))
	})

	t.Run("styled box composites with ANSI", func(t *testing.T) {
		box := Box().Red().String("x")
		got := Canvas().Add(box, 0, 0).String()
		assert.Equal(t, true, strings.Contains(got, "\x1b[31m"))
	})

	t.Run("two styled layers compose correctly", func(t *testing.T) {
		red := Text().Red().String("RR")
		blue := Text().Blue().String("BB")
		got := Canvas().
			Add(red, 0, 0).
			Add(blue, 4, 0).
			String()
		assert.Equal(t, true, strings.Contains(got, "\x1b[31m"))
		assert.Equal(t, true, strings.Contains(got, "\x1b[34m"))
	})
}

func TestCanvasAutoZ(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("auto z increments so later Add is on top", func(t *testing.T) {
		got := Canvas().
			Add("AAAA", 0, 0).
			Add("BB", 1, 0).
			String()
		assert.Equal(t, "ABBA", got)
	})

	t.Run("AddZ does not break auto-z", func(t *testing.T) {
		got := Canvas().
			AddZ("AAAA", 0, 0, 0).
			Add("BB", 1, 0).
			String()
		assert.Equal(t, "ABBA", got)
	})

	t.Run("AddZ with high z bumps auto counter", func(t *testing.T) {
		got := Canvas().
			AddZ("AAAA", 0, 0, 100).
			Add("BB", 1, 0).
			String()
		assert.Equal(t, "ABBA", got)
	})
}

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
		assert.Equal(t, "", lines[1])
		assert.Equal(t, "bot", lines[2])
	})

	t.Run("layer with trailing newline handled correctly", func(t *testing.T) {
		got := Canvas().Add("abc\n", 0, 0).String()
		assert.Equal(t, "abc", got)
	})
}

func TestCanvasMultipleDepths(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("three layers create cascading depth", func(t *testing.T) {
		box := Box().Border(BorderSimple).String("AB")
		got := Canvas().
			Add(box, 4, 4).
			Add(box, 2, 2).
			Add(box, 0, 0).
			String()

		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.HasPrefix(lines[0], "┌──┐"))
		assert.Equal(t, true, strings.HasPrefix(lines[1], "│AB│"))
		assert.Equal(t, true, strings.HasPrefix(lines[2], "└──┘"))

		assert.Equal(t, 7, len(lines))
	})
}
