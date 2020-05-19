package main

import (
	"fmt"
	"os"
	"reflect"
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
	debug = false
	log   string
)

func main() {
	if os.Getenv("HMDEBUG") != "" {
		debug = true
		fmt.Println("debug mode enabled")
	}
	e := engine.NewEngine()
	e.Init()

	window, err := sdl.CreateWindow(
		"hackerman",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(winW),
		int32(winH),
		sdl.WINDOW_SHOWN,
	)
	checkErr(err)
	window.UpdateSurface()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err)
	defer renderer.Destroy()

	player, err := engine.NewPlayer(renderer)
	checkErr(err)
	// create a gravity effect and assign it to the player
	gravity := func(player *engine.Player) {
		// keep in mind that SDL has a reversed Y coordinate system, where 0,0 is the top [left] of the screen, and winH,0 is the bottom [left] of the screen
		if player.Y < winH-float64(player.SizeY*2) {
			player.Y += 4
		}
	}
	player.Effects = append(player.Effects, gravity)

	// Load in our level asset and generate a random map
	level, err := engine.NewRandomizedLevel("sprites/overworld.bmp", renderer)
	checkErr(err)
	level.CameraX = winW
	level.CameraY = winH
	level.XSize = winW * 10
	level.YSize = winH * 10

	// setup a dummy enemy
	enemy, err := engine.NewEnemy("computer", renderer)
	enemy.LevelX = 8 * 100
	enemy.LevelY = 8 * 20

	if reflect.TypeOf(level.EntityMap[0]) == reflect.TypeOf(nil) {
		panic("wtf")
	}
	level.EntityMap[0][0] = enemy

	// Setup audio
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		checkErr(err)
	}
	defer mix.CloseAudio()

	// Play BG Wav
	chunk, err := engine.QueueWAV("sfx/streets.wav")
	checkErr(err)
	level.Sounds["background"] = append(level.Sounds["background"], chunk)
	//e.PlayWAV(level.Sounds["background"][0])

	entities := []engine.Entity{player, enemy}

	// Set tick rate to 8 FPS
	// 8 looks more natural for our 8 bit style animations
	tick := time.NewTicker(time.Second / 16)

	// lastDebugMsg is used as a cache, to make sure we only print debug
	//   messages if they are new messages. This helps us not flood output
	//   during the main loop. While it is set to a tick rate, multiple messages per tick get our of hand.
	//   log is used to log out the combined output of mainloop in debug mode.
	var lastDebugMsg string
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
					e.Quit()
					sdl.Quit()
					os.Exit(0)

				case *sdl.KeyboardEvent:
					if debug && t.Keysym.Scancode == sdl.SCANCODE_G && t.State == 1 {
						level.Debug = !level.Debug
					}
					// If you want to explore keybinding values, you can comment this out and they will log
					//   to your console.
					/*
						if debug {
							log += fmt.Sprintf("%s[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
								log, t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
						}
					*/

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
			originalX := level.X
			originalY := level.Y
			level.Draw(renderer)
			level.Update()

			for _, e := range entities {
				eX, eY := e.GetLevelCoords()
				isPlayer := reflect.TypeOf(e) == reflect.TypeOf(player)
				inCameraView := eX >= level.X && eX < level.X+winW && eY >= level.Y && eY < level.Y+winH
				if inCameraView || isPlayer {
					e.Draw(renderer, level.X, level.Y)
				}
				e.Update(level.X-originalX, level.Y-originalY)
			}
			renderer.Present()

			if debug {
				if log != lastDebugMsg {
					fmt.Println(log)
				}
				lastDebugMsg = log
				log = ""
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
