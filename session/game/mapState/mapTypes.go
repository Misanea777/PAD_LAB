package mapstate

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
