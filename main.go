package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync"
)

const (
	MAINMENU = 0
	OPTIONS  = 1
	GAME     = 2
	PAUSE    = 3
)

func main() {
	monitorInt := rl.GetCurrentMonitor()
	width, height := 1280, 720
	rl.InitWindow(int32(width), int32(height), "Final Project - hmi083")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(monitorInt) * 2))
	defer rl.CloseWindow()

	var wg sync.WaitGroup

	camera := rl.NewCamera2D(rl.Vector2{0, 0}, rl.Vector2{0, 0}, 0, float32(height/360))

	strokeImgs := []*rl.Image{}
	frostedImgs := []rl.Image{}
	strokeTexts := []rl.Texture2D{}
	for i := 0; i < 1; i++ {
		strokeImgs = append(strokeImgs, rl.LoadImage("resources/images/stroke_test.png"))
		rl.ImageBlurGaussian(strokeImgs[i], 1)
		frostedImgs = append(frostedImgs, *FrostedGlassTexture(strokeImgs[i], 5, 0.5))
		strokeTexts = append(strokeTexts, rl.LoadTextureFromImage(&frostedImgs[i]))
	}

	gameState := MAINMENU
	blurTimer := float64(0.0)

	for !rl.WindowShouldClose() {
		blurTimer += float64(rl.GetFrameTime())
		if blurTimer > 0.1 {
			blurTimer = 0.0

			for i := 0; i < len(strokeTexts); i++ {
				frostedImgs[i] = *FrostedGlassTexture(strokeImgs[i], 5, 0.5)
				colors := rl.LoadImageColors(&frostedImgs[i])
				rl.UpdateTexture(strokeTexts[i], colors)
				rl.UnloadImageColors(colors)
			}
		}

		switch gameState {
		case MAINMENU:
			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(rl.RayWhite)

			for i := 0; i < len(strokeTexts); i++ {
				rl.DrawTexture(strokeTexts[i], 0, 0, rl.White)
			}

			rl.EndMode2D()

			rl.DrawFPS(10, 10)
			rl.EndDrawing()
		case OPTIONS:
			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			rl.EndDrawing()
		case GAME:
			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			rl.EndDrawing()
		case PAUSE:
			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			rl.EndDrawing()
		}
		wg.Wait()
	}
}
