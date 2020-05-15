// This is incomplete. Only checking in so I can come back to this.
package main

import (
	"fmt"
	"os"

	"github.com/ryanhartje/gogome/pkg/engine"
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
}
