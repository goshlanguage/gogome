package engine

import (
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Menu is an abstract type for presenting the user with a menu
// It gives the developer an ability to break out of the main loop by providing a looping function of its own,
// as well as a keymapping to help compose the menu and what menu keys will do
type Menu struct {
	BGImagePath      string
	BGSizeX, BGSizeY int
	// Trigger for breaking out of a menu loop
	Break      bool
	Components []Text
	// Cycle provides a hook for animation cycles, or otherwise
	Cycle      int
	Debug      bool
	KeyMapping map[sdl.Keycode]func(*Menu)
	WinH, WinW int // Winow height and width for draw surface
}

// Draw renders the menu to the screen
func (menu *Menu) Draw(renderer *sdl.Renderer, x, y int) {
	image, err := img.Load(menu.BGImagePath)
	checkErr(err)

	texture, err := renderer.CreateTextureFromSurface(image)
	checkErr(err)

	renderer.Copy(
		texture,
		&sdl.Rect{X: 0, Y: 0, W: int32(menu.BGSizeX), H: int32(menu.BGSizeY)},
		&sdl.Rect{X: 0, Y: 0, W: int32(menu.WinW), H: int32(menu.WinH)},
	)
	for _, i := range menu.Components {
		i.Draw(renderer, x, y)
	}
}

// Update should bind a controller and update the menu according to user input
func (menu *Menu) Update(x, y int) {
	for _, component := range menu.Components {
		component.Update(x, y)
	}
}

// Loop iterates through a tick until the user exits the menu
func (menu *Menu) Loop(renderer *sdl.Renderer) {
	// setup a consistent tick rate
	tick := time.NewTicker(time.Second / 30) // 30 FPS
	for menu.Break == false {
		select {
		case <-tick.C:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.KeyboardEvent:
					if menu.KeyMapping[t.Keysym.Sym] != nil {
						menu.KeyMapping[t.Keysym.Sym](menu)
					}
				}
			}

			menu.Draw(renderer, 0, 0)
			menu.Update(0, 0)
			renderer.Present()
		}
	}
}
