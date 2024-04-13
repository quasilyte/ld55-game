package game

type BotProg struct {
	MovementThread *ProgThread
	Weapon1Thread  *ProgThread
	Weapon2Thread  *ProgThread
	DefThread      *ProgThread
}

func NewBotProg() *BotProg {
	return &BotProg{
		MovementThread: &ProgThread{},
		Weapon1Thread:  &ProgThread{},
		Weapon2Thread:  &ProgThread{},
		DefThread:      &ProgThread{},
	}
}

func (p *BotProg) NumThreads() int {
	return 4
}

func (p *BotProg) EachThread(f func(i int, t *ProgThread)) {
	f(0, p.MovementThread)
	f(1, p.Weapon1Thread)
	f(2, p.Weapon2Thread)
	f(3, p.DefThread)
}

type ThreadKind int

const (
	UnknownThread ThreadKind = iota
	MovementThread
	Weapon1Thread
	Weapon2Thread
	DefThread
)

type ProgThread struct {
	Kind     ThreadKind
	Branches []*ProgBranch
}

type ProgBranch struct {
	Instructions []ProgInstruction
}

type InstructionKind int

const (
	UnknownInstruction InstructionKind = iota

	// Common instructions.

	RandomPosInstruction
	RandomOffsetInstruction
	VesselPosInstruction
	TargetPosInstruction
	CenterPosInstruction
	ChanceInstruction
	IsLtInstruction
	IsGtInstruction
	EnergyPercentInstruction
	HealthPercentInstruction
	DistanceToInstruction

	// Movement instructions.

	RotateToInstruction
	MoveForwardInstruction
	MoveAndRotateInstruction

	// Weapon instructions.

	SnapShotInstruction
	NormalShotInstruction
)

type ProgInstructionInfo struct {
	Kind InstructionKind

	Cond bool
}

var ProgInstInfoTab = func() []*ProgInstructionInfo {
	insts := []*ProgInstructionInfo{
		UnknownInstruction: {},

		RandomPosInstruction:     {},
		RandomOffsetInstruction:  {},
		VesselPosInstruction:     {},
		TargetPosInstruction:     {},
		CenterPosInstruction:     {},
		ChanceInstruction:        {Cond: true},
		IsLtInstruction:          {Cond: true},
		IsGtInstruction:          {Cond: true},
		DistanceToInstruction:    {},
		HealthPercentInstruction: {},
		EnergyPercentInstruction: {},

		RotateToInstruction:      {},
		MoveForwardInstruction:   {},
		MoveAndRotateInstruction: {},

		SnapShotInstruction:   {},
		NormalShotInstruction: {},
	}

	for kind, inst := range insts {
		inst.Kind = InstructionKind(kind)
	}

	return insts
}()

type ProgInstruction struct {
	Info *ProgInstructionInfo

	Params []any
}

func MakeInst(kind InstructionKind, params ...any) ProgInstruction {
	return ProgInstruction{
		Info:   ProgInstInfoTab[kind],
		Params: params,
	}
}
