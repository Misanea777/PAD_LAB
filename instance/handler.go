package instance

type Pos struct {
	posX float32
	posY float32
}

type Player struct {
	name   string
	pos    Pos
	status bool
}

func newPlayer(name string) *Player {
	p := Player{pos: Pos{posX: 0, posY: 0}}
	p.name = name
	return &p
}

type SessionMap struct {
}

type State struct {
}

type Session struct {
	Id          int
	maxPlayers  int
	players     []Player
	mapInstance SessionMap
	gameState   State
}

func newSession(maxPlayers int) *Session {
	s := Session{maxPlayers: maxPlayers}
	return &s
}

func (s Session) addPlayer(p *Player) {
	s.players = append(s.players, *p)
}

func Create(usnm string, maxPlayers int) *Session {
	s := newSession(maxPlayers)
	p := newPlayer(usnm)
	s.addPlayer(p)
	return s
}
