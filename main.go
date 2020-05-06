package main

import (
	"os"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winW = 800
	winH = 600
)

var (
	fonts = []string{"fonts/monogram.ttf"}
)

func main() {
	e := NewEngine()
	e.Init()

	window, err := sdl.CreateWindow(
		"GoGnome",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winW,
		winH,
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

	// Load in our level asset and generate a plain map
	level, err := NewLevel("sprites/overworld.bmp", renderer)
	checkErr(err)
	grass := Tile{x0: 0, x1: 16, y0: 0, y1: 16}
	mapping := map[int]map[int]Tile{}
	// Generate a plain grass map
	for x := 0; x < winW; x += 16 {
		mapping[x] = make(map[int]Tile)
		for y := 0; y < winH; y += 16 {
			mapping[x][y] = grass
		}
	}
	level.TileMap = mapping

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
			// Setup ESC to exit keybinding
			keys := sdl.GetKeyboardState()
			if keys[sdl.SCANCODE_ESCAPE] == 1 {
				os.Exit(0)
			}
		}

		// Render tilemap
		for x := 0; x < winW; x += 16 {
			for y := 0; y < winH; y += 16 {
				tile := level.TileMap[x][y]
				width := tile.x1 - tile.x0
				height := tile.y1 - tile.y0
				// Render the background of the level
				renderer.Copy(
					level.Texture,
					&sdl.Rect{X: tile.x0, Y: tile.y0, W: width, H: height},
					&sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height},
				)
			}
		}

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
