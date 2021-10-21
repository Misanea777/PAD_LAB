package main

import (
	"fmt"
	"log"
	"net/http"
	"session/cache"
	"session/conn"
	"session/instance"
	"session/myTypes"
	"time"

	"encoding/json"

	"session/db"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func main() {
	db.Connect()
	cache.Init()
	time.Sleep(time.Second * 10)
	conn.PingEureka()

	test()
}

func test() {
	rtr := mux.NewRouter()
	s := rtr.PathPrefix("/session").Subrouter()

	port := 8081
	s.HandleFunc("/create", createHandler)
	s.HandleFunc("/join", joinHandl).Methods("POST")
	s.Handle("/ss", http.RedirectHandler("https://9gag.com/", 302))
	s.HandleFunc("/echo", echo)
	http.ListenAndServe(fmt.Sprintf(":%v", port), rtr)
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("reached")

	i := instance.NewGameState(2)
	ids := []myTypes.PosAsID{{PosX: 1, PosY: 0}, {PosX: 0, PosY: 2}}
	c, _ := i.GetChunks(ids)
	db.SaveState(i.WrldMap.Id, i.Players)

	js, e := json.Marshal(c)
	if e == nil {
		w.Write([]byte(js))
	}
}

type JoinReq struct {
	Usnm      string
	SessionId uint64
}

func joinHandl(w http.ResponseWriter, r *http.Request) {
	var user JoinReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "good")
}
