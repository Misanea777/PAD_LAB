package main

import (
	"log"
	"session/cache"
	"session/conn"
	"session/db"

	"session/mux"

	"time"

	"session/game/actions"
)

func main() {
	db.Connect()
	cache.Init()
	time.Sleep(time.Second * 10)
	conn.PingEureka()

	// test()

	mux.Init()

}

func test() {
	gm := actions.NewGameState()
	actions.JoinPlayer("misa", gm)
	actions.JoinPlayer("grishsa", gm)
	db.UpdateSt(gm)
	res, err := db.GetSt(gm.Id)
	if err != nil {
		log.Default().Println(err)
	}
	log.Default().Println(len(res.Players))
}
