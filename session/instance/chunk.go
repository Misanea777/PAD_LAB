package instance

import (
	"session/myTypes"
)

func genRandChunk(size int, id myTypes.PosAsID) *myTypes.Chunk {
	var t []myTypes.Tile
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			t = append(t, *GenRandTile(myTypes.PosAsID{int64(i), int64(j)}, 2))
		}
	}
	c := myTypes.Chunk{Size: size, Id: id, Tiles: t}
	return &c
}
