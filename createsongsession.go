package main

import (
	"bufio"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
)

type CreateSongSession struct {
	Music     rl.Music
	name      string
	log       string
	Time      float32
	Frames    uint32
	ChalkDown bool
}

func NewCreateSongSession(filename string) *CreateSongSession {
	music := rl.LoadMusicStream(fmt.Sprintf("resources/audio/%s", filename))
	music.Looping = false
	rl.PlayMusicStream(music)

	// chalkdown is set to true since the mouse is *very* likely being pressed entering the session.
	return &CreateSongSession{music, filename, "", 0, 0, true}
}

func (ss *CreateSongSession) UpdateSongSessionRune(inputType rune) {
	ss.log += fmt.Sprintf("%c,%f\n", inputType, ss.Time)
}

func (ss *CreateSongSession) UpdateSongSessionMouse(mouse rl.Vector2) {
	ss.log += fmt.Sprintf("M,%f,%f,%f\n", mouse.X, mouse.Y, ss.Time)
}

func (ss *CreateSongSession) UpdateSongSession(frameTime float32) {
	ss.Time += frameTime
	ss.Frames++
}

func (ss *CreateSongSession) CompleteSession() {
	file, err := os.Create(fmt.Sprintf("resources/sheets/%s.sheet", ss.name))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(ss.log)
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
