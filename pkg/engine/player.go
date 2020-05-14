package engine

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const speed = 0.75

// Player holds all things relevant to make the Player model self sufficient.
type Player struct {
	// Frame tracks what Frame of the player animation we're on
	Frame      int32
	FrameLimit int32
	// store the Renderer pointer so we can render through a method
	Renderer *sdl.Renderer
	// size x and y pertain to what the standard size of the enemy is
	SizeX, SizeY int32
	// The player sprites are chunked into 16x32 Frames
	// SpriteXPos and SpriteYPos is a Frame reference eg: [0, 1, 2, 3] for 4 Frames of animation
	SpriteXPos, SpriteYPos int32
	// Texture holds the bitmap used to render the Player
	Texture *sdl.Texture
	// x and y are the x and y coordinates for the Player
	X, Y float64
}

// NewPlayer is a Player factory that sets it's defaults and returns it
func NewPlayer(Renderer *sdl.Renderer) (*Player, error) {

	Texture, err := img.LoadTexture(Renderer, "Sprites/character.png")
	checkErr(err)

	return &Player{
		Frame:      0,
		FrameLimit: 3,
		Renderer:   Renderer,
		SizeX:      16,
		SizeY:      32,
		SpriteXPos: 1,
		SpriteYPos: 1,
		Texture:    Texture,
		// Place the user in the middle of the screen, assuming 800x600 minus half the Sprite size
		X: 396.0,
		Y: 272.0,
	}, nil
}

// Draw render's the Player Sprite to the screen
func (player *Player) Draw() {
	player.Renderer.Copy(
		player.Texture,
		&sdl.Rect{X: player.SpriteXPos * 16, Y: player.SpriteYPos * 32, W: 16, H: 32},
		&sdl.Rect{X: int32(player.X), Y: int32(player.Y), W: 32, H: 64},
	)
}

// Update checks for keystrokes and calls the appropriate method based on the user input
func (player *Player) Update() {
	keys := sdl.GetKeyboardState()
	moving := false
	// UP
	if keys[sdl.SCANCODE_W] == 1 {
		player.move(0, -1)
		moving = true
	}
	// DOWN
	if keys[sdl.SCANCODE_S] == 1 {
		player.move(0, 1)
		moving = true
	}
	// LEFT
	if keys[sdl.SCANCODE_A] == 1 {
		player.move(-1, 0)
		moving = true
	}
	// RIGHT
	if keys[sdl.SCANCODE_D] == 1 {
		player.move(1, 0)
		moving = true
	}

	// If we've stopped moving, reset our animation to our still Frame
	if !moving {
		player.Frame = 0
	}

}

func (player *Player) move(x float64, y float64) {
	// Don't let player move beyond bounds, but DO update their animation
	if player.X >= 392 && player.X <= 400 {
		player.X += x * speed
	} else {
		if player.X < 392 {
			player.X = 392
		}
		if player.X > 400 {
			player.X = 400
		}
	}
	if player.Y >= 268 && player.Y <= 276 {
		player.Y += y * speed
	} else {
		if player.Y < 268 {
			player.Y = 268
		}
		if player.Y > 276 {
			player.Y = 276
		}
	}
	player.Frame++
	// If we've iterated past our number of Frames, reset to 0
	if player.Frame > player.FrameLimit {
		player.Frame = 0
	}

	// If x is the coordinate that we're changing, set the Frame animation.
	//   Our original player Sprite has it's right and left facing Sprites in the
	//   1st and 3rd Y positions, so use those respectively
	//   Repeat this logic for movements along the Y axis, setting their Y positions in the
	//   Sprite sheet.
	if x != 0 {
		if x > 0 {
			player.SpriteYPos = 1
		} else {
			player.SpriteYPos = 3
		}
		player.SpriteXPos++
		if player.SpriteXPos > player.FrameLimit {
			player.SpriteXPos = 0
		}
	}
	if y != 0 {
		if y > 0 {
			player.SpriteYPos = 0
		} else {
			player.SpriteYPos = 2
		}
		player.SpriteXPos++
		if player.SpriteXPos > player.FrameLimit {
			player.SpriteXPos = 0
		}
	}
}

// SetX sets the player X coordinate
func (player *Player) SetX(x float64) {
	player.X = x
}

// SetY sets the player Y coordinate
func (player *Player) SetY(y float64) {
	player.Y = y
}

// Size returns x and y values of the sprite's size
func (player *Player) Size() (int32, int32) {
	return player.SizeX, player.SizeY
}
