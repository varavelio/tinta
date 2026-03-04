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

Tinta is built around three composable primitives:

- `Text`: ANSI colors and modifiers
- `Box`: structured frames with spacing, alignment, corners, and title/footer
- `Canvas`: layered 2D composition with z-index ordering

## Install

```sh
go get github.com/varavelio/tinta
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/varavelio/tinta"
)

func main() {
	tinta.Text().Red().Bold().Println("error: file not found")

	panel := tinta.Box().
		Border(tinta.BorderRounded).
		PaddingX(1).
		Title("Status", tinta.AlignLeft).
		Footer("ok", tinta.AlignRight).
		String("service started")

	fmt.Println(panel)
}
```

## Text

`Text()` is immutable and safe to branch.

```go
base := tinta.Text().Bold()
base.Red().Println("error")
base.Green().Println("ok")
```

## Box

`Box()` supports:

- independent border sides (`Top`, `Left`, `Right`, `Bottom` and corners)
- side visibility controls (`DisableTop`, `DisableBottom`, `DisableLeft`, `DisableRight`)
- independent corner controls (`DisableTopLeftCorner`, etc.)
- top/bottom border labels (`Title`, `Footer`) with `AlignLeft`, `AlignCenter`, `AlignRight`

```go
custom := tinta.Border{
	TopLeft: "+",
	TopRight: "+",
	BottomLeft: "+",
	BottomRight: "+",
	Top: "-",
	Left: "|",
	Right: "!",
	Bottom: "~",
}

tinta.Box().
	Border(custom).
	PaddingX(1).
	Title("Example", tinta.AlignCenter).
	Println("custom frame")
```

Corner behavior is explicit: corners render as long as they are not explicitly disabled and at least one adjacent side is visible.

All these borders are already included:

- `tinta.BorderSimple`
- `tinta.BorderDashed`
- `tinta.BorderDotted`
- `tinta.BorderRounded`
- `tinta.BorderRoundedDashed`
- `tinta.BorderRoundedDotted`
- `tinta.BorderDouble`
- `tinta.BorderHeavy`
- `tinta.BorderASCII`
- `tinta.BorderBlock`
- `tinta.BorderBlockHalf`
- `tinta.BorderBlockLight`
- `tinta.BorderBlockMedium`
- `tinta.BorderBlockDark`

## Canvas

`Canvas()` composites pre-rendered strings (commonly boxes) into a layered output.

```go
front := tinta.Box().Border(tinta.BorderHeavy).PaddingX(2).String("front")
back := tinta.Box().Border(tinta.BorderRounded).PaddingX(2).String("back")

out := tinta.Canvas().
	Add(back, -1, 1).
	Add(front, 0, 0).
	String()
```

By default, negative `x/y` coordinates expand the canvas to fit all content. Fixed `Width(...)` and `Height(...)` apply cropping.

## Output and Color Control

- `Print*` methods write to the package output writer (`os.Stdout` by default)
- `SetOutput(w)` redirects default output
- `ForceColors(true|false)` overrides automatic color detection

## API Reference

For the complete API surface and method-level documentation, see:

- [pkg.go.dev/github.com/varavelio/tinta](https://pkg.go.dev/github.com/varavelio/tinta)

## License

This project is released under the MIT License. See [LICENSE](LICENSE).
