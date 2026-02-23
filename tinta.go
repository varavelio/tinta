// Package tinta provides a minimal, chainable terminal text styler.
//
// Configure first, write last, e.g.:
//
//	tinta.Red().Bold().Println("error: something broke")
//	tinta.White().OnBlue().Printf("status: %d", code)
//	msg := tinta.Green().Bold().String("ok")
//
// The default output is os.Stdout. Change it with [SetOutput].
//
// Color support is detected automatically. Override with [ForceColors].
package tinta

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const cReset = "\x1b[0m"

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

	cOnBlack   = "40"
	cOnRed     = "41"
	cOnGreen   = "42"
	cOnYellow  = "43"
	cOnBlue    = "44"
	cOnMagenta = "45"
	cOnCyan    = "46"
	cOnWhite   = "47"

	cOnBrightBlack   = "100"
	cOnBrightRed     = "101"
	cOnBrightGreen   = "102"
	cOnBrightYellow  = "103"
	cOnBrightBlue    = "104"
	cOnBrightMagenta = "105"
	cOnBrightCyan    = "106"
	cOnBrightWhite   = "107"
)

// Package-level state, protected by mu.
var (
	mu      sync.RWMutex
	output  io.Writer = os.Stdout
	enabled           = detectColor()
)

// SetOutput changes the default writer used by [Style.Print], [Style.Println]
// and [Style.Printf]. It is safe for concurrent use.
func SetOutput(w io.Writer) {
	mu.Lock()
	output = w
	mu.Unlock()
}

// ForceColors overrides automatic detection. It is safe for concurrent use.
func ForceColors(on bool) {
	mu.Lock()
	enabled = on
	mu.Unlock()
}

func getOutput() io.Writer {
	mu.RLock()
	w := output
	mu.RUnlock()
	return w
}

func isEnabled() bool {
	mu.RLock()
	v := enabled
	mu.RUnlock()
	return v
}

// style is intentionally unexported. Users interact with it through the
// exported Style type alias below and the package-level constructors.
// This prevents manual instantiation (e.g. Style{}) which would be meaningless.
type style struct {
	codes []string
}

// Style is the public handle returned by every constructor and modifier.
// The underlying struct is opaque; users cannot create one manually.
type Style = *style

// newStyle allocates a style with a single initial code.
func newStyle(code string) Style {
	return &style{codes: []string{code}}
}

// with returns a new style that has all existing codes plus one more.
// It copies the slice to guarantee immutability of the source.
func (s *style) with(code string) Style {
	cp := make([]string, len(s.codes)+1)
	copy(cp, s.codes)
	cp[len(s.codes)] = code
	return &style{codes: cp}
}

// --- Foreground constructors (package entry points) ---

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

// --- Background (On*) ---

func (s *style) OnBlack() Style   { return s.with(cOnBlack) }
func (s *style) OnRed() Style     { return s.with(cOnRed) }
func (s *style) OnGreen() Style   { return s.with(cOnGreen) }
func (s *style) OnYellow() Style  { return s.with(cOnYellow) }
func (s *style) OnBlue() Style    { return s.with(cOnBlue) }
func (s *style) OnMagenta() Style { return s.with(cOnMagenta) }
func (s *style) OnCyan() Style    { return s.with(cOnCyan) }
func (s *style) OnWhite() Style   { return s.with(cOnWhite) }

func (s *style) OnBrightBlack() Style   { return s.with(cOnBrightBlack) }
func (s *style) OnBrightRed() Style     { return s.with(cOnBrightRed) }
func (s *style) OnBrightGreen() Style   { return s.with(cOnBrightGreen) }
func (s *style) OnBrightYellow() Style  { return s.with(cOnBrightYellow) }
func (s *style) OnBrightBlue() Style    { return s.with(cOnBrightBlue) }
func (s *style) OnBrightMagenta() Style { return s.with(cOnBrightMagenta) }
func (s *style) OnBrightCyan() Style    { return s.with(cOnBrightCyan) }
func (s *style) OnBrightWhite() Style   { return s.with(cOnBrightWhite) }

// --- Modifiers ---

func (s *style) Bold() Style      { return s.with(cBold) }
func (s *style) Dim() Style       { return s.with(cDim) }
func (s *style) Italic() Style    { return s.with(cItalic) }
func (s *style) Underline() Style { return s.with(cUnderline) }
func (s *style) Invert() Style    { return s.with(cInvert) }
func (s *style) Hidden() Style    { return s.with(cHidden) }
func (s *style) Strike() Style    { return s.with(cStrike) }

// --- Terminal methods ---

// String returns the styled text.
func (s *style) String(text string) string {
	return s.render(text)
}

// Sprintf formats and returns the styled text.
func (s *style) Sprintf(format string, a ...any) string {
	return s.render(fmt.Sprintf(format, a...))
}

// Print writes the styled text to the default output.
func (s *style) Print(text string) {
	_, _ = fmt.Fprint(getOutput(), s.render(text))
}

// Printf formats and writes the styled text to the default output.
func (s *style) Printf(format string, a ...any) {
	_, _ = fmt.Fprint(getOutput(), s.render(fmt.Sprintf(format, a...)))
}

// Println writes the styled text followed by a newline to the default output.
func (s *style) Println(text string) {
	_, _ = fmt.Fprintln(getOutput(), s.render(text))
}

// Fprint writes the styled text to w.
func (s *style) Fprint(w io.Writer, text string) (int, error) {
	return fmt.Fprint(w, s.render(text))
}

// Fprintf formats and writes the styled text to w.
func (s *style) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, s.render(fmt.Sprintf(format, a...)))
}

// Fprintln writes the styled text followed by a newline to w.
func (s *style) Fprintln(w io.Writer, text string) (int, error) {
	return fmt.Fprintln(w, s.render(text))
}

// --- internals ---

func (s *style) render(text string) string {
	if !isEnabled() || len(s.codes) == 0 {
		return text
	}

	// Compute exact size: \x1b[ + code1;code2;... + m + text + \x1b[0m
	size := 2 // \x1b[
	for i, c := range s.codes {
		if i > 0 {
			size++ // ;
		}
		size += len(c)
	}
	size++ // m
	size += len(text)
	size += len(cReset)

	var b strings.Builder
	b.Grow(size)
	b.WriteString("\x1b[")
	for i, c := range s.codes {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(c)
	}
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
