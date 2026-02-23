package tinta

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/varavelio/tinta/internal/assert"
)

func init() {
	ForceColors(true)
}

// --- String output ---

func TestString(t *testing.T) {
	t.Run("foreground only", func(t *testing.T) {
		assert.Equal(t, "\x1b[31merror\x1b[0m", Text().Red().String("error"))
	})

	t.Run("foreground and bold", func(t *testing.T) {
		assert.Equal(t, "\x1b[31;1merror\x1b[0m", Text().Red().Bold().String("error"))
	})

	t.Run("foreground background and modifier", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;44;1minfo\x1b[0m", Text().White().OnBlue().Bold().String("info"))
	})

	t.Run("bright foreground", func(t *testing.T) {
		assert.Equal(t, "\x1b[92mok\x1b[0m", Text().BrightGreen().String("ok"))
	})

	t.Run("bright background", func(t *testing.T) {
		assert.Equal(t, "\x1b[30;103mwarn\x1b[0m", Text().Black().OnBrightYellow().String("warn"))
	})

	t.Run("multiple modifiers", func(t *testing.T) {
		assert.Equal(t, "\x1b[36;1;4;3mfancy\x1b[0m", Text().Cyan().Bold().Underline().Italic().String("fancy"))
	})
}

func TestSprintf(t *testing.T) {
	t.Run("formats with args", func(t *testing.T) {
		assert.Equal(t, "\x1b[31mcount: 42\x1b[0m", Text().Red().Sprintf("count: %d", 42))
	})
}

// --- Text() entry point ---

func TestTextEntry(t *testing.T) {
	t.Run("text with no codes returns plain text", func(t *testing.T) {
		assert.Equal(t, "hello", Text().String("hello"))
	})

	t.Run("text chained with color", func(t *testing.T) {
		assert.Equal(t, "\x1b[31mhello\x1b[0m", Text().Red().String("hello"))
	})

	t.Run("text chained with color and modifier", func(t *testing.T) {
		assert.Equal(t, "\x1b[31;1mhello\x1b[0m", Text().Red().Bold().String("hello"))
	})

	t.Run("text has all foreground methods", func(t *testing.T) {
		base := Text()
		assert.Equal(t, "\x1b[30mx\x1b[0m", base.Black().String("x"))
		assert.Equal(t, "\x1b[31mx\x1b[0m", base.Red().String("x"))
		assert.Equal(t, "\x1b[32mx\x1b[0m", base.Green().String("x"))
		assert.Equal(t, "\x1b[33mx\x1b[0m", base.Yellow().String("x"))
		assert.Equal(t, "\x1b[34mx\x1b[0m", base.Blue().String("x"))
		assert.Equal(t, "\x1b[35mx\x1b[0m", base.Magenta().String("x"))
		assert.Equal(t, "\x1b[36mx\x1b[0m", base.Cyan().String("x"))
		assert.Equal(t, "\x1b[37mx\x1b[0m", base.White().String("x"))
	})

	t.Run("text has all bright foreground methods", func(t *testing.T) {
		base := Text()
		assert.Equal(t, "\x1b[90mx\x1b[0m", base.BrightBlack().String("x"))
		assert.Equal(t, "\x1b[91mx\x1b[0m", base.BrightRed().String("x"))
		assert.Equal(t, "\x1b[92mx\x1b[0m", base.BrightGreen().String("x"))
		assert.Equal(t, "\x1b[93mx\x1b[0m", base.BrightYellow().String("x"))
		assert.Equal(t, "\x1b[94mx\x1b[0m", base.BrightBlue().String("x"))
		assert.Equal(t, "\x1b[95mx\x1b[0m", base.BrightMagenta().String("x"))
		assert.Equal(t, "\x1b[96mx\x1b[0m", base.BrightCyan().String("x"))
		assert.Equal(t, "\x1b[97mx\x1b[0m", base.BrightWhite().String("x"))
	})
}

// --- Disabled ---

func TestDisabled(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("string returns plain text", func(t *testing.T) {
		assert.Equal(t, "plain", Text().Red().Bold().String("plain"))
	})

	t.Run("sprintf returns plain text", func(t *testing.T) {
		assert.Equal(t, "n=5", Text().Green().Sprintf("n=%d", 5))
	})
}

// --- Modifiers ---

func TestModifiers(t *testing.T) {
	t.Run("dim", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;2mdim\x1b[0m", Text().White().Dim().String("dim"))
	})

	t.Run("italic", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;3mitalic\x1b[0m", Text().White().Italic().String("italic"))
	})

	t.Run("underline", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;4munderline\x1b[0m", Text().White().Underline().String("underline"))
	})

	t.Run("invert", func(t *testing.T) {
		assert.Equal(t, "\x1b[31;47;7minv\x1b[0m", Text().Red().OnWhite().Invert().String("inv"))
	})

	t.Run("hidden", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;8mhidden\x1b[0m", Text().White().Hidden().String("hidden"))
	})

	t.Run("strike", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;9mstrike\x1b[0m", Text().White().Strike().String("strike"))
	})
}

// --- Modifier-only styles ---

func TestBoldOnly(t *testing.T) {
	t.Run("bold without color", func(t *testing.T) {
		assert.Equal(t, "\x1b[1mbold\x1b[0m", Text().Bold().String("bold"))
	})
}

func TestUnderlineOnly(t *testing.T) {
	t.Run("underline without color", func(t *testing.T) {
		assert.Equal(t, "\x1b[4mlink\x1b[0m", Text().Underline().String("link"))
	})
}

// --- Writer methods ---

func TestFprint(t *testing.T) {
	t.Run("fprint writes styled text", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Text().Red().Bold().Fprint(&buf, "err")
		assert.Equal(t, nil, err)
		assert.Equal(t, "\x1b[31;1merr\x1b[0m", buf.String())
	})

	t.Run("fprintln appends newline", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Text().Green().Fprintln(&buf, "ok")
		assert.Equal(t, nil, err)
		assert.Equal(t, "\x1b[32mok\x1b[0m\n", buf.String())
	})

	t.Run("fprintf formats with args", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Text().Blue().Fprintf(&buf, "n=%d", 7)
		assert.Equal(t, nil, err)
		assert.Equal(t, "\x1b[34mn=7\x1b[0m", buf.String())
	})
}

// --- Immutability ---

func TestImmutability(t *testing.T) {
	t.Run("branching from same style does not corrupt", func(t *testing.T) {
		base := Text().Red()
		bold := base.Bold()
		underline := base.Underline()

		assert.Equal(t, "\x1b[31;1mbold\x1b[0m", bold.String("bold"))
		assert.Equal(t, "\x1b[31;4munderline\x1b[0m", underline.String("underline"))
		assert.Equal(t, "\x1b[31mbase\x1b[0m", base.String("base"))
	})

	t.Run("deep chain does not affect parent", func(t *testing.T) {
		a := Text().White().OnBlue()
		b := a.Bold()
		c := a.Italic()

		assert.Equal(t, "\x1b[37;44mplain\x1b[0m", a.String("plain"))
		assert.Equal(t, "\x1b[37;44;1mbold\x1b[0m", b.String("bold"))
		assert.Equal(t, "\x1b[37;44;3mitalic\x1b[0m", c.String("italic"))
	})
}

// --- All backgrounds ---

func TestAllBackgrounds(t *testing.T) {
	cases := []struct {
		name string
		fn   func(TextStyle) TextStyle
		code string
	}{
		{"OnBlack", TextStyle.OnBlack, "40"},
		{"OnRed", TextStyle.OnRed, "41"},
		{"OnGreen", TextStyle.OnGreen, "42"},
		{"OnYellow", TextStyle.OnYellow, "43"},
		{"OnBlue", TextStyle.OnBlue, "44"},
		{"OnMagenta", TextStyle.OnMagenta, "45"},
		{"OnCyan", TextStyle.OnCyan, "46"},
		{"OnWhite", TextStyle.OnWhite, "47"},
		{"OnBrightBlack", TextStyle.OnBrightBlack, "100"},
		{"OnBrightRed", TextStyle.OnBrightRed, "101"},
		{"OnBrightGreen", TextStyle.OnBrightGreen, "102"},
		{"OnBrightYellow", TextStyle.OnBrightYellow, "103"},
		{"OnBrightBlue", TextStyle.OnBrightBlue, "104"},
		{"OnBrightMagenta", TextStyle.OnBrightMagenta, "105"},
		{"OnBrightCyan", TextStyle.OnBrightCyan, "106"},
		{"OnBrightWhite", TextStyle.OnBrightWhite, "107"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.fn(Text().White())
			expected := "\x1b[37;" + tc.code + "mx\x1b[0m"
			assert.Equal(t, expected, s.String("x"))
		})
	}
}

// --- All foreground colors ---

func TestAllForegrounds(t *testing.T) {
	cases := []struct {
		name string
		fn   func(TextStyle) TextStyle
		code string
	}{
		{"Black", TextStyle.Black, "30"},
		{"Red", TextStyle.Red, "31"},
		{"Green", TextStyle.Green, "32"},
		{"Yellow", TextStyle.Yellow, "33"},
		{"Blue", TextStyle.Blue, "34"},
		{"Magenta", TextStyle.Magenta, "35"},
		{"Cyan", TextStyle.Cyan, "36"},
		{"White", TextStyle.White, "37"},
		{"BrightBlack", TextStyle.BrightBlack, "90"},
		{"BrightRed", TextStyle.BrightRed, "91"},
		{"BrightGreen", TextStyle.BrightGreen, "92"},
		{"BrightYellow", TextStyle.BrightYellow, "93"},
		{"BrightBlue", TextStyle.BrightBlue, "94"},
		{"BrightMagenta", TextStyle.BrightMagenta, "95"},
		{"BrightCyan", TextStyle.BrightCyan, "96"},
		{"BrightWhite", TextStyle.BrightWhite, "97"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.fn(Text())
			expected := "\x1b[" + tc.code + "mx\x1b[0m"
			assert.Equal(t, expected, s.String("x"))
		})
	}
}

// --- Render buffer exactness ---

func TestRenderBufferExact(t *testing.T) {
	t.Run("buffer has no excess capacity", func(t *testing.T) {
		s := Text().White().OnBlue().Bold()
		got := s.String("hello")

		// \x1b[37;44;1mhello\x1b[0m
		//  2 + 2+1+2+1+1 + 1 + 5 + 4 = 19
		expected := "\x1b[37;44;1mhello\x1b[0m"
		assert.Equal(t, expected, got)
		assert.Equal(t, 19, len(got))
	})
}

// --- Print/Println/Printf with SetOutput ---

func TestPrintMethods(t *testing.T) {
	t.Run("print does not append newline", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Text().Red().Print("a")
		Text().Red().Print("b")
		assert.Equal(t, "\x1b[31ma\x1b[0m\x1b[31mb\x1b[0m", buf.String())
	})

	t.Run("println appends exactly one newline", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Text().Blue().Println("line")
		assert.Equal(t, "\x1b[34mline\x1b[0m\n", buf.String())
	})

	t.Run("printf formats correctly", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Text().Green().Printf("val=%s num=%d", "ok", 7)
		assert.Equal(t, "\x1b[32mval=ok num=7\x1b[0m", buf.String())
	})
}

// --- ForceColors round-trip ---

func TestForceColorsRoundTrip(t *testing.T) {
	t.Run("enable disable enable", func(t *testing.T) {
		ForceColors(false)
		assert.Equal(t, "x", Text().Red().String("x"))

		ForceColors(true)
		assert.Equal(t, "\x1b[31mx\x1b[0m", Text().Red().String("x"))

		ForceColors(false)
		assert.Equal(t, "x", Text().Red().String("x"))

		ForceColors(true) // restore
	})
}

// --- Fprint methods with disabled colors ---

func TestFprintDisabled(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("fprint writes plain text", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Text().Red().Bold().Fprint(&buf, "plain")
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "plain", buf.String())
	})

	t.Run("fprintln writes plain text with newline", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Text().Green().Fprintln(&buf, "line")
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "line\n", buf.String())
	})

	t.Run("fprintf writes plain formatted text", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Text().Blue().Fprintf(&buf, "n=%d", 42)
		assert.Equal(t, nil, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, "n=42", buf.String())
	})
}

// --- Edge cases ---

func TestEmptyText(t *testing.T) {
	t.Run("empty string still wraps in ANSI", func(t *testing.T) {
		assert.Equal(t, "\x1b[31m\x1b[0m", Text().Red().String(""))
	})

	t.Run("disabled returns empty string as-is", func(t *testing.T) {
		ForceColors(false)
		defer ForceColors(true)
		assert.Equal(t, "", Text().Red().String(""))
	})
}

func TestSprintfNoArgs(t *testing.T) {
	t.Run("sprintf with no format args", func(t *testing.T) {
		assert.Equal(t, "\x1b[32mliteral\x1b[0m", Text().Green().Sprintf("literal"))
	})
}

func TestLongChain(t *testing.T) {
	t.Run("many modifiers chained", func(t *testing.T) {
		s := Text().White().OnRed().Bold().Dim().Italic().Underline().Invert().Hidden().Strike()
		got := s.String("x")
		assert.Equal(t, "\x1b[37;41;1;2;3;4;7;8;9mx\x1b[0m", got)
	})
}

func TestSprintfMultipleVerbs(t *testing.T) {
	t.Run("multiple format arguments", func(t *testing.T) {
		got := Text().Yellow().Sprintf("%s: %d/%d (%.1f%%)", "progress", 3, 10, 30.0)
		assert.Equal(t, "\x1b[33mprogress: 3/10 (30.0%)\x1b[0m", got)
	})
}

// --- Concurrency ---

func TestConcurrentString(t *testing.T) {
	t.Run("shared style used from many goroutines", func(t *testing.T) {
		s := Text().Red().Bold()
		expected := "\x1b[31;1mhello\x1b[0m"

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				got := s.String("hello")
				if got != expected {
					t.Errorf("expected %q, got %q", expected, got)
				}
			}()
		}
		wg.Wait()
	})
}

func TestConcurrentBranching(t *testing.T) {
	t.Run("concurrent branching from same base", func(t *testing.T) {
		base := Text().White().OnBlue()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				bold := base.Bold()
				italic := base.Italic()

				gotBold := bold.String("b")
				gotItalic := italic.String("i")
				gotBase := base.String("x")

				if gotBold != "\x1b[37;44;1mb\x1b[0m" {
					t.Errorf("bold: got %q", gotBold)
				}
				if gotItalic != "\x1b[37;44;3mi\x1b[0m" {
					t.Errorf("italic: got %q", gotItalic)
				}
				if gotBase != "\x1b[37;44mx\x1b[0m" {
					t.Errorf("base: got %q", gotBase)
				}
			}()
		}
		wg.Wait()
	})
}

func TestConcurrentPrint(t *testing.T) {
	t.Run("concurrent Print to discard does not panic or race", func(t *testing.T) {
		SetOutput(io.Discard)
		defer SetOutput(nil)

		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Text().Red().Print("x")
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent Fprint to independent buffers", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var buf bytes.Buffer
				_, err := Text().Red().Fprint(&buf, "x")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				got := buf.String()
				if got != "\x1b[31mx\x1b[0m" {
					t.Errorf("expected styled output, got %q", got)
				}
			}()
		}
		wg.Wait()
	})
}

func TestConcurrentForceColors(t *testing.T) {
	t.Run("toggle ForceColors while rendering", func(t *testing.T) {
		defer ForceColors(true) // restore

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			i := i // capture for Go <1.22
			wg.Add(1)
			go func() {
				defer wg.Done()
				if i%2 == 0 {
					ForceColors(false)
				} else {
					ForceColors(true)
				}
				// Should not panic or race; result depends on timing.
				_ = Text().Red().String("x")
			}()
		}
		wg.Wait()
	})
}

func TestConcurrentSetOutput(t *testing.T) {
	t.Run("switch output while printing", func(t *testing.T) {
		defer SetOutput(nil)

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var buf bytes.Buffer
				SetOutput(&buf)
			}()
		}
		wg.Wait()
	})
}

func TestConcurrentFprintIndependent(t *testing.T) {
	t.Run("concurrent Fprint to separate buffers", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			i := i // capture for Go <1.22
			wg.Add(1)
			go func() {
				defer wg.Done()
				var buf bytes.Buffer
				_, err := Text().Red().Fprint(&buf, fmt.Sprintf("item-%d", i))
				if err != nil {
					t.Errorf("goroutine %d: Fprint error: %v", i, err)
					return
				}
				got := buf.String()
				expected := fmt.Sprintf("\x1b[31mitem-%d\x1b[0m", i)
				if got != expected {
					t.Errorf("goroutine %d: expected %q, got %q", i, expected, got)
				}
			}()
		}
		wg.Wait()
	})
}
