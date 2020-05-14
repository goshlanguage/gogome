package engine

import "github.com/veandco/go-sdl2/sdl"

// Enemy holds all things necessary for the Enemy to make their moves
type Enemy struct {
	// Frame tracks what frame of the player animation we're on
	frame      int32
	frameLimit int32
	// Enemy health out of 1.0 representing 100%
	health float64
	// renderer allows us to draw this object to screen
	renderer *sdl.Renderer
	// size x and y pertain to what the standard size of the enemy is
	sizeX, sizeY int32
	// spriteXPos and spriteYPos is a frame reference eg: [0, 1, 2, 3] for 4 frames of animation
	spriteXPos, spriteYPos int32
	//the texture object for the enemy
	texture *sdl.Texture
	// The x and y coordiates for the enemy on a tileMap
	x, y float64
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
		frameLimit: 1,
		health:     1.0,
		renderer:   renderer,
		sizeX:      32,
		sizeY:      32,
		texture:    texture,
		x:          x,
		y:          y,
	}, nil
}

// Draw renders the enemy to the screen
func (enemy *Enemy) Draw() {
	enemy.renderer.Copy(
		enemy.texture,
		&sdl.Rect{X: enemy.spriteXPos * enemy.sizeX, Y: enemy.spriteYPos * enemy.sizeY, W: 32, H: 32},
		&sdl.Rect{X: int32(enemy.x), Y: int32(enemy.y), W: 32, H: 32},
	)
}

// Update advances the enemy animation
func (enemy *Enemy) Update() {
	enemy.frame++
	// If we've iterated past our number of frames, reset to 0
	if enemy.frame > enemy.frameLimit {
		enemy.frame = 0
	}

	enemy.spriteXPos++
	if enemy.spriteXPos > enemy.frameLimit {
		enemy.spriteXPos = 0
	}
}
