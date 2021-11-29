package main

import (
	"session/cache"
	"session/conn"
	"session/db"

	"session/mux"

	"time"
)

func main() {
	db.Connect()
	cache.Init()
	time.Sleep(time.Second * 20)
	conn.PingEureka()

	// test()

	mux.Init()

}

func test() {
	// db.UpdateSt()
}
