package game

type Session struct {
	Level int

	Prog *BotProg

	Weapons      []*WeaponDesign
	VesselDesign *VesselDesign
}

type SavedVessel struct {
	Weapons      []string
	VesselDesign string
	Prog         *BotProg
}

func NewSession() *Session {
	return &Session{
		Prog: NewBotProg(),
	}
}
