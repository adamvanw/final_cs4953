package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextButton struct {
	Text     string
	position rl.Vector2
	fontSize int32
	size     rl.Vector2
	Selected bool
}

func NewTextButton(text string, position rl.Vector2, fontSize int32) *TextButton {
	size := rl.Vector2{float32(rl.MeasureText(text, fontSize)), float32(fontSize)}
	return &TextButton{text, position, fontSize, size, false}
}

func (tb *TextButton) CheckSelected(mousePos rl.Vector2) bool {
	tb.Selected = rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(tb.position.X, tb.position.Y, tb.size.X, tb.size.Y))
	return tb.Selected
}

func (tb *TextButton) Draw(selectedColor rl.Color, notSelectedColor rl.Color) {
	var color = notSelectedColor
	if tb.Selected {
		color = selectedColor
	}
	rl.DrawText(tb.Text, int32(tb.position.X), int32(tb.position.Y), tb.fontSize, color)
}
