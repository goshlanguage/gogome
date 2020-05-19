package engine

import (
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
	//bush := Tile{Name: "bush", X0: 32, X1: 48, Y0: 224, Y1: 240}
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
			if x == 0 || y == 0 {
				mapping[x][y] = append(mapping[x][y], grass2)
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

			// Here's a little experiment, if level.X > 0, let's try to draw the last tile at a negative x value and see what we get
			if (level.X > 0 || level.Y > 0) && x == 0 && y == 0 {
				xOffset := level.X % level.TileSize
				tileX := level.X - xOffset
				yOffset := level.Y % level.TileSize
				tileY := level.Y - yOffset

				tiles := level.TileMap[tileX][tileY]
				// Now that we've processed offsets and our tile for this iteration, range throuhg tiles and draw them
				for _, tile := range tiles {
					width := tile.X1 - tile.X0
					height := tile.Y1 - tile.Y0

					// Render the background of the level
					renderer.Copy(
						level.Texture,
						&sdl.Rect{X: tile.X0 + int32(xOffset/2), Y: tile.Y0, W: int32(width), H: int32(height)},
						&sdl.Rect{X: int32(-level.TileSize + xOffset), Y: int32(-level.TileSize + yOffset), W: int32(level.TileSize), H: int32(level.TileSize)},
					)
				}
			}

			// Here is the information we need to know where to render all tiles
			var xOffset int
			var yOffset int

			tileX := level.X + x
			tileY := level.Y + y

			if x == 0 && level.X%level.TileSize != 0 {
				xOffset = level.X % level.TileSize
				tileX = level.X - xOffset
			} else if y == 0 && level.Y%level.TileSize != 0 {
				yOffset = level.Y % level.TileSize
				tileY = level.Y - yOffset
			}

			tiles := level.TileMap[tileX][tileY]
			// Now that we've processed offsets and our tile for this iteration, range throuhg tiles and draw them
			for _, tile := range tiles {
				width := tile.X1 - tile.X0
				height := tile.Y1 - tile.Y0

				// Render the background of the level
				renderer.Copy(
					level.Texture,
					&sdl.Rect{X: tile.X0 + int32(xOffset/2), Y: tile.Y0 + int32(yOffset/2), W: int32(width), H: int32(height)},
					&sdl.Rect{X: int32(x), Y: int32(y), W: int32(level.TileSize), H: int32(level.TileSize)},
				)
			}

		}
	}

	// Render a grid to the screen if debug is on
	if level.Debug {
		// Render a grid. Draw lines from 0 to the width/height of the screen along the X and Y axis, incrementing by our tilesize
		for i := 0; i < winW; i += level.TileSize {
			gfx.LineRGBA(renderer, int32(0), int32(i), int32(winW), int32(i), 100, 0, 0, 100)
			gfx.LineRGBA(renderer, int32(i), int32(0), int32(i), int32(winH), 100, 0, 0, 100)
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
