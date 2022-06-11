package nolang

import (
	"fmt"
	"time"
)

func defaultLib(s *Scope) {
	defaultStrLib(s)
	defaultMemLib(s)
	defaultMathLib(s)
	defaultBoolLib(s)
	defaultJumpLibrary(s)
	defaultTimeLibrary(s)
	defaultComparingLibrary(s)
}

func defaultMemLib(s *Scope) {
	s.Mem["set"] = NoFunc(func(s *Scope) (any, error) {
		val, err := s.Step()
		if err != nil {
			return nil, err
		}
		name, ok := val.(string)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		newVal, err := s.Step()
		if err != nil {
			return nil, err
		}
		s.Mem[name] = newVal
		return nil, nil
	})
}

func defaultStrLib(s *Scope) {
	s.Mem["print"] = NoFunc(func(s *Scope) (any, error) {
		str, err := s.Step()
		if err != nil {
			return nil, err
		}
		fmt.Println(str)
		return nil, nil
	})
	s.Mem["concat"] = NoFunc(func(s *Scope) (any, error) {
		a1, err := s.Step()
		if err != nil {
			return nil, err
		}
		a2, err := s.Step()
		if err != nil {
			return nil, err
		}
		s1, ok := a1.(string)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		s2, ok := a2.(string)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		return s1 + s2, nil
	})
}

func defaultMathLib(s *Scope) {
	s.Mem["add"] = defaultLibMathOpFunc(func(a, b float64) float64 { return a + b })
	s.Mem["sub"] = defaultLibMathOpFunc(func(a, b float64) float64 { return a - b })
	s.Mem["mul"] = defaultLibMathOpFunc(func(a, b float64) float64 { return a * b })
	s.Mem["div"] = defaultLibMathOpFunc(func(a, b float64) float64 { return a / b })
	s.Mem["mod"] = defaultLibMathOpFunc(func(a, b float64) float64 { return float64(int(a) % int(b)) })
}

func defaultBoolLib(s *Scope) {
	s.Mem["and"] = defaultLibBoolOpFunc(func(a, b bool) bool { return a && b })
	s.Mem["or"] = defaultLibBoolOpFunc(func(a, b bool) bool { return a || b })
	s.Mem["not"] = NoFunc(func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		b, ok := v.(bool)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		return !b, nil
	})
}

func defaultLibBoolOpFunc(f func(a, b bool) bool) NoFunc {
	return NoFunc(func(s *Scope) (any, error) {
		v1, err := s.Step()
		if err != nil {
			return nil, err
		}
		b1, ok := v1.(bool)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		v2, err := s.Step()
		if err != nil {
			return nil, err
		}
		b2, ok := v2.(bool)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		return f(b1, b2), nil
	})
}

func defaultLibMathOpFunc(f func(a, b float64) float64) NoFunc {
	return NoFunc(func(s *Scope) (any, error) {
		o1, err := s.Step()
		if err != nil {
			return nil, err
		}
		o2, err := s.Step()
		if err != nil {
			return nil, err
		}
		n1, ok := NumberToFloat(o1)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		n2, ok := NumberToFloat(o2)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		return f(n1, n2), nil
	})
}

func defaultJumpLibrary(s *Scope) {
	s.Mem["goto"] = NoFunc(func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		n, ok := v.(float64)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		s.Pos = int(n)
		return nil, nil
	})
	s.Mem["ret"] = NoFunc(func(s *Scope) (any, error) {
		pos, ok := s.CallStack.Pop(0)
		if !ok {
			return nil, newError(ErrCodeEnd, s.LastLine)
		}
		s.Pos = pos
		return nil, nil
	})
	s.Mem["call"] = NoFunc(func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		n, ok := v.(float64)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		newPos := int(n)
		s.CallStack.Push(s.Pos)
		s.Pos = newPos
		return nil, nil
	})
	s.Mem["if"] = NoFunc(func(s *Scope) (any, error) {
		b, err := StepAndCast(s, false)
		if err != nil {
			return nil, err
		}
		f, err := StepAndCast[float64](s, 0)
		if err != nil {
			return nil, err
		}
		pos := int(f)
		if b {
			s.CallStack.Push(s.Pos)
			s.Pos = pos
		}
		return nil, nil
	})
	s.Mem["if-else"] = NoFunc(func(s *Scope) (any, error) {
		b, err := StepAndCast(s, false)
		if err != nil {
			return nil, err
		}
		f1, err := StepAndCast[float64](s, 0)
		if err != nil {
			return nil, err
		}
		f2, err := StepAndCast[float64](s, 0)
		if err != nil {
			return nil, err
		}
		pos1 := int(f1)
		pos2 := int(f2)
		s.CallStack.Push(s.Pos)
		if b {
			s.Pos = pos1
		} else {
			s.Pos = pos2
		}
		return nil, nil
	})
}

func defaultTimeLibrary(s *Scope) {
	s.Mem["sleep"] = NoFunc(func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		n, ok := v.(float64)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		time.Sleep(time.Duration(float64(time.Second) * n))
		return nil, nil
	})
}

func defaultComparingLibrary(s *Scope) {
	s.Mem["<"] = defaultComparingLibraryNumbersFunc(func(a, b float64) bool { return a < b })
	s.Mem[">"] = defaultComparingLibraryNumbersFunc(func(a, b float64) bool { return a > b })
	s.Mem["<="] = defaultComparingLibraryNumbersFunc(func(a, b float64) bool { return a <= b })
	s.Mem[">="] = defaultComparingLibraryNumbersFunc(func(a, b float64) bool { return a >= b })
	s.Mem["=="] = defaultComparingLibraryNumbersFunc(func(a, b float64) bool { return a == b })
}

func defaultComparingLibraryNumbersFunc(f func(a, b float64) bool) NoFunc {
	return func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		a, ok := NumberToFloat(v)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		v, err = s.Step()
		if err != nil {
			return nil, err
		}
		b, ok := NumberToFloat(v)
		if !ok {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		return f(a, b), nil
	}
}
