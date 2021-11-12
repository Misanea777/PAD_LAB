package actions

import (
	"errors"
	"session/cache"
	"session/db"
	state "session/game/gameState"
	mapSt "session/game/mapState"
)

func NewGameStateWithId(id uint32) *state.GameState {
	gms := state.GameState{Id: id}
	cache.StoreSt(&gms)
	return &gms
}

func NewGameState() *state.GameState {
	gms := NewGameStateWithId(12345)
	return gms
}

func JoinPlayer(usnm string, gms *state.GameState) (state.Pos, error) {
	for index, player := range gms.Players {
		if usnm == player.Usnm {
			gms.Players[index].Status = true
			cache.StoreSt(gms)
			return player.PlayerPos, nil
		}
	}
	newPlayer := state.Player{Usnm: usnm, PlayerPos: state.Pos{X: 0, Y: 0}, Status: true}
	gms.Players = append(gms.Players, newPlayer)
	cache.StoreSt(gms)
	return newPlayer.PlayerPos, nil
}

func LeavePlayer(usnm string, gms *state.GameState) error {
	for index, player := range gms.Players {
		if usnm == player.Usnm {
			player.Status = false
			gms.Players[index] = player
			cache.StoreSt(gms)
			checkIfAnyPlayersPresent(gms)
			return nil
		}
	}
	return errors.New("player not found")
}

func checkIfAnyPlayersPresent(gms *state.GameState) bool {
	for _, player := range gms.Players {
		if player.Status {
			return true
		}
	}
	cache.DeleteSt(gms.Id)
	db.UpdateSt(gms)
	return false
}

func FindPlayer(usnm string, gms *state.GameState) (state.Player, error) {
	var p state.Player
	for _, player := range gms.Players {
		if usnm == player.Usnm {
			// player.Status = true
			return player, nil
		}
	}
	return p, errors.New("player not found")
}

func GetAllPlayersExcept(usnm string, gms *state.GameState) []state.Player {
	var p []state.Player
	for _, player := range gms.Players {
		if usnm == player.Usnm {
			continue
		}
		p = append(p, player)
	}
	return p
}

func ModifyPlayer(usnm string, gms *state.GameState, newPl state.Player) error {
	for index, player := range gms.Players {
		if usnm == player.Usnm {
			gms.Players[index] = newPl
			cache.StoreSt(gms)
			return nil
		}
	}
	return errors.New("player not found")
}

func GetChunks(ids []mapSt.PosAsID, gms *state.GameState) ([]mapSt.Chunk, error) {
	var reqChunks []mapSt.Chunk
	for _, id := range ids {
		c := *GetChunk(gms.Id, id)
		reqChunks = append(reqChunks, c)
	}
	return reqChunks, nil
}

func Getst(id uint32) (state.GameState, error) {
	gm, err := cache.GetSt(id)
	if err == nil {
		return gm, nil
	}

	gm, err = db.GetSt(id)
	if err == nil {
		return gm, nil
	}

	return gm, err
}
