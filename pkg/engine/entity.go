package engine

import "github.com/veandco/go-sdl2/sdl"

// Entity is an interface that aides the engine in having commmon functionality
type Entity interface {
	// Draw is given the renderer, and an X,Y coordinate to draw the entity to the screen
	Draw(*sdl.Renderer, int, int)
	// GetLevelCoords returns the X,Y coordinate pair for the entity as it relates to the level map
	GetLevelCoords() (int, int)
	// SetX,SetY take a precise pixel coordinate pair to set the entity's sprite drawing to
	SetX(float64)
	SetY(float64)
	// Size returns the number of pixels wide and high the sprite/hitbox should be
	Size() (int32, int32)
	// Update takes a level X,Y coordinate pair so it can update itself relative to the camera
	Update(int, int)
}
