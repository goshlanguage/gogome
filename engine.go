package main

import (
	"fmt"
	"io/ioutil"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Engine holds all of the assets necessary to run a 2D engine
type Engine struct {
	Fonts    []*ttf.Font
	Entities *sdl.Surface
}

// NewEngine creates and instanciates our engine
func NewEngine() *Engine {
	e := &Engine{}
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

// QueueWAV uses the initialized mixer and plays a WAV on the first channel
func QueueWAV(filepath string) (*mix.Chunk, error) {
	// Load in BG wav
	data, err := ioutil.ReadFile("sfx/streets.wav")
	if err != nil {
		return &mix.Chunk{}, err
	}

	chunk, err := mix.QuickLoadWAV(data)
	if err != nil {
		return &mix.Chunk{}, err
	}
	// defer chunk.Free()
	return chunk, nil
}

// PlayWAV is blocking, please call it as a go routine
// Increments channels first so we don't inadvertently stop playing of another chunk
// This logic is really poor and doesn't garbage collect.
func (e *Engine) PlayWAV(chunk *mix.Chunk) {
	chunk.Play(-1, 0)
}
