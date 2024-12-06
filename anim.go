package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type Anim struct {
	texture       rl.Texture2D
	frames        int
	hold          bool
	Time          float32
	frameDuration float32
	color         rl.Color
}

func NewAnim(filename string, frames int, hold bool, frameDuration float32, color rl.Color) *Anim {
	return &Anim{rl.LoadTexture(fmt.Sprintf("resources/sprites/%s", filename)), frames, hold, 0.0, frameDuration, color}
}

func (a *Anim) Draw(position rl.Vector2, frameTime float32) {
	a.Time += frameTime
	if a.Time >= a.frameDuration*float32(a.frames) && !a.hold {
		a.Time = 0
	}
	index := int32(math.Floor(float64(a.Time / a.frameDuration)))
	if a.hold && index > int32(a.frames-1) {
		index = int32(a.frames) - 1
	}
	rl.DrawTexturePro(a.texture, rl.Rectangle{float32(index * (a.texture.Width / int32(a.frames))), 0, float32(a.texture.Width / int32(a.frames)), float32(a.texture.Height)}, rl.Rectangle{0, 0, float32((a.texture.Width / int32(a.frames)) / 2), float32(a.texture.Height / 2)}, position, 0, a.color)
}
