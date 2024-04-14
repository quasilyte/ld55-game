package game

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
)

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
	NopInstruction InstructionKind = iota

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

	Icon resource.ImageID

	MaxParam float64

	Param bool
	Cond  bool
}

var ProgInstInfoTab = func() []*ProgInstructionInfo {
	insts := []*ProgInstructionInfo{
		NopInstruction: {},

		RandomPosInstruction:     {Icon: assets.ImageIconRandomPos},
		RandomOffsetInstruction:  {Param: true, MaxParam: 600},
		VesselPosInstruction:     {},
		TargetPosInstruction:     {Icon: assets.ImageIconTargetPos},
		CenterPosInstruction:     {},
		ChanceInstruction:        {Param: true, MaxParam: 100, Cond: true},
		IsLtInstruction:          {Icon: assets.ImageIconIsLt, Param: true, Cond: true},
		IsGtInstruction:          {Icon: assets.ImageIconIsGt, Param: true, Cond: true},
		DistanceToInstruction:    {Icon: assets.ImageIconDistanceTo},
		HealthPercentInstruction: {},
		EnergyPercentInstruction: {},

		RotateToInstruction:      {Icon: assets.ImageIconRotateTo},
		MoveForwardInstruction:   {Param: true, MaxParam: 1000, Icon: assets.ImageIconMoveForward},
		MoveAndRotateInstruction: {Param: true, MaxParam: 1000},

		SnapShotInstruction:   {Icon: assets.ImageIconSnapShot},
		NormalShotInstruction: {},
	}

	for kind, inst := range insts {
		inst.Kind = InstructionKind(kind)
	}

	return insts
}()

type ProgInstruction struct {
	Info *ProgInstructionInfo

	Param float64
}

func MakeInst(kind InstructionKind, param float64) ProgInstruction {
	inst := ProgInstruction{
		Info: ProgInstInfoTab[kind],
	}
	inst.SetParam(param)
	return inst
}

func (inst *ProgInstruction) SetParam(v float64) {
	inst.Param = gmath.Clamp(v, 0, inst.Info.MaxParam)
}
