package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"session/cache"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func main() {
	cache.Init()
	t := cache.Store("loh", "misanea blea")
	fmt.Print(t)
	// r := cache.Get("loh")
	// fmt.Print(r)
	test()
}

func test() {
	port := 8080
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/join", joinHandl)
	http.Handle("/ss", http.RedirectHandler("https://9gag.com/", 302))
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
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
	// s := instance.Create("loh", 4)
	// i := *s
	temp, err := cache.Get("loh")
	if err != nil {
		log.Default().Println("Not found")
	}
	log.Default().Println(temp)
	w.Write([]byte(temp))
	// fmt.Fprintf(w, temp)
}
func joinHandl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pozno")
}
