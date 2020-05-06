package main

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

type Level struct {
	// BGFile is the filepath to the background
	BGFile  string
	Sounds  map[string][]*mix.Chunk
	Texture *sdl.Texture
	// TileMap is a matrix representing the map
	// TileMap[x][y]
	TileMap map[int]map[int]Tile
}

// Tile represents a tile in a tilemap. This might be a 16x16 sprite or a 16x128 tile.
type Tile struct {
	x0 int32
	x1 int32
	y0 int32
	y1 int32
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
		Sounds:  make(map[string][]*mix.Chunk),
	}, nil
}
