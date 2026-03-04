package tinta

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/varavelio/tinta/internal/assert"
)

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

func TestBoxBorderStyles(t *testing.T) {
	t.Run("rounded", func(t *testing.T) {
		got := Box().Border(BorderRounded).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "╭─╮", lines[0])
		assert.Equal(t, "│x│", lines[1])
		assert.Equal(t, "╰─╯", lines[2])
	})

	t.Run("double", func(t *testing.T) {
		got := Box().Border(BorderDouble).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "╔═╗", lines[0])
		assert.Equal(t, "║x║", lines[1])
		assert.Equal(t, "╚═╝", lines[2])
	})

	t.Run("heavy", func(t *testing.T) {
		got := Box().Border(BorderHeavy).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┏━┓", lines[0])
		assert.Equal(t, "┃x┃", lines[1])
		assert.Equal(t, "┗━┛", lines[2])
	})

	t.Run("simple is default", func(t *testing.T) {
		got := Box().Border(BorderSimple).String("x")
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

func TestBoxPadding(t *testing.T) {
	t.Run("uniform padding", func(t *testing.T) {
		got := Box().Padding(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│   │", lines[1])
		assert.Equal(t, "│ x │", lines[2])
		assert.Equal(t, "│   │", lines[3])
		assert.Equal(t, "└───┘", lines[4])
	})

	t.Run("horizontal padding only", func(t *testing.T) {
		got := Box().PaddingX(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "│  x  │", lines[1])
		assert.Equal(t, "└─────┘", lines[2])
	})

	t.Run("individual side padding", func(t *testing.T) {
		got := Box().PaddingTop(1).PaddingRight(2).PaddingBottom(0).PaddingLeft(3).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──────┐", lines[0])
		assert.Equal(t, "│      │", lines[1])
		assert.Equal(t, "│   x  │", lines[2])
		assert.Equal(t, "└──────┘", lines[3])
	})
}

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
		assert.Equal(t, "", lines[0])
		assert.Equal(t, true, strings.HasPrefix(lines[1], " ┌"))
		assert.Equal(t, true, strings.HasSuffix(lines[1], "┐ "))
	})
}

func TestBoxColors(t *testing.T) {
	t.Run("box with border color", func(t *testing.T) {
		ForceColors(true)
		got := Box().Red().String("x")
		lines := strings.Split(got, "\n")
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

func TestBoxANSIContent(t *testing.T) {
	t.Run("styled content does not break width calculation", func(t *testing.T) {
		styledHello := Text().Red().String("hello")
		got := Box().String(styledHello)
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, true, strings.Contains(lines[1], "hello"))
		assert.Equal(t, "└─────┘", lines[2])
	})

	t.Run("mixed styled and plain lines", func(t *testing.T) {
		line1 := Text().Red().String("hi")
		line2 := "world"
		got := Box().String(line1 + "\n" + line2)
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "└─────┘", lines[3])
	})
}

func TestBoxImmutability(t *testing.T) {
	t.Run("border change does not affect original", func(t *testing.T) {
		base := Box()
		rounded := base.Border(BorderRounded)
		heavy := base.Border(BorderHeavy)

		baseLines := strings.Split(base.String("x"), "\n")
		roundedLines := strings.Split(rounded.String("x"), "\n")
		heavyLines := strings.Split(heavy.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, "╭─╮", roundedLines[0])
		assert.Equal(t, "┏━┓", heavyLines[0])
	})

	t.Run("padding change does not affect original", func(t *testing.T) {
		base := Box()
		padded := base.Padding(2)

		baseLines := strings.Split(base.String("x"), "\n")
		paddedLines := strings.Split(padded.String("x"), "\n")

		assert.Equal(t, 3, len(baseLines))
		assert.Equal(t, true, len(paddedLines) > len(baseLines))
	})

	t.Run("color change does not affect original", func(t *testing.T) {
		base := Box()
		colored := base.Red()

		baseGot := base.String("x")
		coloredGot := colored.String("x")

		assert.Equal(t, false, strings.Contains(baseGot, "\x1b["))
		assert.Equal(t, true, strings.Contains(coloredGot, "\x1b[31m"))
	})
}

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

func TestBoxSprintf(t *testing.T) {
	t.Run("sprintf formats and boxes", func(t *testing.T) {
		got := Box().Sprintf("count: %d", 42)
		assert.Equal(t, true, strings.Contains(got, "count: 42"))
		assert.Equal(t, true, strings.Contains(got, "┌"))
	})
}

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

func TestBoxConcurrent(t *testing.T) {
	t.Run("shared box used from many goroutines", func(t *testing.T) {
		b := Box().Border(BorderRounded).Padding(1)

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
				rounded := base.Border(BorderRounded)
				heavy := base.Border(BorderHeavy)

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

func TestBoxPaddingExplicit(t *testing.T) {
	t.Run("PaddingX sets left and right", func(t *testing.T) {
		got := Box().PaddingX(3).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌───────┐", lines[0])
		assert.Equal(t, "│   x   │", lines[1])
		assert.Equal(t, "└───────┘", lines[2])
	})

	t.Run("PaddingY sets top and bottom", func(t *testing.T) {
		got := Box().PaddingY(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│ │", lines[1])
		assert.Equal(t, "│x│", lines[2])
		assert.Equal(t, "│ │", lines[3])
		assert.Equal(t, "└─┘", lines[4])
	})

	t.Run("PaddingTop only", func(t *testing.T) {
		got := Box().PaddingTop(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│ │", lines[1])
		assert.Equal(t, "│ │", lines[2])
		assert.Equal(t, "│x│", lines[3])
		assert.Equal(t, "└─┘", lines[4])
	})

	t.Run("PaddingBottom only", func(t *testing.T) {
		got := Box().PaddingBottom(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, "│x│", lines[1])
		assert.Equal(t, "│ │", lines[2])
		assert.Equal(t, "│ │", lines[3])
		assert.Equal(t, "└─┘", lines[4])
	})

	t.Run("PaddingLeft only", func(t *testing.T) {
		got := Box().PaddingLeft(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│  x│", lines[1])
		assert.Equal(t, "└───┘", lines[2])
	})

	t.Run("PaddingRight only", func(t *testing.T) {
		got := Box().PaddingRight(2).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│x  │", lines[1])
		assert.Equal(t, "└───┘", lines[2])
	})
}

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
		assert.Equal(t, "", lines[0])
		assert.Equal(t, "┌─┐", lines[1])
		assert.Equal(t, true, len(lines) >= 4)
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

func TestBoxCenter(t *testing.T) {
	t.Run("center aligns shorter lines", func(t *testing.T) {
		got := Box().Center().String("hello\nhi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─────┐", lines[0])
		assert.Equal(t, "│hello│", lines[1])
		assert.Equal(t, "│ hi  │", lines[2])
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
		assert.Equal(t, "┌───────┐", lines[0])
		assert.Equal(t, "│ hello │", lines[1])
		assert.Equal(t, "│  hi   │", lines[2])
		assert.Equal(t, "└───────┘", lines[3])
	})
}

func TestBoxCenterTrim(t *testing.T) {
	t.Run("trims whitespace before centering", func(t *testing.T) {
		got := Box().CenterTrim().String("  hello  \n  hi  ")
		lines := strings.Split(got, "\n")
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

func TestBoxCenterImmutability(t *testing.T) {
	t.Run("Center does not affect original", func(t *testing.T) {
		base := Box()
		centered := base.Center()

		baseLines := strings.Split(base.String("hello\nhi"), "\n")
		assert.Equal(t, "│hi   │", baseLines[2])

		centeredLines := strings.Split(centered.String("hello\nhi"), "\n")
		assert.Equal(t, "│ hi  │", centeredLines[2])
	})

	t.Run("CenterTrim does not affect original", func(t *testing.T) {
		base := Box()
		trimmed := base.CenterTrim()

		baseLines := strings.Split(base.String("  hi"), "\n")
		assert.Equal(t, "│  hi│", baseLines[1])

		trimmedLines := strings.Split(trimmed.String("  hi"), "\n")
		assert.Equal(t, "│hi│", trimmedLines[1])
	})
}

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

func TestBoxDisableTop(t *testing.T) {
	t.Run("hides top border row", func(t *testing.T) {
		got := Box().DisableTop().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌hi┐", lines[0])
		assert.Equal(t, "└──┘", lines[1])
	})

	t.Run("with padding still hides top border", func(t *testing.T) {
		got := Box().DisableTop().Padding(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌   ┐", lines[0])
		assert.Equal(t, "│ x │", lines[1])
		assert.Equal(t, "│   │", lines[2])
		assert.Equal(t, "└───┘", lines[3])
	})
}

func TestBoxDisableBottom(t *testing.T) {
	t.Run("hides bottom border row", func(t *testing.T) {
		got := Box().DisableBottom().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "└hi┘", lines[1])
		assert.Equal(t, 3, len(lines))
	})

	t.Run("with padding still hides bottom border", func(t *testing.T) {
		got := Box().DisableBottom().Padding(1).String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌───┐", lines[0])
		assert.Equal(t, "│   │", lines[1])
		assert.Equal(t, "│ x │", lines[2])
		assert.Equal(t, "└   ┘", lines[3])
		assert.Equal(t, 5, len(lines))
	})
}

func TestBoxDisableLeft(t *testing.T) {
	t.Run("hides left border glyphs", func(t *testing.T) {
		got := Box().DisableLeft().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, " hi│", lines[1])
		assert.Equal(t, "└──┘", lines[2])
	})
}

func TestBoxDisableRight(t *testing.T) {
	t.Run("hides right border glyphs", func(t *testing.T) {
		got := Box().DisableRight().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "│hi ", lines[1])
		assert.Equal(t, "└──┘", lines[2])
	})
}

func TestBoxDisableCombined(t *testing.T) {
	t.Run("disable top and bottom leaves only content rows", func(t *testing.T) {
		got := Box().DisableTop().DisableBottom().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌hi┐", lines[0])
		assert.Equal(t, 2, len(lines))
	})

	t.Run("disable left and right removes verticals", func(t *testing.T) {
		got := Box().DisableLeft().DisableRight().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, " hi ", lines[1])
		assert.Equal(t, "└──┘", lines[2])
	})

	t.Run("disable all four sides", func(t *testing.T) {
		got := Box().DisableTop().DisableBottom().DisableLeft().DisableRight().PaddingX(1).String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "  hi  ", lines[0])
		assert.Equal(t, 2, len(lines))
	})

	t.Run("disable top and left for open corner", func(t *testing.T) {
		got := Box().DisableTop().DisableLeft().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, " hi┐", lines[0])
		assert.Equal(t, "└──┘", lines[1])
	})

	t.Run("disable bottom and right for shadow offset", func(t *testing.T) {
		got := Box().DisableBottom().DisableRight().String("hi")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", lines[0])
		assert.Equal(t, "└hi ", lines[1])
		assert.Equal(t, 3, len(lines))
	})
}

func TestBoxCornerCapsWithHiddenTopBottom(t *testing.T) {
	t.Run("top hidden plus left hidden keeps top-right corner", func(t *testing.T) {
		got := Box().DisableTop().DisableLeft().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, " x┐", lines[0])
		assert.Equal(t, "└─┘", lines[1])
	})

	t.Run("top hidden plus right hidden keeps top-left corner", func(t *testing.T) {
		got := Box().DisableTop().DisableRight().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌x ", lines[0])
		assert.Equal(t, "└─┘", lines[1])
	})

	t.Run("bottom hidden plus left hidden keeps bottom-right corner", func(t *testing.T) {
		got := Box().DisableBottom().DisableLeft().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, "┌─┐", lines[0])
		assert.Equal(t, " x┘", lines[1])
	})

	t.Run("explicit corner disable wins for top-right cap", func(t *testing.T) {
		got := Box().DisableTop().DisableLeft().DisableTopRightCorner().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, " x│", lines[0])
		assert.Equal(t, "└─┘", lines[1])
	})
}

func TestBoxDisableImmutability(t *testing.T) {
	t.Run("DisableTop does not affect original", func(t *testing.T) {
		base := Box()
		noTop := base.DisableTop()

		baseLines := strings.Split(base.String("x"), "\n")
		noTopLines := strings.Split(noTop.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, "┌x┐", noTopLines[0])
		assert.Equal(t, 3, len(baseLines))
		assert.Equal(t, 2, len(noTopLines))
	})

	t.Run("DisableLeft does not affect original", func(t *testing.T) {
		base := Box()
		noLeft := base.DisableLeft()

		baseLines := strings.Split(base.String("x"), "\n")
		noLeftLines := strings.Split(noLeft.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, "┌─┐", noLeftLines[0])
	})

	t.Run("DisableRight does not affect original", func(t *testing.T) {
		base := Box()
		noRight := base.DisableRight()

		baseLines := strings.Split(base.String("x"), "\n")
		noRightLines := strings.Split(noRight.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseLines[0])
		assert.Equal(t, "┌─┐", noRightLines[0])
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

func TestBoxDisableWithColors(t *testing.T) {
	t.Run("disabled sides still apply style to visible parts", func(t *testing.T) {
		ForceColors(true)
		got := Box().DisableTop().Red().String("x")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[0], "\x1b[31m"))
		assert.Equal(t, true, strings.Contains(lines[0], "x"))
	})
}

func TestBoxNestedColorRobustness(t *testing.T) {
	t.Run("inner box resets do not corrupt outer border", func(t *testing.T) {
		ForceColors(true)
		inner := Box().Border(BorderRounded).Blue().String("hello")
		outer := Box().Border(BorderDouble).Red().String(inner)

		rows := strings.Split(outer, "\n")
		for _, row := range rows {
			if row == "" {
				continue
			}
			assert.Equal(t, true, strings.Contains(row, "\x1b[31m"))
			if strings.Contains(row, "hello") {
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
			assert.Equal(t, true, strings.Contains(row, "\x1b[32m"))
			if strings.Contains(row, "styled") {
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
			assert.Equal(t, true, strings.Contains(row, "\x1b[47;31m"))
		}
	})

	t.Run("no style box with styled content passes through cleanly", func(t *testing.T) {
		ForceColors(true)
		styledContent := Text().Red().String("red text")
		outer := Box().String(styledContent)

		rows := strings.Split(outer, "\n")
		for _, row := range rows {
			if strings.Contains(row, "red text") {
				assert.Equal(t, true, strings.Contains(row, "\x1b[31m"))
			}
		}
	})
}

func TestBoxCenterLine(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("center specific line by index", func(t *testing.T) {
		got := Box().CenterLine(0).String("hi\nworld!")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[1], "│  hi  │"))
		assert.Equal(t, true, strings.Contains(lines[2], "│world!│"))
	})

	t.Run("center multiple specific lines", func(t *testing.T) {
		got := Box().CenterLine(0).CenterLine(2).String("a\nbb\nc")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[1], "│a │"))
		assert.Equal(t, true, strings.Contains(lines[2], "│bb│"))
		assert.Equal(t, true, strings.Contains(lines[3], "│c │"))
	})

	t.Run("out of bounds index is silently ignored", func(t *testing.T) {
		got := Box().CenterLine(99).String("hello")
		lines := strings.Split(got, "\n")
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
		assert.Equal(t, true, strings.Contains(lines[1], "│long line│"))
		assert.Equal(t, true, strings.Contains(lines[2], "│   hi    │"))
	})
}

func TestBoxCenterFirstLine(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("centers only the first line", func(t *testing.T) {
		got := Box().CenterFirstLine().String("Title\nLeft-aligned content")
		lines := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(lines[1], "│       Title        │"))
		assert.Equal(t, true, strings.Contains(lines[2], "│Left-aligned content│"))
	})

	t.Run("single line still centers", func(t *testing.T) {
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

		baseLines := strings.Split(baseGot, "\n")
		assert.Equal(t, true, strings.Contains(baseLines[1], "│a │"))

		centeredLines := strings.Split(centeredGot, "\n")
		assert.Equal(t, true, strings.Contains(centeredLines[1], "│a │"))
	})

	t.Run("CenterLine does not modify original with wider content", func(t *testing.T) {
		base := Box()
		centered := base.CenterLine(0)

		baseGot := base.String("hi\nworld!")
		centeredGot := centered.String("hi\nworld!")

		baseLines := strings.Split(baseGot, "\n")
		assert.Equal(t, true, strings.Contains(baseLines[1], "│hi    │"))

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

func TestBoxCornerControls(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("disabling left side keeps corners", func(t *testing.T) {
		got := Box().Border(BorderRounded).DisableLeft().String("x")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "╭─╮", rows[0])
		assert.Equal(t, " x│", rows[1])
		assert.Equal(t, "╰─╯", rows[2])
	})

	t.Run("disable individual corners", func(t *testing.T) {
		got := Box().DisableTopLeftCorner().DisableBottomRightCorner().String("x")
		rows := strings.Split(got, "\n")
		assert.Equal(t, " ─┐", rows[0])
		assert.Equal(t, "│x│", rows[1])
		assert.Equal(t, "└─ ", rows[2])
	})

	t.Run("disable all corners", func(t *testing.T) {
		got := Box().DisableCorners().String("x")
		rows := strings.Split(got, "\n")
		assert.Equal(t, " ─ ", rows[0])
		assert.Equal(t, "│x│", rows[1])
		assert.Equal(t, " ─ ", rows[2])
	})
}

func TestBoxCornerControlImmutability(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("corner changes do not modify original", func(t *testing.T) {
		base := Box()
		changed := base.DisableCorners()

		baseRows := strings.Split(base.String("x"), "\n")
		changedRows := strings.Split(changed.String("x"), "\n")

		assert.Equal(t, "┌─┐", baseRows[0])
		assert.Equal(t, " ─ ", changedRows[0])
	})
}

func TestBoxTitle(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("title left aligned", func(t *testing.T) {
		got := Box().Title("Hi", AlignLeft).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌─Hi────┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└───────┘", rows[2])
	})

	t.Run("title center aligned", func(t *testing.T) {
		got := Box().Title("Hi", AlignCenter).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌──Hi───┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└───────┘", rows[2])
	})

	t.Run("title right aligned", func(t *testing.T) {
		got := Box().Title("Hi", AlignRight).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌────Hi─┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└───────┘", rows[2])
	})

	t.Run("title with padding", func(t *testing.T) {
		got := Box().PaddingX(1).Title("T", AlignLeft).String("hi")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌─T──┐", rows[0])
		assert.Equal(t, "│ hi │", rows[1])
		assert.Equal(t, "└────┘", rows[2])
	})

	t.Run("title widens box when text is longer than content", func(t *testing.T) {
		got := Box().Title("LongTitle", AlignLeft).String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌─LongTitle─┐", rows[0])
		assert.Equal(t, "│ab         │", rows[1])
		assert.Equal(t, "└───────────┘", rows[2])
	})

	t.Run("title with rounded border", func(t *testing.T) {
		got := Box().Border(BorderRounded).Title("OK", AlignCenter).String("hi")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "╭─OK─╮", rows[0])
		assert.Equal(t, "│hi  │", rows[1])
		assert.Equal(t, "╰────╯", rows[2])
	})

	t.Run("title with hidden top does not show title", func(t *testing.T) {
		got := Box().DisableTop().Title("Hi", AlignLeft).String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌ab  ┐", rows[0])
		assert.Equal(t, "└────┘", rows[1])
	})

	t.Run("empty title has no effect", func(t *testing.T) {
		got := Box().Title("", AlignLeft).String("hi")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌──┐", rows[0])
	})
}

func TestBoxFooter(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("footer left aligned", func(t *testing.T) {
		got := Box().Footer("Ok", AlignLeft).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌───────┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└─Ok────┘", rows[2])
	})

	t.Run("footer center aligned", func(t *testing.T) {
		got := Box().Footer("Ok", AlignCenter).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌───────┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└──Ok───┘", rows[2])
	})

	t.Run("footer right aligned", func(t *testing.T) {
		got := Box().Footer("Ok", AlignRight).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌───────┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└────Ok─┘", rows[2])
	})

	t.Run("footer widens box when text is longer than content", func(t *testing.T) {
		got := Box().Footer("LongFooter", AlignLeft).String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌────────────┐", rows[0])
		assert.Equal(t, "│ab          │", rows[1])
		assert.Equal(t, "└─LongFooter─┘", rows[2])
	})

	t.Run("footer with hidden bottom does not show footer", func(t *testing.T) {
		got := Box().DisableBottom().Footer("Ok", AlignLeft).String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌────┐", rows[0])
		assert.Equal(t, "└ab  ┘", rows[1])
	})
}

func TestBoxTitleAndFooter(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("title and footer together", func(t *testing.T) {
		got := Box().
			Title("Name", AlignLeft).
			Footer("Done", AlignRight).
			String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌─Name──┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└──Done─┘", rows[2])
	})

	t.Run("title and footer both center aligned", func(t *testing.T) {
		got := Box().
			Title("Top", AlignCenter).
			Footer("Bot", AlignCenter).
			String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌──Top──┐", rows[0])
		assert.Equal(t, "│content│", rows[1])
		assert.Equal(t, "└──Bot──┘", rows[2])
	})

	t.Run("title wider than footer expands box", func(t *testing.T) {
		got := Box().
			Title("VeryLongTitle", AlignLeft).
			Footer("X", AlignLeft).
			String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌─VeryLongTitle─┐", rows[0])
		assert.Equal(t, "│ab             │", rows[1])
		assert.Equal(t, "└─X─────────────┘", rows[2])
	})

	t.Run("footer wider than title expands box", func(t *testing.T) {
		got := Box().
			Title("X", AlignLeft).
			Footer("VeryLongFooter", AlignLeft).
			String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌─X──────────────┐", rows[0])
		assert.Equal(t, "│ab              │", rows[1])
		assert.Equal(t, "└─VeryLongFooter─┘", rows[2])
	})
}

func TestBoxTitleImmutability(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("title does not modify original", func(t *testing.T) {
		base := Box()
		withTitle := base.Title("Hi", AlignLeft)

		baseRows := strings.Split(base.String("ab"), "\n")
		titleRows := strings.Split(withTitle.String("ab"), "\n")

		assert.Equal(t, "┌──┐", baseRows[0])
		assert.Equal(t, "┌─Hi─┐", titleRows[0])
	})

	t.Run("footer does not modify original", func(t *testing.T) {
		base := Box()
		withFooter := base.Footer("Ok", AlignLeft)

		baseRows := strings.Split(base.String("ab"), "\n")
		footerRows := strings.Split(withFooter.String("ab"), "\n")

		assert.Equal(t, "└──┘", baseRows[2])
		assert.Equal(t, "└─Ok─┘", footerRows[2])
	})
}

func TestBoxTitleWithCornerControls(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("title with hidden top-left corner", func(t *testing.T) {
		got := Box().DisableTopLeftCorner().Title("Hi", AlignLeft).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, " ─Hi────┐", rows[0])
	})

	t.Run("title with hidden top-right corner", func(t *testing.T) {
		got := Box().DisableTopRightCorner().Title("Hi", AlignRight).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, "┌────Hi─ ", rows[0])
	})

	t.Run("footer with hidden bottom corners", func(t *testing.T) {
		got := Box().DisableCorners().Footer("Ok", AlignCenter).String("content")
		rows := strings.Split(got, "\n")
		assert.Equal(t, " ──Ok─── ", rows[2])
	})
}

func TestBoxTitleWithColors(t *testing.T) {
	ForceColors(true)

	t.Run("title is wrapped in box ANSI codes", func(t *testing.T) {
		got := Box().Red().Title("Hi", AlignLeft).String("ab")
		rows := strings.Split(got, "\n")
		assert.Equal(t, true, strings.Contains(rows[0], "Hi"))
		assert.Equal(t, true, strings.Contains(rows[0], "\x1b["))
		assert.Equal(t, true, strings.Contains(rows[0], "\x1b[0m"))
	})
}

func TestBox3DBorderEffect(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("nested box with canvas creates 3D effect", func(t *testing.T) {
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
