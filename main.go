package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MAINMENU       = 0
	GAME           = 1
	PAUSE          = 2
	SONG_SELECTION = 3
	OPTIONS        = 4
	CREATE_SONG    = 5
)

func initializeMainMenu() []*TextButton {
	return []*TextButton{NewTextButton("Song Selection", rl.Vector2{float32(630 - rl.MeasureText("Song Selection", 10)), 155}, 10),
		NewTextButton("Options", rl.Vector2{float32(630 - rl.MeasureText("Options", 10)), 175}, 10),
		NewTextButton("Exit to Desktop", rl.Vector2{float32(630 - rl.MeasureText("Exit to Desktop", 10)), 195}, 10)}
}

func initializeSongsMenu() []*TextButton {
	return []*TextButton{
		NewTextButton("25th Hour - Schlatt and Lud's Musical Emporium", rl.Vector2{10, 175}, 10),
		NewTextButton("Back to Main Menu", rl.Vector2{float32(630 - rl.MeasureText("Back to Main Menu", 10)), 340}, 10),
	}
}

func initializeOptionsMenu() []*TextButton {
	return []*TextButton{
		NewTextButton("Back to Main Menu", rl.Vector2{float32(630 - rl.MeasureText("Back to Main Menu", 10)), 340}, 10),
		NewTextButton("Resolution: 1920x1080", rl.Vector2{10, 30}, 10),
	}
}

func initializeGame(filename string) *Game {
	return NewGame(filename)
}

func main() {
	color_bg := rl.Color{69, 57, 71, 255}
	color_text := rl.Color{239, 214, 239, 255}
	color_text_selected := rl.Color{249, 161, 174, 255}

	var resolutionIndex = 1
	resolutions := []rl.Vector2{rl.Vector2{640, 360}, rl.Vector2{1280, 720}, rl.Vector2{1920, 1080}, rl.Vector2{2560, 1440}, rl.Vector2{3840, 2160}}
	var width int32 = int32(resolutions[resolutionIndex].X)
	var height int32 = int32(resolutions[resolutionIndex].Y)
	rl.InitWindow(width, height, "Final - hmi083")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())) * 2)
	defer rl.CloseWindow()
	rl.InitAudioDevice()

	// disables CloseWindow() for ESC
	rl.SetExitKey(0)

	gameState := MAINMENU
	var game Game
	var creation CreateSongSession

	bgImg := rl.LoadImage("resources/sprites/bg.png")
	backg := rl.LoadTextureFromImage(bgImg)
	rl.UnloadImage(bgImg)

	camera := rl.NewCamera2D(rl.Vector2{0, 0}, rl.Vector2{0, 0}, 0, float32(width/640))

	menuButtons := initializeMainMenu()

	for !rl.WindowShouldClose() {
		switch gameState {
		case MAINMENU:
			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].CheckSelected(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
			}
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				for i := 0; i < len(menuButtons); i++ {
					if menuButtons[i].Selected {
						switch i {
						case 0:
							menuButtons = initializeSongsMenu()
							gameState = SONG_SELECTION
							break
						case 1:
							menuButtons = initializeOptionsMenu()
							menuButtons[1].Text = fmt.Sprintf("Resolution: %dx%d", width, height)
							gameState = OPTIONS
							break
						case 2:
							rl.CloseAudioDevice()
							rl.CloseWindow()
							break
						}
					}
				}
			}

			if gameState != MAINMENU {
				break
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].Draw(color_text_selected, color_text)
			}

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()

			break
		case GAME:
			if game.Time > rl.GetMusicTimeLength(game.Music) {

			}
			rl.UpdateMusicStream(game.Music)

			if rl.IsMusicStreamPlaying(game.Music) {
				game.UpdateGame(rl.GetFrameTime())
			}
			if rl.IsKeyPressed(rl.KeyEscape) {
				gameState = PAUSE
				break
			}
			if rl.IsKeyPressed(rl.KeyA) {
				game.HandleInputRune('A')
			}
			if rl.IsKeyPressed(rl.KeyD) {
				game.HandleInputRune('D')
			}
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				fmt.Printf("%f,%f\n", rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom).X, rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom).Y)
				game.HandleInputMouse(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			rl.DrawTexture(backg, 0, 0, rl.White)
			game.Draw()
			game.DrawScore()

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()

			break
		case SONG_SELECTION:
			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].CheckSelected(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
			}
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				for i := 0; i < len(menuButtons); i++ {
					if menuButtons[i].Selected {
						switch i {
						case 0:
							menuButtons = []*TextButton{}
							gameState = GAME
							if rl.IsKeyDown(rl.KeyLeftShift) {
								rl.SetTargetFPS(1000)
								gameState = CREATE_SONG
								creation = *NewCreateSongSession("25th_hour.wav")
							} else {
								game = *initializeGame("25th_hour.wav")
							}
							break
						case 1:
							menuButtons = initializeMainMenu()
							gameState = MAINMENU
							break
						}
					}
				}
			}
			if gameState != SONG_SELECTION {
				break
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].Draw(color_text_selected, color_text)
			}

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()
			break
		case PAUSE:
			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(rl.RayWhite)

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()
			break
		case OPTIONS:
			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].CheckSelected(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
			}
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				for i := 0; i < len(menuButtons); i++ {
					if menuButtons[i].Selected {
						switch i {
						case 0:
							menuButtons = initializeMainMenu()
							gameState = MAINMENU
							break
						case 1:
							resolutionIndex++
							if resolutionIndex >= len(resolutions) {
								resolutionIndex = 0
							}
							width = int32(resolutions[resolutionIndex].X)
							height = int32(resolutions[resolutionIndex].Y)
							rl.SetWindowSize(int(width), int(height))
							camera.Zoom = float32(width / 640)
							menuButtons[1].Text = fmt.Sprintf("Resolution: %dx%d", width, height)
						}
					}
				}
			}
			if gameState != OPTIONS {
				break
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].Draw(color_text_selected, color_text)
			}

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()
			break
		case CREATE_SONG:
			if creation.Time >= rl.GetMusicTimeLength(creation.Music) {
				rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())) * 2)
				creation.CompleteSession()
				gameState = SONG_SELECTION
				menuButtons = initializeSongsMenu()
				creation = CreateSongSession{}
				break
			}
			creation.UpdateSongSession(rl.GetFrameTime())

			rl.UpdateMusicStream(creation.Music)
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !creation.ChalkDown {
				creation.UpdateSongSessionMouse(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
				creation.ChalkDown = true
			}
			if rl.IsMouseButtonUp(rl.MouseButtonLeft) && creation.ChalkDown {
				creation.ChalkDown = false
			}
			if rl.IsKeyPressed(rl.KeyA) {
				creation.UpdateSongSessionRune('A')
			}
			if rl.IsKeyPressed(rl.KeyD) {
				creation.UpdateSongSessionRune('D')
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			rl.DrawTexture(backg, 0, 0, rl.White)

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()
		}
	}

}
