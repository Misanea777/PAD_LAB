package instance

import (
	"session/myTypes"
)

type GameState struct {
	WrldMap WorldMap
	Players []myTypes.Player
}

func NewGameState(id uint64) *GameState {
	gms := GameState{WrldMap: WorldMap{Id: id}}
	return &gms
}

func (gms GameState) JoinPlayer(usnm string) error {
	for _, player := range gms.Players {
		if usnm == player.Usnm {
			player.Status = true
			return nil
		}
	}
	newPlayer := myTypes.Player{Usnm: usnm, PlayerPos: myTypes.Pos{X: 0, Y: 0}, Status: true}
	gms.Players = append(gms.Players, newPlayer)
	return nil
}

func (gms GameState) addPlayer(pl myTypes.Player) error {
	gms.Players = append(gms.Players, pl)
	return nil
}

func (gms GameState) GetChunks(ids []myTypes.PosAsID) ([]myTypes.Chunk, error) {
	var reqChunks []myTypes.Chunk
	for _, id := range ids {
		c := *gms.WrldMap.GetChunk(id)
		reqChunks = append(reqChunks, c)
	}
	return reqChunks, nil
}
