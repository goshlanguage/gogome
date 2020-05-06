package engine

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text exists to render text to the screen.
type Text struct {
	font     *ttf.Font
	renderer *sdl.Renderer
	text     string
	x, y     float64
}

func NewText(renderer *sdl.Renderer, text string, x float64, y float64) *Text {
	font, err := ttf.OpenFont("fonts/monogram.ttf", 64)
	checkErr(err)
	return &Text{
		font:     font,
		renderer: renderer,
		text:     text,
		x:        x,
		y:        y,
	}
}

func (t *Text) Draw() {
	gfx.StringRGBA(t.renderer, int32(t.x), int32(t.y), t.text, 255, 255, 255, 255)
	surface, err := t.font.RenderUTF8Solid(t.text, sdl.Color{255, 255, 255, 255})
	checkErr(err)
	box := sdl.Rect{int32(t.x), int32(t.y), surface.W, surface.H}
	var texture *sdl.Texture
	if texture, err = t.renderer.CreateTextureFromSurface(surface); err != nil {
		panic(err)
	}
	t.renderer.Copy(texture, nil, &box)

}

func (t *Text) Update() {

}
