package engine

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// Enemy holds all things necessary for the Enemy to make their moves
type Enemy struct {
	// Frame tracks what Frame of the enemy animation we're on
	Frame      int32
	FrameLimit int32
	// Enemy health out of 1.0 representing 100%
	Health float64
	//LevelX, LevelY are used to track where on the level map the enemy is
	LevelX, LevelY int
	// This primitive logger exists so we can print out any helpful debug information
	Log string
	// Name tracks enemy objects
	Name string
	// size x and y pertain to what the standard size of the enemy is
	SizeX, SizeY int32
	// SpriteXPos and SpriteYPos is a Frame reference eg: [0, 1, 2, 3] for 4 Frames of animation
	// This is used to coordinate where in it's bitmap file it is
	SpriteXPos, SpriteYPos int32
	//the texture object for the enemy
	Texture *sdl.Texture
	// The x and y coordiates for the enemy on a tileMap
	X, Y float64
}

var debug = os.Getenv("HMDEBUG") == ""

// NewEnemy constructs a basic terminal object
func NewEnemy(name string, renderer *sdl.Renderer) (*Enemy, error) {
	img, err := sdl.LoadBMP("sprites/terminal.bmp")
	if err != nil {
		return &Enemy{}, err
	}

	texture, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return &Enemy{}, err
	}
	return &Enemy{
		FrameLimit: 1,
		Health:     1.0,
		Name:       name,
		SizeX:      32,
		SizeY:      32,
		Texture:    texture,
	}, nil
}

// GetLevelCoords returns the X and Y coordinates on the Level where
//    the enemy is supposed to be.
func (enemy *Enemy) GetLevelCoords() (x int, y int) {
	return enemy.LevelX, enemy.LevelY
}

// SetCoords helps when rendering the entity map to set where to Draw the enemy on the screen
// You might also think of this as setting the enemy coordinates relative to the camera overlooking our tileMap array
func (enemy *Enemy) SetCoords(x float64, y float64) {
	enemy.X = x
	enemy.Y = y
}

// Draw renders the enemy to the screen
func (enemy *Enemy) Draw(renderer *sdl.Renderer) {
	if debug {
		fmt.Printf("drawing enemy %s to %v,%v\n", enemy.Name, enemy.X, enemy.Y)
	}
	renderer.Copy(
		enemy.Texture,
		&sdl.Rect{X: enemy.SpriteXPos * enemy.SizeX, Y: enemy.SpriteYPos * enemy.SizeY, W: 32, H: 32},
		&sdl.Rect{X: int32(enemy.X), Y: int32(enemy.Y), W: 32, H: 32},
	)
}

// Update advances the enemy animation and updates x/y coords for the enemy
func (enemy *Enemy) Update(levelX int, levelY int) {
	enemy.Frame++
	// If we've iterated past our number of Frames, reset to 0
	if enemy.Frame > enemy.FrameLimit {
		enemy.Frame = 0
	}

	enemy.SpriteXPos++
	if enemy.SpriteXPos > enemy.FrameLimit {
		enemy.SpriteXPos = 0
	}
}

// SetX sets the enemy X coordinate
func (enemy *Enemy) SetX(x float64) {
	enemy.X = x
}

// SetY sets the enemy Y coordinate
func (enemy *Enemy) SetY(y float64) {
	enemy.Y = y
}

// Size returns the enemy size
func (enemy *Enemy) Size() (int32, int32) {
	return enemy.SizeX, enemy.SizeY
}
