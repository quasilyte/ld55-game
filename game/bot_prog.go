package game

type BotProg struct {
	Threads []ProgThread
}

type ThreadKind int

const (
	UnknownThread ThreadKind = iota
	MovementThread
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

	// Movement instructions.

	RotateToInstruction
	MoveForwardInstruction
	MoveAndRotateInstruction
)

type ProgInstructionInfo struct {
	Kind InstructionKind
}

var ProgInstInfoTab = func() []*ProgInstructionInfo {
	insts := []*ProgInstructionInfo{
		UnknownInstruction: {},

		RandomPosInstruction:    {},
		RandomOffsetInstruction: {},
		VesselPosInstruction:    {},
		TargetPosInstruction:    {},

		RotateToInstruction:      {},
		MoveForwardInstruction:   {},
		MoveAndRotateInstruction: {},
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
