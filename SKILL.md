---
name: tinta-expert
description: Generate Go CLI output with Tinta (github.com/varavelio/tinta). Use when building terminal styled text, bordered boxes, spacing layouts, shadows, and any terminal style related task.
---

# Tinta Expert

Use this skill to produce correct, idiomatic Go code with `github.com/varavelio/tinta`.

## When to use

Use this skill when the task involves:

- ANSI-styled CLI text (status lines, warnings, headings)
- Boxed layouts with spacing/alignment and border customization
- Layered composition of rendered strings via canvas
- Stable output assertions in tests (`ForceColors` and buffers)

## Instructions

1. Start from the right primitive:
   - `tinta.Text()` for inline styling
   - `tinta.Box()` for framed layout
   - `tinta.Canvas()` for multi-layer composition
2. Keep chaining immutable. Every method returns a new style; never rely on mutation.
3. Prefer readable output first; style should clarify hierarchy, not replace it.
4. In tests, force deterministic color behavior when needed.

## API essentials

### Text

- Foreground: `Black..White`, `BrightBlack..BrightWhite`
- Background: `OnBlack..OnWhite`, `OnBrightBlack..OnBrightWhite`
- Modifiers: `Bold`, `Dim`, `Italic`, `Underline`, `Invert`, `Hidden`, `Strike`
- Output: `String`, `Sprintf`, `Print`, `Printf`, `Println`, `Fprint`, `Fprintf`, `Fprintln`

### Box

- Border setup:
  - Presets: `tinta.BorderSimple`, `tinta.BorderDashed`, `tinta.BorderDotted`, `tinta.BorderRounded`, `tinta.BorderRoundedDashed`, `tinta.BorderRoundedDotted`, `tinta.BorderDouble`, `tinta.BorderHeavy`, `tinta.BorderASCII`, `tinta.BorderBlock`, `tinta.BorderBlockHalf`, `tinta.BorderBlockLight`, `tinta.BorderBlockMedium`, `tinta.BorderBlockDark`
  - Custom struct fields:
    - corners: `TopLeft`, `TopRight`, `BottomLeft`, `BottomRight`
    - sides: `Top`, `Left`, `Right`, `Bottom`
  - Apply with `Box().Border(borderValue)`
- Spacing: `Padding*`, `Margin*` (`Padding`, `PaddingX`, `PaddingY`, etc.)
- Content alignment: `Center`, `CenterTrim`, `CenterLine`, `CenterFirstLine`, `CenterLastLine`
- Side visibility: `DisableTop`, `DisableBottom`, `DisableLeft`, `DisableRight`
- Corner visibility: `DisableCorners`, `DisableTopLeftCorner`, `DisableTopRightCorner`, `DisableBottomLeftCorner`, `DisableBottomRightCorner`
- Border labels:
  - `Title(text, align)` on top border row
  - `Footer(text, align)` on bottom border row
  - `align`: `AlignLeft`, `AlignCenter`, `AlignRight`
- Colors/modifiers: same color set as `Text`, plus `Bold`, `Dim`
- Output: same method family as `Text`

### Canvas

- `Canvas()` creates an empty immutable compositor
- `Add(s, x, y)` appends a layer with auto z
- `AddZ(s, x, y, z)` appends with explicit z
- `Width(w)` / `Height(h)` set fixed output dimensions (`0` means auto)
- `String()` composites layers

Compositing behavior:

- Layers render by z ascending, then insertion order
- Layer cells are opaque (overwrite underlying cells)
- Negative `x/y` expands auto-sized canvas to fit all content
- Fixed width/height applies cropping after expansion

## Output and color control

- `SetOutput(w)` changes default writer for `Print*`
- `ForceColors(true|false)` overrides auto-detection
- Auto-detection honors typical env flags (`NO_COLOR`, `FORCE_COLOR`, `CLICOLOR`, `TERM=dumb`)

## Testing guidance

- For plain-text assertions: `ForceColors(false)` and restore with `defer ForceColors(true)`
- For ANSI assertions: `ForceColors(true)` and assert escape sequences deliberately
- Prefer `String()` for deterministic snapshots
- Use `bytes.Buffer` with `Fprint/Fprintf/Fprintln` when writer behavior is under test

## Common pitfalls

- Do not measure styled width with `len(s)`; ANSI is non-visible
- Do not assume mutable style objects

## Examples

User: "Create a titled status panel with custom borders"

```go
panelBorder := tinta.Border{
  TopLeft: "+", 
  TopRight: "+", 
  BottomLeft: "+", 
  BottomRight: "+",
  Top: "-", 
  Left: "|", 
  Right: "!", 
  Bottom: "~",
}

out := tinta.Box().
    Border(panelBorder).
    PaddingX(1).
    Title("Status", tinta.AlignCenter).
    Footer("ok", tinta.AlignRight).
    String("service started")
```

User: "Compose two boxes with depth"

```go
front := tinta.Box().Border(tinta.BorderHeavy).PaddingX(2).String("front")
back := tinta.Box().Border(tinta.BorderRounded).PaddingX(2).String("back")

scene := tinta.Canvas().
    Add(back, -1, 1).
    Add(front, 0, 0).
    String()
```
