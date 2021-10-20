package instance

import (
	"log"
	"session/cache"
	"session/myTypes"
)

type WorldMap struct {
	Id uint64
}

func NewWorldMap() *WorldMap {
	wm := WorldMap{Id: 0}
	return &wm
}

func (wm WorldMap) GetChunk(id myTypes.PosAsID) *myTypes.Chunk {
	key := cache.ChunkKey{WordlId: wm.Id, PosX: id.PosX, PosY: id.PosY}
	res, err := cache.Get(key)
	if err == nil {
		log.Default().Println("found")
		return &res
	}
	log.Default().Println("generated")
	newChnk := genRandChunk(32, id)
	cache.Store(key, *newChnk)

	return newChnk
}
