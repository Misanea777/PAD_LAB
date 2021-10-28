package actions

import (
	"log"
	"math/rand"
	"session/cache"
	mapSt "session/game/mapState"
)

func GetChunk(sessionId uint32, id mapSt.PosAsID) *mapSt.Chunk {
	key := cache.ChunkKey{WordlId: sessionId, PosX: id.PosX, PosY: id.PosY}
	res, err := cache.GetChnk(key)
	if err == nil {
		log.Default().Println("found in cache")
		return &res
	}

	log.Default().Println("generated")
	newChnk := genRandChunk(32, id)
	cache.StoreChnk(key, *newChnk)

	return newChnk
}

func GenRandTile(id mapSt.PosAsID, max int) *mapSt.Tile {
	t := mapSt.Tile{TileType: rand.Intn(max), Id: id}
	return &t
}

func genRandChunk(size int, id mapSt.PosAsID) *mapSt.Chunk {
	var t []mapSt.Tile
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			t = append(t, *GenRandTile(mapSt.PosAsID{int64(i), int64(j)}, 2))
		}
	}
	c := mapSt.Chunk{Size: size, Id: id, Tiles: t}
	return &c
}
