package mux

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"session/cache"
	"session/game/actions"
	gamestate "session/game/gameState"
	mapstate "session/game/mapState"
)

var upgrader = websocket.Upgrader{} // use default options

func Init() {
	rtr := mux.NewRouter()
	s := rtr.PathPrefix("/session").Subrouter()

	port := 8081
	s.HandleFunc("/create", createHandler).Methods("POST")
	s.HandleFunc("/join", joinHandl).Methods("POST")
	s.HandleFunc("/get/chunks", getChnkHandler).Methods("POST")
	s.HandleFunc("/get/players", getPlHandle).Methods("POST")
	s.HandleFunc("/update/player", updatePlHandle).Methods("POST")
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

type CreateReq struct {
	Usnm string
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var cReq CreateReq
	err := json.NewDecoder(r.Body).Decode(&cReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gm := actions.NewGameState()
	pos, _ := actions.JoinPlayer(cReq.Usnm, gm)
	resp, _ := json.Marshal(map[string]interface{}{"id": gm.Id, "posX": pos.X, "posY": pos.Y})
	w.Write([]byte(resp))
}

type JoinReq struct {
	Usnm      string
	SessionId uint64
}

func joinHandl(w http.ResponseWriter, r *http.Request) {
	var joinReq JoinReq
	err := json.NewDecoder(r.Body).Decode(&joinReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gm, err := cache.GetSt(joinReq.SessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pos, _ := actions.JoinPlayer(joinReq.Usnm, &gm)
	resp, _ := json.Marshal(map[string]interface{}{"posX": pos.X, "posY": pos.Y})
	w.Write([]byte(resp))
}

type GetChnkReq struct {
	Usnm      string
	SessionId uint64
	ChnksIds  []mapstate.PosAsID
}

func getChnkHandler(w http.ResponseWriter, r *http.Request) {
	var getReq GetChnkReq
	err := json.NewDecoder(r.Body).Decode(&getReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gm, err := cache.GetSt(getReq.SessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chnks, _ := actions.GetChunks(getReq.ChnksIds, &gm)
	resp, _ := json.Marshal(chnks)
	w.Write([]byte(resp))
}

type GetPlReq struct {
	Usnm      string
	SessionId uint64
}

func getPlHandle(w http.ResponseWriter, r *http.Request) {
	var getReq GetPlReq
	err := json.NewDecoder(r.Body).Decode(&getReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gm, err := cache.GetSt(getReq.SessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player, err := actions.FindPlayer(getReq.Usnm, &gm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	restPlayers := actions.GetAllPlayersExcept(getReq.Usnm, &gm)
	resp, _ := json.Marshal(map[string]interface{}{"yourPlayer": player, "restPlayers": restPlayers})
	w.Write([]byte(resp))
}

type UpdatePlReq struct {
	Usnm      string
	SessionId uint64
	Pos       gamestate.Pos
}

func updatePlHandle(w http.ResponseWriter, r *http.Request) {
	var upReq UpdatePlReq
	err := json.NewDecoder(r.Body).Decode(&upReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gm, err := cache.GetSt(upReq.SessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newPl := gamestate.Player{Usnm: upReq.Usnm, PlayerPos: upReq.Pos, Status: true}
	err = actions.ModifyPlayer(upReq.Usnm, &gm, newPl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "good")
}
