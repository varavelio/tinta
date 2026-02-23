package tinta

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/varavelio/tinta/internal/assert"
)

// --- Box basics ---

func TestBoxString(t *testing.T) {
	t.Run("simple border wraps content", func(t *testing.T) {
		got := Box().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "│hi│", lines[1])
		assert.Equal(t, "└──┘", lines[2])
	})

	t.Run("multiline content", func(t *testing.T) {
		got := Box().String("ab\ncd")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "│ab│", lines[1])
		assert.Equal(t, "│cd│", lines[2])
		assert.Equal(t, "└──┘", lines[3])
	})

	t.Run("multiline uneven lines pads shorter lines", func(t *testing.T) {
		got := Box().String("hello\nhi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "│hello│", lines[1])
		assert.Equal(t, "│hi   │", lines[2])
		assert.Equal(t, "└─────┘", lines[3])
	})

	t.Run("empty content", func(t *testing.T) {
		got := Box().String("")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌┐", lines[0])
		assert.Equal(t, "││", lines[1])
		assert.Equal(t, "└┘", lines[2])
	})
}

// --- Border styles ---

func TestBoxBorderStyles(t *testing.T) {
	t.Run("rounded", func(t *testing.T) {
		got := Box().BorderRounded().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "╭─╮", lines[0])
		assert.Equal(t, "│x│", lines[1])
		assert.Equal(t, "╰─╯", lines[2])
	})

	t.Run("double", func(t *testing.T) {
		got := Box().BorderDouble().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "╔═╗", lines[0])
		assert.Equal(t, "║x║", lines[1])
		assert.Equal(t, "╚═╝", lines[2])
	})

	t.Run("heavy", func(t *testing.T) {
		got := Box().BorderHeavy().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┏━┓", lines[0])
		assert.Equal(t, "┃x┃", lines[1])
		assert.Equal(t, "┗━┛", lines[2])
	})

	t.Run("simple is default", func(t *testing.T) {
		got := Box().BorderSimple().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│x│", lines[1])
		assert.Equal(t, "└─┘", lines[2])
	})

	t.Run("custom border", func(t *testing.T) {
		custom := Border{
			TopLeft: "+", TopRight: "+", BottomLeft: "+", BottomRight: "+",
			Horizontal: "-", Vertical: "|",
		}
		got := Box().Border(custom).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "+-+", lines[0])
		assert.Equal(t, "|x|", lines[1])
		assert.Equal(t, "+-+", lines[2])
	})
}

// --- Padding ---

func TestBoxPadding(t *testing.T) {
	t.Run("uniform padding", func(t *testing.T) {
		got := Box().Padding(1).String("x")
		lines := strings.Split(got, "\n")
		// Padding 1 all sides: inner width = 1(content) + 1(left) + 1(right) = 3
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│   │", lines[1]) // top pad
		assert.Equal(t, "│ x │", lines[2]) // content with left/right pad
		assert.Equal(t, "│   │", lines[3]) // bottom pad
		assert.Equal(t, "└───┘", lines[4])
	})

	t.Run("horizontal padding only", func(t *testing.T) {
		got := Box().PaddingX(2).String("x")
		lines := strings.Split(got, "\n")
		// innerW = 1 + 2 + 2 = 5
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "│  x  │", lines[1])
		assert.Equal(t, "└─────┘", lines[2])
	})

	t.Run("individual side padding", func(t *testing.T) {
		got := Box().PaddingTop(1).PaddingRight(2).PaddingBottom(0).PaddingLeft(3).String("x")
		lines := strings.Split(got, "\n")
		// innerW = 1 + 3(left) + 2(right) = 6
		assert.Equal(t, "┌──────┐", lines[0])
		assert.Equal(t, "│      │", lines[1]) // top pad
		assert.Equal(t, "│   x  │", lines[2]) // left=3, content=1, right-fill=2
		assert.Equal(t, "└──────┘", lines[3]) // no bottom pad
	})
}

// --- Margin ---

func TestBoxMargin(t *testing.T) {
	t.Run("left margin adds spaces", func(t *testing.T) {
		got := Box().MarginLeft(3).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.HasPrefix(lines[0], "   ┌"))
		assert.Equal(t, true, strings.HasPrefix(lines[1], "   │"))
	})

	t.Run("top margin adds blank lines", func(t *testing.T) {
		got := Box().MarginTop(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "", lines[0])
		assert.Equal(t, "", lines[1])
		assert.Equal(t, "┌─┐", lines[2])
	})

	t.Run("bottom margin adds trailing newlines", func(t *testing.T) {
		got := Box().MarginBottom(2).String("x")
		assert.Equal(t, true, strings.HasSuffix(got, "└─┘\n\n"))
	})

	t.Run("uniform margin", func(t *testing.T) {
		got := Box().Margin(1).String("x")
		lines := strings.Split(got, "\n")
		// line 0: top margin (empty)
		assert.Equal(t, "", lines[0])
		// line 1: top border with left margin and right margin
		assert.Equal(t, true, strings.HasPrefix(lines[1], " ┌"))
		assert.Equal(t, true, strings.HasSuffix(lines[1], "┐ "))
	})
}

// --- Box with colors ---

func TestBoxColors(t *testing.T) {
	t.Run("box with border color", func(t *testing.T) {
		ForceColors(true)
		got := Box().Red().String("x")
		lines := strings.Split(got, "\n")
		// Border should be wrapped in red ANSI codes
		assert.Equal(t, true, strings.Contains(lines[0], "\x1b[31m"))
		assert.Equal(t, true, strings.Contains(lines[0], cReset))
	})

	t.Run("box with background", func(t *testing.T) {
		got := Box().OnBlue().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[0], "\x1b[44m"))
	})

	t.Run("box disabled colors", func(t *testing.T) {
		ForceColors(false)
		defer ForceColors(true)
		got := Box().Red().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│x│", lines[1])
		assert.Equal(t, "└─┘", lines[2])
	})
}

// --- Box with styled content (ANSI-aware width) ---

func TestBoxANSIContent(t *testing.T) {
	t.Run("styled content does not break width calculation", func(t *testing.T) {
		styledHello := Text().Red().String("hello")
		got := Box().String(styledHello)
		lines := strings.Split(got, "\n")
		// "hello" is 5 visible chars, so border is 7 wide
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, true, strings.Contains(lines[1], "hello"))
		assert.Equal(t, "└─────┘", lines[2])
	})

	t.Run("mixed styled and plain lines", func(t *testing.T) {
		line1 := Text().Red().String("hi") // 2 visible
		line2 := "world"                   // 5 visible
		got := Box().String(line1 + "\n" + line2)
		lines := strings.Split(got, "\n")
		// Width determined by longest visible line (5)
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "└─────┘", lines[3])
	})
}

// --- Box immutability ---

func TestBoxImmutability(t *testing.T) {
	t.Run("border change does not affect original", func(t *testing.T) {
		base := Box()
		rounded := base.BorderRounded()
		heavy := base.BorderHeavy()

		baseLines := strings.Split(base.String("x"), "\n")
		roundedLines := strings.Split(rounded.String("x"), "\n")
		heavyLines := strings.Split(heavy.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])    // simple unchanged
		assert.Equal(t, "╭─╮", roundedLines[0]) // rounded
		assert.Equal(t, "┏━┓", heavyLines[0])   // heavy
	})

	t.Run("padding change does not affect original", func(t *testing.T) {
		base := Box()
		padded := base.Padding(2)

		baseLines := strings.Split(base.String("x"), "\n")
		paddedLines := strings.Split(padded.String("x"), "\n")

		assert.Equal(t, 3, len(baseLines))                       // top, content, bottom
		assert.Equal(t, true, len(paddedLines) > len(baseLines)) // has padding rows
	})

	t.Run("color change does not affect original", func(t *testing.T) {
		base := Box()
		colored := base.Red()

		baseGot := base.String("x")
		coloredGot := colored.String("x")

		// Base should have no ANSI codes
		assert.Equal(t, false, strings.Contains(baseGot, "\x1b["))
		// Colored should have ANSI codes
		assert.Equal(t, true, strings.Contains(coloredGot, "\x1b[31m"))
	})
}

// --- Box Fprint methods ---

func TestBoxFprint(t *testing.T) {
	t.Run("fprint writes box to buffer", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Box().Fprint(&buf, "hi")
		assert.Equal(t, nil, err)
		assert.Equal(t, true, n > 0)
		assert.Equal(t, true, strings.Contains(buf.String(), "┌──┐"))
	})

	t.Run("fprintln appends newline", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Box().Fprintln(&buf, "hi")
		assert.Equal(t, nil, err)
		assert.Equal(t, true, strings.HasSuffix(buf.String(), "└──┘\n"))
	})

	t.Run("fprintf formats content", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Box().Fprintf(&buf, "n=%d", 42)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, strings.Contains(buf.String(), "n=42"))
	})
}

// --- Box Sprintf ---

func TestBoxSprintf(t *testing.T) {
	t.Run("sprintf formats and boxes", func(t *testing.T) {
		got := Box().Sprintf("count: %d", 42)
		assert.Equal(t, true, strings.Contains(got, "count: 42"))
		assert.Equal(t, true, strings.Contains(got, "┌"))
	})
}

// --- Box Print methods ---

func TestBoxPrintMethods(t *testing.T) {
	t.Run("print writes to default output", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Box().Print("hi")
		assert.Equal(t, true, strings.Contains(buf.String(), "│hi│"))
	})

	t.Run("println writes with trailing newline", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Box().Println("hi")
		assert.Equal(t, true, strings.HasSuffix(buf.String(), "└──┘\n"))
	})

	t.Run("printf formats content", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Box().Printf("n=%d", 42)
		assert.Equal(t, true, strings.Contains(buf.String(), "n=42"))
	})
}

// --- Box concurrent safety ---

func TestBoxConcurrent(t *testing.T) {
	t.Run("shared box used from many goroutines", func(t *testing.T) {
		b := Box().BorderRounded().Padding(1)

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				got := b.String("hello")
				if !strings.Contains(got, "╭") {
					t.Errorf("expected rounded border, got %q", got)
				}
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent branching from same base", func(t *testing.T) {
		base := Box()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rounded := base.BorderRounded()
				heavy := base.BorderHeavy()

				rLines := strings.Split(rounded.String("x"), "\n")
				hLines := strings.Split(heavy.String("x"), "\n")

				if !strings.HasPrefix(rLines[0], "╭") {
					t.Errorf("rounded: got %q", rLines[0])
				}
				if !strings.HasPrefix(hLines[0], "┏") {
					t.Errorf("heavy: got %q", hLines[0])
				}
			}()
		}
		wg.Wait()
	})
}

// --- Box all border color methods ---

func TestBoxAllBorderColors(t *testing.T) {
	fgCases := []struct {
		name string
		fn   func(BoxStyle) BoxStyle
		code string
	}{
		{"Black", BoxStyle.Black, "30"},
		{"Red", BoxStyle.Red, "31"},
		{"Green", BoxStyle.Green, "32"},
		{"Yellow", BoxStyle.Yellow, "33"},
		{"Blue", BoxStyle.Blue, "34"},
		{"Magenta", BoxStyle.Magenta, "35"},
		{"Cyan", BoxStyle.Cyan, "36"},
		{"White", BoxStyle.White, "37"},
		{"BrightBlack", BoxStyle.BrightBlack, "90"},
		{"BrightRed", BoxStyle.BrightRed, "91"},
		{"BrightGreen", BoxStyle.BrightGreen, "92"},
		{"BrightYellow", BoxStyle.BrightYellow, "93"},
		{"BrightBlue", BoxStyle.BrightBlue, "94"},
		{"BrightMagenta", BoxStyle.BrightMagenta, "95"},
		{"BrightCyan", BoxStyle.BrightCyan, "96"},
		{"BrightWhite", BoxStyle.BrightWhite, "97"},
	}
	for _, tc := range fgCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn(Box()).String("x")
			assert.Equal(t, true, strings.Contains(got, "\x1b["+tc.code+"m"))
		})
	}
}

func TestBoxAllBackgrounds(t *testing.T) {
	bgCases := []struct {
		name string
		fn   func(BoxStyle) BoxStyle
		code string
	}{
		{"OnBlack", BoxStyle.OnBlack, "40"},
		{"OnRed", BoxStyle.OnRed, "41"},
		{"OnGreen", BoxStyle.OnGreen, "42"},
		{"OnYellow", BoxStyle.OnYellow, "43"},
		{"OnBlue", BoxStyle.OnBlue, "44"},
		{"OnMagenta", BoxStyle.OnMagenta, "45"},
		{"OnCyan", BoxStyle.OnCyan, "46"},
		{"OnWhite", BoxStyle.OnWhite, "47"},
		{"OnBrightBlack", BoxStyle.OnBrightBlack, "100"},
		{"OnBrightRed", BoxStyle.OnBrightRed, "101"},
		{"OnBrightGreen", BoxStyle.OnBrightGreen, "102"},
		{"OnBrightYellow", BoxStyle.OnBrightYellow, "103"},
		{"OnBrightBlue", BoxStyle.OnBrightBlue, "104"},
		{"OnBrightMagenta", BoxStyle.OnBrightMagenta, "105"},
		{"OnBrightCyan", BoxStyle.OnBrightCyan, "106"},
		{"OnBrightWhite", BoxStyle.OnBrightWhite, "107"},
	}
	for _, tc := range bgCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn(Box()).String("x")
			assert.Equal(t, true, strings.Contains(got, "\x1b["+tc.code+"m"))
		})
	}
}

// --- Box modifiers ---

func TestBoxModifiers(t *testing.T) {
	t.Run("bold border", func(t *testing.T) {
		got := Box().Bold().String("x")
		assert.Equal(t, true, strings.Contains(got, "\x1b[1m"))
	})

	t.Run("dim border", func(t *testing.T) {
		got := Box().Dim().String("x")
		assert.Equal(t, true, strings.Contains(got, "\x1b[2m"))
	})
}

// --- Explicit padding methods ---

func TestBoxPaddingExplicit(t *testing.T) {
	t.Run("PaddingX sets left and right", func(t *testing.T) {
		got := Box().PaddingX(3).String("x")
		lines := strings.Split(got, "\n")
		// innerW = 1 + 3 + 3 = 7
		assert.Equal(t, "┌───────┐", lines[0])
		assert.Equal(t, "│   x   │", lines[1])
		assert.Equal(t, "└───────┘", lines[2])
	})

	t.Run("PaddingY sets top and bottom", func(t *testing.T) {
		got := Box().PaddingY(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│ │", lines[1]) // top pad
		assert.Equal(t, "│x│", lines[2]) // content
		assert.Equal(t, "│ │", lines[3]) // bottom pad
		assert.Equal(t, "└─┘", lines[4])
	})

	t.Run("PaddingTop only", func(t *testing.T) {
		got := Box().PaddingTop(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│ │", lines[1]) // top pad 1
		assert.Equal(t, "│ │", lines[2]) // top pad 2
		assert.Equal(t, "│x│", lines[3]) // content
		assert.Equal(t, "└─┘", lines[4])
	})

	t.Run("PaddingBottom only", func(t *testing.T) {
		got := Box().PaddingBottom(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│x│", lines[1]) // content
		assert.Equal(t, "│ │", lines[2]) // bottom pad 1
		assert.Equal(t, "│ │", lines[3]) // bottom pad 2
		assert.Equal(t, "└─┘", lines[4])
	})

	t.Run("PaddingLeft only", func(t *testing.T) {
		got := Box().PaddingLeft(2).String("x")
		lines := strings.Split(got, "\n")
		// innerW = 1 + 2(left) + 0(right) = 3
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│  x│", lines[1])
		assert.Equal(t, "└───┘", lines[2])
	})

	t.Run("PaddingRight only", func(t *testing.T) {
		got := Box().PaddingRight(2).String("x")
		lines := strings.Split(got, "\n")
		// innerW = 1 + 0(left) + 2(right) = 3
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│x  │", lines[1])
		assert.Equal(t, "└───┘", lines[2])
	})
}

// --- Explicit margin methods ---

func TestBoxMarginExplicit(t *testing.T) {
	t.Run("MarginX sets left and right", func(t *testing.T) {
		got := Box().MarginX(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.HasPrefix(lines[0], "  ┌"))
		assert.Equal(t, true, strings.HasSuffix(lines[0], "┐  "))
	})

	t.Run("MarginY sets top and bottom", func(t *testing.T) {
		got := Box().MarginY(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "", lines[0])          // top margin
		assert.Equal(t, "┌─┐", lines[1])       // border
		assert.Equal(t, true, len(lines) >= 4) // content + bottom border + trailing
		assert.Equal(t, true, strings.HasSuffix(got, "└─┘\n"))
	})

	t.Run("MarginLeft only", func(t *testing.T) {
		got := Box().MarginLeft(4).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.HasPrefix(lines[0], "    ┌"))
	})

	t.Run("MarginRight only", func(t *testing.T) {
		got := Box().MarginRight(3).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.HasSuffix(lines[0], "┐   "))
	})

	t.Run("MarginTop only", func(t *testing.T) {
		got := Box().MarginTop(3).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "", lines[0])
		assert.Equal(t, "", lines[1])
		assert.Equal(t, "", lines[2])
		assert.Equal(t, "┌─┐", lines[3])
	})

	t.Run("MarginBottom only", func(t *testing.T) {
		got := Box().MarginBottom(3).String("x")
		assert.Equal(t, true, strings.HasSuffix(got, "└─┘\n\n\n"))
	})
}

// --- Center ---

func TestBoxCenter(t *testing.T) {
	t.Run("center aligns shorter lines", func(t *testing.T) {
		got := Box().Center().String("hello\nhi")
		lines := strings.Split(got, "\n")
		// maxW = 5 ("hello"), "hi" is 2 -> pad total = 3, left = 1, right = 2
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "│hello│", lines[1])
		assert.Equal(t, "│ hi  │", lines[2]) // centered: 1 left, 2 right
		assert.Equal(t, "└─────┘", lines[3])
	})

	t.Run("center with single line is no-op", func(t *testing.T) {
		got := Box().Center().String("hello")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "│hello│", lines[1])
	})

	t.Run("center with padding", func(t *testing.T) {
		got := Box().Center().PaddingX(1).String("hello\nhi")
		lines := strings.Split(got, "\n")
		// innerW = 5 + 1 + 1 = 7
		// availW = 7 - 1 - 1 = 5
		// "hi" vis=2, total=3, leftPad=1, rightPad=2
		// row = │ + pad(1+1) + hi + pad(2+1) + │
		assert.Equal(t, "┌───────┐", lines[0])
		assert.Equal(t, "│ hello │", lines[1])
		assert.Equal(t, "│  hi   │", lines[2])
		assert.Equal(t, "└───────┘", lines[3])
	})
}

// --- CenterTrim ---

func TestBoxCenterTrim(t *testing.T) {
	t.Run("trims whitespace before centering", func(t *testing.T) {
		got := Box().CenterTrim().String("  hello  \n  hi  ")
		lines := strings.Split(got, "\n")
		// After trim: "hello" and "hi", maxW = 5
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "│hello│", lines[1])
		assert.Equal(t, "│ hi  │", lines[2])
		assert.Equal(t, "└─────┘", lines[3])
	})

	t.Run("CenterTrim with no extra whitespace behaves like Center", func(t *testing.T) {
		trimGot := Box().CenterTrim().String("hello\nhi")
		centerGot := Box().Center().String("hello\nhi")
		assert.Equal(t, centerGot, trimGot)
	})
}

// --- Center/CenterTrim immutability ---

func TestBoxCenterImmutability(t *testing.T) {
	t.Run("Center does not affect original", func(t *testing.T) {
		base := Box()
		centered := base.Center()

		// base should left-align (default)
		baseLines := strings.Split(base.String("hello\nhi"), "\n")
		assert.Equal(t, "│hi   │", baseLines[2]) // left-aligned, padded right

		// centered should center
		centeredLines := strings.Split(centered.String("hello\nhi"), "\n")
		assert.Equal(t, "│ hi  │", centeredLines[2]) // centered
	})

	t.Run("CenterTrim does not affect original", func(t *testing.T) {
		base := Box()
		trimmed := base.CenterTrim()

		// base should not trim
		baseLines := strings.Split(base.String("  hi"), "\n")
		assert.Equal(t, "│  hi│", baseLines[1]) // untrimmed

		// trimmed should trim and center
		trimmedLines := strings.Split(trimmed.String("  hi"), "\n")
		assert.Equal(t, "│hi│", trimmedLines[1]) // trimmed
	})
}

// --- Padding/Margin explicit method immutability ---

func TestBoxPaddingMarginImmutability(t *testing.T) {
	t.Run("PaddingX does not affect original", func(t *testing.T) {
		base := Box()
		padded := base.PaddingX(2)

		baseLines := strings.Split(base.String("x"), "\n")
		paddedLines := strings.Split(padded.String("x"), "\n")

		assert.Equal(t, "│x│", baseLines[1])
		assert.Equal(t, "│  x  │", paddedLines[1])
	})

	t.Run("MarginLeft does not affect original", func(t *testing.T) {
		base := Box()
		margined := base.MarginLeft(3)

		baseLines := strings.Split(base.String("x"), "\n")
		marginedLines := strings.Split(margined.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, true, strings.HasPrefix(marginedLines[0], "   ┌"))
	})
}

// --- DisableTop ---

func TestBoxDisableTop(t *testing.T) {
	t.Run("hides top border row", func(t *testing.T) {
		got := Box().DisableTop().String("hi")
		lines := strings.Split(got, "\n")
		// No top border — first line is content.
		assert.Equal(t, "│hi│", lines[0])
		assert.Equal(t, "└──┘", lines[1])
	})

	t.Run("with padding still hides top border", func(t *testing.T) {
		got := Box().DisableTop().Padding(1).String("x")
		lines := strings.Split(got, "\n")
		// First line is top padding row (no top border).
		assert.Equal(t, "│   │", lines[0])
		assert.Equal(t, "│ x │", lines[1])
		assert.Equal(t, "│   │", lines[2])
		assert.Equal(t, "└───┘", lines[3])
	})
}

// --- DisableBottom ---

func TestBoxDisableBottom(t *testing.T) {
	t.Run("hides bottom border row", func(t *testing.T) {
		got := Box().DisableBottom().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "│hi│", lines[1])
		// No bottom border — string ends after content row + newline.
		assert.Equal(t, 3, len(lines)) // top, content, trailing empty
	})

	t.Run("with padding still hides bottom border", func(t *testing.T) {
		got := Box().DisableBottom().Padding(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│   │", lines[1])
		assert.Equal(t, "│ x │", lines[2])
		assert.Equal(t, "│   │", lines[3])
		// No bottom border.
		assert.Equal(t, 5, len(lines)) // 4 rows + trailing empty
	})
}

// --- DisableLeft ---

func TestBoxDisableLeft(t *testing.T) {
	t.Run("hides left border glyphs", func(t *testing.T) {
		got := Box().DisableLeft().String("hi")
		lines := strings.Split(got, "\n")
		// Left corners/verticals replaced with space.
		assert.Equal(t, " ──┐", lines[0])
		assert.Equal(t, " hi│", lines[1])
		assert.Equal(t, " ──┘", lines[2])
	})
}

// --- DisableRight ---

func TestBoxDisableRight(t *testing.T) {
	t.Run("hides right border glyphs", func(t *testing.T) {
		got := Box().DisableRight().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌── ", lines[0])
		assert.Equal(t, "│hi ", lines[1])
		assert.Equal(t, "└── ", lines[2])
	})
}

// --- Combined disable ---

func TestBoxDisableCombined(t *testing.T) {
	t.Run("disable top and bottom leaves only content rows", func(t *testing.T) {
		got := Box().DisableTop().DisableBottom().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "│hi│", lines[0])
		assert.Equal(t, 2, len(lines)) // content + trailing empty
	})

	t.Run("disable left and right removes verticals", func(t *testing.T) {
		got := Box().DisableLeft().DisableRight().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, " ── ", lines[0])
		assert.Equal(t, " hi ", lines[1])
		assert.Equal(t, " ── ", lines[2])
	})

	t.Run("disable all four sides", func(t *testing.T) {
		got := Box().DisableTop().DisableBottom().DisableLeft().DisableRight().PaddingX(1).String("hi")
		lines := strings.Split(got, "\n")
		// Only content row remains, no border glyphs.
		assert.Equal(t, "  hi  ", lines[0])
		assert.Equal(t, 2, len(lines))
	})

	t.Run("disable top and left for open corner", func(t *testing.T) {
		got := Box().DisableTop().DisableLeft().String("hi")
		lines := strings.Split(got, "\n")
		// No top row. Content row has space instead of left vertical.
		assert.Equal(t, " hi│", lines[0])
		assert.Equal(t, " ──┘", lines[1])
	})

	t.Run("disable bottom and right for shadow offset", func(t *testing.T) {
		got := Box().DisableBottom().DisableRight().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌── ", lines[0])
		assert.Equal(t, "│hi ", lines[1])
		assert.Equal(t, 3, len(lines))
	})
}

// --- Disable immutability ---

func TestBoxDisableImmutability(t *testing.T) {
	t.Run("DisableTop does not affect original", func(t *testing.T) {
		base := Box()
		noTop := base.DisableTop()

		baseLines := strings.Split(base.String("x"), "\n")
		noTopLines := strings.Split(noTop.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])  // base has top border
		assert.Equal(t, "│x│", noTopLines[0]) // no top has content first
		assert.Equal(t, 3, len(baseLines))    // top + content + bottom
		assert.Equal(t, 2, len(noTopLines))   // content + bottom
	})

	t.Run("DisableLeft does not affect original", func(t *testing.T) {
		base := Box()
		noLeft := base.DisableLeft()

		baseLines := strings.Split(base.String("x"), "\n")
		noLeftLines := strings.Split(noLeft.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, " ─┐", noLeftLines[0])
	})

	t.Run("DisableRight does not affect original", func(t *testing.T) {
		base := Box()
		noRight := base.DisableRight()

		baseLines := strings.Split(base.String("x"), "\n")
		noRightLines := strings.Split(noRight.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, "┌─ ", noRightLines[0])
	})

	t.Run("DisableBottom does not affect original", func(t *testing.T) {
		base := Box()
		noBottom := base.DisableBottom()

		baseGot := base.String("x")
		noBottomGot := noBottom.String("x")

		assert.Equal(t, true, strings.Contains(baseGot, "└─┘"))
		assert.Equal(t, false, strings.Contains(noBottomGot, "└─┘"))
	})
}

// --- Disable with colors ---

func TestBoxDisableWithColors(t *testing.T) {
	t.Run("disabled sides still apply style to visible parts", func(t *testing.T) {
		ForceColors(true)
		got := Box().DisableTop().Red().String("x")
		lines := strings.Split(got, "\n")
		// Content row should still have ANSI codes.
		assert.Equal(t, true, strings.Contains(lines[0], "\x1b[31m"))
		// No top border row at all.
		assert.Equal(t, true, strings.Contains(lines[0], "x"))
	})
}
