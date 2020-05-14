package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ryanhartje/gogome/pkg/engine"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winW = 800
	winH = 600
)

var (
	fonts = []string{"fonts/monogram.ttf"}
	grid  = false
	debug = false
)

func main() {
	e := engine.NewEngine()
	e.Init()

	window, err := sdl.CreateWindow(
		"GoGnome",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(winW),
		int32(winH),
		sdl.WINDOW_OPENGL,
	)
	checkErr(err)
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
	for x := 0; x < (winW * 10); x += 16 {
		mapping[x] = make(map[int]engine.Tile)
		for y := 0; y < (winH * 10); y += 16 {
			mapping[x][y] = grass
			// half the time, give us different grass
			if rand.Intn(10) > 5 {
				mapping[x][y] = grass2
			}
		}
	}
	level.XSize = winW * 10
	level.YSize = winH * 10
	level.TileMap = mapping

	fmt.Println("MAP: %s", level.TileMap)

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		checkErr(err)
	}
	defer mix.CloseAudio()

	// Play BG Wav
	chunk, err := engine.QueueWAV("sfx/streets.wav")
	checkErr(err)
	level.Sounds["background"] = append(level.Sounds["background"], chunk)
	//e.PlayWAV(level.Sounds["background"][0])

	// setup a dummy enemy for demo
	enemy, err := engine.NewEnemy(384, 150, renderer)
	checkErr(err)

	entities := []engine.Entity{player, enemy}

	// Set tick rate to 8 FPS
	// 8 looks more natural for our 8 bit style animations
	tick := time.NewTicker(time.Second / 16)

	count := 0
	for {
		renderer.Clear()
		// Setup ESC to exit keybinding
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_ESCAPE] == 1 {
			os.Exit(0)
		}

		// Setup a tick rate
		select {
		case <-tick.C:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					os.Exit(0)

				case *sdl.KeyboardEvent:
					if debug {
						fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
							t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
					}
					if t.Keysym.Scancode == sdl.SCANCODE_G && t.State == 1 {
						grid = !grid
					}

					/*
						case *sdl.MouseButtonEvent:
									x, y, state := sdl.GetMouseState()
									if state == 1 {
										coords := fmt.Sprintf("(%d, %d)", x, y)
										text := engine.NewText(renderer, coords, float64(x), float64(y))
										entities = append(entities, text)
									}
					*/

				}
			}

			// Render level to window
			for x := 0; x < winW; x += 16 {
				for y := 0; y < winH; y += 16 {
					tile := level.TileMap[level.X+x][level.Y+y]
					width := tile.X1 - tile.X0
					height := tile.Y1 - tile.Y0
					// Render the background of the level
					renderer.Copy(
						level.Texture,
						&sdl.Rect{X: tile.X0, Y: tile.Y0, W: width, H: height},
						&sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height},
					)
					if grid {
						// draw vertical grid lines
						gfx.LineRGBA(renderer, int32(x), 0, int32(x), winH, 100, 0, 0, 100)
						// draw horizontal line
						gfx.LineRGBA(renderer, 0, int32(y), winW, int32(y), 100, 0, 0, 100)
					}
				}
			}
			level.Update()
			for _, e := range entities {
				e.Draw()
				e.Update()
			}
			renderer.Present()

			if debug {
				count++
				if count > 3 {
					count = 0
				}
				if count == 3 {
					fmt.Printf("level: (%d,%d)\n", level.X, level.Y)
				}
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
