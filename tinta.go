// Package tinta provides a minimal, chainable terminal styling system
// built on two pillars: content ([Text]) and structure ([Box]).
//
// # Text: ANSI colors and modifiers
//
// Use [Text] as the single entry point for styled text:
//
//	tinta.Text().Red().Bold().Println("error: something broke")
//	msg := tinta.Text().Green().Bold().String("ok")
//
// # Box: layout, borders and spacing
//
//	tinta.Box().Rounded().Padding(1).PaddingX(2).Println("hello")
//	tinta.Box().Double().Blue().PaddingX(1).Println("status")
//
// The default output is [os.Stdout]. Change it with [SetOutput].
// Color support is detected automatically. Override with [ForceColors].
package tinta
