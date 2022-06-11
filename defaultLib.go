package nolang

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AldieNightStar/golexem"
)

var stopValue = &struct{}{}

func defaultLib(s *Scope) {
	defaultValues(s)
	defaultStrLib(s)
	defaultMemLib(s)
	defaultMathLib(s)
	defaultBoolLib(s)
	defaultJumpLibrary(s)
	defaultTimeLibrary(s)
	defaultComparingLibrary(s)
	defaultStackLib(s)
	defaultArrayLib(s)
}

func defaultValues(s *Scope) {
	s.Mem["nil"] = nil
	s.Mem["!!"] = stopValue
}

func defaultArrayLib(s *Scope) {
	s.Mem["arr-new"] = NoFunc(func(s *Scope) (any, error) {
		return make([]any, 0, 32), nil
	})
	s.Mem["arr-add"] = NoFunc(func(s *Scope) (any, error) {
		arr, err := StepAndCast[[]any](s, nil)
		if err != nil {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
		return arr, nil
	})
	s.Mem["arr-get"] = NoFunc(func(s *Scope) (any, error) {
		arr, err := StepAndCast[[]any](s, nil)
		if err != nil {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		idF, err := StepAndCast[float64](s, 0)
		if err != nil {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		id := int(idF)
		if id < 0 || id > len(arr) {
			return nil, nil
		}
		return arr[id], nil
	})
	s.Mem["arr-len"] = NoFunc(func(s *Scope) (any, error) {
		arr, err := StepAndCast[[]any](s, nil)
		if err != nil {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		return float64(len(arr)), nil
	})
	s.Mem["arr-all"] = NoFunc(func(s *Scope) (any, error) {
		arr, err := StepAndCast[[]any](s, nil)
		if err != nil {
			return nil, err
		}
		for {
			v, err := s.Step()
			if err != nil {
				return nil, err
			}
			if v == stopValue {
				break
			}
			arr = append(arr, v)
		}
		return arr, nil
	})
}

func defaultStackLib(s *Scope) {
	s.Mem["stack-new"] = NoFunc(func(s *Scope) (any, error) {
		return NewStack[any](1024), nil
	})
	s.Mem["stack-push"] = NoFunc(func(s *Scope) (any, error) {
		stack, err := StepAndCast[*Stack[any]](s, nil)
		if err != nil {
			return nil, err
		}
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		stack.Push(v)
		return nil, nil
	})
	s.Mem["stack-pop"] = NoFunc(func(s *Scope) (any, error) {
		stack, err := StepAndCast[*Stack[any]](s, nil)
		if err != nil {
			return nil, err
		}
		v, _ := stack.Pop(nil)
		return v, nil
	})
}

func defaultMemLib(s *Scope) {
	s.Mem["set"] = NoFunc(func(s *Scope) (any, error) {
		nameToken := s.Next()
		nameEtc, etcOk := nameToken.(golexem.ETC)
		if !etcOk {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		name := nameEtc.Value
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
	s.Mem["num"] = NoFunc(func(s *Scope) (any, error) {
		str, err := StepAndCast(s, "")
		if err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	})
	s.Mem["str"] = NoFunc(func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		return fmt.Sprint(v), nil
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
	s.Mem["!if"] = NoFunc(func(s *Scope) (any, error) {
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
		if b {
			s.Pos = pos1
		} else {
			s.Pos = pos2
		}
		return nil, nil
	})
	s.Mem["!if-else"] = NoFunc(func(s *Scope) (any, error) {
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
	s.Mem["loop"] = NoFunc(func(s *Scope) (any, error) {
		name := NextAndGetName(s)
		if name == "" {
			return nil, newError(ErrWrongType, s.LastLine)
		}
		posF, err := StepAndCast[float64](s, 0)
		if err != nil {
			return nil, err
		}
		pos := int(posF)
		val, ok := s.Mem[name]
		if !ok {
			return nil, nil
		}
		num, ok := val.(float64)
		if !ok {
			return nil, nil
		}
		num -= 1
		if num > 0 {
			s.Pos = pos
		}
		s.Mem[name] = num
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
