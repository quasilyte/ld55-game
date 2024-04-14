package game

import (
	"encoding/json"

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

func (p *BotProg) Compact() *BotProg {
	return &BotProg{
		MovementThread: p.MovementThread.Compact(),
		Weapon1Thread:  p.Weapon1Thread.Compact(),
		Weapon2Thread:  p.Weapon2Thread.Compact(),
		DefThread:      p.DefThread.Compact(),
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

func (t *ProgThread) Compact() *ProgThread {
	t2 := &ProgThread{Kind: t.Kind}
	for _, b := range t.Branches {
		b2 := b.Compact()
		if len(b2.Instructions) == 0 {
			continue
		}
		t2.Branches = append(t2.Branches, b2)
	}
	return t2
}

type ProgBranch struct {
	Instructions []ProgInstruction
}

func (b *ProgBranch) Compact() *ProgBranch {
	b2 := &ProgBranch{}
	for _, inst := range b.Instructions {
		if inst.Info.Kind == NopInstruction {
			continue
		}
		b2.Instructions = append(b2.Instructions, inst)
	}
	return b2
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
	IsOutBoundsInstruction
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

	// Instruction added after saves were introduced.
	// They should be somewhere above, but I don't
	// have time to fix save/load.

	WaitInstruction
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
		IsOutBoundsInstruction:   {Icon: assets.ImageIconIsOutOfBounds, Cond: true, Mask: anyThreadMask, StackChange: -1, StackInType: "vec2"},
		DistanceToInstruction:    {Icon: assets.ImageIconDistanceTo, Mask: anyThreadMask, StackChange: 0, StackInType: "vec2", StackOutType: "float"},
		HealthPercentInstruction: {Icon: assets.ImageIconSelfHealthPercent, Mask: anyThreadMask, StackChange: +1, StackOutType: "float"},
		EnergyPercentInstruction: {Icon: assets.ImageIconSelfEnergyPercent, Mask: anyThreadMask, StackChange: +1, StackOutType: "float"},
		WaitInstruction:          {Icon: assets.ImageIconWait, Mask: anyThreadMask, Param: true, MaxParam: 9999, DefaultParam: 30},

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

type jsonProgInstruction struct {
	Tag   int
	Param float64
}

func (inst ProgInstruction) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonProgInstruction{
		Tag:   int(inst.Info.Kind),
		Param: inst.Param,
	})
}

func (inst *ProgInstruction) UnmarshalJSON(data []byte) error {
	var tmp jsonProgInstruction
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	inst.Info = ProgInstInfoTab[tmp.Tag]
	inst.Param = tmp.Param
	return nil
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
