package main

import (
	"session/cache"
	"session/conn"

	"session/mux"

	"time"
)

func main() {
	// db.Connect()
	cache.Init()
	time.Sleep(time.Second * 10)
	conn.PingEureka()

	mux.Init()

	// test()
}

// func test() {
// 	gm := actions.NewGameState(12345)
// 	actions.JoinPlayer("misa", gm)
// 	chnk := actions.GetChunk(gm.Id, mapSt.PosAsID{0,0})
// }
