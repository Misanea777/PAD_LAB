package main

import (
	"log"
	"session/cache"
	"session/conn"
	"session/db"

	"session/mux"

	"time"

	"session/game/actions"
	mapstate "session/game/mapState"
)

func main() {
	db.Connect()
	cache.Init()
	time.Sleep(time.Second * 10)
	conn.PingEureka()

	test()

	mux.Init()

}

func test() {
	gm := actions.NewGameState()
	actions.JoinPlayer("misa", gm)
	chnk := actions.GetChunk(gm.Id, mapstate.PosAsID{1, 0})
	db.UpdateChnk(gm.Id, *chnk)
	db.GetChnk(gm.Id, mapstate.PosAsID{0, 0})
	log.Default().Println("HEREEEEEEEEE")
}
