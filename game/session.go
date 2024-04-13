package game

type Session struct {
	Level int

	Prog *BotProg
}

func NewSession() *Session {
	return &Session{
		Prog: NewBotProg(),
	}
}
