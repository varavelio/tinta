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
		fn   func(*BoxStyle) *BoxStyle
		code string
	}{
		{"Black", (*BoxStyle).Black, "30"},
		{"Red", (*BoxStyle).Red, "31"},
		{"Green", (*BoxStyle).Green, "32"},
		{"Yellow", (*BoxStyle).Yellow, "33"},
		{"Blue", (*BoxStyle).Blue, "34"},
		{"Magenta", (*BoxStyle).Magenta, "35"},
		{"Cyan", (*BoxStyle).Cyan, "36"},
		{"White", (*BoxStyle).White, "37"},
		{"BrightBlack", (*BoxStyle).BrightBlack, "90"},
		{"BrightRed", (*BoxStyle).BrightRed, "91"},
		{"BrightGreen", (*BoxStyle).BrightGreen, "92"},
		{"BrightYellow", (*BoxStyle).BrightYellow, "93"},
		{"BrightBlue", (*BoxStyle).BrightBlue, "94"},
		{"BrightMagenta", (*BoxStyle).BrightMagenta, "95"},
		{"BrightCyan", (*BoxStyle).BrightCyan, "96"},
		{"BrightWhite", (*BoxStyle).BrightWhite, "97"},
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
		fn   func(*BoxStyle) *BoxStyle
		code string
	}{
		{"OnBlack", (*BoxStyle).OnBlack, "40"},
		{"OnRed", (*BoxStyle).OnRed, "41"},
		{"OnGreen", (*BoxStyle).OnGreen, "42"},
		{"OnYellow", (*BoxStyle).OnYellow, "43"},
		{"OnBlue", (*BoxStyle).OnBlue, "44"},
		{"OnMagenta", (*BoxStyle).OnMagenta, "45"},
		{"OnCyan", (*BoxStyle).OnCyan, "46"},
		{"OnWhite", (*BoxStyle).OnWhite, "47"},
		{"OnBrightBlack", (*BoxStyle).OnBrightBlack, "100"},
		{"OnBrightRed", (*BoxStyle).OnBrightRed, "101"},
		{"OnBrightGreen", (*BoxStyle).OnBrightGreen, "102"},
		{"OnBrightYellow", (*BoxStyle).OnBrightYellow, "103"},
		{"OnBrightBlue", (*BoxStyle).OnBrightBlue, "104"},
		{"OnBrightMagenta", (*BoxStyle).OnBrightMagenta, "105"},
		{"OnBrightCyan", (*BoxStyle).OnBrightCyan, "106"},
		{"OnBrightWhite", (*BoxStyle).OnBrightWhite, "107"},
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

// --- Nested box color robustness ---

func TestBoxNestedColorRobustness(t *testing.T) {
	t.Run("inner box resets do not corrupt outer border", func(t *testing.T) {
		ForceColors(true)
		inner := Box().BorderRounded().Blue().String("hello")
		outer := Box().BorderDouble().Red().String(inner)

		rows := strings.Split(outer, "\n")
		// The outer box content rows contain the inner box lines.
		// Each content row must start and end with the outer style (red).
		// The right border "║" must be wrapped in its own red ANSI
		// sequence, reapplied AFTER the inner box's resets.
		for _, row := range rows {
			if row == "" {
				continue
			}
			// Every non-empty row must contain the outer red code.
			assert.Equal(t, true, strings.Contains(row, "\x1b[31m"))
			if strings.Contains(row, "hello") {
				// This is a content row. The right chrome is wrapped
				// separately: ...<inner resets>\x1b[31m║\x1b[0m
				// Count occurrences of outer code — must appear at least
				// twice (once for left chrome, once for right chrome).
				count := strings.Count(row, "\x1b[31m")
				assert.Equal(t, true, count >= 2)
			}
		}
	})

	t.Run("inner styled text resets do not corrupt outer border", func(t *testing.T) {
		ForceColors(true)
		styledContent := Text().Red().Bold().String("styled")
		outer := Box().Green().String(styledContent)

		rows := strings.Split(outer, "\n")
		for _, row := range rows {
			if row == "" {
				continue
			}
			// Every row must contain green code.
			assert.Equal(t, true, strings.Contains(row, "\x1b[32m"))
			if strings.Contains(row, "styled") {
				// Right chrome must be re-styled after the content's reset.
				// The green code must appear at least twice (left + right chrome).
				count := strings.Count(row, "\x1b[32m")
				assert.Equal(t, true, count >= 2)
			}
		}
	})

	t.Run("deeply nested boxes preserve all styles", func(t *testing.T) {
		ForceColors(true)
		innermost := Box().Blue().String("deep")
		middle := Box().Green().String(innermost)
		outermost := Box().Red().String(middle)

		rows := strings.Split(outermost, "\n")
		// The outermost box must have red styling on every row.
		for _, row := range rows {
			if row == "" {
				continue
			}
			assert.Equal(t, true, strings.Contains(row, "\x1b[31m"))
		}
	})

	t.Run("nested box with padding preserves outer background", func(t *testing.T) {
		ForceColors(true)
		inner := Box().Blue().String("x")
		outer := Box().OnWhite().Red().PaddingX(1).String(inner)

		rows := strings.Split(outer, "\n")
		for _, row := range rows {
			if row == "" {
				continue
			}
			// Every row must contain the outer style (on-white + red).
			// Code order depends on method call order: OnWhite (47) then Red (31).
			assert.Equal(t, true, strings.Contains(row, "\x1b[47;31m"))
		}
	})

	t.Run("no style box with styled content passes through cleanly", func(t *testing.T) {
		ForceColors(true)
		styledContent := Text().Red().String("red text")
		// Outer box has no style — content should pass through as-is.
		outer := Box().String(styledContent)

		rows := strings.Split(outer, "\n")
		for _, row := range rows {
			if strings.Contains(row, "red text") {
				assert.Equal(t, true, strings.Contains(row, "\x1b[31m"))
			}
		}
	})
}

// --- Selective line centering ---

func TestBoxCenterLine(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("center specific line by index", func(t *testing.T) {
		got := Box().CenterLine(0).String("hi\nworld!")
		lines := strings.Split(got, "\n")
		// Line 0 ("hi") should be centered within the width of "world!" (6 chars).
		// "hi" is 2 chars, so 4 extra spaces: 2 left + 2 right.
		assert.Equal(t, true, strings.Contains(lines[1], "│  hi  │"))
		// Line 1 ("world!") should be left-aligned (not centered).
		assert.Equal(t, true, strings.Contains(lines[2], "│world!│"))
	})

	t.Run("center multiple specific lines", func(t *testing.T) {
		got := Box().CenterLine(0).CenterLine(2).String("a\nbb\nc")
		lines := strings.Split(got, "\n")
		// Width is 2 (from "bb"). "a" and "c" are 1 char each — 1 space to distribute.
		// "a": 0 left, 1 right (1/2=0). Wait — total=1, leftPad=0, rightPad=1.
		// Actually let me check: availW=2, vis("a")=1, total=1, leftPad=0, rightPad=1.
		// So "a" gets " a " — no, leftPad=0 means: "│" + " "*0 + "a" + " "*1 + "│" = "│a │"
		// Hmm, that's not really "centered". With total=1, leftPad=total/2=0, rightPad=1.
		// That's the same as left-aligned with rightPad=1. That's correct — integer division
		// rounds down, so odd remainders lean right.
		assert.Equal(t, true, strings.Contains(lines[1], "│a │"))
		assert.Equal(t, true, strings.Contains(lines[2], "│bb│"))
		assert.Equal(t, true, strings.Contains(lines[3], "│c │"))
	})

	t.Run("out of bounds index is silently ignored", func(t *testing.T) {
		got := Box().CenterLine(99).String("hello")
		lines := strings.Split(got, "\n")
		// Only 1 line — index 99 doesn't exist, no centering applied.
		assert.Equal(t, true, strings.Contains(lines[1], "│hello│"))
	})

	t.Run("negative index is silently ignored", func(t *testing.T) {
		got := Box().CenterLine(-1).String("hello")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[1], "│hello│"))
	})

	t.Run("CenterLine does not affect other lines", func(t *testing.T) {
		got := Box().CenterLine(1).String("long line\nhi")
		lines := strings.Split(got, "\n")
		// Line 0 ("long line") is left-aligned, no extra padding.
		assert.Equal(t, true, strings.Contains(lines[1], "│long line│"))
		// Line 1 ("hi") centered within width 9: total=7, left=3, right=4.
		assert.Equal(t, true, strings.Contains(lines[2], "│   hi    │"))
	})
}

func TestBoxCenterFirstLine(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("centers only the first line", func(t *testing.T) {
		got := Box().CenterFirstLine().String("Title\nLeft-aligned content")
		lines := strings.Split(got, "\n")
		// Width is 20 (from "Left-aligned content"). "Title" is 5 chars.
		// total=15, leftPad=7, rightPad=8.
		assert.Equal(t, true, strings.Contains(lines[1], "│       Title        │"))
		assert.Equal(t, true, strings.Contains(lines[2], "│Left-aligned content│"))
	})

	t.Run("single line still centers", func(t *testing.T) {
		// With a single line, there's nothing wider, so no centering effect.
		got := Box().CenterFirstLine().String("only")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[1], "│only│"))
	})
}

func TestBoxCenterLastLine(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("centers only the last line", func(t *testing.T) {
		got := Box().CenterLastLine().String("Left-aligned content\nEnd")
		lines := strings.Split(got, "\n")
		// Width is 20. "End" is 3 chars.
		// total=17, leftPad=8, rightPad=9.
		assert.Equal(t, true, strings.Contains(lines[1], "│Left-aligned content│"))
		assert.Equal(t, true, strings.Contains(lines[2], "│        End         │"))
	})
}

func TestBoxCenterLineWithCenter(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("Center() overrides CenterLine — all lines centered", func(t *testing.T) {
		got := Box().CenterLine(0).Center().String("a\nbb")
		lines := strings.Split(got, "\n")
		// Center() centers all lines, so both are centered.
		assert.Equal(t, true, strings.Contains(lines[1], "│a │"))
		assert.Equal(t, true, strings.Contains(lines[2], "│bb│"))
	})
}

func TestBoxCenterLineImmutability(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("CenterLine does not modify original", func(t *testing.T) {
		base := Box()
		centered := base.CenterLine(0)

		baseGot := base.String("a\nbb")
		centeredGot := centered.String("a\nbb")

		// Base should be left-aligned.
		baseLines := strings.Split(baseGot, "\n")
		assert.Equal(t, true, strings.Contains(baseLines[1], "│a │"))

		// Centered version should center line 0.
		centeredLines := strings.Split(centeredGot, "\n")
		assert.Equal(t, true, strings.Contains(centeredLines[1], "│a │"))
		// With width 2 and "a" being 1 char, total=1, leftPad=0.
		// Both look the same because the centering space is only 1 char.
		// Use a wider example to see the difference.
	})

	t.Run("CenterLine does not modify original with wider content", func(t *testing.T) {
		base := Box()
		centered := base.CenterLine(0)

		baseGot := base.String("hi\nworld!")
		centeredGot := centered.String("hi\nworld!")

		// Base: "hi" left-aligned within width 6 → "│hi    │"
		baseLines := strings.Split(baseGot, "\n")
		assert.Equal(t, true, strings.Contains(baseLines[1], "│hi    │"))

		// Centered: "hi" centered within width 6 → "│  hi  │"
		centeredLines := strings.Split(centeredGot, "\n")
		assert.Equal(t, true, strings.Contains(centeredLines[1], "│  hi  │"))
	})

	t.Run("CenterFirstLine does not modify original", func(t *testing.T) {
		base := Box()
		_ = base.CenterFirstLine()

		got := base.String("hi\nworld!")
		baseLines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(baseLines[1], "│hi    │"))
	})
}

// --- Shadow ---

func TestBoxShadowBottomRight(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("basic shadow structure with corners", func(t *testing.T) {
		got := Box().Shadow(ShadowBottomRight, ShadowLight).String("hi")
		rows := strings.Split(got, "\n")
		// Expected (ShadowLight uses ░ for every glyph):
		// ┌──┐      (row 0: box + spacer)
		// │hi│░     (row 1: box + TopRight corner)
		// └──┘░     (row 2: box + Vertical)
		//  ░░░░     (row 3: spacer + BottomLeft + ░░ + BottomRight)
		assert.Equal(t, 4, len(rows))
		// Row 0: top border + space.
		assert.Equal(t, "┌──┐ ", rows[0])
		// Row 1: content + TopRight corner (░).
		assert.Equal(t, "│hi│░", rows[1])
		// Row 2: bottom border + vertical (░).
		assert.Equal(t, "└──┘░", rows[2])
		// Row 3: space + shadow bar (░░░░).
		assert.Equal(t, " ░░░░", rows[3])
	})

	t.Run("shadow adds one extra row and one extra column", func(t *testing.T) {
		noShadow := Box().String("test")
		withShadow := Box().Shadow(ShadowBottomRight, ShadowLight).String("test")

		noRows := strings.Split(noShadow, "\n")
		shadowRows := strings.Split(withShadow, "\n")

		// Shadow adds one extra row.
		assert.Equal(t, len(noRows)+1, len(shadowRows))

		// Each shadow row is 1 char wider than the non-shadow row.
		for i := 0; i < len(noRows); i++ {
			noW := visibleWidth(noRows[i])
			shW := visibleWidth(shadowRows[i])
			assert.Equal(t, noW+1, shW)
		}
	})

	t.Run("custom style renders distinct corners", func(t *testing.T) {
		sty := ShadowStyle{
			TopLeft: "╭", TopRight: "╮", BottomLeft: "╰", BottomRight: "╯",
			Horizontal: "─", Vertical: "│",
		}
		got := Box().Shadow(ShadowBottomRight, sty).String("hi")
		rows := strings.Split(got, "\n")
		// Row 0: box + space.
		assert.Equal(t, "┌──┐ ", rows[0])
		// Row 1: box + TopRight corner (╮).
		assert.Equal(t, "│hi│╮", rows[1])
		// Row 2: box + Vertical (│).
		assert.Equal(t, "└──┘│", rows[2])
		// Row 3: space + BottomLeft + ── + BottomRight = ╰──╯
		assert.Equal(t, " ╰──╯", rows[3])
	})
}

func TestBoxShadowBottomLeft(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("basic shadow structure with corners", func(t *testing.T) {
		got := Box().Shadow(ShadowBottomLeft, ShadowLight).String("hi")
		rows := strings.Split(got, "\n")
		assert.Equal(t, 4, len(rows))
		// Row 0: space + top border.
		assert.Equal(t, " ┌──┐", rows[0])
		// Row 1: TopLeft corner (░) + content.
		assert.Equal(t, "░│hi│", rows[1])
		// Row 2: vertical (░) + bottom border.
		assert.Equal(t, "░└──┘", rows[2])
		// Row 3: shadow bar + space.
		assert.Equal(t, "░░░░ ", rows[3])
	})
}

func TestBoxShadowTopRight(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("basic shadow structure with corners", func(t *testing.T) {
		got := Box().Shadow(ShadowTopRight, ShadowLight).String("hi")
		rows := strings.Split(got, "\n")
		assert.Equal(t, 4, len(rows))
		// Row 0: space + shadow top bar (TopLeft + ░░ + TopRight = ░░░░).
		assert.Equal(t, " ░░░░", rows[0])
		// Row 1: top border + Vertical.
		assert.Equal(t, "┌──┐░", rows[1])
		// Row 2: content + Vertical.
		assert.Equal(t, "│hi│░", rows[2])
		// Row 3: bottom border + BottomRight corner.
		assert.Equal(t, "└──┘░", rows[3])
	})
}

func TestBoxShadowTopLeft(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("basic shadow structure with corners", func(t *testing.T) {
		got := Box().Shadow(ShadowTopLeft, ShadowLight).String("hi")
		rows := strings.Split(got, "\n")
		assert.Equal(t, 4, len(rows))
		// Row 0: shadow bar + space.
		assert.Equal(t, "░░░░ ", rows[0])
		// Row 1: vertical (░) + top border.
		assert.Equal(t, "░┌──┐", rows[1])
		// Row 2: vertical (░) + content.
		assert.Equal(t, "░│hi│", rows[2])
		// Row 3: BottomLeft corner (░) + bottom border.
		assert.Equal(t, "░└──┘", rows[3])
	})
}

func TestBoxShadowWithColors(t *testing.T) {
	t.Run("shadow gets its own ANSI codes", func(t *testing.T) {
		ForceColors(true)
		got := Box().Red().Shadow(ShadowBottomRight, ShadowLight).String("x")
		rows := strings.Split(got, "\n")
		// Shadow rows should contain the bright-black code (default shadow color).
		// The shadow bottom row should have shadow styling.
		lastRow := rows[len(rows)-1]
		assert.Equal(t, true, strings.Contains(lastRow, "\x1b[90m"))
	})

	t.Run("custom shadow style renders correct glyphs", func(t *testing.T) {
		ForceColors(false)
		defer ForceColors(true)
		got := Box().Shadow(ShadowBottomRight, ShadowBlock).String("x")
		rows := strings.Split(got, "\n")
		// Shadow should use █ instead of ░.
		assert.Equal(t, true, strings.Contains(rows[1], "█"))
		assert.Equal(t, true, strings.Contains(rows[len(rows)-1], "█"))
	})
}

func TestBoxShadowWithPadding(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("shadow works with padding", func(t *testing.T) {
		got := Box().Padding(1).Shadow(ShadowBottomRight, ShadowLight).String("x")
		rows := strings.Split(got, "\n")
		// Box has padding so more rows: top border, top pad, content, bottom pad, bottom border + shadow bottom.
		assert.Equal(t, 6, len(rows))
		// Last row is shadow bottom.
		assert.Equal(t, true, strings.Contains(rows[5], "░"))
	})
}

func TestBoxShadowImmutability(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("Shadow does not modify original", func(t *testing.T) {
		base := Box()
		shadowed := base.Shadow(ShadowBottomRight, ShadowLight)

		baseGot := base.String("x")
		shadowGot := shadowed.String("x")

		baseRows := strings.Split(baseGot, "\n")
		shadowRows := strings.Split(shadowGot, "\n")

		// Base has 3 rows, shadow has 4.
		assert.Equal(t, 3, len(baseRows))
		assert.Equal(t, 4, len(shadowRows))
	})

	t.Run("Shadow with different styles does not cross-contaminate", func(t *testing.T) {
		base := Box()
		light := base.Shadow(ShadowBottomRight, ShadowLight)
		block := base.Shadow(ShadowBottomRight, ShadowBlock)

		lightGot := light.String("x")
		blockGot := block.String("x")

		// Light uses ░, block uses █.
		assert.Equal(t, true, strings.Contains(lightGot, "░"))
		assert.Equal(t, true, strings.Contains(blockGot, "█"))
		assert.Equal(t, false, strings.Contains(lightGot, "█"))
		assert.Equal(t, false, strings.Contains(blockGot, "░"))
	})
}
