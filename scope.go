package nolang

import (
	"fmt"

	"github.com/AldieNightStar/golexem"
)

type Scope struct {
	Pos       int
	Code      []any
	Mem       map[string]any
	LastLine  int
	CallStack *Stack[int]
}

func (s *Scope) Next() any {
	if s.Pos >= len(s.Code) {
		return nil
	}
	res := s.Code[s.Pos]
	s.Pos += 1
	return res
}

func (s *Scope) Step() (any, error) {
	op := s.Next()
	if op == nil {
		return nil, newError(ErrCodeEnd, s.LastLine)
	}
	if str, ok := op.(golexem.STRING); ok {
		s.LastLine = str.LineNumber
		return str.Value, nil
	} else if num, ok := op.(golexem.NUMBER); ok {
		s.LastLine = num.LineNumber
		return num.ValueNumber, nil
	} else if etc, ok := op.(golexem.ETC); ok {
		name := etc.Value
		s.LastLine = etc.LineNumber
		if name == "true" {
			return true, nil
		} else if name == "false" {
			return false, nil
		}
		obj, found := s.Mem[name]
		if !found {
			return nil, fmt.Errorf("unknown '%s' Line: %d", name, etc.LineNumber)
		}
		if fn, isFunc := obj.(NoFunc); isFunc {
			return fn(s)
		}
		return obj, nil
	}
	return op, nil
}

func (s *Scope) Run() error {
	for {
		_, err := s.Step()
		if err != nil {
			if isErrEq(err, ErrCodeEnd) {
				return nil
			} else {
				return err
			}
		}
	}
}

func (s *Scope) RunLocal(pos int) error {
	oldCallStack := s.CallStack
	curPos := s.Pos

	s.Pos = pos
	s.CallStack = NewStack[int](256)

	err := s.Run()

	s.CallStack = oldCallStack
	s.Pos = curPos

	if err != nil && !isCodeEndErr(err) {
		return err
	}
	return nil
}

func NewScope(code []any) *Scope {
	return &Scope{Pos: 0, Code: code, Mem: make(map[string]any, 64), LastLine: 0, CallStack: NewStack[int](1024)}
}
