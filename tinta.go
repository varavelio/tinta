// Package tinta provides a minimal, chainable terminal text styler.
//
// Configure first, write last:
//
//	tinta.Red().Bold().Println("error: something broke")
//	tinta.White().BgBlue().Printf("status: %d", code)
//	msg := tinta.Green().Bold().String("ok")
package tinta

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const cReset = "\x1b[0m"

// ANSI codes (unexported, flat strings for zero-alloc append).
const (
	cBold      = "1"
	cDim       = "2"
	cItalic    = "3"
	cUnderline = "4"
	cInvert    = "7"
	cHidden    = "8"
	cStrike    = "9"

	cBlack   = "30"
	cRed     = "31"
	cGreen   = "32"
	cYellow  = "33"
	cBlue    = "34"
	cMagenta = "35"
	cCyan    = "36"
	cWhite   = "37"

	cBrightBlack   = "90"
	cBrightRed     = "91"
	cBrightGreen   = "92"
	cBrightYellow  = "93"
	cBrightBlue    = "94"
	cBrightMagenta = "95"
	cBrightCyan    = "96"
	cBrightWhite   = "97"

	cBgBlack   = "40"
	cBgRed     = "41"
	cBgGreen   = "42"
	cBgYellow  = "43"
	cBgBlue    = "44"
	cBgMagenta = "45"
	cBgCyan    = "46"
	cBgWhite   = "47"

	cBgBrightBlack   = "100"
	cBgBrightRed     = "101"
	cBgBrightGreen   = "102"
	cBgBrightYellow  = "103"
	cBgBrightBlue    = "104"
	cBgBrightMagenta = "105"
	cBgBrightCyan    = "106"
	cBgBrightWhite   = "107"
)

// enabled caches whether stdout supports ANSI colors. Evaluated once at init.
var enabled = detectColor()

// ForceColors overrides automatic detection and always enables or disables colors.
func ForceColors(on bool) { enabled = on }

// Style accumulates ANSI codes. Build with constructors, write with terminal methods.
type Style struct {
	codes []string
}

func newStyle(code string) Style {
	return Style{codes: []string{code}}
}

// --- Foreground constructors ---

func Black() Style   { return newStyle(cBlack) }
func Red() Style     { return newStyle(cRed) }
func Green() Style   { return newStyle(cGreen) }
func Yellow() Style  { return newStyle(cYellow) }
func Blue() Style    { return newStyle(cBlue) }
func Magenta() Style { return newStyle(cMagenta) }
func Cyan() Style    { return newStyle(cCyan) }
func White() Style   { return newStyle(cWhite) }

func BrightBlack() Style   { return newStyle(cBrightBlack) }
func BrightRed() Style     { return newStyle(cBrightRed) }
func BrightGreen() Style   { return newStyle(cBrightGreen) }
func BrightYellow() Style  { return newStyle(cBrightYellow) }
func BrightBlue() Style    { return newStyle(cBrightBlue) }
func BrightMagenta() Style { return newStyle(cBrightMagenta) }
func BrightCyan() Style    { return newStyle(cBrightCyan) }
func BrightWhite() Style   { return newStyle(cBrightWhite) }

// Modifier-only constructors (no foreground color).

func Bold() Style      { return newStyle(cBold) }
func Underline() Style { return newStyle(cUnderline) }

// --- Background ---

func (s Style) BgBlack() Style   { return s.with(cBgBlack) }
func (s Style) BgRed() Style     { return s.with(cBgRed) }
func (s Style) BgGreen() Style   { return s.with(cBgGreen) }
func (s Style) BgYellow() Style  { return s.with(cBgYellow) }
func (s Style) BgBlue() Style    { return s.with(cBgBlue) }
func (s Style) BgMagenta() Style { return s.with(cBgMagenta) }
func (s Style) BgCyan() Style    { return s.with(cBgCyan) }
func (s Style) BgWhite() Style   { return s.with(cBgWhite) }

func (s Style) BgBrightBlack() Style   { return s.with(cBgBrightBlack) }
func (s Style) BgBrightRed() Style     { return s.with(cBgBrightRed) }
func (s Style) BgBrightGreen() Style   { return s.with(cBgBrightGreen) }
func (s Style) BgBrightYellow() Style  { return s.with(cBgBrightYellow) }
func (s Style) BgBrightBlue() Style    { return s.with(cBgBrightBlue) }
func (s Style) BgBrightMagenta() Style { return s.with(cBgBrightMagenta) }
func (s Style) BgBrightCyan() Style    { return s.with(cBgBrightCyan) }
func (s Style) BgBrightWhite() Style   { return s.with(cBgBrightWhite) }

// --- Modifiers ---

func (s Style) Bold() Style      { return s.with(cBold) }
func (s Style) Dim() Style       { return s.with(cDim) }
func (s Style) Italic() Style    { return s.with(cItalic) }
func (s Style) Underline() Style { return s.with(cUnderline) }
func (s Style) Invert() Style    { return s.with(cInvert) }
func (s Style) Hidden() Style    { return s.with(cHidden) }
func (s Style) Strike() Style    { return s.with(cStrike) }

// --- Terminal methods (these produce output) ---

// String returns the styled text.
func (s Style) String(text string) string {
	return s.render(text)
}

// Sprintf formats and returns the styled text.
func (s Style) Sprintf(format string, a ...any) string {
	return s.render(fmt.Sprintf(format, a...))
}

// Print writes the styled text to stdout.
func (s Style) Print(text string) {
	fmt.Fprint(os.Stdout, s.render(text))
}

// Printf formats and writes the styled text to stdout.
func (s Style) Printf(format string, a ...any) {
	fmt.Fprint(os.Stdout, s.render(fmt.Sprintf(format, a...)))
}

// Println writes the styled text followed by a newline to stdout.
func (s Style) Println(text string) {
	fmt.Fprintln(os.Stdout, s.render(text))
}

// Fprint writes the styled text to w.
func (s Style) Fprint(w io.Writer, text string) {
	fmt.Fprint(w, s.render(text))
}

// Fprintf formats and writes the styled text to w.
func (s Style) Fprintf(w io.Writer, format string, a ...any) {
	fmt.Fprint(w, s.render(fmt.Sprintf(format, a...)))
}

// Fprintln writes the styled text followed by a newline to w.
func (s Style) Fprintln(w io.Writer, text string) {
	fmt.Fprintln(w, s.render(text))
}

// --- internals ---

func (s Style) with(code string) Style {
	s.codes = append(s.codes, code)
	return s
}

func (s Style) render(text string) string {
	if !enabled || len(s.codes) == 0 {
		return text
	}
	var b strings.Builder
	b.Grow(4 + len(s.codes)*4 + len(text) + len(cReset))
	b.WriteString("\x1b[")
	b.WriteString(strings.Join(s.codes, ";"))
	b.WriteByte('m')
	b.WriteString(text)
	b.WriteString(cReset)
	return b.String()
}

func detectColor() bool {
	return colorEnabled(os.Getenv)
}

func colorEnabled(getenv func(string) string) bool {
	if getenv("NO_COLOR") != "" || getenv("NO_COLORS") != "" || getenv("DISABLE_COLORS") != "" {
		return false
	}
	if getenv("FORCE_COLOR") != "" || getenv("CLICOLOR_FORCE") != "" {
		return true
	}
	if getenv("CLICOLOR") == "0" {
		return false
	}
	if strings.EqualFold(getenv("TERM"), "dumb") {
		return false
	}
	return isTerminal(os.Stdout)
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok || f == nil {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice != 0
}
