package main

import "github.com/veandco/go-sdl2/sdl"

type Level struct {
	// BGFile is the filepath to the background
	BGFile  string
	Texture *sdl.Texture
}

// NewLevel takes in the filepath of a level's background, and a renderer
func NewLevel(filepath string, renderer *sdl.Renderer) (*Level, error) {
	img, err := sdl.LoadBMP(filepath)
	checkErr(err)
	defer img.Free()

	bgTexture, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return &Level{}, err
	}
	return &Level{
		BGFile:  filepath,
		Texture: bgTexture,
	}, nil
}
