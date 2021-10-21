package instance

import (
	"session/myTypes"
)

type GameState struct {
	wrldMap WorldMap
	players []Player
}

func NewGameState(id uint64) *GameState {
	gms := GameState{wrldMap: WorldMap{Id: id}}
	return &gms
}

func (gms GameState) JoinPlayer(usnm string) error {
	for _, player := range gms.players {
		if usnm == player.Usnm {
			player.status = true
			return nil
		}
	}
	newPlayer := Player{Usnm: usnm, PlayerPos: Pos{X: 0, Y: 0}, status: true}
	gms.players = append(gms.players, newPlayer)
	return nil
}

func (gms GameState) addPlayer(pl Player) error {
	gms.players = append(gms.players, pl)
	return nil
}

func (gms GameState) GetChunks(ids []myTypes.PosAsID) ([]myTypes.Chunk, error) {
	var reqChunks []myTypes.Chunk
	for _, id := range ids {
		c := *gms.wrldMap.GetChunk(id)
		reqChunks = append(reqChunks, c)
	}
	return reqChunks, nil
}
