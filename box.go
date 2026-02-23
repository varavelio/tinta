package tinta

import (
	"fmt"
	"io"
	"strings"
)

// Border defines the glyphs used to draw a box frame.
type Border struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

// Predefined border styles.
var (
	BorderSimple  = Border{"┌", "┐", "└", "┘", "─", "│"}
	BorderRounded = Border{"╭", "╮", "╰", "╯", "─", "│"}
	BorderDouble  = Border{"╔", "╗", "╚", "╝", "═", "║"}
	BorderHeavy   = Border{"┏", "┓", "┗", "┛", "━", "┃"}
)

// box is intentionally unexported. Users interact with it through the
// exported BoxStyle type alias and the [Box] constructor.
type box struct {
	border       Border
	codes        []string // ANSI SGR codes for the border/background
	padTop       int
	padRight     int
	padBottom    int
	padLeft      int
	marginTop    int
	marginRight  int
	marginBottom int
	marginLeft   int
	center       bool // center content lines horizontally
	centerTrim   bool // trim whitespace from lines before centering
}

// BoxStyle is the public handle returned by [Box] and every chaining method.
// The underlying struct is opaque; users cannot create one manually.
type BoxStyle = *box

// Box returns a new BoxStyle with a simple border and no padding or margin.
func Box() BoxStyle {
	return &box{border: BorderSimple}
}

// copyBox returns a deep copy of the box, including the codes slice.
func copyBox(b *box) *box {
	cp := *b
	if len(b.codes) > 0 {
		cp.codes = make([]string, len(b.codes))
		copy(cp.codes, b.codes)
	}
	return &cp
}

// withCode returns a new box with an additional ANSI code appended.
func (b *box) withCode(code string) BoxStyle {
	cp := copyBox(b)
	cp.codes = append(cp.codes, code)
	return cp
}

// --- Border style setters ---

// Simple sets the border to a simple line style (┌─┐).
func (b *box) Simple() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderSimple
	return cp
}

// Rounded sets the border to a rounded style (╭─╮).
func (b *box) Rounded() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderRounded
	return cp
}

// Double sets the border to a double line style (╔═╗).
func (b *box) Double() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderDouble
	return cp
}

// Heavy sets the border to a heavy line style (┏━┓).
func (b *box) Heavy() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderHeavy
	return cp
}

// --- Layout: Padding ---

// Padding sets equal padding on all four sides.
func (b *box) Padding(n int) BoxStyle {
	cp := copyBox(b)
	cp.padTop = n
	cp.padRight = n
	cp.padBottom = n
	cp.padLeft = n
	return cp
}

// PaddingTop sets the top padding.
func (b *box) PaddingTop(n int) BoxStyle {
	cp := copyBox(b)
	cp.padTop = n
	return cp
}

// PaddingBottom sets the bottom padding.
func (b *box) PaddingBottom(n int) BoxStyle {
	cp := copyBox(b)
	cp.padBottom = n
	return cp
}

// PaddingLeft sets the left padding.
func (b *box) PaddingLeft(n int) BoxStyle {
	cp := copyBox(b)
	cp.padLeft = n
	return cp
}

// PaddingRight sets the right padding.
func (b *box) PaddingRight(n int) BoxStyle {
	cp := copyBox(b)
	cp.padRight = n
	return cp
}

// PaddingX sets the left and right padding.
func (b *box) PaddingX(n int) BoxStyle {
	cp := copyBox(b)
	cp.padLeft = n
	cp.padRight = n
	return cp
}

// PaddingY sets the top and bottom padding.
func (b *box) PaddingY(n int) BoxStyle {
	cp := copyBox(b)
	cp.padTop = n
	cp.padBottom = n
	return cp
}

// --- Layout: Margin ---

// Margin sets equal margin on all four sides.
func (b *box) Margin(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginTop = n
	cp.marginRight = n
	cp.marginBottom = n
	cp.marginLeft = n
	return cp
}

// MarginTop sets the top margin.
func (b *box) MarginTop(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginTop = n
	return cp
}

// MarginBottom sets the bottom margin.
func (b *box) MarginBottom(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginBottom = n
	return cp
}

// MarginLeft sets the left margin.
func (b *box) MarginLeft(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginLeft = n
	return cp
}

// MarginRight sets the right margin.
func (b *box) MarginRight(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginRight = n
	return cp
}

// MarginX sets the left and right margin.
func (b *box) MarginX(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginLeft = n
	cp.marginRight = n
	return cp
}

// MarginY sets the top and bottom margin.
func (b *box) MarginY(n int) BoxStyle {
	cp := copyBox(b)
	cp.marginTop = n
	cp.marginBottom = n
	return cp
}

// --- Content alignment ---

// Center enables horizontal centering of content lines within the box.
// Shorter lines are padded equally on both sides to match the widest line.
func (b *box) Center() BoxStyle {
	cp := copyBox(b)
	cp.center = true
	return cp
}

// CenterTrim enables horizontal centering and trims leading/trailing
// whitespace from each line before centering. This is useful when content
// has inconsistent indentation that should be ignored.
func (b *box) CenterTrim() BoxStyle {
	cp := copyBox(b)
	cp.center = true
	cp.centerTrim = true
	return cp
}

// --- Colors (border + background) ---

func (b *box) OnBlack() BoxStyle   { return b.withCode(cOnBlack) }
func (b *box) OnRed() BoxStyle     { return b.withCode(cOnRed) }
func (b *box) OnGreen() BoxStyle   { return b.withCode(cOnGreen) }
func (b *box) OnYellow() BoxStyle  { return b.withCode(cOnYellow) }
func (b *box) OnBlue() BoxStyle    { return b.withCode(cOnBlue) }
func (b *box) OnMagenta() BoxStyle { return b.withCode(cOnMagenta) }
func (b *box) OnCyan() BoxStyle    { return b.withCode(cOnCyan) }
func (b *box) OnWhite() BoxStyle   { return b.withCode(cOnWhite) }

func (b *box) OnBrightBlack() BoxStyle   { return b.withCode(cOnBrightBlack) }
func (b *box) OnBrightRed() BoxStyle     { return b.withCode(cOnBrightRed) }
func (b *box) OnBrightGreen() BoxStyle   { return b.withCode(cOnBrightGreen) }
func (b *box) OnBrightYellow() BoxStyle  { return b.withCode(cOnBrightYellow) }
func (b *box) OnBrightBlue() BoxStyle    { return b.withCode(cOnBrightBlue) }
func (b *box) OnBrightMagenta() BoxStyle { return b.withCode(cOnBrightMagenta) }
func (b *box) OnBrightCyan() BoxStyle    { return b.withCode(cOnBrightCyan) }
func (b *box) OnBrightWhite() BoxStyle   { return b.withCode(cOnBrightWhite) }

// Foreground colors for the border glyphs.

func (b *box) Black() BoxStyle   { return b.withCode(cBlack) }
func (b *box) Red() BoxStyle     { return b.withCode(cRed) }
func (b *box) Green() BoxStyle   { return b.withCode(cGreen) }
func (b *box) Yellow() BoxStyle  { return b.withCode(cYellow) }
func (b *box) Blue() BoxStyle    { return b.withCode(cBlue) }
func (b *box) Magenta() BoxStyle { return b.withCode(cMagenta) }
func (b *box) Cyan() BoxStyle    { return b.withCode(cCyan) }
func (b *box) White() BoxStyle   { return b.withCode(cWhite) }

func (b *box) BrightBlack() BoxStyle   { return b.withCode(cBrightBlack) }
func (b *box) BrightRed() BoxStyle     { return b.withCode(cBrightRed) }
func (b *box) BrightGreen() BoxStyle   { return b.withCode(cBrightGreen) }
func (b *box) BrightYellow() BoxStyle  { return b.withCode(cBrightYellow) }
func (b *box) BrightBlue() BoxStyle    { return b.withCode(cBrightBlue) }
func (b *box) BrightMagenta() BoxStyle { return b.withCode(cBrightMagenta) }
func (b *box) BrightCyan() BoxStyle    { return b.withCode(cBrightCyan) }
func (b *box) BrightWhite() BoxStyle   { return b.withCode(cBrightWhite) }

// Modifiers for the border style.

func (b *box) Bold() BoxStyle { return b.withCode(cBold) }
func (b *box) Dim() BoxStyle  { return b.withCode(cDim) }

// --- Output methods ---

// String renders the box around the given content and returns the result.
func (b *box) String(content string) string {
	return b.render(content)
}

// Sprintf formats the content, renders it inside the box, and returns the result.
func (b *box) Sprintf(format string, a ...any) string {
	return b.render(fmt.Sprintf(format, a...))
}

// Print renders the box and writes it to the default output.
func (b *box) Print(content string) {
	_, _ = fmt.Fprint(getOutput(), b.render(content))
}

// Printf formats the content, renders it inside the box, and writes to the default output.
func (b *box) Printf(format string, a ...any) {
	_, _ = fmt.Fprint(getOutput(), b.render(fmt.Sprintf(format, a...)))
}

// Println renders the box and writes it followed by a newline to the default output.
func (b *box) Println(content string) {
	_, _ = fmt.Fprintln(getOutput(), b.render(content))
}

// Fprint renders the box and writes it to w.
func (b *box) Fprint(w io.Writer, content string) (int, error) {
	return fmt.Fprint(w, b.render(content))
}

// Fprintf formats the content, renders it inside the box, and writes to w.
func (b *box) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, b.render(fmt.Sprintf(format, a...)))
}

// Fprintln renders the box and writes it followed by a newline to w.
func (b *box) Fprintln(w io.Writer, content string) (int, error) {
	return fmt.Fprintln(w, b.render(content))
}

// --- Internals ---

// wrapStyle wraps s in ANSI codes if colors are enabled and codes are present.
func (b *box) wrapStyle(s string) string {
	if !isEnabled() || len(b.codes) == 0 {
		return s
	}
	size := 2
	for i, c := range b.codes {
		if i > 0 {
			size++
		}
		size += len(c)
	}
	size++ // m
	size += len(s)
	size += len(cReset)

	var buf strings.Builder
	buf.Grow(size)
	buf.WriteString("\x1b[")
	for i, c := range b.codes {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteString(c)
	}
	buf.WriteByte('m')
	buf.WriteString(s)
	buf.WriteString(cReset)
	return buf.String()
}

// render builds the full box frame around content.
func (b *box) render(content string) string {
	lines := strings.Split(content, "\n")

	// Apply trim if CenterTrim is active.
	if b.centerTrim {
		for i, line := range lines {
			lines[i] = strings.TrimSpace(line)
		}
	}

	// Find the widest visible line.
	maxW := 0
	for _, line := range lines {
		w := visibleWidth(line)
		if w > maxW {
			maxW = w
		}
	}

	// Inner width = content width + horizontal padding.
	innerW := maxW + b.padLeft + b.padRight

	marginLeft := strings.Repeat(" ", b.marginLeft)
	marginRight := strings.Repeat(" ", b.marginRight)

	var out strings.Builder

	// Top margin.
	for i := 0; i < b.marginTop; i++ {
		out.WriteByte('\n')
	}

	// Top border: ┌───┐
	topBar := b.border.TopLeft + strings.Repeat(b.border.Horizontal, innerW) + b.border.TopRight
	out.WriteString(marginLeft)
	out.WriteString(b.wrapStyle(topBar))
	out.WriteString(marginRight)
	out.WriteByte('\n')

	// Top padding rows.
	for i := 0; i < b.padTop; i++ {
		padLine := b.border.Vertical + strings.Repeat(" ", innerW) + b.border.Vertical
		out.WriteString(marginLeft)
		out.WriteString(b.wrapStyle(padLine))
		out.WriteString(marginRight)
		out.WriteByte('\n')
	}

	// Content rows.
	for _, line := range lines {
		vis := visibleWidth(line)
		availW := innerW - b.padLeft - b.padRight

		var leftPad, rightPad int
		if b.center && vis < availW {
			// Center within the available content area.
			total := availW - vis
			leftPad = total / 2
			rightPad = total - leftPad
		} else {
			rightPad = availW - vis
			if rightPad < 0 {
				rightPad = 0
			}
		}

		row := b.border.Vertical +
			strings.Repeat(" ", b.padLeft+leftPad) +
			line +
			strings.Repeat(" ", rightPad+b.padRight) +
			b.border.Vertical
		out.WriteString(marginLeft)
		out.WriteString(b.wrapStyle(row))
		out.WriteString(marginRight)
		out.WriteByte('\n')
	}

	// Bottom padding rows.
	for i := 0; i < b.padBottom; i++ {
		padLine := b.border.Vertical + strings.Repeat(" ", innerW) + b.border.Vertical
		out.WriteString(marginLeft)
		out.WriteString(b.wrapStyle(padLine))
		out.WriteString(marginRight)
		out.WriteByte('\n')
	}

	// Bottom border: └───┘
	botBar := b.border.BottomLeft + strings.Repeat(b.border.Horizontal, innerW) + b.border.BottomRight
	out.WriteString(marginLeft)
	out.WriteString(b.wrapStyle(botBar))
	out.WriteString(marginRight)

	// Bottom margin.
	for i := 0; i < b.marginBottom; i++ {
		out.WriteByte('\n')
	}

	return out.String()
}
