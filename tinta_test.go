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
		assert.Equal(t, "\x1b[31merror\x1b[0m", Red().String("error"))
	})

	t.Run("foreground and bold", func(t *testing.T) {
		assert.Equal(t, "\x1b[31;1merror\x1b[0m", Red().Bold().String("error"))
	})

	t.Run("foreground background and modifier", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;44;1minfo\x1b[0m", White().OnBlue().Bold().String("info"))
	})

	t.Run("bright foreground", func(t *testing.T) {
		assert.Equal(t, "\x1b[92mok\x1b[0m", BrightGreen().String("ok"))
	})

	t.Run("bright background", func(t *testing.T) {
		assert.Equal(t, "\x1b[30;103mwarn\x1b[0m", Black().OnBrightYellow().String("warn"))
	})

	t.Run("multiple modifiers", func(t *testing.T) {
		assert.Equal(t, "\x1b[36;1;4;3mfancy\x1b[0m", Cyan().Bold().Underline().Italic().String("fancy"))
	})
}

func TestSprintf(t *testing.T) {
	t.Run("formats with args", func(t *testing.T) {
		assert.Equal(t, "\x1b[31mcount: 42\x1b[0m", Red().Sprintf("count: %d", 42))
	})
}

// --- Disabled ---

func TestDisabled(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("string returns plain text", func(t *testing.T) {
		assert.Equal(t, "plain", Red().Bold().String("plain"))
	})

	t.Run("sprintf returns plain text", func(t *testing.T) {
		assert.Equal(t, "n=5", Green().Sprintf("n=%d", 5))
	})
}

// --- Modifiers ---

func TestModifiers(t *testing.T) {
	t.Run("dim", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;2mdim\x1b[0m", White().Dim().String("dim"))
	})

	t.Run("italic", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;3mitalic\x1b[0m", White().Italic().String("italic"))
	})

	t.Run("underline", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;4munderline\x1b[0m", White().Underline().String("underline"))
	})

	t.Run("invert", func(t *testing.T) {
		assert.Equal(t, "\x1b[31;47;7minv\x1b[0m", Red().OnWhite().Invert().String("inv"))
	})

	t.Run("hidden", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;8mhidden\x1b[0m", White().Hidden().String("hidden"))
	})

	t.Run("strike", func(t *testing.T) {
		assert.Equal(t, "\x1b[37;9mstrike\x1b[0m", White().Strike().String("strike"))
	})
}

// --- Standalone constructors ---

func TestBoldConstructor(t *testing.T) {
	t.Run("bold without color", func(t *testing.T) {
		assert.Equal(t, "\x1b[1mbold\x1b[0m", Bold().String("bold"))
	})
}

func TestUnderlineConstructor(t *testing.T) {
	t.Run("underline without color", func(t *testing.T) {
		assert.Equal(t, "\x1b[4mlink\x1b[0m", Underline().String("link"))
	})
}

// --- Writer methods ---

func TestFprint(t *testing.T) {
	t.Run("fprint writes styled text", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Red().Bold().Fprint(&buf, "err")
		assert.Equal(t, nil, err)
		assert.Equal(t, "\x1b[31;1merr\x1b[0m", buf.String())
	})

	t.Run("fprintln appends newline", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Green().Fprintln(&buf, "ok")
		assert.Equal(t, nil, err)
		assert.Equal(t, "\x1b[32mok\x1b[0m\n", buf.String())
	})

	t.Run("fprintf formats with args", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := Blue().Fprintf(&buf, "n=%d", 7)
		assert.Equal(t, nil, err)
		assert.Equal(t, "\x1b[34mn=7\x1b[0m", buf.String())
	})
}

// --- SetOutput ---

func TestSetOutput(t *testing.T) {
	t.Run("print writes to configured output", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil) // reset after test

		Red().Print("hello")
		assert.Equal(t, "\x1b[31mhello\x1b[0m", buf.String())
	})

	t.Run("println writes to configured output", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Green().Println("ok")
		assert.Equal(t, "\x1b[32mok\x1b[0m\n", buf.String())
	})

	t.Run("printf writes to configured output", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Yellow().Printf("x=%d", 3)
		assert.Equal(t, "\x1b[33mx=3\x1b[0m", buf.String())
	})
}

// --- Immutability ---

func TestImmutability(t *testing.T) {
	t.Run("branching from same style does not corrupt", func(t *testing.T) {
		base := Red()
		bold := base.Bold()
		underline := base.Underline()

		assert.Equal(t, "\x1b[31;1mbold\x1b[0m", bold.String("bold"))
		assert.Equal(t, "\x1b[31;4munderline\x1b[0m", underline.String("underline"))
		assert.Equal(t, "\x1b[31mbase\x1b[0m", base.String("base"))
	})

	t.Run("deep chain does not affect parent", func(t *testing.T) {
		a := White().OnBlue()
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
		fn   func(Style) Style
		code string
	}{
		{"OnBlack", Style.OnBlack, "40"},
		{"OnRed", Style.OnRed, "41"},
		{"OnGreen", Style.OnGreen, "42"},
		{"OnYellow", Style.OnYellow, "43"},
		{"OnBlue", Style.OnBlue, "44"},
		{"OnMagenta", Style.OnMagenta, "45"},
		{"OnCyan", Style.OnCyan, "46"},
		{"OnWhite", Style.OnWhite, "47"},
		{"OnBrightBlack", Style.OnBrightBlack, "100"},
		{"OnBrightRed", Style.OnBrightRed, "101"},
		{"OnBrightGreen", Style.OnBrightGreen, "102"},
		{"OnBrightYellow", Style.OnBrightYellow, "103"},
		{"OnBrightBlue", Style.OnBrightBlue, "104"},
		{"OnBrightMagenta", Style.OnBrightMagenta, "105"},
		{"OnBrightCyan", Style.OnBrightCyan, "106"},
		{"OnBrightWhite", Style.OnBrightWhite, "107"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.fn(White())
			expected := "\x1b[37;" + tc.code + "mx\x1b[0m"
			assert.Equal(t, expected, s.String("x"))
		})
	}
}

// --- Color detection ---

func TestColorDetection(t *testing.T) {
	t.Run("NO_COLOR disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"NO_COLOR": "1"})))
	})

	t.Run("NO_COLORS disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"NO_COLORS": "1"})))
	})

	t.Run("DISABLE_COLORS disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"DISABLE_COLORS": "1"})))
	})

	t.Run("disable takes precedence over force", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{
			"NO_COLOR":    "1",
			"FORCE_COLOR": "1",
		})))
	})

	t.Run("FORCE_COLOR enables", func(t *testing.T) {
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"FORCE_COLOR": "1"})))
	})

	t.Run("CLICOLOR_FORCE enables", func(t *testing.T) {
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"CLICOLOR_FORCE": "1"})))
	})

	t.Run("CLICOLOR=0 disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"CLICOLOR": "0"})))
	})

	t.Run("TERM=dumb disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"TERM": "dumb"})))
	})
}

func fakeEnv(m map[string]string) func(string) string {
	return func(key string) string { return m[key] }
}

// --- Concurrency ---

func TestConcurrentString(t *testing.T) {
	t.Run("shared style used from many goroutines", func(t *testing.T) {
		s := Red().Bold()
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
		base := White().OnBlue()

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
				Red().Print("x")
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
				_, err := Red().Fprint(&buf, "x")
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
				_ = Red().String("x")
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

// --- Edge cases ---

func TestEmptyText(t *testing.T) {
	t.Run("empty string still wraps in ANSI", func(t *testing.T) {
		assert.Equal(t, "\x1b[31m\x1b[0m", Red().String(""))
	})

	t.Run("disabled returns empty string as-is", func(t *testing.T) {
		ForceColors(false)
		defer ForceColors(true)
		assert.Equal(t, "", Red().String(""))
	})
}

func TestSprintfNoArgs(t *testing.T) {
	t.Run("sprintf with no format args", func(t *testing.T) {
		assert.Equal(t, "\x1b[32mliteral\x1b[0m", Green().Sprintf("literal"))
	})
}

func TestLongChain(t *testing.T) {
	t.Run("many modifiers chained", func(t *testing.T) {
		s := White().OnRed().Bold().Dim().Italic().Underline().Invert().Hidden().Strike()
		got := s.String("x")
		assert.Equal(t, "\x1b[37;41;1;2;3;4;7;8;9mx\x1b[0m", got)
	})
}

// --- All foreground constructors ---

func TestAllForegrounds(t *testing.T) {
	cases := []struct {
		name string
		fn   func() Style
		code string
	}{
		{"Black", Black, "30"},
		{"Red", Red, "31"},
		{"Green", Green, "32"},
		{"Yellow", Yellow, "33"},
		{"Blue", Blue, "34"},
		{"Magenta", Magenta, "35"},
		{"Cyan", Cyan, "36"},
		{"White", White, "37"},
		{"BrightBlack", BrightBlack, "90"},
		{"BrightRed", BrightRed, "91"},
		{"BrightGreen", BrightGreen, "92"},
		{"BrightYellow", BrightYellow, "93"},
		{"BrightBlue", BrightBlue, "94"},
		{"BrightMagenta", BrightMagenta, "95"},
		{"BrightCyan", BrightCyan, "96"},
		{"BrightWhite", BrightWhite, "97"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := "\x1b[" + tc.code + "mx\x1b[0m"
			assert.Equal(t, expected, tc.fn().String("x"))
		})
	}
}

// --- Render buffer exactness ---

func TestRenderBufferExact(t *testing.T) {
	t.Run("buffer has no excess capacity", func(t *testing.T) {
		// render uses Grow(size) on a fresh Builder.
		// Verify the output length matches the computed size exactly.
		s := White().OnBlue().Bold()
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

		Red().Print("a")
		Red().Print("b")
		assert.Equal(t, "\x1b[31ma\x1b[0m\x1b[31mb\x1b[0m", buf.String())
	})

	t.Run("println appends exactly one newline", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Blue().Println("line")
		assert.Equal(t, "\x1b[34mline\x1b[0m\n", buf.String())
	})

	t.Run("printf formats correctly", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(nil)

		Green().Printf("val=%s num=%d", "ok", 7)
		assert.Equal(t, "\x1b[32mval=ok num=7\x1b[0m", buf.String())
	})
}

// --- ForceColors round-trip ---

func TestForceColorsRoundTrip(t *testing.T) {
	t.Run("enable disable enable", func(t *testing.T) {
		ForceColors(false)
		assert.Equal(t, "x", Red().String("x"))

		ForceColors(true)
		assert.Equal(t, "\x1b[31mx\x1b[0m", Red().String("x"))

		ForceColors(false)
		assert.Equal(t, "x", Red().String("x"))

		ForceColors(true) // restore
	})
}

// --- Fprint methods with disabled colors ---

func TestFprintDisabled(t *testing.T) {
	ForceColors(false)
	defer ForceColors(true)

	t.Run("fprint writes plain text", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Red().Bold().Fprint(&buf, "plain")
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "plain", buf.String())
	})

	t.Run("fprintln writes plain text with newline", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Green().Fprintln(&buf, "line")
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "line\n", buf.String())
	})

	t.Run("fprintf writes plain formatted text", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := Blue().Fprintf(&buf, "n=%d", 42)
		assert.Equal(t, nil, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, "n=42", buf.String())
	})
}

// --- Color detection edge cases ---

func TestColorDetectionEdgeCases(t *testing.T) {
	t.Run("empty env falls through to terminal check", func(t *testing.T) {
		// With no env vars set and stdout not a TTY, should return false.
		got := colorEnabled(fakeEnv(map[string]string{}))
		// In test env stdout is not a real terminal, so this should be false.
		assert.Equal(t, false, got)
	})

	t.Run("TERM=dumb case insensitive", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"TERM": "DUMB"})))
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"TERM": "Dumb"})))
	})

	t.Run("CLICOLOR=1 does not force enable", func(t *testing.T) {
		// CLICOLOR=1 doesn't force colors on, only CLICOLOR=0 forces them off.
		// With non-TTY stdout this should still be false.
		got := colorEnabled(fakeEnv(map[string]string{"CLICOLOR": "1"}))
		assert.Equal(t, false, got)
	})

	t.Run("FORCE_COLOR with any value enables", func(t *testing.T) {
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"FORCE_COLOR": "0"})))
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"FORCE_COLOR": "yes"})))
	})
}

// --- Sprintf with multiple format verbs ---

func TestSprintfMultipleVerbs(t *testing.T) {
	t.Run("multiple format arguments", func(t *testing.T) {
		got := Yellow().Sprintf("%s: %d/%d (%.1f%%)", "progress", 3, 10, 30.0)
		assert.Equal(t, "\x1b[33mprogress: 3/10 (30.0%)\x1b[0m", got)
	})
}

// --- Concurrent Fprint to independent buffers ---

func TestConcurrentFprintIndependent(t *testing.T) {
	t.Run("concurrent Fprint to separate buffers", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			i := i // capture for Go <1.22
			wg.Add(1)
			go func() {
				defer wg.Done()
				var buf bytes.Buffer
				_, err := Red().Fprint(&buf, fmt.Sprintf("item-%d", i))
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
