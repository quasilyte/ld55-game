package game

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
)

const (
	MaxBranches     = 3
	MaxInstructions = 10
)

type BotProg struct {
	MovementThread *ProgThread
	Weapon1Thread  *ProgThread
	Weapon2Thread  *ProgThread
	DefThread      *ProgThread
}

func NewBotProg() *BotProg {
	return &BotProg{
		MovementThread: &ProgThread{Kind: MovementThread},
		Weapon1Thread:  &ProgThread{Kind: Weapon1Thread},
		Weapon2Thread:  &ProgThread{Kind: Weapon2Thread},
		DefThread:      &ProgThread{Kind: DefThread},
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
	MovementThread ThreadKind = 1 << iota
	Weapon1Thread
	Weapon2Thread
	DefThread
)

const (
	anyThreadMask = MovementThread | Weapon1Thread | Weapon2Thread | DefThread
	anyWeaponMask = Weapon1Thread | Weapon2Thread
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

	Mask ThreadKind

	MaxParam     float64
	DefaultParam float64

	StackInType  string
	StackOutType string
	StackChange  int
	Param        bool
	Cond         bool
}

var ProgInstInfoTab = func() []*ProgInstructionInfo {
	insts := []*ProgInstructionInfo{
		NopInstruction: {},

		RandomPosInstruction:     {Icon: assets.ImageIconRandomPos, Mask: anyThreadMask, StackChange: +1, StackOutType: "vec2"},
		RandomOffsetInstruction:  {Icon: assets.ImageIconRandomOffset, Param: true, MaxParam: 999, DefaultParam: 32, Mask: anyThreadMask, StackInType: "vec2", StackOutType: "vec2"},
		VesselPosInstruction:     {Icon: assets.ImageIconSelfPos, Mask: anyThreadMask, StackChange: +1, StackOutType: "vec2"},
		TargetPosInstruction:     {Icon: assets.ImageIconTargetPos, Mask: anyThreadMask, StackChange: +1, StackOutType: "vec2"},
		CenterPosInstruction:     {Icon: assets.ImageIconCenterPos, Mask: anyThreadMask, StackChange: +1, StackOutType: "vec2"},
		ChanceInstruction:        {Icon: assets.ImageIconRand, Param: true, MaxParam: 100, DefaultParam: 50, Cond: true, Mask: anyThreadMask},
		IsLtInstruction:          {Icon: assets.ImageIconIsLt, Param: true, MaxParam: 9999, DefaultParam: 50, Cond: true, Mask: anyThreadMask, StackChange: -1, StackInType: "float"},
		IsGtInstruction:          {Icon: assets.ImageIconIsGt, Param: true, MaxParam: 9999, DefaultParam: 50, Cond: true, Mask: anyThreadMask, StackChange: -1, StackInType: "float"},
		DistanceToInstruction:    {Icon: assets.ImageIconDistanceTo, Mask: anyThreadMask, StackChange: 0, StackInType: "vec2", StackOutType: "float"},
		HealthPercentInstruction: {Icon: assets.ImageIconSelfHealthPercent, Mask: anyThreadMask, StackChange: +1, StackOutType: "float"},
		EnergyPercentInstruction: {Icon: assets.ImageIconSelfEnergyPercent, Mask: anyThreadMask, StackChange: +1, StackOutType: "float"},

		RotateToInstruction:      {Icon: assets.ImageIconRotateTo, Mask: MovementThread, StackChange: -1, StackInType: "vec2"},
		MoveForwardInstruction:   {Icon: assets.ImageIconMoveForward, Param: true, MaxParam: 999, DefaultParam: 100, Mask: MovementThread},
		MoveAndRotateInstruction: {Icon: assets.ImageIconMoveAndRotate, Param: true, MaxParam: 999, DefaultParam: 100, Mask: MovementThread, StackChange: -1, StackInType: "vec2"},

		SnapShotInstruction:   {Icon: assets.ImageIconSnapShot, Mask: anyWeaponMask},
		NormalShotInstruction: {Icon: assets.ImageIconNormalShot, Mask: anyWeaponMask},
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
