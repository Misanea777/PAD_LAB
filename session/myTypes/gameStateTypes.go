package myTypes

type WorldMap struct {
	Id uint64
}

type PosAsID struct {
	PosX int64
	PosY int64
}

type Chunk struct {
	Size  int
	Id    PosAsID
	Tiles []Tile
}

type Tile struct {
	TileType int
	Id       PosAsID
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
