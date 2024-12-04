package main

import (
	"bufio"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"os"
)

type Game struct {
	Music        rl.Music
	ATimes       []float32
	DTimes       []float32
	Drawings     []rl.Vector2
	DrawingTimes []float32
	Time         float32
	Footsteps    [2]rl.Texture2D

	AScores       []float32
	DScores       []float32
	DrawingScores []float32
}

func NewGame(filename string) *Game {
	music := rl.LoadMusicStream(fmt.Sprintf("resources/audio/%s", filename))
	file, err := os.Open(fmt.Sprintf("resources/sheets/%s.sheet", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var ATimes []float32
	var DTimes []float32
	var Drawings []rl.Vector2
	var DrawingTimes []float32

	reader := bufio.NewReader(file)
	var readString string
	readString, err = reader.ReadString('\n')
	fmt.Println(readString)
	for err == nil {
		if readString[0] == 'M' {
			var drawing rl.Vector2
			var time float32
			fmt.Sscanf(readString, "M,%f,%f,%f", &drawing.X, &drawing.Y, &time)
			Drawings = append(Drawings, drawing)
			DrawingTimes = append(DrawingTimes, time)
			fmt.Printf("DrawingRead\n")
		} else {
			var inputType rune
			var time float32
			fmt.Sscanf(readString, "%c,%f", &inputType, &time)
			if inputType == 'A' {
				ATimes = append(ATimes, time)
			} else {
				DTimes = append(DTimes, time)
			}
			fmt.Printf("RuneRead\n")
		}
		readString, err = reader.ReadString('\n')
	}

	fmt.Printf("Loaded Game. Here's the amount of inputs for the game:\n")
	fmt.Printf("A: %d\n", len(ATimes))
	fmt.Printf("D: %d\n", len(DTimes))
	fmt.Printf("Chalk inputs: %d\n", len(Drawings))

	rl.PlayMusicStream(music)

	footstep1 := rl.LoadImage("resources/sprites/footstep_1.png")
	footstep2 := rl.LoadImage("resources/sprites/footstep_2.png")

	footsteps := [2]rl.Texture2D{rl.LoadTextureFromImage(footstep1), rl.LoadTextureFromImage(footstep2)}

	rl.UnloadImage(footstep1)
	rl.UnloadImage(footstep2)

	AScores := make([]float32, len(ATimes))
	for i := 0; i < len(ATimes); i++ {
		AScores[i] = -1.0
	}
	DScores := make([]float32, len(DTimes))
	for i := 0; i < len(DTimes); i++ {
		DScores[i] = -1.0
	}
	DrawingScores := make([]float32, len(DrawingTimes))
	for i := 0; i < len(DrawingTimes); i++ {
		DrawingScores[i] = -1.0
	}

	return &Game{music, ATimes, DTimes, Drawings, DrawingTimes, 0.0, footsteps, AScores, DScores, DrawingScores}
}

func (g *Game) UpdateGame(frameTime float32) {
	g.Time += frameTime
}

func (g *Game) Draw() {
	for i := 0; i < len(g.ATimes); i++ {
		if g.AScores[i] != -1 {
			continue
		}
		if g.ATimes[i] <= g.Time+0.85 && g.ATimes[i] > g.Time-0.75 {
			color := rl.White
			if g.ATimes[i] <= g.Time+0.25 && g.ATimes[i] > g.Time-0.25 {
				color = rl.Gray
			}
			rl.DrawTexture(g.Footsteps[0], 208, int32(300+((g.Time-g.ATimes[i])/0.75)*400), color)
		}
	}
	for i := 0; i < len(g.DTimes); i++ {
		if g.DScores[i] != -1 {
			continue
		}
		if g.DTimes[i] <= g.Time+0.85 && g.DTimes[i] > g.Time-0.75 {
			color := rl.White
			if g.DTimes[i] <= g.Time+0.25 && g.DTimes[i] > g.Time-0.25 {
				color = rl.Gray
			}
			rl.DrawTexture(g.Footsteps[1], 300, int32(300+((g.Time-g.DTimes[i])/0.75)*400), color)
		}
	}

	for i := 0; i < len(g.DrawingTimes); i++ {
		if g.DrawingScores[i] != -1 {
			continue
		}
		if g.DrawingTimes[i] <= g.Time+1 && g.DrawingTimes[i] > g.Time-0.25 {
			var opacity uint8
			if g.DrawingTimes[i] < g.Time {
				opacity = uint8(math.Abs(float64(g.Time-g.DrawingTimes[i])) * 255)
			} else if g.DrawingTimes[i] > g.Time {
				opacity = uint8((math.Abs(float64(g.DrawingTimes[i]-g.Time)) / 0.25) * 255)
			} else {
				opacity = 255
			}
			color := rl.Color{255, 255, 255, opacity}
			rl.DrawCircleV(g.Drawings[i], 15, color)
		}
	}
}

func (g *Game) HandleInputRune(inputType rune) {
	var times []float32
	if inputType == 'A' {
		times = g.ATimes
	} else {
		times = g.DTimes
	}

	var closestTime float32 = times[0]
	var index int = 0

	for i := 1; i < len(times); i++ {
		if math.Abs(float64(times[i]-g.Time)) < math.Abs(float64(closestTime-g.Time)) {
			closestTime = times[i]
			index = i
		}
	}

	fmt.Printf("%f\n", math.Abs(float64(closestTime-g.Time)))
	if math.Abs(float64(closestTime-g.Time)) <= 0.25 {
		if inputType == 'A' {
			g.AScores[index] = float32(math.Abs(float64(closestTime-g.Time)) / 0.25)
		} else {
			g.DScores[index] = float32(math.Abs(float64(closestTime-g.Time)) / 0.25)
		}
	}
}

func (g *Game) DrawScore() {
	dScores := 0
	for i := 0; i < len(g.DScores); i++ {
		if g.DScores[i] >= 0.0 && g.DScores[i] <= 1.0 {
			dScores++
		}
	}
	aScores := 0
	for i := 0; i < len(g.AScores); i++ {
		if g.AScores[i] >= 0.0 && g.AScores[i] <= 1.0 {
			aScores++
		}
	}
	drawingScores := 0
	for i := 0; i < len(g.DrawingScores); i++ {
		if g.DrawingScores[i] >= 0.0 && g.DrawingScores[i] <= 1.0 {
			drawingScores++
		}
	}
	rl.DrawText(fmt.Sprintf("Left Footsteps: %d", aScores), 500, 280, 10, rl.White)
	rl.DrawText(fmt.Sprintf("Right Footsteps: %d", dScores), 500, 300, 10, rl.White)
	rl.DrawText(fmt.Sprintf("Chalk Marks: %d", drawingScores), 500, 320, 10, rl.White)
}

func (g *Game) HandleInputMouse(mousePos rl.Vector2) {

}
