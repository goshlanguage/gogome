package main

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/ryanhartje/gogome/pkg/engine"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	winW = 800
	winH = 600
)

var (
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
		sdl.WINDOW_OPENGL,
	)
	checkErr(err)
	window.UpdateSurface()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err)
	defer renderer.Destroy()

	player, err := engine.NewPlayer(renderer, "assets/sprites/character.png")
	checkErr(err)

	// Load in our level asset and generate a random map
	level, err := engine.NewRandomizedLevel("assets/sprites/overworld.bmp", renderer)
	checkErr(err)
	level.CameraX = winW
	level.CameraY = winH
	level.XSize = winW * 10
	level.YSize = winH * 10

	// Attempted to provide a lighting effect by providing an alpha layer over the viewport/camera
	// It had weird implementation effects
	/*
		var lightcycle float64
		sineLighting := func(level *engine.Level) {
			// 10 is the max value in alpha we're looking to promote. Will have a max value of 55, minimum value of 0 as we only use 0-pi in the loop
			alphaModifier := math.Sin(lightcycle) * 10
			lightcycle += (math.Pi / 3200)
			if lightcycle >= math.Pi {
				lightcycle = 0
			}
			if debug {
				fmt.Printf("a: %v\tlc: %v\n", alphaModifier, lightcycle)
			}
			// gfx.BoxColor(renderer, 0, 0, winW, winH, sdl.Color{0, 0, 0, uint8(alpha)})
			gfx.BoxRGBA(renderer, 0, 0, winW, winH, 0, 0, 0, uint8(alphaModifier))
		}
		if debug {
			level.Lighting = sineLighting
		}
	*/

	// setup a dummy enemy
	enemy, err := engine.NewEnemy("computer", renderer)
	enemy.LevelX = 8 * 100
	enemy.LevelY = 8 * 20

	if reflect.TypeOf(level.EntityMap[0]) == reflect.TypeOf(nil) {
		panic("EntityMap is nil after assignment. Can't render entities")
	}
	level.EntityMap[0][0] = enemy

	// Setup audio
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		checkErr(err)
	}
	defer mix.CloseAudio()

	// Play BG Wav
	// chunk, err := engine.QueueWAV("assets/sfx/streets.wav")
	// checkErr(err)
	// level.Sounds["background"] = append(level.Sounds["background"], chunk)
	//e.PlayWAV(level.Sounds["background"][0])

	font64, err := ttf.OpenFont("assets/fonts/monogram.ttf", 64)
	checkErr(err)
	font32, err := ttf.OpenFont("assets/fonts/monogram.ttf", 32)
	checkErr(err)

	var sineCounter float64
	sineFunc := func(text *engine.Text) {
		sineCounter += math.Pi / 16
		if sineCounter >= math.Pi {
			sineCounter = 0
		}
		alpha := math.Sin(sineCounter) * 255
		text.Color = sdl.Color{R: 255, G: 255, B: 255, A: uint8(alpha)}
		text.X -= 10
	}
	var menuText []engine.Text
	menuText = append(menuText, engine.Text{
		Color: sdl.Color{R: 255, G: 255, B: 255, A: 240},
		Font:  font64,
		Text:  "Tyler the Tiler",
		X:     200,
		Y:     450,
	})
	menuText = append(menuText, engine.Text{
		Color:   sdl.Color{R: 255, G: 255, B: 255, A: 240},
		Effects: []func(*engine.Text){sineFunc},
		Font:    font32,
		Text:    "Press spacebar to continue...",
		X:       210,
		Y:       520,
	})

	// Setup a main menu, real retro like
	keyMapFuncs := make(map[sdl.Keycode]func(*engine.Menu))
	keyMapFuncs[sdl.K_SPACE] = func(menu *engine.Menu) {
		menu.Break = true
	}
	keyMapFuncs[sdl.K_ESCAPE] = func(menu *engine.Menu) {
		e.Quit()
		sdl.Quit()
		os.Exit(0)
	}
	menu := engine.Menu{
		BGImagePath: "assets/backgrounds/main.png",
		BGSizeX:     384,
		BGSizeY:     244,
		Components:  menuText,
		Debug:       debug,
		KeyMapping:  keyMapFuncs,
		WinH:        winH,
		WinW:        winW,
	}
	menu.Loop(renderer)

	// Set tick rate to 8 FPS
	// 8 looks more natural for our 8 bit style animations
	tick := time.NewTicker(time.Second / 16)

	// lastDebugMsg is used as a cache, to make sure we only print debug
	//   messages if they are new messages. This helps us not flood output
	//   during the main loop. While it is set to a tick rate, multiple messages per tick get our of hand.
	//   log is used to log out the combined output of mainloop in debug mode.
	var lastDebugMsg string
	entities := []engine.Entity{player, enemy}
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
