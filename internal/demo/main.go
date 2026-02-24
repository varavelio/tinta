package main

import (
	"fmt"

	t "github.com/varavelio/tinta"
)

func main() {
	t.ForceColors(true)

	// ── Header ──
	t.Box().BorderRounded().BrightWhite().OnBlue().PaddingX(2).Println("tinta demo")
	fmt.Println()

	// ── Text: Base colors ──
	section("Text: Base colors")
	t.Text().Black().OnWhite().Println("  Black  ")
	t.Text().Red().Println("  Red  ")
	t.Text().Green().Println("  Green  ")
	t.Text().Yellow().Println("  Yellow  ")
	t.Text().Blue().Println("  Blue  ")
	t.Text().Magenta().Println("  Magenta  ")
	t.Text().Cyan().Println("  Cyan  ")
	t.Text().White().Println("  White  ")

	fmt.Println()
	section("Text: Bright colors")
	t.Text().BrightBlack().Println("  BrightBlack  ")
	t.Text().BrightRed().Println("  BrightRed  ")
	t.Text().BrightGreen().Println("  BrightGreen  ")
	t.Text().BrightYellow().Println("  BrightYellow  ")
	t.Text().BrightBlue().Println("  BrightBlue  ")
	t.Text().BrightMagenta().Println("  BrightMagenta  ")
	t.Text().BrightCyan().Println("  BrightCyan  ")
	t.Text().BrightWhite().Println("  BrightWhite  ")

	fmt.Println()
	section("Text: Backgrounds")
	t.Text().White().OnRed().Println("  White on Red  ")
	t.Text().Black().OnBrightYellow().Println("  Black on BrightYellow  ")
	t.Text().BrightWhite().OnBlue().Println("  BrightWhite on Blue  ")

	fmt.Println()
	section("Text: Modifiers")
	t.Text().Cyan().Bold().Println("  Bold  ")
	t.Text().Cyan().Dim().Println("  Dim  ")
	t.Text().Cyan().Italic().Println("  Italic  ")
	t.Text().Cyan().Underline().Println("  Underline  ")
	t.Text().Cyan().Strike().Println("  Strike  ")
	t.Text().Cyan().Hidden().Println("  Hidden  ")

	fmt.Println()
	section("Text: Invert")
	t.Text().Red().OnWhite().Invert().Println("  Red on White (inverted)  ")

	fmt.Println()
	section("Text: Chaining")
	t.Text().Red().Bold().Underline().Println("  Red + Bold + Underline  ")
	t.Text().White().OnBlue().Bold().Println("  White on Blue + Bold  ")

	fmt.Println()
	section("Text: Immutability")
	base := t.Text().Red()
	base.Bold().Println("  bold branch  ")
	base.Underline().Println("  underline branch (base unaffected)  ")

	fmt.Println()
	section("Text: String() / Sprintf()")
	msg := t.Text().Green().Bold().String("ok")
	fmt.Printf("  built string: %s\n", msg)
	t.Text().Yellow().Printf("  items: %d, status: %s\n", 42, "ready")

	fmt.Println()
	section("Text: Bold / Underline only")
	t.Text().Bold().Println("  Bold without color  ")
	t.Text().Underline().Println("  Underline without color  ")

	// ── Box demos ──

	fmt.Println()
	section("Box: Border styles")
	t.Box().Println("Simple (default)")
	fmt.Println()
	t.Box().BorderRounded().Println("Rounded")
	fmt.Println()
	t.Box().BorderDouble().Println("Double")
	fmt.Println()
	t.Box().BorderHeavy().Println("Heavy")

	fmt.Println()
	section("Box: Custom border")
	custom := t.Border{
		TopLeft: "*", TopRight: "*", BottomLeft: "*", BottomRight: "*",
		Horizontal: "~", Vertical: "!",
	}
	t.Box().Border(custom).PaddingX(1).Println("Custom border")

	fmt.Println()
	section("Box: Padding")
	t.Box().BorderRounded().PaddingY(1).PaddingX(3).Println("Padded content")

	fmt.Println()
	section("Box: Margin")
	t.Box().BorderRounded().MarginLeft(4).Println("Left margin = 4")

	fmt.Println()
	section("Box: Colored borders")
	t.Box().BorderRounded().Red().Println("Red border")
	fmt.Println()
	t.Box().BorderDouble().Blue().Bold().Println("Blue bold border")
	fmt.Println()
	t.Box().BorderHeavy().Green().OnBlack().PaddingX(1).Println("Green on black")

	fmt.Println()
	section("Box: Styled content inside a box")
	styled := t.Text().Red().Bold().String("Error:") + " " + t.Text().White().String("something broke")
	t.Box().BorderRounded().Yellow().PaddingX(1).Println(styled)

	fmt.Println()
	section("Box: Multiline content")
	t.Box().BorderDouble().Cyan().PaddingY(1).PaddingX(2).Println("Line 1: Hello\nLine 2: World\nLine 3: Tinta!")

	fmt.Println()
	section("Box: Center")
	t.Box().BorderRounded().Center().PaddingX(1).Println("Hello, World!\nhi\nTinta!")

	fmt.Println()
	section("Box: CenterTrim")
	t.Box().BorderRounded().CenterTrim().PaddingX(1).Println("  Hello  \n  hi  \n  Tinta!  ")

	fmt.Println()
	section("Box: Immutability")
	boxBase := t.Box().BorderRounded()
	boxBase.Red().PaddingX(1).Println("Red rounded")
	fmt.Println()
	boxBase.Blue().PaddingX(1).Println("Blue rounded (base unaffected)")

	// ── Disabled sides ──

	fmt.Println()
	section("Box: DisableTop")
	t.Box().BorderRounded().DisableTop().Cyan().PaddingX(1).Println("No top border")

	fmt.Println()
	section("Box: DisableBottom")
	t.Box().BorderRounded().DisableBottom().Magenta().PaddingX(1).Println("No bottom border")

	fmt.Println()
	section("Box: DisableLeft")
	t.Box().BorderRounded().DisableLeft().Yellow().PaddingX(1).Println("No left border")

	fmt.Println()
	section("Box: DisableRight")
	t.Box().BorderRounded().DisableRight().Green().PaddingX(1).Println("No right border")

	fmt.Println()
	section("Box: DisableTop + DisableBottom (horizontal rule effect)")
	t.Box().DisableTop().DisableBottom().Cyan().PaddingX(2).Println("Sandwiched content")

	fmt.Println()
	section("Box: DisableLeft + DisableRight (top/bottom only)")
	t.Box().BorderDouble().DisableLeft().DisableRight().Blue().PaddingX(1).Println("Horizontal frame")

	// ── Shadow effect ──
	// Native 3D shadow with configurable glyphs and position.

	fmt.Println()
	section("Box: Shadow (bottom-right)")
	t.Box().BorderRounded().BrightWhite().PaddingX(2).
		Shadow(t.ShadowBottomRight, t.ShadowLight).
		Println("Shadow box")

	fmt.Println()
	section("Box: Shadow (bottom-left)")
	t.Box().BorderRounded().Cyan().PaddingX(2).
		Shadow(t.ShadowBottomLeft, t.ShadowLight).
		Println("Shadow left")

	fmt.Println()
	section("Box: Shadow (top-right)")
	t.Box().BorderDouble().Magenta().PaddingX(1).
		Shadow(t.ShadowTopRight, t.ShadowLight).
		Println("Top-right")

	fmt.Println()
	section("Box: Shadow (top-left)")
	t.Box().BorderHeavy().Yellow().PaddingX(1).
		Shadow(t.ShadowTopLeft, t.ShadowLight).
		Println("Top-left")

	fmt.Println()
	section("Box: Shadow with custom style (dark blocks)")
	t.Box().BorderRounded().Green().PaddingX(2).
		Shadow(t.ShadowBottomRight, t.ShadowDark).
		Println("Dark shadow")

	fmt.Println()
	section("Box: Shadow with custom style (full blocks)")
	t.Box().BorderRounded().Red().PaddingX(2).
		Shadow(t.ShadowBottomRight, t.ShadowBlock).
		Println("Block shadow")

	fmt.Println()
	section("Box: Shadow with custom style (Rounded corners)")
	t.Box().BorderRounded().Green().PaddingX(2).
		Shadow(t.ShadowBottomRight, t.ShadowStyle{
			TopLeft:     t.BorderRounded.TopLeft,
			TopRight:    t.BorderRounded.TopRight,
			BottomLeft:  t.BorderRounded.BottomLeft,
			BottomRight: t.BorderRounded.BottomRight,
			Horizontal:  t.BorderRounded.Horizontal,
			Vertical:    t.BorderRounded.Vertical,
		}).
		Println("Custom rounded shadow")

	// ── Selective line centering ──

	fmt.Println()
	section("Box: CenterFirstLine (title centering)")
	t.Box().BorderDouble().Blue().PaddingX(1).CenterFirstLine().
		Println("Title\nLeft-aligned body\ncontinues here")

	fmt.Println()
	section("Box: CenterLastLine")
	t.Box().BorderDouble().Green().PaddingX(1).CenterLastLine().
		Println("Body content here\ncontinues...\n-- The End --")

	fmt.Println()
	section("Box: CenterLine(1) — center only middle line")
	t.Box().BorderRounded().Cyan().PaddingX(1).CenterLine(1).
		Println("Top\nCentered\nBottom")

	// ── Nested box (color-safe) ──

	fmt.Println()
	section("Box: Nested boxes (color-safe)")

	// Inner box: a small styled box.
	inner := t.Box().BorderRounded().Green().PaddingX(1).String(
		t.Text().Green().Bold().String("Inner box") + "\n" +
			t.Text().White().String("with content"),
	)

	// Outer box wraps the inner box as its content.
	// The inner box's ANSI resets do NOT corrupt the outer box's styling.
	t.Box().BorderDouble().Blue().PaddingY(1).PaddingX(2).Println(inner)

	// ── Advanced nested: multiple inner boxes ──

	fmt.Println()
	section("Box: Multiple nested boxes")

	box1 := t.Box().BorderRounded().Red().PaddingX(1).String("Alert")
	box2 := t.Box().BorderRounded().Green().PaddingX(1).String("Success")
	box3 := t.Box().BorderRounded().Yellow().PaddingX(1).String("Warning")

	combined := box1 + "\n" + box2 + "\n" + box3
	t.Box().BorderHeavy().Cyan().PaddingY(1).PaddingX(2).Println(combined)

	// ── Open corner effect ──

	fmt.Println()
	section("Box: Open corner (DisableTop + DisableLeft)")
	t.Box().BorderHeavy().Red().DisableTop().DisableLeft().PaddingX(1).Println("Open top-left\ncorner effect")

	fmt.Println()
	section("Box: L-shape (DisableTop + DisableRight)")
	t.Box().BorderDouble().Magenta().DisableTop().DisableRight().PaddingX(1).Println("L-shape border\neffect")

	// ── Quote / blockquote effect ──
	// Only left border visible: disable top, bottom, and right.

	fmt.Println()
	section("Box: Blockquote (left border only)")
	t.Box().BorderHeavy().BrightCyan().
		DisableTop().DisableBottom().DisableRight().
		PaddingLeft(1).
		Println("The best way to predict the\nfuture is to invent it.\n— Alan Kay")

	// ── Underline / heading effect ──
	// Only bottom border visible.

	fmt.Println()
	section("Box: Heading underline (bottom border only)")
	t.Box().BorderHeavy().BrightYellow().
		DisableTop().DisableLeft().DisableRight().
		Println("Section Title")
}

func section(label string) {
	t.Text().BrightBlue().Bold().Printf("== %s ==\n", label)
}
