package instance

type Pos struct {
	X float64
	Y float64
}

type Player struct {
	Usnm      string
	PlayerPos Pos
	status    bool
}
