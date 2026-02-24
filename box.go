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

// ShadowStyle defines the glyphs used to draw a shadow around a box.
// The structure mirrors [Border]: corners, horizontal, and vertical pieces.
type ShadowStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

// Predefined shadow styles.
var (
	// ShadowLight uses light shade characters (░).
	ShadowLight = ShadowStyle{"░", "░", "░", "░", "░", "░"}

	// ShadowMedium uses medium shade characters (▒).
	ShadowMedium = ShadowStyle{"▒", "▒", "▒", "▒", "▒", "▒"}

	// ShadowDark uses dark shade characters (▓).
	ShadowDark = ShadowStyle{"▓", "▓", "▓", "▓", "▓", "▓"}

	// ShadowBlock uses full block characters (█).
	ShadowBlock = ShadowStyle{"█", "█", "█", "█", "█", "█"}
)

// ShadowPosition determines the direction in which the shadow is cast.
type ShadowPosition int

const (
	// ShadowBottomRight places the shadow below and to the right of the box.
	ShadowBottomRight ShadowPosition = iota
	// ShadowBottomLeft places the shadow below and to the left.
	ShadowBottomLeft
	// ShadowTopRight places the shadow above and to the right.
	ShadowTopRight
	// ShadowTopLeft places the shadow above and to the left.
	ShadowTopLeft
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
	center       bool             // center all content lines horizontally
	centerTrim   bool             // trim whitespace from lines before centering
	centerLines  map[int]struct{} // specific line indices to center (0-based)
	centerFirst  bool             // center the first content line
	centerLast   bool             // center the last content line
	hideTop      bool             // hide the top border row
	hideBottom   bool             // hide the bottom border row
	hideLeft     bool             // hide the left vertical border
	hideRight    bool             // hide the right vertical border
	shadow       *ShadowStyle     // nil means no shadow
	shadowPos    ShadowPosition   // direction of the shadow
	shadowCodes  []string         // ANSI SGR codes for the shadow glyphs
}

// BoxStyle is the public handle returned by [Box] and every chaining method.
// The underlying struct is opaque; users cannot create one manually.
type BoxStyle = *box

// Box returns a new BoxStyle with a simple border and no padding or margin.
func Box() BoxStyle {
	return &box{border: BorderSimple}
}

// copyBox returns a deep copy of the box, including the codes slice,
// centerLines map, and shadow configuration.
func copyBox(b *box) *box {
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
	if b.shadow != nil {
		s := *b.shadow
		cp.shadow = &s
	}
	if len(b.shadowCodes) > 0 {
		cp.shadowCodes = make([]string, len(b.shadowCodes))
		copy(cp.shadowCodes, b.shadowCodes)
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

// Border sets a custom border using the provided [Border] struct.
func (b *box) Border(border Border) BoxStyle {
	cp := copyBox(b)
	cp.border = border
	return cp
}

// BorderSimple sets the border to a simple line style (┌─┐).
func (b *box) BorderSimple() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderSimple
	return cp
}

// BorderRounded sets the border to a rounded style (╭─╮).
func (b *box) BorderRounded() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderRounded
	return cp
}

// BorderDouble sets the border to a double line style (╔═╗).
func (b *box) BorderDouble() BoxStyle {
	cp := copyBox(b)
	cp.border = BorderDouble
	return cp
}

// BorderHeavy sets the border to a heavy line style (┏━┓).
func (b *box) BorderHeavy() BoxStyle {
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

// CenterLine marks the line at index n (0-based) for horizontal centering.
// If n is out of bounds at render time, the call is silently ignored.
// This can be called multiple times to center several specific lines.
func (b *box) CenterLine(n int) BoxStyle {
	cp := copyBox(b)
	if cp.centerLines == nil {
		cp.centerLines = make(map[int]struct{})
	}
	cp.centerLines[n] = struct{}{}
	return cp
}

// CenterFirstLine centers the first content line (index 0).
// This is a convenience shortcut useful for centering titles.
func (b *box) CenterFirstLine() BoxStyle {
	cp := copyBox(b)
	cp.centerFirst = true
	return cp
}

// CenterLastLine centers the last content line.
// The line count is determined at render time.
func (b *box) CenterLastLine() BoxStyle {
	cp := copyBox(b)
	cp.centerLast = true
	return cp
}

// --- Side visibility ---

// DisableTop hides the top border row. The vertical borders on content
// rows remain unchanged.
func (b *box) DisableTop() BoxStyle {
	cp := copyBox(b)
	cp.hideTop = true
	return cp
}

// DisableBottom hides the bottom border row.
func (b *box) DisableBottom() BoxStyle {
	cp := copyBox(b)
	cp.hideBottom = true
	return cp
}

// DisableLeft hides the left vertical border on all rows (top corner,
// content, padding, and bottom corner).
func (b *box) DisableLeft() BoxStyle {
	cp := copyBox(b)
	cp.hideLeft = true
	return cp
}

// DisableRight hides the right vertical border on all rows.
func (b *box) DisableRight() BoxStyle {
	cp := copyBox(b)
	cp.hideRight = true
	return cp
}

// --- Shadow ---

// Shadow enables a shadow effect in the given direction using the
// provided [ShadowStyle]. Shadow glyphs are rendered with bright-black
// color by default; use [ShadowDim], [ShadowBlack], or
// [ShadowBrightBlack] to change the color.
func (b *box) Shadow(pos ShadowPosition, sty ShadowStyle) BoxStyle {
	cp := copyBox(b)
	cp.shadow = &sty
	cp.shadowPos = pos
	if len(cp.shadowCodes) == 0 {
		cp.shadowCodes = []string{cBrightBlack}
	}
	return cp
}

// ShadowDim applies the dim modifier to the shadow.
func (b *box) ShadowDim() BoxStyle {
	cp := copyBox(b)
	cp.shadowCodes = append(cp.shadowCodes, cDim)
	return cp
}

// ShadowBlack sets the shadow foreground to black.
func (b *box) ShadowBlack() BoxStyle {
	cp := copyBox(b)
	cp.shadowCodes = append(cp.shadowCodes, cBlack)
	return cp
}

// ShadowBrightBlack sets the shadow foreground to bright black.
func (b *box) ShadowBrightBlack() BoxStyle {
	cp := copyBox(b)
	cp.shadowCodes = append(cp.shadowCodes, cBrightBlack)
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
func (b *box) wrapStyle(s string) string {
	return wrapCodes(s, b.codes)
}

// wrapShadow wraps s in the box's shadow ANSI codes.
func (b *box) wrapShadow(s string) string {
	return wrapCodes(s, b.shadowCodes)
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

	// Top border: ┌───┐
	if !b.hideTop {
		tl := b.border.TopLeft
		tr := b.border.TopRight
		if b.hideLeft {
			tl = strings.Repeat(" ", visibleWidth(b.border.TopLeft))
		}
		if b.hideRight {
			tr = strings.Repeat(" ", visibleWidth(b.border.TopRight))
		}
		topBar := tl + strings.Repeat(b.border.Horizontal, innerW) + tr
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

	// Bottom border: └───┘
	if !b.hideBottom {
		bl := b.border.BottomLeft
		br := b.border.BottomRight
		if b.hideLeft {
			bl = strings.Repeat(" ", visibleWidth(b.border.BottomLeft))
		}
		if b.hideRight {
			br = strings.Repeat(" ", visibleWidth(b.border.BottomRight))
		}
		botBar := bl + strings.Repeat(b.border.Horizontal, innerW) + br
		boxRows = append(boxRows, b.wrapStyle(botBar))
	}

	// Compute visible width of the box (from the first row).
	boxVisW := 0
	if len(boxRows) > 0 {
		boxVisW = visibleWidth(boxRows[0])
	}

	// Track the index of the bottom border row (if present) before shadow
	// is applied. This determines trailing-newline behavior.
	bottomBorderIdx := -1
	if !b.hideBottom {
		bottomBorderIdx = len(boxRows) - 1
	}

	// Apply shadow if enabled.
	hasShadow := b.shadow != nil
	if hasShadow {
		boxRows = b.applyShadow(boxRows, boxVisW)
	}

	// Assemble final output with margins.
	marginLeft := strings.Repeat(" ", b.marginLeft)
	marginRight := strings.Repeat(" ", b.marginRight)

	var out strings.Builder

	for i := 0; i < b.marginTop; i++ {
		out.WriteByte('\n')
	}

	lastRow := len(boxRows) - 1
	for i := 0; i < len(boxRows); i++ {
		out.WriteString(marginLeft)
		out.WriteString(boxRows[i])
		out.WriteString(marginRight)
		// The final rendered row never gets a trailing \n.
		// Without shadow: the bottom border (if present) is final.
		// With shadow: the shadow's last row is final.
		// When the bottom border is hidden and there's no shadow,
		// the last content/padding row still gets \n (legacy behavior).
		if hasShadow {
			// With shadow, the very last row is the shadow's last row.
			if i < lastRow {
				out.WriteByte('\n')
			}
		} else if bottomBorderIdx >= 0 && i == bottomBorderIdx {
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

// applyShadow adds shadow glyphs to the collected box rows based on
// the shadow position. It returns a new slice with the shadow applied.
//
// The shadow forms an L-shape with three visible corners. For example,
// ShadowBottomRight produces:
//
//	┌────┐
//	│ hi │╮   ← TopRight corner
//	└────┘│   ← Vertical
//	 ╰────╯   ← BottomLeft, Horizontals, BottomRight
func (b *box) applyShadow(rows []string, boxVisW int) []string {
	s := b.shadow
	n := len(rows)
	if n == 0 {
		return rows
	}

	// Styled glyphs.
	shadowV := b.wrapShadow(s.Vertical)

	// Visible width of the vertical glyph (usually 1).
	vertW := visibleWidth(s.Vertical)
	spacer := strings.Repeat(" ", vertW)

	// Horizontal fill count for the shadow bar (between the two corners).
	// The bar spans boxVisW visible characters total: two corners (each
	// vertW wide) plus horizontal fill in between.
	horzW := visibleWidth(s.Horizontal)
	hFill := boxVisW - 2*vertW
	if hFill < 0 {
		hFill = 0
	}
	hCount := 0
	if horzW > 0 {
		hCount = hFill / horzW
	}

	switch b.shadowPos {
	case ShadowBottomRight:
		// Right column: TopRight corner at top, Vertical for the rest.
		// Bottom row:   BottomLeft corner, Horizontals, BottomRight corner.
		result := make([]string, 0, n+1)
		result = append(result, rows[0]+spacer)
		if n > 1 {
			result = append(result, rows[1]+b.wrapShadow(s.TopRight))
			for i := 2; i < n; i++ {
				result = append(result, rows[i]+shadowV)
			}
		}
		hBar := s.BottomLeft + strings.Repeat(s.Horizontal, hCount) + s.BottomRight
		result = append(result, spacer+b.wrapShadow(hBar))
		return result

	case ShadowBottomLeft:
		// Left column: TopLeft corner at top, Vertical for the rest.
		// Bottom row:  BottomLeft corner, Horizontals, BottomRight corner.
		result := make([]string, 0, n+1)
		result = append(result, spacer+rows[0])
		if n > 1 {
			result = append(result, b.wrapShadow(s.TopLeft)+rows[1])
			for i := 2; i < n; i++ {
				result = append(result, shadowV+rows[i])
			}
		}
		hBar := s.BottomLeft + strings.Repeat(s.Horizontal, hCount) + s.BottomRight
		result = append(result, b.wrapShadow(hBar)+spacer)
		return result

	case ShadowTopRight:
		// Top row:     TopLeft corner, Horizontals, TopRight corner.
		// Right column: Vertical for most rows, BottomRight corner at bottom.
		result := make([]string, 0, n+1)
		hBar := s.TopLeft + strings.Repeat(s.Horizontal, hCount) + s.TopRight
		result = append(result, spacer+b.wrapShadow(hBar))
		for i := 0; i < n-1; i++ {
			result = append(result, rows[i]+shadowV)
		}
		if n > 1 {
			result = append(result, rows[n-1]+b.wrapShadow(s.BottomRight))
		} else {
			result = append(result, rows[0]+b.wrapShadow(s.BottomRight))
		}
		return result

	case ShadowTopLeft:
		// Top row:      TopLeft corner, Horizontals, TopRight corner.
		// Left column:  Vertical for most rows, BottomLeft corner at bottom.
		result := make([]string, 0, n+1)
		hBar := s.TopLeft + strings.Repeat(s.Horizontal, hCount) + s.TopRight
		result = append(result, b.wrapShadow(hBar)+spacer)
		for i := 0; i < n-1; i++ {
			result = append(result, shadowV+rows[i])
		}
		if n > 1 {
			result = append(result, b.wrapShadow(s.BottomLeft)+rows[n-1])
		} else {
			result = append(result, b.wrapShadow(s.BottomLeft)+rows[0])
		}
		return result
	}

	return rows
}
