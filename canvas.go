package tinta

import (
	"sort"
	"strings"
)

// cell represents a single terminal position: one visible rune plus any
// ANSI style prefix.
type cell struct {
	r     rune   // visible character
	style string // ANSI escape prefix (everything before the rune, e.g. "\x1b[31m")
}

// layer holds a parsed grid of cells at an (x, y) offset with a z-index.
type layer struct {
	grid [][]cell // rows of cells
	x, y int
	z    int
	seq  int // insertion order for stable sort
}

// CanvasStyle holds layers and compositing settings. Create one with
// [Canvas] and chain Add/AddZ/Width/Height methods. Call [CanvasStyle.String]
// to composite all layers into a final string.
//
// All methods return a new CanvasStyle to preserve immutability.
type CanvasStyle struct {
	layers []layer
	width  int // 0 = auto (derived from layer extents)
	height int // 0 = auto
	nextZ  int // auto-incrementing z counter
}

// Canvas returns a new empty [CanvasStyle].
func Canvas() *CanvasStyle {
	return &CanvasStyle{}
}

// copyCanvas returns a deep copy of the CanvasStyle, including the layers slice.
func copyCanvas(c *CanvasStyle) *CanvasStyle {
	cp := *c
	if len(c.layers) > 0 {
		cp.layers = make([]layer, len(c.layers))
		copy(cp.layers, c.layers)
	}
	return &cp
}

// Add places a rendered string at position (x, y) on the canvas. The
// z-index is auto-incremented: each successive Add gets a higher z than
// the previous one.
func (c *CanvasStyle) Add(s string, x, y int) *CanvasStyle {
	return c.AddZ(s, x, y, c.nextZ)
}

// AddZ places a rendered string at position (x, y) with an explicit z-index.
// Higher z values are drawn on top of lower ones. When two layers share the
// same z, insertion order wins (later Add calls draw on top).
func (c *CanvasStyle) AddZ(s string, x, y, z int) *CanvasStyle {
	cp := copyCanvas(c)
	grid := parseGrid(s)
	cp.layers = append(cp.layers, layer{
		grid: grid,
		x:    x,
		y:    y,
		z:    z,
		seq:  len(cp.layers),
	})
	if z >= cp.nextZ {
		cp.nextZ = z + 1
	}
	return cp
}

// Width sets a fixed canvas width. If zero (default), the width is
// derived from the rightmost visible cell across all layers.
func (c *CanvasStyle) Width(w int) *CanvasStyle {
	cp := copyCanvas(c)
	cp.width = w
	return cp
}

// Height sets a fixed canvas height. If zero (default), the height is
// derived from the bottommost visible cell across all layers.
func (c *CanvasStyle) Height(h int) *CanvasStyle {
	cp := copyCanvas(c)
	cp.height = h
	return cp
}

// String composites all layers and returns the final rendered string.
// Layers are drawn in z-order (ascending), then by insertion order.
// Each layer is fully opaque: every cell in a layer's grid overwrites
// whatever is below it. Positions not covered by any layer are rendered
// as plain spaces. The result has no trailing newline on the last row.
func (c *CanvasStyle) String() string {
	if len(c.layers) == 0 {
		return ""
	}

	// Sort layers by (z, seq) ascending.
	sorted := make([]layer, len(c.layers))
	copy(sorted, c.layers)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].z != sorted[j].z {
			return sorted[i].z < sorted[j].z
		}
		return sorted[i].seq < sorted[j].seq
	})

	// Determine canvas dimensions.
	w := c.width
	h := c.height
	if w == 0 || h == 0 {
		autoW, autoH := 0, 0
		for _, ly := range sorted {
			for rowIdx, row := range ly.grid {
				cy := ly.y + rowIdx
				if cy < 0 {
					continue
				}
				for colIdx := range row {
					cx := ly.x + colIdx
					if cx < 0 {
						continue
					}
					if cx+1 > autoW {
						autoW = cx + 1
					}
				}
				if cy+1 > autoH {
					autoH = cy + 1
				}
			}
		}
		if w == 0 {
			w = autoW
		}
		if h == 0 {
			h = autoH
		}
	}

	if w <= 0 || h <= 0 {
		return ""
	}

	// Build the composited grid. Positions start as plain spaces (the
	// default background). Layers paint over these.
	grid := make([][]cell, h)
	for i := range grid {
		grid[i] = make([]cell, w)
		for j := range grid[i] {
			grid[i][j] = cell{r: ' '} // transparent
		}
	}

	// Paint layers in order (lowest z first). Each layer's cells fully
	// overwrite whatever is below — all cells within a layer's grid are
	// opaque. Positions outside a layer's grid are not touched.
	for _, ly := range sorted {
		for rowIdx, row := range ly.grid {
			cy := ly.y + rowIdx
			if cy < 0 || cy >= h {
				continue
			}
			for colIdx, cl := range row {
				cx := ly.x + colIdx
				if cx < 0 || cx >= w {
					continue
				}
				grid[cy][cx] = cl
			}
		}
	}

	// Render the composited grid into a string with proper ANSI sequences.
	var buf strings.Builder
	for rowIdx, row := range grid {
		if rowIdx > 0 {
			buf.WriteByte('\n')
		}

		// Find the last non-blank cell to avoid trailing whitespace.
		// A "blank" trailing cell is an unstyled space.
		lastVisible := len(row) - 1
		for lastVisible >= 0 && row[lastVisible].r == ' ' && row[lastVisible].style == "" {
			lastVisible--
		}

		// Track the last style applied so we can avoid redundant resets/sets.
		lastStyle := ""
		for colIdx := 0; colIdx <= lastVisible; colIdx++ {
			cl := row[colIdx]
			if cl.style != lastStyle {
				// Close previous style if any.
				if lastStyle != "" {
					buf.WriteString(cReset)
				}
				// Open new style if any.
				if cl.style != "" {
					buf.WriteString(cl.style)
				}
				lastStyle = cl.style
			}
			buf.WriteRune(cl.r)
		}
		// Close any open style at end of row.
		if lastStyle != "" {
			buf.WriteString(cReset)
		}
	}

	return buf.String()
}

// parseGrid converts a rendered string (possibly containing ANSI escape
// sequences) into a 2D grid of cells. Each row corresponds to a line of
// the input (split by \n). Within each row, ANSI sequences are accumulated
// into a "current style" prefix that gets attached to each subsequent
// visible rune.
//
// A reset sequence (\x1b[0m) clears the current style.
func parseGrid(s string) [][]cell {
	// Strip trailing newline to avoid a phantom empty row.
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if s == "" {
		return nil
	}

	lines := strings.Split(s, "\n")
	grid := make([][]cell, len(lines))

	for i, line := range lines {
		grid[i] = parseLine(line)
	}
	return grid
}

// parseLine converts a single line of text (with possible ANSI sequences)
// into a slice of cells.
func parseLine(line string) []cell {
	var cells []cell
	var styleBuf strings.Builder // accumulates ANSI sequences between visible runes

	j := 0
	runes := []byte(line)
	n := len(runes)

	for j < n {
		if runes[j] == '\x1b' {
			// Start of an escape sequence. Capture the entire sequence.
			start := j
			if j+1 < n {
				switch runes[j+1] {
				case '[': // CSI sequence
					j += 2
					for j < n {
						if runes[j] >= 0x40 && runes[j] <= 0x7E {
							j++
							break
						}
						j++
					}
				case ']': // OSC sequence
					j += 2
					for j < n {
						if runes[j] == '\x07' {
							j++
							break
						}
						if runes[j] == '\x1b' && j+1 < n && runes[j+1] == '\\' {
							j += 2
							break
						}
						j++
					}
				default:
					j += 2
				}
			} else {
				j++
			}
			seq := line[start:j]

			// Check if this is a reset. If so, clear the style buffer.
			if seq == cReset {
				styleBuf.Reset()
			} else {
				styleBuf.WriteString(seq)
			}
			continue
		}

		// Visible character. We need to consume one UTF-8 rune.
		// Since we're working with byte offsets but need rune boundaries,
		// let's extract the rune.
		r, size := decodeRune(line[j:])
		cells = append(cells, cell{
			r:     r,
			style: styleBuf.String(),
		})
		j += size
	}

	return cells
}

// decodeRune decodes the first UTF-8 rune from s and returns it along with
// its byte length. This avoids importing unicode/utf8 to keep dependencies
// at zero.
func decodeRune(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}
	b := s[0]
	if b < 0x80 {
		return rune(b), 1
	}
	// Determine byte length from leading bits.
	var size int
	switch {
	case b>>5 == 0x06:
		size = 2
	case b>>4 == 0x0E:
		size = 3
	case b>>3 == 0x1E:
		size = 4
	default:
		return rune(b), 1 // invalid leading byte, treat as single byte
	}
	if len(s) < size {
		return rune(b), 1
	}
	var r rune
	switch size {
	case 2:
		r = rune(b&0x1F)<<6 | rune(s[1]&0x3F)
	case 3:
		r = rune(b&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F)
	case 4:
		r = rune(b&0x07)<<18 | rune(s[1]&0x3F)<<12 | rune(s[2]&0x3F)<<6 | rune(s[3]&0x3F)
	}
	return r, size
}
