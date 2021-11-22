package cntrl

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"encoding/json"

	"context"
	"time"

	"github.com/gorilla/mux"
)

const taskLimit = 2

var timeoutAt = 75000 * time.Millisecond

type limitHandler struct {
	connc   chan struct{}
	handler http.Handler
}

func (h *limitHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	select {
	case <-h.connc:
		defer func() { h.connc <- struct{}{} }()
		ctx, cancel := context.WithTimeout(context.Background(), timeoutAt)
		defer cancel()
		h.handler.ServeHTTP(w, req.WithContext(ctx))

	default:
		http.Error(w, "503 too busy", http.StatusServiceUnavailable)
	}
}

func NewLimitHandler(maxConns int, handler http.Handler) http.Handler {
	h := &limitHandler{
		connc:   make(chan struct{}, maxConns),
		handler: handler,
	}
	for i := 0; i < maxConns; i++ {
		h.connc <- struct{}{}
	}
	return h
}

type Args struct {
	A, B int
}

func Init() {
	rtr := mux.NewRouter()
	s := rtr.PathPrefix("/lobby").Subrouter()

	port := 8084
	s.HandleFunc("/games/active", gamesActiveHandler).Methods("GET")

	lmHandl := NewLimitHandler(taskLimit, rtr)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.TimeoutHandler(lmHandl, timeoutAt, "timeout!!!"),
		Addr:         fmt.Sprintf(":%v", port),
	}

	srv.ListenAndServe()
}

func gamesActiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("getActive")

	resp, _ := json.Marshal(map[string]interface{}{"resp": rpcCall()})
	w.Write([]byte(resp))
}

func rpcCall() int {
	client, err := rpc.DialHTTP("tcp", "session:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.GetAll", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)
	return reply
}
