package game

type Session struct {
	Level int
}

func NewSession() *Session {
	return &Session{}
}
