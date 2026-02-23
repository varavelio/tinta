package tinta

import (
	"fmt"
	"io"
	"strings"
)

// ANSI SGR constants.
const (
	cReset = "\x1b[0m"

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

// text is intentionally unexported. Users interact with it through the
// exported TextStyle type alias and the package-level constructors.
type text struct {
	codes []string
}

// TextStyle is the public handle returned by [Text] and every chaining method.
// The underlying struct is opaque; users cannot create one manually.
type TextStyle = *text

// Text returns a new TextStyle with no codes. Use it as the single entry
// point for building styled output: tinta.Text().Red().Bold().Println("hello").
func Text() TextStyle {
	return &text{}
}

// with returns a new text that has all existing codes plus one more.
// It copies the slice to guarantee immutability of the source.
func (t *text) with(code string) TextStyle {
	cp := make([]string, len(t.codes)+1)
	copy(cp, t.codes)
	cp[len(t.codes)] = code
	return &text{codes: cp}
}

// --- Foreground colors ---

func (t *text) Black() TextStyle   { return t.with(cBlack) }
func (t *text) Red() TextStyle     { return t.with(cRed) }
func (t *text) Green() TextStyle   { return t.with(cGreen) }
func (t *text) Yellow() TextStyle  { return t.with(cYellow) }
func (t *text) Blue() TextStyle    { return t.with(cBlue) }
func (t *text) Magenta() TextStyle { return t.with(cMagenta) }
func (t *text) Cyan() TextStyle    { return t.with(cCyan) }
func (t *text) White() TextStyle   { return t.with(cWhite) }

func (t *text) BrightBlack() TextStyle   { return t.with(cBrightBlack) }
func (t *text) BrightRed() TextStyle     { return t.with(cBrightRed) }
func (t *text) BrightGreen() TextStyle   { return t.with(cBrightGreen) }
func (t *text) BrightYellow() TextStyle  { return t.with(cBrightYellow) }
func (t *text) BrightBlue() TextStyle    { return t.with(cBrightBlue) }
func (t *text) BrightMagenta() TextStyle { return t.with(cBrightMagenta) }
func (t *text) BrightCyan() TextStyle    { return t.with(cBrightCyan) }
func (t *text) BrightWhite() TextStyle   { return t.with(cBrightWhite) }

// --- Background colors (On*) ---

func (t *text) OnBlack() TextStyle   { return t.with(cOnBlack) }
func (t *text) OnRed() TextStyle     { return t.with(cOnRed) }
func (t *text) OnGreen() TextStyle   { return t.with(cOnGreen) }
func (t *text) OnYellow() TextStyle  { return t.with(cOnYellow) }
func (t *text) OnBlue() TextStyle    { return t.with(cOnBlue) }
func (t *text) OnMagenta() TextStyle { return t.with(cOnMagenta) }
func (t *text) OnCyan() TextStyle    { return t.with(cOnCyan) }
func (t *text) OnWhite() TextStyle   { return t.with(cOnWhite) }

func (t *text) OnBrightBlack() TextStyle   { return t.with(cOnBrightBlack) }
func (t *text) OnBrightRed() TextStyle     { return t.with(cOnBrightRed) }
func (t *text) OnBrightGreen() TextStyle   { return t.with(cOnBrightGreen) }
func (t *text) OnBrightYellow() TextStyle  { return t.with(cOnBrightYellow) }
func (t *text) OnBrightBlue() TextStyle    { return t.with(cOnBrightBlue) }
func (t *text) OnBrightMagenta() TextStyle { return t.with(cOnBrightMagenta) }
func (t *text) OnBrightCyan() TextStyle    { return t.with(cOnBrightCyan) }
func (t *text) OnBrightWhite() TextStyle   { return t.with(cOnBrightWhite) }

// --- Modifiers ---

func (t *text) Bold() TextStyle      { return t.with(cBold) }
func (t *text) Dim() TextStyle       { return t.with(cDim) }
func (t *text) Italic() TextStyle    { return t.with(cItalic) }
func (t *text) Underline() TextStyle { return t.with(cUnderline) }
func (t *text) Invert() TextStyle    { return t.with(cInvert) }
func (t *text) Hidden() TextStyle    { return t.with(cHidden) }
func (t *text) Strike() TextStyle    { return t.with(cStrike) }

// --- Output methods ---

// String returns the styled text.
func (t *text) String(s string) string {
	return t.render(s)
}

// Sprintf formats and returns the styled text.
func (t *text) Sprintf(format string, a ...any) string {
	return t.render(fmt.Sprintf(format, a...))
}

// Print writes the styled text to the default output.
func (t *text) Print(s string) {
	_, _ = fmt.Fprint(getOutput(), t.render(s))
}

// Printf formats and writes the styled text to the default output.
func (t *text) Printf(format string, a ...any) {
	_, _ = fmt.Fprint(getOutput(), t.render(fmt.Sprintf(format, a...)))
}

// Println writes the styled text followed by a newline to the default output.
func (t *text) Println(s string) {
	_, _ = fmt.Fprintln(getOutput(), t.render(s))
}

// Fprint writes the styled text to w.
func (t *text) Fprint(w io.Writer, s string) (int, error) {
	return fmt.Fprint(w, t.render(s))
}

// Fprintf formats and writes the styled text to w.
func (t *text) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, t.render(fmt.Sprintf(format, a...)))
}

// Fprintln writes the styled text followed by a newline to w.
func (t *text) Fprintln(w io.Writer, s string) (int, error) {
	return fmt.Fprintln(w, t.render(s))
}

// --- Internals ---

func (t *text) render(s string) string {
	if !isEnabled() || len(t.codes) == 0 {
		return s
	}

	// Compute exact size: \x1b[ + code1;code2;... + m + text + \x1b[0m
	size := 2 // \x1b[
	for i, c := range t.codes {
		if i > 0 {
			size++ // ;
		}
		size += len(c)
	}
	size++ // m
	size += len(s)
	size += len(cReset)

	var b strings.Builder
	b.Grow(size)
	b.WriteString("\x1b[")
	for i, c := range t.codes {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(c)
	}
	b.WriteByte('m')
	b.WriteString(s)
	b.WriteString(cReset)
	return b.String()
}
