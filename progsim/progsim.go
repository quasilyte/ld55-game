package progsim

import (
	"fmt"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/game"
)

type Executor struct {
	prog   *runningProg
	world  *game.World
	vessel *game.Vessel

	rand *gmath.Rand

	commands VesselCommands

	delta float64

	tick int
}

type ExecutorConfig struct {
	Prog   *game.BotProg
	World  *game.World
	Vessel *game.Vessel
}

type VesselCommands struct {
	MoveForward bool
	RotateLeft  bool
	RotateRight bool

	FireCommands []VesselFireCommand
}

func (c *VesselCommands) Reset() {
	c.MoveForward = false
	c.RotateLeft = false
	c.RotateRight = false

	// Re-use memory for the fire commands.
	c.FireCommands = c.FireCommands[:0]
}

type VesselFireCommand struct {
	WeaponIndex      uint
	ReloadMultiplier float64
	TargetPos        gmath.Vec
}

type branchStatus int

const (
	branchRunning branchStatus = iota
	branchFinished
	branchCancelled
)

type instStatus int

const (
	instRunning instStatus = iota
	instFinished
	instCancelled
)

type runningProg struct {
	src     *game.BotProg
	threads []*runningThread
}

type runningThread struct {
	t             *game.ProgThread
	currentBranch int
	stack         *threadValueStack
	branches      []*runningBranch
}

func (t *runningThread) weaponIndex() int {
	switch t.t.Kind {
	case game.Weapon1Thread:
		return 0
	case game.Weapon2Thread:
		return 1
	default:
		return -1
	}
}

type runningBranch struct {
	currentInst int
	insts       []*runningInst
}

type runningInst struct {
	firstTick bool
	rotation  gmath.Rad
	amount    float64
	game.ProgInstruction
}

func NewExecutor(config ExecutorConfig) *Executor {
	e := &Executor{
		world:  config.World,
		vessel: config.Vessel,
		rand:   &config.World.Rand,
	}
	e.prog = &runningProg{
		src:     config.Prog,
		threads: make([]*runningThread, config.Prog.NumThreads()),
	}
	config.Prog.EachThread(func(i int, t *game.ProgThread) {
		rt := &runningThread{
			t: t,
			stack: &threadValueStack{
				values: make([]stackValue, 0, 4),
			},
			branches: make([]*runningBranch, len(t.Branches)),
		}
		for j, b := range t.Branches {
			rb := &runningBranch{
				insts: make([]*runningInst, len(b.Instructions)),
			}
			rt.branches[j] = rb
			for k, inst := range b.Instructions {
				rb.insts[k] = &runningInst{
					ProgInstruction: inst,
					firstTick:       true,
				}
			}
		}
		e.prog.threads[i] = rt
	})
	return e
}

func (e *Executor) RunTick(delta float64) VesselCommands {
	e.commands.Reset()

	e.delta = delta

	for _, t := range e.prog.threads {
		e.runThread(t)
	}

	e.tick++
	return e.commands
}

func (e *Executor) runThread(t *runningThread) {
	for i := t.currentBranch; i < len(t.branches); i++ {
		b := t.branches[i]
		if len(t.stack.values) >= 12 {
			fmt.Printf("warning: thread stack len is %d\n", len(t.stack.values))
		}

		s := e.runBranch(t, b)
		if s == branchFinished {
			// This will cause the current branch to
			// become 0.
			t.stack.Clear()
			break
		}
		if s == branchRunning {
			return
		}
		t.currentBranch++
	}

	t.currentBranch = 0
}

func (e *Executor) runBranch(t *runningThread, b *runningBranch) branchStatus {
	for i := b.currentInst; i < len(b.insts); i++ {
		inst := b.insts[i]
		if inst.firstTick {
			b.currentInst = i
		}
		s := e.runInst(t, inst)
		if inst.firstTick {
			inst.firstTick = false
		}

		// Not using a switch to avoid having to
		// place an explicit label above the for loop.
		if s == instFinished {
			inst.firstTick = true
			continue // Can run another instruction
		}
		if s == instCancelled {
			b.currentInst = 0
			return branchCancelled
		}
		if s == instRunning {
			return branchRunning
		}
	}

	b.currentInst = 0
	return branchFinished
}

func (e *Executor) runInst(t *runningThread, inst *runningInst) instStatus {
	switch inst.Info.Kind {
	case game.NopInstruction:
		return instFinished

	case game.HealthPercentInstruction:
		percent := e.vessel.Health / e.vessel.Design.MaxHealth
		t.stack.Push(stackValue{
			value: 100 * percent,
			tag:   "health percent",
		})

	case game.EnergyPercentInstruction:
		percent := e.vessel.Energy / e.vessel.Design.MaxEnergy
		t.stack.Push(stackValue{
			value: 100 * percent,
			tag:   "energy percent",
		})

	case game.RandomPosInstruction:
		p := e.world.Size.Sub(gmath.Vec{
			X: e.rand.FloatRange(0, e.world.Size.X),
			Y: e.rand.FloatRange(0, e.world.Size.Y),
		})
		t.stack.Push(stackValue{
			value: p,
			tag:   "random pos",
		})

	case game.RandomOffsetInstruction:
		p := t.stack.PopVec()
		r := inst.Param
		newPos := p.Add(e.rand.Offset(-r, r))
		t.stack.Push(stackValue{
			value: newPos,
			tag:   "random offset",
		})

	case game.VesselPosInstruction:
		t.stack.Push(stackValue{
			value: e.vessel.Pos,
			tag:   "vessel pos",
		})

	case game.TargetPosInstruction:
		t.stack.Push(stackValue{
			value: e.vessel.Target.Pos,
			tag:   "target pos",
		})

	case game.TargetSpeedInstruction:
		t.stack.Push(stackValue{
			value: e.vessel.Target.Velocity().Len(),
			tag:   "target speed",
		})

	case game.CenterPosInstruction:
		t.stack.Push(stackValue{
			value: e.world.Size.Mulf(0.5),
			tag:   "center pos",
		})

	case game.ChanceInstruction:
		chance := inst.Param
		if !e.rand.Chance(chance * 0.01) {
			return instCancelled
		}

	case game.IsLtInstruction:
		x := inst.Param
		y := t.stack.PopFloat()
		if y >= x {
			return instCancelled
		}

	case game.IsGtInstruction:
		x := inst.Param
		y := t.stack.PopFloat()
		if y <= x {
			return instCancelled
		}

	case game.IsOutBoundsInstruction:
		p := t.stack.PopVec()
		rect := gmath.Rect{Max: e.world.Size}
		if rect.Contains(p) {
			return instCancelled
		}

	case game.DistanceToInstruction:
		p := t.stack.PopVec()
		dist := p.DistanceTo(e.vessel.Pos)
		t.stack.Push(stackValue{
			value: dist,
			tag:   "distance to",
		})

	case game.WaitInstruction:
		if inst.firstTick {
			inst.amount = inst.Param / 60.0
		}
		inst.amount -= e.delta
		if inst.amount > 0 {
			return instRunning
		}

	case game.MoveAndRotateInstruction:
		if inst.firstTick {
			p := t.stack.PopVec()
			inst.rotation = e.vessel.Pos.AngleToPoint(p).Normalized()
			inst.amount = inst.Param
		}
		vesselRotation := e.vessel.Rotation.Normalized()
		delta := vesselRotation.AngleDelta(inst.rotation)
		absDelta := delta.Abs()
		if absDelta >= float64(e.vessel.Design.RotationSpeed)*e.delta {
			if delta < 0 {
				e.commands.RotateLeft = true
			} else {
				e.commands.RotateRight = true
			}
		}
		if inst.amount > 0 {
			inst.amount -= e.delta * e.vessel.Design.MaxSpeed
			e.commands.MoveForward = true
			return instRunning
		}

	case game.RotateToInstruction:
		if inst.firstTick {
			p := t.stack.PopVec()
			inst.rotation = e.vessel.Pos.AngleToPoint(p).Normalized()
		}
		vesselRotation := e.vessel.Rotation.Normalized()
		delta := vesselRotation.AngleDelta(inst.rotation)
		absDelta := delta.Abs()
		if absDelta < float64(e.vessel.Design.RotationSpeed)*e.delta {
			return instFinished
		}
		if delta < 0 {
			e.commands.RotateLeft = true
		} else {
			e.commands.RotateRight = true
		}
		return instRunning

	case game.MoveForwardInstruction:
		if inst.firstTick {
			inst.amount = inst.Param
		}
		if inst.amount > 0 {
			inst.amount -= e.delta * e.vessel.Design.MaxSpeed
			e.commands.MoveForward = true
			return instRunning
		}

	case game.SnapShotInstruction:
		// TODO: check whether this weapon can be used?
		weaponIndex := t.weaponIndex()
		e.commands.FireCommands = append(e.commands.FireCommands, VesselFireCommand{
			ReloadMultiplier: 0.8,
			WeaponIndex:      uint(weaponIndex),
			TargetPos:        e.vessel.Target.Pos.Add(e.rand.Offset(-32, 32)),
		})

	case game.NormalShotInstruction:
		weaponIndex := t.weaponIndex()
		e.commands.FireCommands = append(e.commands.FireCommands, VesselFireCommand{
			ReloadMultiplier: 1.0,
			WeaponIndex:      uint(weaponIndex),
			TargetPos:        e.vessel.Target.Pos,
		})

	case game.AimShotInstruction:
		weaponIndex := t.weaponIndex()
		w := e.vessel.Weapons[weaponIndex]
		targetVelocity := e.vessel.Target.Velocity().Mulf(e.rand.FloatRange(0.9, 1.1))
		targetPos := snipePos(w.Design.ProjectileSpeed, e.vessel.Pos, e.vessel.Target.Pos, targetVelocity)
		e.commands.FireCommands = append(e.commands.FireCommands, VesselFireCommand{
			ReloadMultiplier: 2.0,
			WeaponIndex:      uint(weaponIndex),
			TargetPos:        targetPos,
		})
	}

	return instFinished
}
