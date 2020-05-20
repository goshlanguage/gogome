package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text exists to render text to the screen.
type Text struct {
	Color sdl.Color
	// Effect allows the developer a way to mutate the object
	Effects []func(*Text)
	Font    *ttf.Font
	Surface *sdl.Surface
	Text    string
	X, Y    float64
}

// NewText is a helper factory for Text on screen
// Ex:
// font, _ := ttf.OpenFont(exampleFont, 32)
// t := Text{}
// surface, err := font.RenderUTF8Solid("example", sdl.Color{R: 255, G: 255, B: 255, A:255})
// t.Font = font
// t.Surface = surface
// t.X = 400
// t.Y = 400
func NewText(font *ttf.Font, text string, x float64, y float64) *Text {
	return &Text{
		Font: font,
		Text: text,
		X:    x,
		Y:    y,
	}
}

// Draw for Text requires creating the surface, then copying the surface object onto a rectangle,
//   then rendering that rectangle to the screen
func (t *Text) Draw(renderer *sdl.Renderer, x, y int) {
	surface, err := t.Font.RenderUTF8Solid(t.Text, t.Color)
	checkErr(err)
	t.Surface = surface
	box := sdl.Rect{X: int32(t.X), Y: int32(t.Y), W: t.Surface.W, H: t.Surface.H}
	var texture *sdl.Texture
	if texture, err = renderer.CreateTextureFromSurface(surface); err != nil {
		panic(err)
	}
	renderer.Copy(texture, nil, &box)
}

// Update exists to fulfill the entities interface contract
func (t *Text) Update(x, y int) {
	for _, effect := range t.Effects {
		effect(t)
	}
}
