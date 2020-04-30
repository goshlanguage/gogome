package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("Initializing SDL")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"GoGome",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_OPENGL,
	)
	errHelper(err)
	defer window.Destroy()

	window.UpdateSurface()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	errHelper(err)
	defer renderer.Destroy()

	player, err := NewPlayer(renderer)
	errHelper(err)

	level, err := NewLevel("sprites/background.bmp", renderer)
	errHelper(err)

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		errHelper(err)
	}
	defer mix.CloseAudio()

	// Load in BG wav
	data, err := ioutil.ReadFile("sfx/streets.wav")
	errHelper(err)

	chunk, err := mix.QuickLoadWAV(data)
	errHelper(err)
	defer chunk.Free()

	chunk.Play(1, 0)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				if string(t.Keysym.Sym) == "s" && t.State == 1 {
					if player.y <= 536 {
						player.move(0, 8)
					}
				}
				if string(t.Keysym.Sym) == "w" && t.State == 1 {
					if player.y >= 0 {
						player.move(0, -8)
					}
				}
				if string(t.Keysym.Sym) == "d" && t.State == 1 {
					if player.x < 800 {
						player.move(8, 0)
					}
				}
				if string(t.Keysym.Sym) == "a" && t.State == 1 {
					if player.x > 0 {
						player.move(-8, 0)
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
		renderer.Copy(
			player.texture,
			&sdl.Rect{X: player.spriteXPos * 16, Y: player.spriteYPos * 32, W: 16, H: 32},
			&sdl.Rect{X: player.x, Y: player.y, W: 32, H: 64},
		)

		renderer.Present()
	}
}

func errHelper(err error) {
	if err != nil {
		panic(err)
	}
}
