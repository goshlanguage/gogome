package engine

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var grid = false

// Level provides us a way to scroll and update the level
type Level struct {
	// BGFile is the filepath to the background
	BGFile string
	// Give level a debug passthrough
	Debug            bool
	CameraX, CameraY int
	// EntityMap entities and draws them on the map when they're in focus
	EntityMap map[int]map[int]Entity
	// represents how many pixels we scroll per cycle. default 16
	ScrollSpeed int
	Sounds      map[string][]*mix.Chunk
	Texture     *sdl.Texture
	// TileMap is a matrix representing the map
	// A z index is arbitrarily establish by the slice ordering
	// the 0th object is drawn first, with each subsequent object drawn on top of the previous one
	// TileMap[x][y]
	TileMap  map[int]map[int][]Tile
	TileSize int
	// Current level's coordinates
	X, Y int
	// Size coords to stop scrolling approriately
	XSize int
	YSize int
}

// Tile represents a tile in a tilemap. This might be a 16x16 sprite or a 16x128 tile.
type Tile struct {
	Name string
	X0   int32
	X1   int32
	Y0   int32
	Y1   int32
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
		BGFile:      filepath,
		CameraX:     640,
		CameraY:     480,
		Texture:     bgTexture,
		TileSize:    32,
		ScrollSpeed: 8,
		Sounds:      make(map[string][]*mix.Chunk),
	}, nil
}

// NewRandomizedLevel takes in the filepath of a level's background, and a renderer
func NewRandomizedLevel(filepath string, renderer *sdl.Renderer) (*Level, error) {
	img, err := sdl.LoadBMP(filepath)

	checkErr(err)
	defer img.Free()

	bgTexture, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return &Level{}, err
	}

	level := &Level{
		BGFile:      filepath,
		CameraX:     640,
		CameraY:     480,
		Texture:     bgTexture,
		TileSize:    32,
		ScrollSpeed: 8,
		Sounds:      make(map[string][]*mix.Chunk),
	}

	grass := Tile{Name: "grass", X0: 0, X1: 16, Y0: 0, Y1: 16}
	grass2 := Tile{Name: "grass2", X0: 272, X1: 303, Y0: 464, Y1: 495}
	bush := Tile{Name: "bush", X0: 32, X1: 48, Y0: 224, Y1: 240}
	mapping := map[int]map[int][]Tile{}
	level.EntityMap = map[int]map[int]Entity{}

	// Bootstrap TileMap for the background, and entity map to render entities on top of
	// iterate by the x and y values of the sprite's width and height, so that you don't
	// draw over other tiles.
	for x := 0; x < (winW * 10); x += level.TileSize {
		mapping[x] = make(map[int][]Tile)
		level.EntityMap[x] = make(map[int]Entity)
		for y := 0; y < (winH * 10); y += level.TileSize {
			// populate map with Tiles
			mapping[x][y] = []Tile{}
			mapping[x][y] = append(mapping[x][y], grass)
			// half the time, give us different grass
			if rand.Intn(10) > 5 {
				mapping[x][y] = append(mapping[x][y], grass2)
			}
			if rand.Intn(10) == 0 {
				mapping[x][y] = append(mapping[x][y], bush)
			}
			level.EntityMap[x][y] = nil
		}
	}
	level.TileMap = mapping

	return level, nil
}

// Draw takes the camera viewport and renders it to the screen
func (level *Level) Draw(renderer *sdl.Renderer) {
	// Render level to window tile by tile
	for x := 0; x < level.CameraX; x += level.ScrollSpeed {
		for y := 0; y < level.CameraY; y += level.ScrollSpeed {
			tiles := level.TileMap[level.X+x][level.Y+y]
			for _, tile := range tiles {
				width := tile.X1 - tile.X0
				height := tile.Y1 - tile.Y0

				// Render the background of the level
				renderer.Copy(
					level.Texture,
					&sdl.Rect{X: tile.X0, Y: tile.Y0, W: width, H: height},
					&sdl.Rect{X: int32(x), Y: int32(y), W: width * 2, H: height * 2},
				)
			}

			// on x=0 and y=0, we need to render partials on the left side of the screen when scrolling
			// we're using 16 here because that's the prerendered tile's size (16x16)
			// This will need to be refactored to support dynamically sized tiles
			// TODO - TileSize handling is VERY fragile, This should be reworked/redesigned to be more resilient
			if x == 0 && level.X%level.TileSize != 0 {
				xOffset := level.X % level.TileSize
				halfSet := xOffset / 2 // half set is used to determine where on the texture map to find the asset

				tiles := level.TileMap[level.X-xOffset][level.Y+y]
				for _, tile := range tiles {
					width := level.TileSize - xOffset
					height := tile.Y1 - tile.Y0

					// Render the background of the level
					renderer.Copy(
						level.Texture,
						&sdl.Rect{X: tile.X0 + int32(halfSet), Y: tile.Y0, W: int32(width), H: height},
						&sdl.Rect{X: int32(x), Y: int32(y), W: int32(width) * 2, H: height * 2},
					)
					if y < 32 {
						fmt.Printf("Rendering %d,%d at level %d,%d from pixel %d,%d\t ", x, y, level.X, level.Y, tile.X0+int32(xOffset), tile.Y0)
						fmt.Printf("Using TileMap[%d][%d]: %v\t", level.X-xOffset, level.Y+y, level.TileMap[level.X-xOffset][level.Y+y])
						fmt.Printf("Placed at %d,%d w: %d h: %d\n", x, y, width, height)
					}
				} // end of tiles for loop
			} // end if for x == 0 edge case

			// on x=0 and y=0, we need to render partials on the left side of the screen when scrolling
			// we're using 16 here because that's the prerendered tile's size (16x16)
			// This will need to be refactored to support dynamically sized tiles
			if y == 0 && level.Y%level.TileSize != 0 {
				yOffset := level.Y % level.TileSize
				halfSet := yOffset / 2

				tiles := level.TileMap[level.X+x][level.Y-yOffset]
				for _, tile := range tiles {
					width := tile.X1 - tile.X0
					height := level.TileSize - yOffset

					// Render the background of the level
					renderer.Copy(
						level.Texture,
						&sdl.Rect{X: tile.X0, Y: tile.Y0 + int32(halfSet), W: width, H: int32(height)},
						&sdl.Rect{X: int32(x), Y: int32(y), W: width * 2, H: int32(height) * 2},
					)
				} // end of tiles for loop
			} // end if for y == 0 edge case

			if level.Debug {
				if x%level.TileSize == 0 && y%level.TileSize == 0 {
					// draw vertical grid lines
					gfx.LineRGBA(renderer, int32(x), 0, int32(x), int32(winH), 100, 0, 0, 100)
					// draw horizontal line
					gfx.LineRGBA(renderer, 0, int32(y), int32(winW), int32(y), 100, 0, 0, 100)
				}
			}
		}
	}
}

// Update watches keybindings and scrolls as necessary
func (level *Level) Update() {
	keys := sdl.GetKeyboardState()

	// UP
	if keys[sdl.SCANCODE_W] == 1 {
		if level.Y > 0 {
			level.Y -= level.ScrollSpeed
		} else {
			level.Y = 0
		}
	}
	// DOWN
	if keys[sdl.SCANCODE_S] == 1 {
		if level.Y < level.YSize {
			level.Y += level.ScrollSpeed
		} else {
			level.Y = level.YSize
		}
	}
	// LEFT
	if keys[sdl.SCANCODE_A] == 1 {
		if level.X > 0 {
			level.X -= level.ScrollSpeed
		} else {
			level.X = 0
		}
	}
	// RIGHT
	if keys[sdl.SCANCODE_D] == 1 {
		if level.X < level.XSize {
			level.X += level.ScrollSpeed
		} else {
			level.X = level.XSize
		}
	}
}
