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
		NewTextButton("Back to Main Menu", rl.Vector2{float32(630 - rl.MeasureText("Back to Main Menu", 10)), 340}, 10),
		NewTextButton("The 25th Hour - Schlatt and Lud's Musical Emporium", rl.Vector2{10, 175}, 10),
		/*
			"I Got a Stick Feat James Gavins" Kevin MacLeod (incompetech.com)
			Licensed under Creative Commons: By Attribution 4.0 License
			http://creativecommons.org/licenses/by/4.0/
		*/
		NewTextButton("I Got a Stick (ft. James Gavins) - Kevin McLeod", rl.Vector2{10, 195}, 10),
	}
}

func initializeOptionsMenu() []*TextButton {
	return []*TextButton{
		NewTextButton("Back to Main Menu", rl.Vector2{float32(630 - rl.MeasureText("Back to Main Menu", 10)), 340}, 10),
		NewTextButton("Resolution: 1920x1080", rl.Vector2{10, 30}, 10),
	}
}

func initializePauseMenu() []*TextButton {
	return []*TextButton{
		NewTextButton("Resume", rl.Vector2{float32(640/2 - rl.MeasureText("Resume", 10)/2), 155}, 10),
		NewTextButton("Return to Song Selection", rl.Vector2{float32(640/2 - rl.MeasureText("Return to Song Selection", 10)/2), 175}, 10),
		NewTextButton("Return to Desktop", rl.Vector2{float32(640/2 - rl.MeasureText("Return to Desktop", 10)/2), 195}, 10),
	}
}

func initializeGame(filename string) *Game {
	return NewGame(filename)
}

// i initially put this in a chalk.go, but it made me sad
type ChalkMark struct {
	Position rl.Vector2
	Time     float32
	Color    rl.Color
}

func main() {
	color_bg := rl.Color{69, 57, 71, 255}
	color_text := rl.Color{239, 214, 239, 255}
	color_text_selected := rl.Color{249, 161, 174, 255}

	var resolutionIndex = 1
	resolutions := []rl.Vector2{rl.Vector2{640, 360}, rl.Vector2{1280, 720}, rl.Vector2{1920, 1080}, rl.Vector2{2560, 1440}, rl.Vector2{3840, 2160}}
	var width int32 = int32(resolutions[resolutionIndex].X)
	var height int32 = int32(resolutions[resolutionIndex].Y)
	rl.InitWindow(width, height, "Skotched Rhythm - hmi083")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())) * 2)
	defer rl.CloseWindow()
	rl.InitAudioDevice()

	// disables CloseWindow() for ESC
	rl.SetExitKey(0)

	gameState := MAINMENU
	var game Game
	var creation CreateSongSession
	var chalkMarks = []ChalkMark{}

	chalkText := rl.LoadTexture("resources/sprites/chalk.png")
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
		default:
			if game.Time > rl.GetMusicTimeLength(game.Music) {

			}
			if gameState == GAME {
				rl.UpdateMusicStream(game.Music)

				if rl.IsMusicStreamPlaying(game.Music) {
					game.UpdateGame(rl.GetFrameTime())
				}
			} else {
				for i := 0; i < len(menuButtons); i++ {
					menuButtons[i].CheckSelected(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
				}
				if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
					for i := 0; i < len(menuButtons); i++ {
						if menuButtons[i].Selected {
							switch i {
							case 0:
								gameState = GAME
								break
							case 1:
								menuButtons = initializeSongsMenu()
								gameState = SONG_SELECTION
								break
							case 2:
								rl.CloseAudioDevice()
								rl.CloseWindow()
								break
							}
						}
					}
				}

			}
			if rl.IsKeyPressed(rl.KeyEscape) {
				if gameState == PAUSE {
					gameState = GAME
				} else {
					gameState = PAUSE
				}
			}
			if rl.IsKeyPressed(rl.KeyA) {
				game.HandleInputRune('A')
			}
			if rl.IsKeyPressed(rl.KeyD) {
				game.HandleInputRune('D')
			}
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				fmt.Printf("%f,%f\n", rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom).X, rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom).Y)
				var good = game.HandleInputMouse(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
				if rl.CheckCollisionPointRec(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom), rl.Rectangle{423, 33, 190, 112}) {
					color := rl.Color{186, 242, 164, 255}
					if !good {
						color = rl.Color{248, 162, 162, 255}
					}
					chalkMarks = append(chalkMarks, ChalkMark{rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom), 1.5, color})
				}
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			rl.DrawTexture(backg, 0, 0, rl.White)
			game.Draw()
			game.DrawScore()

			if gameState == PAUSE {
				rl.DrawRectangle(0, 0, 640, 360, rl.Color{0, 0, 0, 128})
				for i := 0; i < len(menuButtons); i++ {
					menuButtons[i].Draw(color_text_selected, color_text)
				}
			} else {
				for i := 0; i < len(chalkMarks); i++ {
					chalkMarks[i].Time -= rl.GetFrameTime()
					color := chalkMarks[i].Color
					if chalkMarks[i].Time <= 1 {
						color = rl.ColorAlpha(chalkMarks[i].Color, chalkMarks[i].Time)
					}
					rl.DrawTexture(chalkText, int32(chalkMarks[i].Position.X)-chalkText.Width/2, int32(chalkMarks[i].Position.Y)-chalkText.Height/2, color)
				}
			}

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()

			for i := len(chalkMarks) - 1; i >= 0; i-- {
				if chalkMarks[i].Time < 0 {
					chalkMarks = append(chalkMarks[:i], chalkMarks[i+1:]...)
				}
			}

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
							menuButtons = initializeMainMenu()
							gameState = MAINMENU
							break
						case 1:
							menuButtons = initializePauseMenu()
							gameState = GAME
							if rl.IsKeyDown(rl.KeyLeftShift) {
								rl.SetTargetFPS(1000)
								gameState = CREATE_SONG
								creation = *NewCreateSongSession("25th_hour.wav")
							} else {
								game = *initializeGame("25th_hour.wav")
							}
							break
						case 2:
							menuButtons = initializePauseMenu()
							gameState = GAME
							if rl.IsKeyDown(rl.KeyLeftShift) {
								rl.SetTargetFPS(1000)
								gameState = CREATE_SONG
								creation = *NewCreateSongSession("i_got_a_stick.mp3")
							} else {
								game = *initializeGame("i_got_a_stick.mp3")
							}
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
