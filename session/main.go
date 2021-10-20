package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"session/cache"
	"session/conn"
	"session/instance"

	"encoding/json"

	"session/myTypes"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func main() {
	cache.Init()
	conn.PingEureka()

	test()
}

func test() {
	rtr := mux.NewRouter()
	s := rtr.PathPrefix("/session").Subrouter()

	port := 8080
	s.HandleFunc("/create", createHandler)
	s.HandleFunc("/join", joinHandl)
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

	i := instance.NewWorldMap()
	c := *i.GetChunk(myTypes.PosAsID{PosX: 0, PosY: 0})
	js, e := json.Marshal(c)
	if e == nil {
		w.Write([]byte(js))
	}
}

func joinHandl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pozno")
}
