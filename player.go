package main

import "github.com/veandco/go-sdl2/sdl"

type Player struct {
	texture    *sdl.Texture
	frame      int32
	frameLimit int32
	x          int32
	y          int32
	// sprites are chunked into 16x32 frames
	spriteXPos int32
	spriteYPos int32
}

func NewPlayer(renderer *sdl.Renderer) (*Player, error) {
	img, err := sdl.LoadBMP("sprites/npc.bmp")
	errHelper(err)
	defer img.Free()

	playerTexture, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return &Player{}, err
	}

	return &Player{
		texture:    playerTexture,
		x:          2,
		y:          530,
		frame:      0,
		frameLimit: 3,
		spriteXPos: 1,
		spriteYPos: 1,
	}, nil
}

func (player *Player) move(x int32, y int32) {
	player.x += x
	player.y += y
	player.frame += 1
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
}
