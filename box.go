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

// Align controls the horizontal alignment of title and footer text
// within a box border.
type Align int

const (
	// AlignLeft places text at the left edge of the border (after the corner).
	AlignLeft Align = iota
	// AlignCenter places text at the horizontal center of the border.
	AlignCenter
	// AlignRight places text at the right edge of the border (before the corner).
	AlignRight
)

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
	codes        []string
	padTop       int
	padRight     int
	padBottom    int
	padLeft      int
	marginTop    int
	marginRight  int
	marginBottom int
	marginLeft   int
	center       bool
	centerTrim   bool
	centerLines  map[int]struct{}
	centerFirst  bool
	centerLast   bool
	hideTop      bool
	hideBottom   bool
	hideLeft     bool
	hideRight    bool
	hideTopLeft  bool
	hideTopRight bool
	hideBotLeft  bool
	hideBotRight bool
	title        string
	titleAlign   Align
	footer       string
	footerAlign  Align
}

// Box returns a new [BoxStyle] with a simple border and no padding or margin.
func Box() *BoxStyle {
	return &BoxStyle{border: BorderSimple}
}

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

func (b *BoxStyle) withCode(code string) *BoxStyle {
	cp := copyBox(b)
	cp.codes = append(cp.codes, code)
	return cp
}

// Border sets a custom border using the provided [Border] struct.
func (b *BoxStyle) Border(border Border) *BoxStyle {
	cp := copyBox(b)
	cp.border = border
	return cp
}

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

// Title sets text to render inside the top border row. The align
// parameter controls horizontal placement: [AlignLeft], [AlignCenter],
// or [AlignRight]. The title is separated from corners by one
// horizontal glyph on each side for visual balance.
func (b *BoxStyle) Title(text string, align Align) *BoxStyle {
	cp := copyBox(b)
	cp.title = text
	cp.titleAlign = align
	return cp
}

// Footer sets text to render inside the bottom border row. The align
// parameter controls horizontal placement: [AlignLeft], [AlignCenter],
// or [AlignRight]. The footer is separated from corners by one
// horizontal glyph on each side for visual balance.
func (b *BoxStyle) Footer(text string, align Align) *BoxStyle {
	cp := copyBox(b)
	cp.footer = text
	cp.footerAlign = align
	return cp
}

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

func (b *BoxStyle) Bold() *BoxStyle { return b.withCode(cBold) }
func (b *BoxStyle) Dim() *BoxStyle  { return b.withCode(cDim) }

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

func (b *BoxStyle) wrapStyle(s string) string {
	return wrapCodes(s, b.codes)
}

func (b *BoxStyle) buildBorderRow(cornerLeft, cornerRight string, hideLeft, hideRight bool, text string, align Align, frameW int) string {
	cl := cornerLeft
	cr := cornerRight
	if hideLeft {
		cl = strings.Repeat(" ", visibleWidth(cornerLeft))
	}
	if hideRight {
		cr = strings.Repeat(" ", visibleWidth(cornerRight))
	}

	fillW := frameW - visibleWidth(cl) - visibleWidth(cr)
	if fillW < 0 {
		fillW = 0
	}

	horW := visibleWidth(b.border.Horizontal)
	if horW == 0 {
		horW = 1
	}

	if text == "" {
		return cl + strings.Repeat(b.border.Horizontal, fillW/horW) + cr
	}

	textW := visibleWidth(text)
	minNeeded := horW + textW + horW
	if fillW < minNeeded {
		return cl + strings.Repeat(b.border.Horizontal, fillW/horW) + cr
	}

	remaining := fillW - textW
	var leftGlyphs, rightGlyphs int

	switch align {
	case AlignCenter:
		leftGlyphs = remaining / 2 / horW
		rightGlyphs = (remaining - leftGlyphs*horW) / horW
	case AlignRight:
		rightGlyphs = 1
		leftGlyphs = (remaining - rightGlyphs*horW) / horW
	default: // AlignLeft
		leftGlyphs = 1
		rightGlyphs = (remaining - leftGlyphs*horW) / horW
	}

	return cl +
		strings.Repeat(b.border.Horizontal, leftGlyphs) +
		text +
		strings.Repeat(b.border.Horizontal, rightGlyphs) +
		cr
}

func (b *BoxStyle) render(content string) string {
	lines := strings.Split(content, "\n")

	if b.centerTrim {
		for i, line := range lines {
			lines[i] = strings.TrimSpace(line)
		}
	}

	maxW := 0
	for _, line := range lines {
		w := visibleWidth(line)
		if w > maxW {
			maxW = w
		}
	}

	innerW := maxW + b.padLeft + b.padRight

	horW := visibleWidth(b.border.Horizontal)
	if horW == 0 {
		horW = 1
	}
	clW := visibleWidth(b.border.TopLeft)
	crW := visibleWidth(b.border.TopRight)
	if b.title != "" {
		needed := visibleWidth(b.title) + 2*horW
		vertSum := visibleWidth(b.border.Vertical)*2 - clW - crW
		minInner := needed - vertSum
		if minInner > innerW {
			innerW = minInner
		}
	}
	if b.footer != "" {
		blW := visibleWidth(b.border.BottomLeft)
		brW := visibleWidth(b.border.BottomRight)
		needed := visibleWidth(b.footer) + 2*horW
		vertSum := visibleWidth(b.border.Vertical)*2 - blW - brW
		minInner := needed - vertSum
		if minInner > innerW {
			innerW = minInner
		}
	}

	leftVert := b.border.Vertical
	rightVert := b.border.Vertical
	if b.hideLeft {
		leftVert = strings.Repeat(" ", visibleWidth(b.border.Vertical))
	}
	if b.hideRight {
		rightVert = strings.Repeat(" ", visibleWidth(b.border.Vertical))
	}

	var boxRows []string
	totalBodyRows := b.padTop + len(lines) + b.padBottom

	showTopLeftCap := b.hideTop && !b.hideLeft && !b.hideTopLeft
	showTopRightCap := b.hideTop && !b.hideRight && !b.hideTopRight
	showBotLeftCap := b.hideBottom && !b.hideLeft && !b.hideBotLeft
	showBotRightCap := b.hideBottom && !b.hideRight && !b.hideBotRight

	bodyEdgeGlyphs := func(bodyIdx int) (string, string) {
		leftGlyph := leftVert
		rightGlyph := rightVert

		if bodyIdx == 0 {
			if showTopLeftCap {
				leftGlyph = b.border.TopLeft
			}
			if showTopRightCap {
				rightGlyph = b.border.TopRight
			}
		}

		if bodyIdx == totalBodyRows-1 {
			if showBotLeftCap && !(bodyIdx == 0 && showTopLeftCap) {
				leftGlyph = b.border.BottomLeft
			}
			if showBotRightCap && !(bodyIdx == 0 && showTopRightCap) {
				rightGlyph = b.border.BottomRight
			}
		}

		return leftGlyph, rightGlyph
	}

	frameW := visibleWidth(b.border.Vertical) + innerW + visibleWidth(b.border.Vertical)

	if !b.hideTop {
		topBar := b.buildBorderRow(
			b.border.TopLeft, b.border.TopRight,
			b.hideTopLeft, b.hideTopRight,
			b.title, b.titleAlign, frameW,
		)
		boxRows = append(boxRows, b.wrapStyle(topBar))
	}

	for i := 0; i < b.padTop; i++ {
		leftGlyph, rightGlyph := bodyEdgeGlyphs(i)
		padLine := leftGlyph + strings.Repeat(" ", innerW) + rightGlyph
		boxRows = append(boxRows, b.wrapStyle(padLine))
	}

	lastIdx := len(lines) - 1
	for i := 0; i < len(lines); i++ {
		bodyIdx := b.padTop + i
		leftGlyph, rightGlyph := bodyEdgeGlyphs(bodyIdx)

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

		row := b.wrapStyle(leftGlyph+strings.Repeat(" ", b.padLeft+leftPad)) +
			line +
			b.wrapStyle(strings.Repeat(" ", rightPad+b.padRight)+rightGlyph)
		boxRows = append(boxRows, row)
	}

	for i := 0; i < b.padBottom; i++ {
		bodyIdx := b.padTop + len(lines) + i
		leftGlyph, rightGlyph := bodyEdgeGlyphs(bodyIdx)
		padLine := leftGlyph + strings.Repeat(" ", innerW) + rightGlyph
		boxRows = append(boxRows, b.wrapStyle(padLine))
	}

	if !b.hideBottom {
		botBar := b.buildBorderRow(
			b.border.BottomLeft, b.border.BottomRight,
			b.hideBotLeft, b.hideBotRight,
			b.footer, b.footerAlign, frameW,
		)
		boxRows = append(boxRows, b.wrapStyle(botBar))
	}

	bottomBorderIdx := -1
	if !b.hideBottom {
		bottomBorderIdx = len(boxRows) - 1
	}

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
		if bottomBorderIdx >= 0 && i == bottomBorderIdx {
		} else {
			out.WriteByte('\n')
		}
	}

	for i := 0; i < b.marginBottom; i++ {
		out.WriteByte('\n')
	}

	return out.String()
}
