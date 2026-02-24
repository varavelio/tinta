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

// TextStyle holds ANSI SGR codes for styled terminal text.
// Create one with [Text] and chain color/modifier methods.
// All fields are unexported to preserve immutability.
type TextStyle struct {
	codes []string
}

// Text returns a new [TextStyle] with no codes. Use it as the single entry
// point for building styled output: tinta.Text().Red().Bold().Println("hello").
func Text() *TextStyle {
	return &TextStyle{}
}

// with returns a new TextStyle that has all existing codes plus one more.
// It copies the slice to guarantee immutability of the source.
func (t *TextStyle) with(code string) *TextStyle {
	cp := make([]string, len(t.codes)+1)
	copy(cp, t.codes)
	cp[len(t.codes)] = code
	return &TextStyle{codes: cp}
}

// --- Foreground colors ---

func (t *TextStyle) Black() *TextStyle   { return t.with(cBlack) }
func (t *TextStyle) Red() *TextStyle     { return t.with(cRed) }
func (t *TextStyle) Green() *TextStyle   { return t.with(cGreen) }
func (t *TextStyle) Yellow() *TextStyle  { return t.with(cYellow) }
func (t *TextStyle) Blue() *TextStyle    { return t.with(cBlue) }
func (t *TextStyle) Magenta() *TextStyle { return t.with(cMagenta) }
func (t *TextStyle) Cyan() *TextStyle    { return t.with(cCyan) }
func (t *TextStyle) White() *TextStyle   { return t.with(cWhite) }

func (t *TextStyle) BrightBlack() *TextStyle   { return t.with(cBrightBlack) }
func (t *TextStyle) BrightRed() *TextStyle     { return t.with(cBrightRed) }
func (t *TextStyle) BrightGreen() *TextStyle   { return t.with(cBrightGreen) }
func (t *TextStyle) BrightYellow() *TextStyle  { return t.with(cBrightYellow) }
func (t *TextStyle) BrightBlue() *TextStyle    { return t.with(cBrightBlue) }
func (t *TextStyle) BrightMagenta() *TextStyle { return t.with(cBrightMagenta) }
func (t *TextStyle) BrightCyan() *TextStyle    { return t.with(cBrightCyan) }
func (t *TextStyle) BrightWhite() *TextStyle   { return t.with(cBrightWhite) }

// --- Background colors (On*) ---

func (t *TextStyle) OnBlack() *TextStyle   { return t.with(cOnBlack) }
func (t *TextStyle) OnRed() *TextStyle     { return t.with(cOnRed) }
func (t *TextStyle) OnGreen() *TextStyle   { return t.with(cOnGreen) }
func (t *TextStyle) OnYellow() *TextStyle  { return t.with(cOnYellow) }
func (t *TextStyle) OnBlue() *TextStyle    { return t.with(cOnBlue) }
func (t *TextStyle) OnMagenta() *TextStyle { return t.with(cOnMagenta) }
func (t *TextStyle) OnCyan() *TextStyle    { return t.with(cOnCyan) }
func (t *TextStyle) OnWhite() *TextStyle   { return t.with(cOnWhite) }

func (t *TextStyle) OnBrightBlack() *TextStyle   { return t.with(cOnBrightBlack) }
func (t *TextStyle) OnBrightRed() *TextStyle     { return t.with(cOnBrightRed) }
func (t *TextStyle) OnBrightGreen() *TextStyle   { return t.with(cOnBrightGreen) }
func (t *TextStyle) OnBrightYellow() *TextStyle  { return t.with(cOnBrightYellow) }
func (t *TextStyle) OnBrightBlue() *TextStyle    { return t.with(cOnBrightBlue) }
func (t *TextStyle) OnBrightMagenta() *TextStyle { return t.with(cOnBrightMagenta) }
func (t *TextStyle) OnBrightCyan() *TextStyle    { return t.with(cOnBrightCyan) }
func (t *TextStyle) OnBrightWhite() *TextStyle   { return t.with(cOnBrightWhite) }

// --- Modifiers ---

func (t *TextStyle) Bold() *TextStyle      { return t.with(cBold) }
func (t *TextStyle) Dim() *TextStyle       { return t.with(cDim) }
func (t *TextStyle) Italic() *TextStyle    { return t.with(cItalic) }
func (t *TextStyle) Underline() *TextStyle { return t.with(cUnderline) }
func (t *TextStyle) Invert() *TextStyle    { return t.with(cInvert) }
func (t *TextStyle) Hidden() *TextStyle    { return t.with(cHidden) }
func (t *TextStyle) Strike() *TextStyle    { return t.with(cStrike) }

// --- Output methods ---

// String returns the styled text.
func (t *TextStyle) String(s string) string {
	return t.render(s)
}

// Sprintf formats and returns the styled text.
func (t *TextStyle) Sprintf(format string, a ...any) string {
	return t.render(fmt.Sprintf(format, a...))
}

// Print writes the styled text to the default output.
func (t *TextStyle) Print(s string) {
	_, _ = fmt.Fprint(getOutput(), t.render(s))
}

// Printf formats and writes the styled text to the default output.
func (t *TextStyle) Printf(format string, a ...any) {
	_, _ = fmt.Fprint(getOutput(), t.render(fmt.Sprintf(format, a...)))
}

// Println writes the styled text followed by a newline to the default output.
func (t *TextStyle) Println(s string) {
	_, _ = fmt.Fprintln(getOutput(), t.render(s))
}

// Fprint writes the styled text to w.
func (t *TextStyle) Fprint(w io.Writer, s string) (int, error) {
	return fmt.Fprint(w, t.render(s))
}

// Fprintf formats and writes the styled text to w.
func (t *TextStyle) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, t.render(fmt.Sprintf(format, a...)))
}

// Fprintln writes the styled text followed by a newline to w.
func (t *TextStyle) Fprintln(w io.Writer, s string) (int, error) {
	return fmt.Fprintln(w, t.render(s))
}

// --- Internals ---

func (t *TextStyle) render(s string) string {
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
