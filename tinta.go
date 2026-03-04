// Package tinta provides a minimal, chainable terminal styling system
// built on three pillars: content ([Text]), structure ([Box]), and
// composition ([Canvas]).
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
// Use [Box] to add structure and layout to your output:
//
//	tinta.Box().Border(tinta.BorderRounded).Padding(1).PaddingLeft(2).Println("hello")
//	tinta.Box().Border(tinta.BorderDouble).Blue().PaddingX(1).Println("status")
//
// # Canvas: layer compositor
//
// Use [Canvas] to composite multiple rendered strings into layered 2D output:
//
//	front := tinta.Box().Border(tinta.BorderHeavy).PaddingX(3).String("hello")
//	shadow := tinta.Box().Border(tinta.BorderRounded).PaddingX(3).String("hello")
//	tinta.Canvas().Add(shadow, 1, 1).Add(front, 0, 0).String()
//
// The default output is [os.Stdout]. Change it with [SetOutput].
// Color support is detected automatically. Override with [ForceColors].
package tinta
