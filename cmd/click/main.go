package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ryanhartje/gogome/pkg/engine"
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
	e := engine.NewEngine()
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

	player, err := engine.NewPlayer(renderer)
	checkErr(err)

	// Load in our level asset and generate a plain map
	level, err := engine.NewLevel("sprites/overworld.bmp", renderer)
	checkErr(err)
	grass := engine.Tile{X0: 0, X1: 16, Y0: 0, Y1: 16}
	grass2 := engine.Tile{X0: 272, X1: 303, Y0: 464, Y1: 495}
	mapping := map[int]map[int]engine.Tile{}
	// Generate a plain grass map
	for x := 0; x < winW; x += 16 {
		mapping[x] = make(map[int]engine.Tile)
		for y := 0; y < winH; y += 16 {
			mapping[x][y] = grass
			if y > (winH - 64) {
				mapping[x][y] = grass2
			}
		}
	}
	level.TileMap = mapping

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		checkErr(err)
	}
	defer mix.CloseAudio()

	// Play BG Wav
	chunk, err := engine.QueueWAV("sfx/streets.wav")
	checkErr(err)
	level.Sounds["background"] = append(level.Sounds["background"], chunk)
	// e.PlayWAV(level.Sounds["background"][0])

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()

	// setup a dummy enemy for demo
	enemy, err := engine.NewEnemy(384, 150, renderer)
	checkErr(err)

	entities := []engine.Entity{player, enemy}

	// Set tick rate to 8 FPS
	// 8 looks more natural for our 8 bit style animations
	tick := time.NewTicker(time.Second / 8)

	for {

		// Setup a tick rate
		select {
		case <-tick.C:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					os.Exit(0)

				case *sdl.MouseButtonEvent:
					x, y, state := sdl.GetMouseState()
					if state == 1 {
						coords := fmt.Sprintf("(%d, %d)", x, y)
						text := engine.NewText(renderer, coords, float64(x), float64(y))
						entities = append(entities, text)
					}
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
					width := tile.X1 - tile.X0
					height := tile.Y1 - tile.Y0
					// Render the background of the level
					renderer.Copy(
						level.Texture,
						&sdl.Rect{X: tile.X0, Y: tile.Y0, W: width, H: height},
						&sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height},
					)
					// gfx.LineRGBA(renderer, 0, 0, int32(x), int32(y), 100, 0, 0, 100)
				}
			}
			for _, e := range entities {
				e.Draw()
				e.Update()
			}

			renderer.Present()
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
