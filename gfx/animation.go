package gfx

import (
	"math"

	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gsignal"
)

type AnimationMode int

const (
	AnimationForward AnimationMode = iota
	AnimationBackward
)

type Animation struct {
	sprite     *graphics.Sprite
	frameWidth int

	numFrames int
	offsetY   float64

	animationSpan float64
	deltaPerFrame float64

	repeated bool
	Mode     AnimationMode

	frame       int
	frameTicker float64

	EventFrameChanged gsignal.Event[int]
}

func NewRepeatedAnimation(s *graphics.Sprite, numFrames int) *Animation {
	a := NewAnimation(s, numFrames)
	a.repeated = true
	return a
}

func NewAnimation(s *graphics.Sprite, numFrames int) *Animation {
	a := &Animation{}
	a.SetSprite(s, numFrames)
	a.SetSecondsPerFrame(0.05)
	return a
}

func (a *Animation) SetSprite(s *graphics.Sprite, numFrames int) {
	a.sprite = s
	frameWidth, _ := s.GetFrameSize()
	if numFrames < 0 {
		w, _ := s.ImageSize()
		numFrames = int(float64(w) / float64(frameWidth))
	}
	a.frameWidth = frameWidth
	a.numFrames = numFrames
	a.SetSecondsPerFrame(a.deltaPerFrame)
}

func (a *Animation) SetOffsetY(offset float64) {
	a.offsetY = offset
}

func (a *Animation) SetAnimationSpan(value float64) {
	a.animationSpan = value
	a.deltaPerFrame = value / float64(a.numFrames)
}

func (a *Animation) SetSecondsPerFrame(seconds float64) {
	a.animationSpan = seconds * float64(a.numFrames)
	a.deltaPerFrame = seconds
}

func (a *Animation) Sprite() *graphics.Sprite {
	return a.sprite
}

func (a *Animation) IsDisposed() bool {
	return a.sprite.IsDisposed()
}

func (a *Animation) Rewind() {
	a.frameTicker = 0
	a.frame = 0
	var x int
	if a.Mode != AnimationForward {
		x = a.frameWidth * (a.numFrames - 1)
	}
	a.sprite.SetFrameOffset(x, int(a.offsetY))
}

func (a *Animation) RewindTo(value float64) {
	a.frameTicker = 0
	a.frame = -1
	a.Tick(value)
}

func (a *Animation) Tick(delta float64) bool {
	if !a.repeated {
		if a.frameTicker >= a.animationSpan {
			return true
		}
	}

	finished := false
	a.frameTicker += delta
	var frame int
	if a.frameTicker >= a.animationSpan {
		finished = true
		if a.repeated {
			rem := math.Mod(a.frameTicker, a.animationSpan)
			a.frameTicker = rem
			frame = int(a.frameTicker / a.deltaPerFrame)
		} else {
			a.frameTicker = a.animationSpan
			frame = a.numFrames - 1
		}
	} else {
		frame = int(a.frameTicker / a.deltaPerFrame)
	}

	if a.Mode == AnimationBackward {
		frame = (a.numFrames - 1) - frame
	}

	framesDelta := frame - a.frame
	a.frame = frame
	if framesDelta != 0 {
		// A small optimization: don't call Emit if there are no listeners.
		// This is more useful for repeated animations as they're less likely to have
		// any frame event listeners.
		if !a.EventFrameChanged.IsEmpty() {
			a.EventFrameChanged.Emit(framesDelta)
		}
		a.sprite.SetFrameOffset(a.frameWidth*frame, int(a.offsetY))
	}

	return finished
}
