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
	FINISH_SONG    = 6
)

func initializeMainMenu() []*TextButton {
	return []*TextButton{NewTextButton("Song Selection", rl.Vector2{float32(440 - rl.MeasureText("Song Selection", 10)/2), 225}, 10),
		NewTextButton("Options", rl.Vector2{float32(440 - rl.MeasureText("Options", 10)/2), 245}, 10),
		NewTextButton("Exit to Desktop", rl.Vector2{float32(440 - rl.MeasureText("Exit to Desktop", 10)/2), 265}, 10)}
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

func initializeFinishMenu() []*TextButton {
	return []*TextButton{
		NewTextButton("Restart", rl.Vector2{float32(640/2 - rl.MeasureText("Restart", 10)/2), 200}, 10),
		NewTextButton("Return to Song Selection", rl.Vector2{float32(640/2 - rl.MeasureText("Return to Song Selection", 10)/2), 220}, 10),
		NewTextButton("Return to Desktop", rl.Vector2{float32(640/2 - rl.MeasureText("Return to Desktop", 10)/2), 240}, 10),
	}
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
	logoText := rl.LoadTexture("resources/sprites/logo.png")
	rl.UnloadImage(bgImg)

	rightFoot := rl.LoadTexture("resources/sprites/right_on.png")
	leftFoot := rl.LoadTexture("resources/sprites/left_on.png")

	misinputSound := rl.LoadSound("resources/audio/misinput.mp3")
	rl.SetSoundVolume(misinputSound, 0.15)

	idleAnim := NewAnim("idle-Sheet.png", 4, false, 0.1, rl.Color{251, 213, 220, 255})
	// blinkAnim := NewAnim("blink-Sheet.png", 4, false, 0.2)
	leftAnim := NewAnim("left-Sheet.png", 4, true, 0.1, rl.Color{232, 146, 146, 255})
	rightAnim := NewAnim("right-Sheet.png", 4, true, 0.1, rl.Color{146, 177, 232, 255})
	bothAnim := NewAnim("both-Sheet.png", 4, true, 0.1, rl.Color{253, 212, 141, 255})
	missAnim := NewAnim("miss-Sheet.png", 5, true, 0.1, rl.Color{253, 143, 141, 255})

	gameAnim := bothAnim

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
			rl.DrawTexture(logoText, 340, 80, rl.White)
			idleAnim.Draw(rl.Vector2{0, -120}, rl.GetFrameTime())

			rl.EndMode2D()
			rl.DrawFPS(10, 10)
			rl.EndDrawing()

			break
		case FINISH_SONG:
			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].CheckSelected(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
			}
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				for i := 0; i < len(menuButtons); i++ {
					if menuButtons[i].Selected {
						switch i {
						case 0:
							gameState = GAME
							game = *initializeGame(game.Filename)
							gameAnim = bothAnim
							gameAnim.Time = 0.0
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

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			for i := 0; i < len(menuButtons); i++ {
				menuButtons[i].Draw(color_text_selected, color_text)
			}

			rl.EndMode2D()
			rl.EndDrawing()
		default:
			if game.Time > rl.GetMusicTimeLength(game.Music) || (gameState == GAME && !rl.IsMusicStreamPlaying(game.Music)) {
				menuButtons = initializeFinishMenu()
				gameState = FINISH_SONG
				break
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
				if !game.HandleInputRune('A') && gameState == GAME {
					rl.PlaySound(misinputSound)
					game.Misses++
					gameAnim = missAnim
					gameAnim.Time = 0.0
				}
				if game.HandleInputRune('A') && gameState == GAME {
					gameAnim = leftAnim
					gameAnim.Time = 0.0
				}
			}
			if rl.IsKeyPressed(rl.KeyD) {
				if !game.HandleInputRune('D') && gameState == GAME {
					rl.PlaySound(misinputSound)
					game.Misses++
					gameAnim = missAnim
					gameAnim.Time = 0.0
				}
				if game.HandleInputRune('D') && gameState == GAME {
					gameAnim = rightAnim
					gameAnim.Time = 0.0
				}
			}
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				fmt.Printf("%f,%f\n", rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom).X, rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom).Y)
				var good = game.HandleInputMouse(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom))
				if rl.CheckCollisionPointRec(rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom), rl.Rectangle{416, 26, 200, 126}) {
					color := rl.Color{186, 242, 164, 255}
					if !good {
						color = rl.Color{248, 162, 162, 255}
						if gameState == GAME {
							rl.PlaySound(misinputSound)
							game.Misses++
							gameAnim = missAnim
							gameAnim.Time = 0.0
						}
					}
					if gameState == GAME {
						chalkMarks = append(chalkMarks, ChalkMark{rl.Vector2Scale(rl.GetMousePosition(), 1/camera.Zoom), 1.5, color})
					}
				}
			}

			rl.BeginDrawing()
			rl.BeginMode2D(camera)
			rl.ClearBackground(color_bg)

			rl.DrawTexture(backg, 0, 0, rl.White)
			game.Draw()
			game.DrawScore()

			gameAnim.Draw(rl.Vector2{10, -120}, rl.GetFrameTime())

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
				if rl.IsKeyDown(rl.KeyA) {
					rl.DrawTexture(leftFoot, 0, 0, rl.White)
				}
				if rl.IsKeyDown(rl.KeyD) {
					rl.DrawTexture(rightFoot, 0, 0, rl.White)
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
						gameAnim = bothAnim
						gameAnim.Time = 0.0
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
