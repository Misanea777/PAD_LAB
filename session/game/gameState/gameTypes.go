package gamestate

type GameState struct {
	Id      uint32   `bson:"_id"`
	Players []Player `bson:"players"`
}

type Pos struct {
	X float64 `bson:"x"`
	Y float64 `bson:"y"`
}

type Player struct {
	Usnm      string `usnm:"id"`
	PlayerPos Pos
	Status    bool `status:"id"`
}
