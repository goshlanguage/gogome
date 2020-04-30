package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// GameEngine is here in the event we need a testing mock.
type GameEngine interface {
	Init()
	AddFont(fontPath ...string)
}

// Engine holds all of the assets necessary to run a 2D engine
type Engine struct {
	Fonts    []*ttf.Font
	Entities *sdl.Surface
}

// NewEngine creates and instanciates our engine
func NewEngine() *Engine {
	e := &Engine{}
	e.Init()
	return e
}

// Init allows the user to initialize everything necessary for the game engine.
//   First, it initializes sdl with the INIT_EVERYTHING flag. This can likely be optimized.
//   Next, we setup ttf (true type font) to render text onto our window.
//   We also initialize audio so we can render WAV or MP3 through the mixer.
func (e *Engine) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		checkErr(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		checkErr(err)
	}

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		checkErr(err)
	}
	// defer mix.CloseAudio()

	font, err := ttf.OpenFont("fonts/monogram.ttf", 32)
	checkErr(err)

	e.Fonts = append(e.Fonts, font)

}

// AddFont provides a helper to variadically add fonts to the engine
func AddFont(fontPaths ...string) {
	for _, path := range fontPaths {
		fmt.Printf("FONT FOUND: %s", path)
	}
}
