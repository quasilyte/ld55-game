package game

type Session struct {
	Level int

	Prog *BotProg

	Weapons        []*WeaponDesign
	VesselDesign   *VesselDesign
	ArtifactDesign *ArtifactDesign
}

type SavedVessel struct {
	Weapons      []string
	VesselDesign string
	Artifact     string
	Prog         *BotProg
}

func NewSession() *Session {
	return &Session{
		Prog: NewBotProg(),
	}
}
