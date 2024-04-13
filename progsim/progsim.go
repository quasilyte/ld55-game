package progsim

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/battle"
	"github.com/quasilyte/ld55-game/game"
)

type Executor struct {
	prog   *runningProg
	world  *battle.World
	vessel *battle.Vessel

	rand *gmath.Rand

	commands VesselCommands

	delta float64

	tick int
}

type ExecutorConfig struct {
	Prog   *game.BotProg
	World  *battle.World
	Vessel *battle.Vessel
}

type VesselCommands struct {
	MoveForward bool
	RotateLeft  bool
	RotateRight bool
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
	t             game.ProgThread
	currentBranch int
	stack         *threadValueStack
	branches      []*runningBranch
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
		threads: make([]*runningThread, len(config.Prog.Threads)),
	}
	for i, t := range config.Prog.Threads {
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
	}
	return e
}

func (e *Executor) RunTick(delta float64) VesselCommands {
	e.commands = VesselCommands{}
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
		t.stack.Clear()

		s := e.runBranch(t, b)
		if s == branchFinished {
			// This will cause the current branch to
			// become 0.
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
		r := inst.Params[0].(float64)
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

	case game.CenterPosInstruction:
		t.stack.Push(stackValue{
			value: e.world.Size.Mulf(0.5),
			tag:   "center pos",
		})

	case game.ChanceInstruction:
		chance := inst.Params[0].(float64)
		if !e.rand.Chance(chance) {
			return instCancelled
		}

	case game.MoveAndRotateInstruction:
		if inst.firstTick {
			p := t.stack.PopVec()
			inst.rotation = e.vessel.Pos.AngleToPoint(p).Normalized()
			inst.amount = inst.Params[0].(float64)
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
			inst.amount = inst.Params[0].(float64)
		}
		if inst.amount > 0 {
			inst.amount -= e.delta * e.vessel.Design.MaxSpeed
			e.commands.MoveForward = true
			return instRunning
		}
	}

	return instFinished
}
