package game

type BotProg struct {
	Threads []ProgThread
}

type ThreadKind int

const (
	UnknownThread ThreadKind = iota
	MovementThread
	Weapon1Thread
)

type ProgThread struct {
	Kind     ThreadKind
	Branches []ProgBranch
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

		RandomPosInstruction:    {},
		RandomOffsetInstruction: {},
		VesselPosInstruction:    {},
		TargetPosInstruction:    {},
		CenterPosInstruction:    {},
		ChanceInstruction:       {Cond: true},
		IsLtInstruction:         {Cond: true},
		IsGtInstruction:         {Cond: true},
		DistanceToInstruction:   {},

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
