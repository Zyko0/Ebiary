package game

import (
	"math/rand"
)

const (
	RoomSize = 1024
)

func AppendRoomCrates(crates []Crate, roomX, roomY int) []Crate {
	const gridsize = RoomSize / CrateSize

	rng := rand.New(rand.NewSource(int64(roomX) | (int64(roomY) << 32)))
	count := rng.Intn(3) + 1
	for i := 0; i < count; i++ {
	crategen:
		x := rng.Intn(gridsize-4) + 2
		y := rng.Intn(gridsize-4) + 2
		for _, c := range crates {
			if x == int(c.X) && y == int(c.Y) {
				goto crategen
			}
		}
		crates = append(crates, Crate{
			X: float64(x * CrateSize),
			Y: float64(y * CrateSize),
		})
	}

	return crates
}
