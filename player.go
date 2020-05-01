package main

import "github.com/veandco/go-sdl2/sdl"

const speed = 0.75

// Player holds all things relevant to make the Player model self sufficient.
type Player struct {
	// Frame tracks what frame of the player animation we're on
	frame      int32
	frameLimit int32
	// store the renderer pointer so we can render through a method
	renderer *sdl.Renderer
	// sprites are chunked into 16x32 frames
	// spriteXPos and spriteYPos is a frame reference eg: [0, 1, 2, 3] for 4 frames of animation
	spriteXPos, spriteYPos int32
	// texture holds the bitmap used to render the Player
	texture *sdl.Texture
	// x and y are the x and y coordinates for the Player
	x, y float64
}

// NewPlayer is a Player factory that sets it's defaults and returns it
func NewPlayer(renderer *sdl.Renderer) (*Player, error) {
	img, err := sdl.LoadBMP("sprites/npc.bmp")
	checkErr(err)
	defer img.Free()

	playerTexture, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return &Player{}, err
	}

	return &Player{
		frame:      0,
		frameLimit: 3,
		renderer:   renderer,
		spriteXPos: 1,
		spriteYPos: 1,
		texture:    playerTexture,
		// Place the user in the middle of the screen, assuming 800x600 minus half the sprite size
		x: 384.0,
		y: 292.0,
	}, nil
}

// Draw render's the Player sprite to the screen
func (player *Player) Draw() {
	player.renderer.Copy(
		player.texture,
		&sdl.Rect{X: player.spriteXPos * 16, Y: player.spriteYPos * 32, W: 16, H: 32},
		&sdl.Rect{X: int32(player.x), Y: int32(player.y), W: 32, H: 64},
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

	// If we've stopped moving, reset our animation to our still frame
	if !moving {
		player.frame = 0
	}

}

func (player *Player) move(x float64, y float64) {
	// Don't let player move beyond bounds, but DO update their animation
	if player.x >= 0 && player.x <= 800 {
		player.x += x * speed
	} else {
		if player.x < 0 {
			player.x = 0
		}
		if player.x > 800 {
			player.x = 800
		}
	}
	if player.y >= 0 && player.y <= 600 {
		player.y += y * speed
	} else {
		if player.y < 0 {
			player.y = 0
		}
		if player.y > 600 {
			player.y = 600
		}
	}
	player.frame++
	// If we've iterated past our number of frames, reset to 0
	if player.frame > player.frameLimit {
		player.frame = 0
	}

	// If x is the coordinate that we're changing, set the frame animation.
	//   Our original player sprite has it's right and left facing sprites in the
	//   1st and 3rd Y positions, so use those respectively
	//   Repeat this logic for movements along the Y axis, setting their Y positions in the
	//   sprite sheet.
	if x != 0 {
		if x > 0 {
			player.spriteYPos = 1
		} else {
			player.spriteYPos = 3
		}
		player.spriteXPos++
		if player.spriteXPos > player.frameLimit {
			player.spriteXPos = 0
		}
	}
	if y != 0 {
		if y > 0 {
			player.spriteYPos = 0
		} else {
			player.spriteYPos = 2
		}
		player.spriteXPos++
		if player.spriteXPos > player.frameLimit {
			player.spriteXPos = 0
		}
	}
}
