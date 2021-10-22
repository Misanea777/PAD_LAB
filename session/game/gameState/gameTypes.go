package gamestate

type GameState struct {
	Id      uint32
	Players []Player
}

type Pos struct {
	X float64
	Y float64
}

type Player struct {
	Usnm      string
	PlayerPos Pos
	Status    bool
}
