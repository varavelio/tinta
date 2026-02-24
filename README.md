# Tinta

<img src="./gopher-200px.png" height="200px" width="auto">

<br/>

<p>
  <a href="https://github.com/varavelio/tinta/actions">
    <img src="https://github.com/varavelio/tinta/actions/workflows/ci.yaml/badge.svg" alt="CI status"/>
  </a>
  <a href="https://pkg.go.dev/github.com/varavelio/tinta">
    <img src="https://pkg.go.dev/badge/github.com/varavelio/tinta" alt="Go Reference"/>
  </a>
  <a href="https://goreportcard.com/report/varavelio/tinta">
    <img src="https://goreportcard.com/badge/varavelio/tinta" alt="Go Report Card"/>
  </a>
  <a href="https://github.com/varavelio/tinta/releases/latest">
    <img src="https://img.shields.io/github/release/varavelio/tinta.svg" alt="Release Version"/>
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/varavelio/tinta.svg" alt="License"/>
  </a>
  <a href="https://github.com/varavelio/tinta">
    <img src="https://img.shields.io/github/stars/varavelio/tinta?style=flat&label=github+stars"/>
  </a>
</p>

Minimal, chainable terminal styling and layout for Go. Zero dependencies.

Two pillars: **Text** for ANSI colors and modifiers, **Box** for bordered containers with padding, margin, shadows, and content alignment.

## Install

```sh
go get github.com/varavelio/tinta
```

## Quick start

```go
package main

import "github.com/varavelio/tinta"

func main() {
    // Styled text
    tinta.Text().Red().Bold().Println("fatal: file not found")
    tinta.Text().Green().Println("ok")
    tinta.Text().White().OnBlue().Printf("info: %s\n", "service started")

    // Bordered containers
    tinta.Box().BorderRounded().PaddingY(1).PaddingX(2).Println("hello")
    tinta.Box().BorderDouble().Blue().PaddingX(1).Println("status")
}
```

## Text

`tinta.Text()` is the single entry point. Every chaining method returns a new
immutable `TextStyle`, so branching from a shared base is safe:

```go
base := tinta.Text().Bold()
err  := base.Red()    // bold + red
warn := base.Yellow() // bold + yellow -- base is unchanged

err.Println("error")
warn.Println("warning")
```

16 foreground colors (`Red`, `Green`, `Blue`, ..., `BrightRed`, `BrightGreen`, ...),
16 background colors with `On*` prefix (`OnRed`, `OnBlue`, `OnBrightCyan`, ...),
and 7 modifiers: `Bold`, `Dim`, `Italic`, `Underline`, `Invert`, `Hidden`, `Strike`.

## Box

`tinta.Box()` creates a bordered container. Borders, padding, margin, alignment,
side visibility, shadows, and colors are all chainable:

```go
tinta.Box().BorderRounded().PaddingY(1).PaddingX(2).Blue().Println("hello, world")
```

```
╭────────────────╮
│                │
│  hello, world  │
│                │
╰────────────────╯
```

### Borders

Four built-in styles plus fully custom borders:

| Method            | Glyphs                  |
| ----------------- | ----------------------- |
| `BorderSimple()`  | `┌ ─ ┐ │ └ ┘` (default) |
| `BorderRounded()` | `╭ ─ ╮ │ ╰ ╯`           |
| `BorderDouble()`  | `╔ ═ ╗ ║ ╚ ╝`           |
| `BorderHeavy()`   | `┏ ━ ┓ ┃ ┗ ┛`           |

```go
tinta.Box().Border(tinta.Border{
    TopLeft: "+", TopRight: "+", BottomLeft: "+", BottomRight: "+",
    Horizontal: "-", Vertical: "|",
}).PaddingX(1).Println("custom")
```

### Padding and Margin

Explicit single-value methods -- no CSS-shorthand overloads:

`Padding(n)`, `PaddingTop(n)`, `PaddingBottom(n)`, `PaddingLeft(n)`, `PaddingRight(n)`, `PaddingX(n)`, `PaddingY(n)`

Same pattern for margin: `Margin(n)`, `MarginTop(n)`, ..., `MarginX(n)`, `MarginY(n)`.

### Content alignment

```go
tinta.Box().BorderRounded().Center().PaddingX(1).Println("Hello, World!\nhi\nTinta!")
```

```
╭───────────────╮
│ Hello, World! │
│      hi       │
│    Tinta!     │
╰───────────────╯
```

`Center()` centers all lines. `CenterTrim()` trims whitespace first.
`CenterFirstLine()`, `CenterLastLine()`, and `CenterLine(n)` target specific lines.

### Side visibility

Hide individual border sides with `DisableTop()`, `DisableBottom()`, `DisableLeft()`, `DisableRight()`.
Combine them for blockquotes, heading underlines, and open-corner effects:

```go
// Blockquote
tinta.Box().BorderHeavy().BrightCyan().
    DisableTop().DisableBottom().DisableRight().
    PaddingLeft(1).
    Println("The best way to predict the\nfuture is to invent it.\n-- Alan Kay")
```

```
┃ The best way to predict the
┃ future is to invent it.
┃ -- Alan Kay
```

```go
// Heading underline
tinta.Box().BorderHeavy().BrightYellow().
    DisableTop().DisableLeft().DisableRight().
    Println("Section Title")
```

```
Section Title
━━━━━━━━━━━━━
```

### Shadow

Add an L-shaped 3D shadow to any box. `Shadow` takes a position and a style -- both required:

```go
tinta.Box().BorderRounded().PaddingX(2).
    Shadow(tinta.ShadowBottomRight, tinta.ShadowLight).
    Println("Hello")
```

```
╭─────────╮
│  Hello  │░
╰─────────╯░
 ░░░░░░░░░░░
```

**Positions:** `ShadowBottomRight`, `ShadowBottomLeft`, `ShadowTopRight`, `ShadowTopLeft`.

**Predefined styles:** `ShadowLight` (░), `ShadowMedium` (▒), `ShadowDark` (▓), `ShadowBlock` (█).

For full control, pass a custom `ShadowStyle` with individual corner, horizontal, and vertical glyphs.
Shadow color defaults to bright-black; adjust with `ShadowDim()`, `ShadowBlack()`, or `ShadowBrightBlack()`.

### Nested boxes

Boxes are ANSI-aware and color-safe. Inner resets never corrupt the outer box styling:

```go
inner := tinta.Box().BorderRounded().Green().PaddingX(1).String("Inner box")
tinta.Box().BorderDouble().Blue().PaddingY(1).PaddingX(2).Println(inner)
```

```
╔══════════════════╗
║                  ║
║  ╭───────────╮   ║
║  │ Inner box │   ║
║  ╰───────────╯   ║
║                  ║
╚══════════════════╝
```

## Output

Both `TextStyle` and `BoxStyle` share the same output methods:

| Method                 | Returns        |
| ---------------------- | -------------- |
| `String(s)`            | `string`       |
| `Sprintf(fmt, ...)`    | `string`       |
| `Print(s)`             | --             |
| `Println(s)`           | --             |
| `Printf(fmt, ...)`     | --             |
| `Fprint(w, s)`         | `(int, error)` |
| `Fprintln(w, s)`       | `(int, error)` |
| `Fprintf(w, fmt, ...)` | `(int, error)` |

`Print`/`Println`/`Printf` write to the default output (`os.Stdout`).
`SetOutput(w)` redirects it. `ForceColors(bool)` overrides auto-detection.

## Color detection

ANSI codes are only emitted when the terminal supports them.

**Disabled** when `NO_COLOR`, `NO_COLORS`, `DISABLE_COLORS`, `CLICOLOR=0`, or `TERM=dumb` is set, or the writer is not a TTY.

**Forced** when `FORCE_COLOR` or `CLICOLOR_FORCE` is set. Disable always takes precedence.

## Full API reference

See the complete API documentation on [pkg.go.dev](https://pkg.go.dev/github.com/varavelio/tinta).

## License

This project is released under the MIT License, read more at [LICENSE](LICENSE)
