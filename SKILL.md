---
name: tinta-expert
description: Expert usage guide for generating Go CLI output with Tinta (github.com/varavelio/tinta). Use when building terminal styled text, bordered boxes, spacing layouts, shadows, and any terminal style related task.
---

# Tinta Expert

Use this skill to write production-quality Go code that **uses** `github.com/varavelio/tinta` correctly and idiomatically.

## When to use this skill

Use this skill whenever the user asks for Go terminal UX such as:

- colored logs, status lines, banners, prompts, notices
- boxed output with borders, padding, margins, alignment, or shadows
- nested styled output (styled text inside boxes, boxes inside boxes)
- deterministic output to custom writers (`bytes.Buffer`, files, tests)
- explicit color behavior in CI/non-TTY environments

## Core usage model

Always start from one of the two entry points:

- `tinta.Text()` for ANSI text styling
- `tinta.Box()` for structural terminal layout

All chaining methods are immutable: each call returns a new style value.

```go
base := tinta.Text().Bold()
errStyle := base.Red()
okStyle := base.Green()
```

## Text API (complete practical surface)

### Foreground colors

`Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White`

### Bright foreground colors

`BrightBlack`, `BrightRed`, `BrightGreen`, `BrightYellow`, `BrightBlue`, `BrightMagenta`, `BrightCyan`, `BrightWhite`

### Background colors

`OnBlack`, `OnRed`, `OnGreen`, `OnYellow`, `OnBlue`, `OnMagenta`, `OnCyan`, `OnWhite`

### Bright background colors

`OnBrightBlack`, `OnBrightRed`, `OnBrightGreen`, `OnBrightYellow`, `OnBrightBlue`, `OnBrightMagenta`, `OnBrightCyan`, `OnBrightWhite`

### Text modifiers

`Bold`, `Dim`, `Italic`, `Underline`, `Invert`, `Hidden`, `Strike`

### Text output methods

- Returns string: `String(s)`, `Sprintf(format, ...)`
- Writes to default output: `Print`, `Printf`, `Println`
- Writes to explicit writer: `Fprint`, `Fprintf`, `Fprintln`

## Box API (complete practical surface)

### Borders

- Presets: `BorderSimple()` (default), `BorderRounded()`, `BorderDouble()`, `BorderHeavy()`
- Custom: `Border(tinta.Border{...})`

### Padding

`Padding`, `PaddingTop`, `PaddingBottom`, `PaddingLeft`, `PaddingRight`, `PaddingX`, `PaddingY`

### Margin

`Margin`, `MarginTop`, `MarginBottom`, `MarginLeft`, `MarginRight`, `MarginX`, `MarginY`

### Alignment

- Full content: `Center()`
- Trim then center: `CenterTrim()`
- Selective: `CenterLine(n)`, `CenterFirstLine()`, `CenterLastLine()`

### Side visibility

`DisableTop`, `DisableBottom`, `DisableLeft`, `DisableRight`

### Shadow

- Enable: `Shadow(position, style)`
- Positions: `ShadowBottomRight`, `ShadowBottomLeft`, `ShadowTopRight`, `ShadowTopLeft`
- Styles: `ShadowLight`, `ShadowMedium`, `ShadowDark`, `ShadowBlock`, or custom `ShadowStyle`
- Shadow color/modifiers: `ShadowDim`, `ShadowBlack`, `ShadowBrightBlack`

### Box colors and modifiers

Box supports the same foreground/background color methods as `Text`, plus:

- modifiers: `Bold`, `Dim`

### Box output methods

- Returns string: `String(s)`, `Sprintf(format, ...)`
- Writes to default output: `Print`, `Printf`, `Println`
- Writes to explicit writer: `Fprint`, `Fprintf`, `Fprintln`

## Output control and color behavior

Use package-level output controls when needed:

- `tinta.SetOutput(w)` changes the default writer used by `Print*`
- `tinta.ForceColors(true/false)` overrides color auto-detection

Color detection behavior:

- disabled when any of: `NO_COLOR`, `NO_COLORS`, `DISABLE_COLORS`, `CLICOLOR=0`, `TERM=dumb`
- forced when any of: `FORCE_COLOR`, `CLICOLOR_FORCE`
- disable has precedence over force
- when not forced/disabled, color depends on TTY detection

## Recommended coding patterns

### 1) Reusable style constants

```go
var (
    label = tinta.Text().BrightBlack()
    ok    = tinta.Text().Green().Bold()
    warn  = tinta.Text().Yellow().Bold()
    fail  = tinta.Text().Red().Bold()
)
```

### 2) Prebuild messages with `String`/`Sprintf`

```go
status := tinta.Text().White().OnBlue().Sprintf(" service: %s ", "running")
fmt.Println(status)
```

### 3) Structured sections with boxes

```go
panel := tinta.Box().BorderRounded().Blue().PaddingY(1).PaddingX(2)
panel.Println("Deploy Summary\n- service: api\n- result: success")
```

### 4) Safe nested rendering

```go
inner := tinta.Box().BorderRounded().Green().PaddingX(1).String("Inner")
outer := tinta.Box().BorderDouble().Blue().PaddingX(2)
outer.Println(inner)
```

### 5) Test-friendly deterministic output

```go
var buf bytes.Buffer
tinta.ForceColors(true)
_, _ = tinta.Text().Red().Fprint(&buf, "error")
```

## Style composition guidance

- Prefer building a small base style and branching from it.
- Keep command output readable first; style should support hierarchy, not replace it.
- Use bright accents sparingly for warnings/errors and key state transitions.
- For multiline content with inconsistent indentation, prefer `CenterTrim()` over `Center()`.
- For callouts/quotes, use side disabling patterns instead of over-styling text.

## Canonical UI patterns

### Error line

```go
tinta.Text().Red().Bold().Println("error: configuration file not found")
```

### Success summary box

```go
tinta.Box().BorderRounded().Green().PaddingY(1).PaddingX(2).
    Println("Build completed\nArtifacts: 3\nDuration: 12.4s")
```

### Section heading with underline effect

```go
tinta.Box().BorderHeavy().BrightYellow().
    DisableTop().DisableLeft().DisableRight().
    Println("Deployment")
```

### Blockquote effect

```go
tinta.Box().BorderHeavy().BrightCyan().
    DisableTop().DisableBottom().DisableRight().
    PaddingLeft(1).
    Println("The best way to predict the future is to invent it.")
```

### Emphasized card with shadow

```go
tinta.Box().BorderRounded().PaddingX(2).
    Shadow(tinta.ShadowBottomRight, tinta.ShadowLight).
    Println("Release v1.4.2")
```

## Mistakes to avoid

- Do not mutate or reuse style values as mutable state assumptions.
- Do not assume colors always render; CI/non-TTY may disable ANSI.
- Do not rely on visual width via `len(s)` when mixing ANSI output.
- Do not overuse chained modifiers; prioritize legibility.
- Do not hardcode escape sequences directly when Tinta methods exist.

## Agent execution checklist

- [ ] import path is `github.com/varavelio/tinta`
- [ ] style creation uses `Text()`/`Box()` chaining
- [ ] output path uses correct method family (`String` vs `Print` vs `Fprint`)
- [ ] colors are deterministic where required (`ForceColors` in tests/demo contexts)
- [ ] nested styled content is rendered via Tinta, not manual ANSI concatenation
- [ ] final output remains readable with colors disabled

## Minimal reference snippets

User request: "Print warning and success lines"

```go
tinta.Text().Yellow().Bold().Println("warning: retrying request")
tinta.Text().Green().Println("ok: operation completed")
```

User request: "Show centered title card"

```go
tinta.Box().BorderDouble().Center().PaddingX(1).
    Println("Build Pipeline\nrelease")
```

User request: "Render to buffer for assertions"

```go
var buf bytes.Buffer
tinta.ForceColors(true)
_, _ = tinta.Box().Red().Fprintln(&buf, "failed")
```
