package engine

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

// Level provides us a way to scroll and update the level
type Level struct {
	// BGFile is the filepath to the background
	BGFile string
	// EntityMap entities and draws them on the map when they're in focus
	EntityMap map[int]map[int]Entity
	// represents how many pixels we scroll per cycle. default 16
	ScrollSpeed int
	Sounds      map[string][]*mix.Chunk
	Texture     *sdl.Texture
	// TileMap is a matrix representing the map
	// TileMap[x][y]
	TileMap map[int]map[int]Tile
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
		Texture:     bgTexture,
		ScrollSpeed: 16,
		Sounds:      make(map[string][]*mix.Chunk),
	}, nil
}

// Update watches keybindings and scrolls as necessary
func (l *Level) Update() {
	keys := sdl.GetKeyboardState()

	// UP
	if keys[sdl.SCANCODE_W] == 1 {
		if l.Y > 0 {
			l.Y -= l.ScrollSpeed
		} else {
			l.Y = 0
		}
	}
	// DOWN
	if keys[sdl.SCANCODE_S] == 1 {
		if l.Y < l.YSize {
			l.Y += l.ScrollSpeed
		} else {
			l.Y = l.YSize
		}
	}
	// LEFT
	if keys[sdl.SCANCODE_A] == 1 {
		if l.X > 0 {
			l.X -= l.ScrollSpeed
		} else {
			l.X = 0
		}
	}
	// RIGHT
	if keys[sdl.SCANCODE_D] == 1 {
		if l.X < l.XSize {
			l.X += l.ScrollSpeed
		} else {
			l.X = l.XSize
		}
	}
}
