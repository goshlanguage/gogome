package engine

import "github.com/veandco/go-sdl2/sdl"

// Enemy holds all things necessary for the Enemy to make their moves
type Enemy struct {
	// Frame tracks what Frame of the enemy animation we're on
	Frame      int32
	FrameLimit int32
	// Enemy health out of 1.0 representing 100%
	Health float64
	// renderer allows us to draw this object to screen
	Renderer *sdl.Renderer
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

// NewEnemy constructs a basic terminal object
func NewEnemy(x float64, y float64, renderer *sdl.Renderer) (*Enemy, error) {
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
		Renderer:   renderer,
		SizeX:      32,
		SizeY:      32,
		Texture:    texture,
		X:          x,
		Y:          y,
	}, nil
}

// SetCoords helps when rendering the entity map to set where to Draw the enemy on the screen
// You might also think of this as setting the enemy coordinates relative to the camera overlooking our tileMap array
func (enemy *Enemy) SetCoords(x float64, y float64) {
	enemy.X = x
	enemy.Y = y
}

// Draw renders the enemy to the screen
func (enemy *Enemy) Draw() {
	enemy.Renderer.Copy(
		enemy.Texture,
		&sdl.Rect{X: enemy.SpriteXPos * enemy.SizeX, Y: enemy.SpriteYPos * enemy.SizeY, W: 32, H: 32},
		&sdl.Rect{X: int32(enemy.X), Y: int32(enemy.Y), W: 32, H: 32},
	)
}

// Update advances the enemy animation
func (enemy *Enemy) Update() {
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
