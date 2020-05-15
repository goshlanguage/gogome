package engine

import "github.com/veandco/go-sdl2/sdl"

// Entity is an interface that aides the engine in having commmon functionality
type Entity interface {
	Draw(*sdl.Renderer)
	GetLevelCoords() (int, int)
	SetX(float64)
	SetY(float64)
	Size() (int32, int32)
	Update(int, int)
}
