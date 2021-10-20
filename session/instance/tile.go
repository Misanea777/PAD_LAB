package instance

import (
	"math/rand"
	"session/myTypes"
)

func GenRandTile(id myTypes.PosAsID, max int) *myTypes.Tile {
	t := myTypes.Tile{TileType: rand.Intn(max), Id: id}
	return &t
}
