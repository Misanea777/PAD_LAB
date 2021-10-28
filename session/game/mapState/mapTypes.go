package mapstate

type PosAsID struct {
	PosX int64 `bson:"posx"`
	PosY int64 `bson:"posy"`
}

type Chunk struct {
	Size  int `bson:"size"`
	Id    PosAsID
	Tiles []Tile
}

type Tile struct {
	TileType int `bson:"title_type"`
	Id       PosAsID
}
