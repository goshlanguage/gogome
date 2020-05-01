package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	fonts = []string{"fonts/monogram.ttf"}
)

func main() {
	e := NewEngine()
	e.Init()
	text, err := e.Fonts[0].RenderUTF8Shaded(
		"Hello!",
		sdl.Color{200, 200, 200, 255},
		sdl.Color{255, 255, 255, 255},
	)

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

	// Play BG Wav
	chunk, err := QueueWAV("sfx/streets.wav")
	checkErr(err)
	level.Sounds["background"] = append(level.Sounds["background"], chunk)
	e.PlayWAV(level.Sounds["background"][0])

	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				if string(t.Keysym.Sym) == "s" && t.State == 1 {
					if player.y <= 536 {
						player.move(0, 1, 1)
					}
				}
				if string(t.Keysym.Sym) == "w" && t.State == 1 {
					if player.y >= 0 {
						player.move(0, -1, 1)
					}
				}
				if string(t.Keysym.Sym) == "d" && t.State == 1 {
					if player.x < 800 {
						player.move(1, 0, 1)
					}
				}
				if string(t.Keysym.Sym) == "a" && t.State == 1 {
					if player.x > 0 {
						player.move(-1, 0, 1)
					}
				}
			}
		}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()

		// Stage background then player sprite
		renderer.Copy(
			level.Texture,
			&sdl.Rect{X: 0, Y: 0, W: 800, H: 600},
			&sdl.Rect{X: 0, Y: 0, W: 800, H: 600},
		)

		renderer.Copy(nil, &text.ClipRect, &sdl.Rect{X: 400, Y: 400, W: 100, H: 36})

		renderer.Present()
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
