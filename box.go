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
	BorderSimple        = Border{"┌", "┐", "└", "┘", "─", "│"}
	BorderDashed        = Border{"┌", "┐", "└", "┘", "╌", "╎"}
	BorderDotted        = Border{"┌", "┐", "└", "┘", "┈", "┊"}
	BorderRounded       = Border{"╭", "╮", "╰", "╯", "─", "│"}
	BorderRoundedDashed = Border{"╭", "╮", "╰", "╯", "╌", "╎"}
	BorderRoundedDotted = Border{"╭", "╮", "╰", "╯", "┈", "┊"}
	BorderDouble        = Border{"╔", "╗", "╚", "╝", "═", "║"}
	BorderHeavy         = Border{"┏", "┓", "┗", "┛", "━", "┃"}
	BorderASCII         = Border{"+", "+", "+", "+", "-", "|"}
	BorderBlock         = Border{"█", "█", "█", "█", "█", "█"}
	BorderBlockHalf     = Border{"▀", "▀", "▄", "▄", "▀", "▄"}
	BorderBlockLight    = Border{"░", "░", "░", "░", "░", "░"}
	BorderBlockMedium   = Border{"▒", "▒", "▒", "▒", "▒", "▒"}
	BorderBlockDark     = Border{"▓", "▓", "▓", "▓", "▓", "▓"}
)

// BoxStyle holds the configuration for a bordered terminal container.
// Create one with [Box] and chain border, padding, margin, alignment,
// corner controls, and color methods. All fields are unexported to
// preserve immutability; use the provided methods to configure the box.
type BoxStyle struct {
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
	center       bool             // center all content lines horizontally
	centerTrim   bool             // trim whitespace from lines before centering
	centerLines  map[int]struct{} // specific line indices to center (0-based)
	centerFirst  bool             // center the first content line
	centerLast   bool             // center the last content line
	hideTop      bool             // hide the top border row
	hideBottom   bool             // hide the bottom border row
	hideLeft     bool             // hide the left vertical border
	hideRight    bool             // hide the right vertical border
	hideTopLeft  bool             // hide the top-left corner glyph
	hideTopRight bool             // hide the top-right corner glyph
	hideBotLeft  bool             // hide the bottom-left corner glyph
	hideBotRight bool             // hide the bottom-right corner glyph
}

// Box returns a new [BoxStyle] with a simple border and no padding or margin.
func Box() *BoxStyle {
	return &BoxStyle{border: BorderSimple}
}

// copyBox returns a deep copy of the BoxStyle, including the codes slice
// and centerLines map.
func copyBox(b *BoxStyle) *BoxStyle {
	cp := *b
	if len(b.codes) > 0 {
		cp.codes = make([]string, len(b.codes))
		copy(cp.codes, b.codes)
	}
	if len(b.centerLines) > 0 {
		cp.centerLines = make(map[int]struct{}, len(b.centerLines))
		for k, v := range b.centerLines {
			cp.centerLines[k] = v
		}
	}
	return &cp
}

// withCode returns a new BoxStyle with an additional ANSI code appended.
func (b *BoxStyle) withCode(code string) *BoxStyle {
	cp := copyBox(b)
	cp.codes = append(cp.codes, code)
	return cp
}

// --- Border style setters ---

// Border sets a custom border using the provided [Border] struct.
func (b *BoxStyle) Border(border Border) *BoxStyle {
	cp := copyBox(b)
	cp.border = border
	return cp
}

// --- Layout: Padding ---

// Padding sets equal padding on all four sides.
func (b *BoxStyle) Padding(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padTop = n
	cp.padRight = n
	cp.padBottom = n
	cp.padLeft = n
	return cp
}

// PaddingTop sets the top padding.
func (b *BoxStyle) PaddingTop(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padTop = n
	return cp
}

// PaddingBottom sets the bottom padding.
func (b *BoxStyle) PaddingBottom(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padBottom = n
	return cp
}

// PaddingLeft sets the left padding.
func (b *BoxStyle) PaddingLeft(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padLeft = n
	return cp
}

// PaddingRight sets the right padding.
func (b *BoxStyle) PaddingRight(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padRight = n
	return cp
}

// PaddingX sets the left and right padding.
func (b *BoxStyle) PaddingX(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padLeft = n
	cp.padRight = n
	return cp
}

// PaddingY sets the top and bottom padding.
func (b *BoxStyle) PaddingY(n int) *BoxStyle {
	cp := copyBox(b)
	cp.padTop = n
	cp.padBottom = n
	return cp
}

// --- Layout: Margin ---

// Margin sets equal margin on all four sides.
func (b *BoxStyle) Margin(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginTop = n
	cp.marginRight = n
	cp.marginBottom = n
	cp.marginLeft = n
	return cp
}

// MarginTop sets the top margin.
func (b *BoxStyle) MarginTop(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginTop = n
	return cp
}

// MarginBottom sets the bottom margin.
func (b *BoxStyle) MarginBottom(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginBottom = n
	return cp
}

// MarginLeft sets the left margin.
func (b *BoxStyle) MarginLeft(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginLeft = n
	return cp
}

// MarginRight sets the right margin.
func (b *BoxStyle) MarginRight(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginRight = n
	return cp
}

// MarginX sets the left and right margin.
func (b *BoxStyle) MarginX(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginLeft = n
	cp.marginRight = n
	return cp
}

// MarginY sets the top and bottom margin.
func (b *BoxStyle) MarginY(n int) *BoxStyle {
	cp := copyBox(b)
	cp.marginTop = n
	cp.marginBottom = n
	return cp
}

// --- Content alignment ---

// Center enables horizontal centering of content lines within the box.
// Shorter lines are padded equally on both sides to match the widest line.
func (b *BoxStyle) Center() *BoxStyle {
	cp := copyBox(b)
	cp.center = true
	return cp
}

// CenterTrim enables horizontal centering and trims leading/trailing
// whitespace from each line before centering. This is useful when content
// has inconsistent indentation that should be ignored.
func (b *BoxStyle) CenterTrim() *BoxStyle {
	cp := copyBox(b)
	cp.center = true
	cp.centerTrim = true
	return cp
}

// CenterLine marks the line at index n (0-based) for horizontal centering.
// If n is out of bounds at render time, the call is silently ignored.
// This can be called multiple times to center several specific lines.
func (b *BoxStyle) CenterLine(n int) *BoxStyle {
	cp := copyBox(b)
	if cp.centerLines == nil {
		cp.centerLines = make(map[int]struct{})
	}
	cp.centerLines[n] = struct{}{}
	return cp
}

// CenterFirstLine centers the first content line (index 0).
// This is a convenience shortcut useful for centering titles.
func (b *BoxStyle) CenterFirstLine() *BoxStyle {
	cp := copyBox(b)
	cp.centerFirst = true
	return cp
}

// CenterLastLine centers the last content line.
// The line count is determined at render time.
func (b *BoxStyle) CenterLastLine() *BoxStyle {
	cp := copyBox(b)
	cp.centerLast = true
	return cp
}

// --- Side visibility ---

// DisableTop hides the top border row. The vertical borders on content
// rows remain unchanged.
func (b *BoxStyle) DisableTop() *BoxStyle {
	cp := copyBox(b)
	cp.hideTop = true
	return cp
}

// DisableBottom hides the bottom border row.
func (b *BoxStyle) DisableBottom() *BoxStyle {
	cp := copyBox(b)
	cp.hideBottom = true
	return cp
}

// DisableLeft hides the left vertical border on content and padding rows.
// Corner glyphs remain visible unless disabled explicitly.
func (b *BoxStyle) DisableLeft() *BoxStyle {
	cp := copyBox(b)
	cp.hideLeft = true
	return cp
}

// DisableRight hides the right vertical border on all rows.
func (b *BoxStyle) DisableRight() *BoxStyle {
	cp := copyBox(b)
	cp.hideRight = true
	return cp
}

// DisableCorners hides all four corner glyphs.
func (b *BoxStyle) DisableCorners() *BoxStyle {
	cp := copyBox(b)
	cp.hideTopLeft = true
	cp.hideTopRight = true
	cp.hideBotLeft = true
	cp.hideBotRight = true
	return cp
}

// DisableTopLeftCorner hides the top-left corner glyph.
func (b *BoxStyle) DisableTopLeftCorner() *BoxStyle {
	cp := copyBox(b)
	cp.hideTopLeft = true
	return cp
}

// DisableTopRightCorner hides the top-right corner glyph.
func (b *BoxStyle) DisableTopRightCorner() *BoxStyle {
	cp := copyBox(b)
	cp.hideTopRight = true
	return cp
}

// DisableBottomLeftCorner hides the bottom-left corner glyph.
func (b *BoxStyle) DisableBottomLeftCorner() *BoxStyle {
	cp := copyBox(b)
	cp.hideBotLeft = true
	return cp
}

// DisableBottomRightCorner hides the bottom-right corner glyph.
func (b *BoxStyle) DisableBottomRightCorner() *BoxStyle {
	cp := copyBox(b)
	cp.hideBotRight = true
	return cp
}

// --- Colors (border + background) ---

func (b *BoxStyle) OnBlack() *BoxStyle   { return b.withCode(cOnBlack) }
func (b *BoxStyle) OnRed() *BoxStyle     { return b.withCode(cOnRed) }
func (b *BoxStyle) OnGreen() *BoxStyle   { return b.withCode(cOnGreen) }
func (b *BoxStyle) OnYellow() *BoxStyle  { return b.withCode(cOnYellow) }
func (b *BoxStyle) OnBlue() *BoxStyle    { return b.withCode(cOnBlue) }
func (b *BoxStyle) OnMagenta() *BoxStyle { return b.withCode(cOnMagenta) }
func (b *BoxStyle) OnCyan() *BoxStyle    { return b.withCode(cOnCyan) }
func (b *BoxStyle) OnWhite() *BoxStyle   { return b.withCode(cOnWhite) }

func (b *BoxStyle) OnBrightBlack() *BoxStyle   { return b.withCode(cOnBrightBlack) }
func (b *BoxStyle) OnBrightRed() *BoxStyle     { return b.withCode(cOnBrightRed) }
func (b *BoxStyle) OnBrightGreen() *BoxStyle   { return b.withCode(cOnBrightGreen) }
func (b *BoxStyle) OnBrightYellow() *BoxStyle  { return b.withCode(cOnBrightYellow) }
func (b *BoxStyle) OnBrightBlue() *BoxStyle    { return b.withCode(cOnBrightBlue) }
func (b *BoxStyle) OnBrightMagenta() *BoxStyle { return b.withCode(cOnBrightMagenta) }
func (b *BoxStyle) OnBrightCyan() *BoxStyle    { return b.withCode(cOnBrightCyan) }
func (b *BoxStyle) OnBrightWhite() *BoxStyle   { return b.withCode(cOnBrightWhite) }

// Foreground colors for the border glyphs.

func (b *BoxStyle) Black() *BoxStyle   { return b.withCode(cBlack) }
func (b *BoxStyle) Red() *BoxStyle     { return b.withCode(cRed) }
func (b *BoxStyle) Green() *BoxStyle   { return b.withCode(cGreen) }
func (b *BoxStyle) Yellow() *BoxStyle  { return b.withCode(cYellow) }
func (b *BoxStyle) Blue() *BoxStyle    { return b.withCode(cBlue) }
func (b *BoxStyle) Magenta() *BoxStyle { return b.withCode(cMagenta) }
func (b *BoxStyle) Cyan() *BoxStyle    { return b.withCode(cCyan) }
func (b *BoxStyle) White() *BoxStyle   { return b.withCode(cWhite) }

func (b *BoxStyle) BrightBlack() *BoxStyle   { return b.withCode(cBrightBlack) }
func (b *BoxStyle) BrightRed() *BoxStyle     { return b.withCode(cBrightRed) }
func (b *BoxStyle) BrightGreen() *BoxStyle   { return b.withCode(cBrightGreen) }
func (b *BoxStyle) BrightYellow() *BoxStyle  { return b.withCode(cBrightYellow) }
func (b *BoxStyle) BrightBlue() *BoxStyle    { return b.withCode(cBrightBlue) }
func (b *BoxStyle) BrightMagenta() *BoxStyle { return b.withCode(cBrightMagenta) }
func (b *BoxStyle) BrightCyan() *BoxStyle    { return b.withCode(cBrightCyan) }
func (b *BoxStyle) BrightWhite() *BoxStyle   { return b.withCode(cBrightWhite) }

// Modifiers for the border style.

func (b *BoxStyle) Bold() *BoxStyle { return b.withCode(cBold) }
func (b *BoxStyle) Dim() *BoxStyle  { return b.withCode(cDim) }

// --- Output methods ---

// String renders the box around the given content and returns the result.
func (b *BoxStyle) String(content string) string {
	return b.render(content)
}

// Sprintf formats the content, renders it inside the box, and returns the result.
func (b *BoxStyle) Sprintf(format string, a ...any) string {
	return b.render(fmt.Sprintf(format, a...))
}

// Print renders the box and writes it to the default output.
func (b *BoxStyle) Print(content string) {
	_, _ = fmt.Fprint(getOutput(), b.render(content))
}

// Printf formats the content, renders it inside the box, and writes to the default output.
func (b *BoxStyle) Printf(format string, a ...any) {
	_, _ = fmt.Fprint(getOutput(), b.render(fmt.Sprintf(format, a...)))
}

// Println renders the box and writes it followed by a newline to the default output.
func (b *BoxStyle) Println(content string) {
	_, _ = fmt.Fprintln(getOutput(), b.render(content))
}

// Fprint renders the box and writes it to w.
func (b *BoxStyle) Fprint(w io.Writer, content string) (int, error) {
	return fmt.Fprint(w, b.render(content))
}

// Fprintf formats the content, renders it inside the box, and writes to w.
func (b *BoxStyle) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, b.render(fmt.Sprintf(format, a...)))
}

// Fprintln renders the box and writes it followed by a newline to w.
func (b *BoxStyle) Fprintln(w io.Writer, content string) (int, error) {
	return fmt.Fprintln(w, b.render(content))
}

// --- Internals ---

// wrapCodes wraps s in the given ANSI SGR codes. Returns s unchanged if
// colors are disabled or codes is empty.
func wrapCodes(s string, codes []string) string {
	if !isEnabled() || len(codes) == 0 {
		return s
	}
	size := 2
	for i, c := range codes {
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
	for i, c := range codes {
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

// wrapStyle wraps s in the box's border/background ANSI codes.
func (b *BoxStyle) wrapStyle(s string) string {
	return wrapCodes(s, b.codes)
}

// render builds the full box frame around content.
func (b *BoxStyle) render(content string) string {
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

	// Determine glyph replacements for disabled sides.
	leftVert := b.border.Vertical
	rightVert := b.border.Vertical
	if b.hideLeft {
		leftVert = strings.Repeat(" ", visibleWidth(b.border.Vertical))
	}
	if b.hideRight {
		rightVert = strings.Repeat(" ", visibleWidth(b.border.Vertical))
	}

	// Collect box rows (without margin, without trailing \n).
	var boxRows []string

	frameW := visibleWidth(b.border.Vertical) + innerW + visibleWidth(b.border.Vertical)

	// Top border.
	if !b.hideTop {
		tl := b.border.TopLeft
		tr := b.border.TopRight
		if b.hideTopLeft {
			tl = strings.Repeat(" ", visibleWidth(b.border.TopLeft))
		}
		if b.hideTopRight {
			tr = strings.Repeat(" ", visibleWidth(b.border.TopRight))
		}

		topFillW := frameW - visibleWidth(tl) - visibleWidth(tr)
		if topFillW < 0 {
			topFillW = 0
		}
		topCount := topFillW
		horW := visibleWidth(b.border.Horizontal)
		if horW > 0 {
			topCount = topFillW / horW
		}

		topBar := tl + strings.Repeat(b.border.Horizontal, topCount) + tr
		boxRows = append(boxRows, b.wrapStyle(topBar))
	}

	// Top padding rows.
	for i := 0; i < b.padTop; i++ {
		padLine := leftVert + strings.Repeat(" ", innerW) + rightVert
		boxRows = append(boxRows, b.wrapStyle(padLine))
	}

	// Content rows.
	lastIdx := len(lines) - 1
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		vis := visibleWidth(line)
		availW := innerW - b.padLeft - b.padRight

		shouldCenter := b.center
		if !shouldCenter {
			_, inSet := b.centerLines[i]
			shouldCenter = inSet
		}
		if !shouldCenter && b.centerFirst && i == 0 {
			shouldCenter = true
		}
		if !shouldCenter && b.centerLast && i == lastIdx {
			shouldCenter = true
		}

		var leftPad, rightPad int
		if shouldCenter && vis < availW {
			total := availW - vis
			leftPad = total / 2
			rightPad = total - leftPad
		} else {
			rightPad = availW - vis
			if rightPad < 0 {
				rightPad = 0
			}
		}

		// Chrome parts wrapped individually to prevent nested ANSI corruption.
		row := b.wrapStyle(leftVert+strings.Repeat(" ", b.padLeft+leftPad)) +
			line +
			b.wrapStyle(strings.Repeat(" ", rightPad+b.padRight)+rightVert)
		boxRows = append(boxRows, row)
	}

	// Bottom padding rows.
	for i := 0; i < b.padBottom; i++ {
		padLine := leftVert + strings.Repeat(" ", innerW) + rightVert
		boxRows = append(boxRows, b.wrapStyle(padLine))
	}

	// Bottom border.
	if !b.hideBottom {
		bl := b.border.BottomLeft
		br := b.border.BottomRight
		if b.hideBotLeft {
			bl = strings.Repeat(" ", visibleWidth(b.border.BottomLeft))
		}
		if b.hideBotRight {
			br = strings.Repeat(" ", visibleWidth(b.border.BottomRight))
		}

		botFillW := frameW - visibleWidth(bl) - visibleWidth(br)
		if botFillW < 0 {
			botFillW = 0
		}
		botCount := botFillW
		horW := visibleWidth(b.border.Horizontal)
		if horW > 0 {
			botCount = botFillW / horW
		}

		botBar := bl + strings.Repeat(b.border.Horizontal, botCount) + br
		boxRows = append(boxRows, b.wrapStyle(botBar))
	}

	// Track the index of the bottom border row (if present). This determines
	// trailing-newline behavior.
	bottomBorderIdx := -1
	if !b.hideBottom {
		bottomBorderIdx = len(boxRows) - 1
	}

	// Assemble final output with margins.
	marginLeft := strings.Repeat(" ", b.marginLeft)
	marginRight := strings.Repeat(" ", b.marginRight)

	var out strings.Builder

	for i := 0; i < b.marginTop; i++ {
		out.WriteByte('\n')
	}

	for i := 0; i < len(boxRows); i++ {
		out.WriteString(marginLeft)
		out.WriteString(boxRows[i])
		out.WriteString(marginRight)
		// The final rendered row never gets a trailing \n.
		// The bottom border (if present) is final.
		// When the bottom border is hidden,
		// the last content/padding row still gets \n (legacy behavior).
		if bottomBorderIdx >= 0 && i == bottomBorderIdx {
			// Bottom border is present and this IS it — no trailing \n.
		} else {
			out.WriteByte('\n')
		}
	}

	for i := 0; i < b.marginBottom; i++ {
		out.WriteByte('\n')
	}

	return out.String()
}
