package main

import "github.com/veandco/go-sdl2/sdl"

const speed = 16

type Player struct {
	frame      int32
	frameLimit int32
	renderer   *sdl.Renderer
	// sprites are chunked into 16x32 frames
	spriteXPos int32
	spriteYPos int32
	texture    *sdl.Texture
	x          int32
	y          int32
}

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
		x:          2,
		y:          530,
	}, nil
}

func (player *Player) Draw() {
	player.renderer.Copy(
		player.texture,
		&sdl.Rect{X: player.spriteXPos * 16, Y: player.spriteYPos * 32, W: 16, H: 32},
		&sdl.Rect{X: player.x, Y: player.y, W: 32, H: 64},
	)
}

func (player *Player) Update() {
	keys := sdl.GetKeyboardState()

	// UP
	if keys[sdl.SCANCODE_W] == 1 {
		player.move(0, -1, 1)
	}
	// DOWN
	if keys[sdl.SCANCODE_S] == 1 {
		player.move(0, 1, 1)
	}
	// LEFT
	if keys[sdl.SCANCODE_A] == 1 {
		player.move(-1, 0, 1)
	}
	// RIGHT
	if keys[sdl.SCANCODE_D] == 1 {
		player.move(1, 0, 1)
	}

}

func (player *Player) move(x int32, y int32, repeat int) {
	halfStep := int32(speed / 2)
	player.x += x * halfStep
	player.y += y * halfStep
	player.frame++
	if player.frame > player.frameLimit {
		player.frame = 0
	}

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

	player.Draw()
	if repeat > 0 {
		player.move(x, y, repeat-1)
	}
}
