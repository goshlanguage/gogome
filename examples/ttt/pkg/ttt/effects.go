package TTT

import (
	"math/rand"
	"github.com/ryanhartje/gogome/pkg/engine"
)

// Wondering creates a casual wondering effect, the entity will move arbitrarily on each update
func Wondering(player *engine.Player) {
		// 10% of the time, drunkenly step left or right
		if rand.Intn(10) < 1 {
			if rand.Intn(2) == 1 {
				player.Move(0, -1)
			} else {
				player.Move(0, 1)
			}
		}
		// another 10% of the time, drunkenly step forward or backward
		if rand.Intn(10) < 1 {
			if rand.Intn(2) == 1 {
				player.Move(-1, 0)
			} else {
				player.Move(1, 0)
			}
		}
	}