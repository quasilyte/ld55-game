package game

type Session struct {
	Level int

	Prog *BotProg

	Weapons      []*WeaponDesign
	VesselDesign *VesselDesign
}

func NewSession() *Session {
	return &Session{
		Prog: NewBotProg(),
	}
}
