package main

import (
	"fmt"
	"net/http"
	"session/cache"
	"session/instance"
	"strconv"
)

func main() {
	cache.Init()
	cache.Store("loh", 5)
	r := cache.Get("loh")
	fmt.Print(r)
	// test()
}

func test() {
	port := 8080
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/join", joinHandl)
	http.Handle("/ss", http.RedirectHandler("https://9gag.com/", 302))
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	s := instance.Create("loh", 4)
	i := *s
	fmt.Fprintf(w, strconv.Itoa(i.Id))
}
func joinHandl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pozno")
}
