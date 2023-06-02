package main

import (
	"lobby/cntrl"
	"lobby/conn"
	"time"
)

func main() {
	time.Sleep(time.Second * 20)
	conn.PingEureka()

	cntrl.Init()
	//testdssd
}
