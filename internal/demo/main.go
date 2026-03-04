package main

import (
	"fmt"

	t "github.com/varavelio/tinta"
)

func main() {
	t.ForceColors(true)

	t.Box().Border(t.BorderRounded).BrightWhite().OnBlue().PaddingX(2).Println("tinta demo")
	fmt.Println()

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

	fmt.Println()
	section("Box: Border styles")
	t.Box().Println("Simple (default)")
	fmt.Println()
	t.Box().Border(t.BorderRounded).Println("Rounded")
	fmt.Println()
	t.Box().Border(t.BorderDouble).Println("Double")
	fmt.Println()
	t.Box().Border(t.BorderHeavy).Println("Heavy")

	fmt.Println()
	section("Box: Custom border")
	custom := t.Border{
		TopLeft: "*", TopRight: "*", BottomLeft: "*", BottomRight: "*",
		Top: "~", Left: "!", Right: "!", Bottom: "~",
	}
	t.Box().Border(custom).PaddingX(1).Println("Custom border")

	fmt.Println()
	section("Box: Padding")
	t.Box().Border(t.BorderRounded).PaddingY(1).PaddingX(3).Println("Padded content")

	fmt.Println()
	section("Box: Margin")
	t.Box().Border(t.BorderRounded).MarginLeft(4).Println("Left margin = 4")

	fmt.Println()
	section("Box: Colored borders")
	t.Box().Border(t.BorderRounded).Red().Println("Red border")
	fmt.Println()
	t.Box().Border(t.BorderDouble).Blue().Bold().Println("Blue bold border")
	fmt.Println()
	t.Box().Border(t.BorderHeavy).Green().OnBlack().PaddingX(1).Println("Green on black")

	fmt.Println()
	section("Box: Styled content inside a box")
	styled := t.Text().Red().Bold().String("Error:") + " " + t.Text().White().String("something broke")
	t.Box().Border(t.BorderRounded).Yellow().PaddingX(1).Println(styled)

	fmt.Println()
	section("Box: Multiline content")
	t.Box().Border(t.BorderDouble).Cyan().PaddingY(1).PaddingX(2).Println("Line 1: Hello\nLine 2: World\nLine 3: Tinta!")

	fmt.Println()
	section("Box: Center")
	t.Box().Border(t.BorderRounded).Center().PaddingX(1).Println("Hello, World!\nhi\nTinta!")

	fmt.Println()
	section("Box: CenterTrim")
	t.Box().Border(t.BorderRounded).CenterTrim().PaddingX(1).Println("  Hello  \n  hi  \n  Tinta!  ")

	fmt.Println()
	section("Box: Immutability")
	boxBase := t.Box().Border(t.BorderRounded)
	boxBase.Red().PaddingX(1).Println("Red rounded")
	fmt.Println()
	boxBase.Blue().PaddingX(1).Println("Blue rounded (base unaffected)")

	fmt.Println()
	section("Box: DisableTop")
	t.Box().Border(t.BorderRounded).DisableTop().Cyan().PaddingX(1).Println("No top border")

	fmt.Println()
	section("Box: DisableBottom")
	t.Box().Border(t.BorderRounded).DisableBottom().Magenta().PaddingX(1).Println("No bottom border")

	fmt.Println()
	section("Box: DisableLeft")
	t.Box().Border(t.BorderRounded).DisableLeft().Yellow().PaddingX(1).Println("No left border")

	fmt.Println()
	section("Box: DisableRight")
	t.Box().Border(t.BorderRounded).DisableRight().Green().PaddingX(1).Println("No right border")

	fmt.Println()
	section("Box: DisableTop + DisableBottom (horizontal rule effect)")
	t.Box().DisableTop().DisableBottom().Cyan().PaddingX(2).Println("Sandwiched content")

	fmt.Println()
	section("Box: DisableLeft + DisableRight (top/bottom only)")
	t.Box().Border(t.BorderDouble).DisableLeft().DisableRight().Blue().PaddingX(1).Println("Horizontal frame")

	fmt.Println()
	section("Canvas: Faux shadow with layered boxes")
	shadowInner := t.Box().Border(t.BorderRounded).BrightWhite().PaddingX(2).String("Shadow box")
	shadowOuter := t.Box().Border(t.BorderBlockLight).PaddingX(3).PaddingY(1).String("Shadow box")
	fmt.Println(t.Canvas().
		Add(shadowOuter, 1, 1).
		Add(shadowInner, 0, 0).
		String())

	fmt.Println()
	section("Box: Independent corner control")
	t.Box().Border(t.BorderHeavy).
		DisableTopLeftCorner().DisableBottomRightCorner().
		PaddingX(1).
		Println("Custom corner mask")

	fmt.Println()
	section("Box: Title (left / center / right)")
	t.Box().Border(t.BorderRounded).Cyan().PaddingX(1).
		Title("Left Title", t.AlignLeft).
		Println("Content with a left-aligned title")
	fmt.Println()
	t.Box().Border(t.BorderRounded).Green().PaddingX(1).
		Title("Center Title", t.AlignCenter).
		Println("Content with a centered title")
	fmt.Println()
	t.Box().Border(t.BorderRounded).Yellow().PaddingX(1).
		Title("Right Title", t.AlignRight).
		Println("Content with a right-aligned title")

	fmt.Println()
	section("Box: Footer (left / center / right)")
	t.Box().Border(t.BorderDouble).Blue().PaddingX(1).
		Footer("v1.0.0", t.AlignLeft).
		Println("Content with a left-aligned footer")
	fmt.Println()
	t.Box().Border(t.BorderDouble).Magenta().PaddingX(1).
		Footer("Page 1/3", t.AlignCenter).
		Println("Content with a centered footer")
	fmt.Println()
	t.Box().Border(t.BorderDouble).Red().PaddingX(1).
		Footer("Status: OK", t.AlignRight).
		Println("Content with a right-aligned footer")

	fmt.Println()
	section("Box: Title + Footer combined")
	t.Box().Border(t.BorderHeavy).BrightCyan().PaddingX(2).PaddingY(1).
		Title("Deploy Summary", t.AlignCenter).
		Footer("Done", t.AlignRight).
		Println("Title + Footer combined\nservice: api\nresult: success\nduration: 12.4s")

	fmt.Println()
	section("Box: Title widens narrow content")
	t.Box().Border(t.BorderRounded).BrightYellow().PaddingX(1).
		Title("A Very Long Title That Widens The Box", t.AlignCenter).
		Println("hi")

	fmt.Println()
	section("Box: Footer widens narrow content")
	t.Box().Border(t.BorderRounded).BrightYellow().PaddingX(1).
		Footer("A Very Long Footer That Widens The Box", t.AlignCenter).
		Println("hi")

	fmt.Println()
	section("Box: CenterFirstLine (title centering)")
	t.Box().Border(t.BorderDouble).Blue().PaddingX(1).CenterFirstLine().
		Println("Title\nLeft-aligned body\ncontinues here")

	fmt.Println()
	section("Box: CenterLastLine")
	t.Box().Border(t.BorderDouble).Green().PaddingX(1).CenterLastLine().
		Println("Body content here\ncontinues...\n-- The End --")

	fmt.Println()
	section("Box: CenterLine(1) — center only middle line")
	t.Box().Border(t.BorderRounded).Cyan().PaddingX(1).CenterLine(1).
		Println("Top\nCentered\nBottom")

	fmt.Println()
	section("Box: Nested boxes (color-safe)")

	inner := t.Box().Border(t.BorderRounded).Green().PaddingX(1).String(
		t.Text().Green().Bold().String("Inner box") + "\n" +
			t.Text().White().String("with content"),
	)

	t.Box().Border(t.BorderDouble).Blue().PaddingY(1).PaddingX(2).Println(inner)

	fmt.Println()
	section("Box: Multiple nested boxes")

	box1 := t.Box().Border(t.BorderRounded).Red().PaddingX(1).String("Alert")
	box2 := t.Box().Border(t.BorderRounded).Green().PaddingX(1).String("Success")
	box3 := t.Box().Border(t.BorderRounded).Yellow().PaddingX(1).String("Warning")

	combined := box1 + "\n" + box2 + "\n" + box3
	t.Box().Border(t.BorderHeavy).Cyan().PaddingY(1).PaddingX(2).Println(combined)

	fmt.Println()
	section("Box: Open corner (DisableTop + DisableLeft)")
	t.Box().Border(t.BorderHeavy).Red().DisableTop().DisableLeft().PaddingX(1).Println("Open top-left\ncorner effect")

	fmt.Println()
	section("Box: L-shape (DisableTop + DisableRight)")
	t.Box().Border(t.BorderDouble).Magenta().DisableTop().DisableRight().PaddingX(1).Println("L-shape border\neffect")

	fmt.Println()
	section("Box: Blockquote (left border only)")
	t.Box().Border(t.BorderHeavy).BrightCyan().
		DisableTop().DisableBottom().DisableRight().
		PaddingLeft(1).
		Println("The best way to predict the\nfuture is to invent it.\n— Alan Kay")

	fmt.Println()
	section("Box: Heading underline (bottom border only)")
	t.Box().Border(t.BorderHeavy).BrightYellow().
		DisableTop().DisableLeft().DisableRight().
		Println("Section Title")

	fmt.Println()
	section("Canvas: 3D border effect")
	text := "Lorem ipsum"

	front := t.Box().
		Border(t.BorderHeavy).
		BrightYellow().
		PaddingX(5).
		PaddingY(1).
		String(text)
	shadow1 := t.Box().
		Border(t.BorderRounded).
		Red().
		PaddingX(5).
		PaddingY(1).
		String(text)
	shadow2 := t.Box().
		Border(t.BorderRounded).
		Blue().
		PaddingX(5).
		PaddingY(1).
		String(text)
	shadow3 := t.Box().
		Border(t.BorderRounded).
		Yellow().
		PaddingX(5).
		PaddingY(1).
		String(text)
	fmt.Println(t.Canvas().
		Add(shadow3, 3, -1).
		Add(shadow2, 4, 2).
		Add(shadow1, 6, 1).
		Add(front, 5, 0).
		String())
}

func section(label string) {
	t.Text().BrightBlue().Bold().Printf("== %s ==\n", label)
}
