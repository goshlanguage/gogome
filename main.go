package main

import (
	"os"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	fonts = []string{"fonts/monogram.ttf"}
)

func main() {
	e := NewEngine()
	e.Init()

	window, err := sdl.CreateWindow(
		"GoGome",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_OPENGL,
	)
	checkErr(err)
	defer window.Destroy()

	window.UpdateSurface()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err)
	defer renderer.Destroy()

	player, err := NewPlayer(renderer)
	checkErr(err)

	level, err := NewLevel("sprites/background.bmp", renderer)
	checkErr(err)

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		checkErr(err)
	}
	defer mix.CloseAudio()

	// Play BG Wav
	chunk, err := QueueWAV("sfx/streets.wav")
	checkErr(err)
	level.Sounds["background"] = append(level.Sounds["background"], chunk)
	// e.PlayWAV(level.Sounds["background"][0])

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()

	// setup a dummy enemy for demo
	enemy, err := NewEnemy(384, 150, renderer)
	checkErr(err)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			}
		}

		// Render the background of the level
		renderer.Copy(
			level.Texture,
			&sdl.Rect{X: 0, Y: 0, W: 800, H: 600},
			&sdl.Rect{X: 0, Y: 0, W: 800, H: 600},
		)
		enemy.Draw()
		enemy.Update()
		//renderer.Copy(nil, &text.ClipRect, &sdl.Rect{X: 400, Y: 400, W: 100, H: 36})
		player.Draw()
		player.Update()

		renderer.Present()
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
