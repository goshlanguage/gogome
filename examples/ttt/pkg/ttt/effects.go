package ttt

import (
	"math/rand"
	"github.com/ryanhartje/gogome/pkg/engine"
)

// Wondering creates a casual wondering effect, the entity will move arbitrarily on each update
func Wondering(player *engine.Player) {
		// 10% of the time, drunkenly step left or right
	if rand.Intn(10) < 1 {
	// 10% of the time, drunkenly step left or right
		if rand.Intn(2) == 1 {
			player.Move(0, -1)
		} else {
			player.Move(0, 1)
	}
		}
	}
	if rand.Intn(10) < 1 {
		if rand.Intn(2) == 1 {
			player.Move(-1, 0)
		} else {
			player.Move(1, 0)
	}
		}
	}
}

// Gravity emulates falling until the player reaches 0.
// Maybe it shouldn't have bounds, but instead kill the player if y > winH + player.SizeY (once the player is off screen)
func Gravity(player *engine.Player) {
	// keep in mind that SDL has a reversed Y coordinate system, where 0,0 is the top [left] of the screen, and winH,0 is the bottom [left] of the screen
	if player.Y < winH-float64(player.SizeY*2) {
		player.Y += 8
		player.SpriteXPos = 9
	}
}


// Bleeding renders randomized red boxes to the screen
func Bleeding(player *engine.Player) {
	droplets := rand.Intn(10)
	for i := 0; i < droplets; i++ {
		randX := int32(player.X+8) + rand.Int31n(player.SizeX)
		randY := int32(player.Y+16) + rand.Int31n(player.SizeY)
		gfx.RoundedBoxColor(renderer, randX, randY, randX+4, randY+4, 1, sdl.Color{255, 0, 0, uint8(rand.Intn(255))})
	}
}

// Hitbox is a debugging friendly way to visualize where sprites collide. Assign this effect to entities to draw a hitbox over the sprite
func Hitbox(player *engine.Player) {
	gfx.BoxColor(renderer, int32(player.X), int32(player.Y), int32(player.X)+16, int32(player.Y)+32, sdl.Color{255, 255, 255, 255})
}
