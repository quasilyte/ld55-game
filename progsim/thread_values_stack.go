package progsim

import (
	"fmt"

	"github.com/quasilyte/gmath"
)

type threadValueStack struct {
	values []stackValue
}

type stackValue struct {
	value any
	tag   string
}

func (s *threadValueStack) Push(v stackValue) {
	s.values = append(s.values, v)
}

func (s *threadValueStack) PopVec() gmath.Vec {
	rv := s.Pop()
	v, ok := rv.value.(gmath.Vec)
	if !ok {
		panic(fmt.Sprintf("PopVec: expected Vec, found %T (%s)", rv, rv.tag))
	}
	return v
}

func (s *threadValueStack) Pop() stackValue {
	v := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return v
}

func (s *threadValueStack) Clear() {
	s.values = s.values[:0]
}
