package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand/v2"
)

// could not figure out a good GLSL implementation for this.
// this is sparingly used in the game due to this.
func FrostedGlassTexture(d *rl.Image, shift float32, intensity float32) *rl.Image {
	var newColors = make([]byte, d.Width*d.Height*4)
	for j := int32(0); j < d.Height; j++ {
		for i := int32(0); i < d.Width; i++ {
			y := j
			x := i

			if rand.Float32() < intensity {
				randVec := rl.Vector2Rotate(rl.Vector2{1, 0}, rand.Float32()*(2*rl.Pi))
				x += int32(randVec.X * shift)
				y += int32(randVec.Y * shift)
				if x < 0 {
					x = 0
				} else if x >= d.Width {
					x = d.Width - 1
				}
				if y < 0 {
					y = 0
				} else if y >= d.Height {
					y = d.Height - 1
				}
			}

			newColors[4*(i+j*d.Width)] = rl.GetImageColor(*d, x, y).R
			newColors[4*(i+j*d.Width)+1] = rl.GetImageColor(*d, x, y).G
			newColors[4*(i+j*d.Width)+2] = rl.GetImageColor(*d, x, y).B
			newColors[4*(i+j*d.Width)+3] = rl.GetImageColor(*d, x, y).A
		}
	}

	newImage := rl.NewImage(newColors, d.Width, d.Height, d.Mipmaps, d.Format)
	return newImage
}
