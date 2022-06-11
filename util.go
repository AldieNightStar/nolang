package nolang

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AldieNightStar/golexem"
)

func readFile(name string) string {
	f, err := os.Open(name)

	if err != nil {
		fmt.Println("File read ERR: ", err)
		return ""
	}
	dat, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("File read ERR: ", err)
		return ""
	}
	return string(dat)
}

func NumberToFloat(n any) (float64, bool) {
	if _int, ok := n.(int); ok {
		return float64(_int), true
	} else if _int32, ok := n.(int32); ok {
		return float64(_int32), true
	} else if _int64, ok := n.(int64); ok {
		return float64(_int64), true
	} else if _uint, ok := n.(uint); ok {
		return float64(_uint), true
	} else if _uint32, ok := n.(uint32); ok {
		return float64(_uint32), true
	} else if _uint64, ok := n.(uint64); ok {
		return float64(_uint64), true
	} else if _int16, ok := n.(int16); ok {
		return float64(_int16), true
	} else if _uint16, ok := n.(uint16); ok {
		return float64(_uint16), true
	} else if _int8, ok := n.(int8); ok {
		return float64(_int8), true
	} else if _uint8, ok := n.(uint8); ok {
		return float64(_uint8), true
	} else if _float32, ok := n.(float32); ok {
		return float64(_float32), true
	} else if _float64, ok := n.(float64); ok {
		return float64(_float64), true
	}
	return 0, false
}

func NextAndGetName(s *Scope) string {
	t := s.Next()
	if t == nil {
		return ""
	}
	etc, ok := t.(golexem.ETC)
	if !ok {
		return ""
	}
	return etc.Value
}

func StepAndCastInt(s *Scope) (val int, err error) {
	f, err := StepAndCast[float64](s, 0)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

func StepAndCast[T any](s *Scope, def T) (val T, err error) {
	v, err := s.Step()
	if err != nil {
		return def, err
	}
	t, ok := v.(T)
	if !ok {
		return def, newError(ErrWrongType, s.LastLine)
	}
	return t, nil
}

func StepAndCast2[A any, B any](s *Scope, def1 A, def2 B) (a A, b B, err error) {
	val1, err := StepAndCast(s, def1)
	if err != nil {
		return def1, def2, err
	}
	val2, err := StepAndCast(s, def2)
	if err != nil {
		return def1, def2, err
	}
	return val1, val2, nil
}

func StepAndCast3[A any, B any, C any](s *Scope, def1 A, def2 B, def3 C) (a A, b B, c C, err error) {
	val1, val2, err := StepAndCast2(s, def1, def2)
	if err != nil {
		return def1, def2, def3, err
	}
	val3, err := StepAndCast(s, def3)
	if err != nil {
		return def1, def2, def3, err
	}
	return val1, val2, val3, nil
}

func StepAndCast4[A any, B any, C any, D any](s *Scope, def1 A, def2 B, def3 C, def4 D) (a A, b B, c C, d D, err error) {
	val1, val2, err := StepAndCast2(s, def1, def2)
	if err != nil {
		return def1, def2, def3, def4, err
	}
	val3, val4, err := StepAndCast2(s, def3, def4)
	if err != nil {
		return def1, def2, def3, def4, err
	}
	return val1, val2, val3, val4, nil
}

func ValueGetter(f func() (any, error)) NoFunc {
	return func(s *Scope) (any, error) {
		return f()
	}
}

func ValueSetter(f func(any) error) NoFunc {
	return func(s *Scope) (any, error) {
		v, err := s.Step()
		if err != nil {
			return nil, err
		}
		return nil, f(v)
	}
}

func NoFunc1[A any](def A, f func(A) (any, error)) NoFunc {
	return func(s *Scope) (any, error) {
		v, err := StepAndCast(s, def)
		if err != nil {
			return nil, err
		}
		return f(v)
	}
}

func NoFunc2[A any, B any](def1 A, def2 B, f func(A, B) (any, error)) NoFunc {
	return func(s *Scope) (any, error) {
		v1, v2, err := StepAndCast2(s, def1, def2)
		if err != nil {
			return nil, err
		}
		return f(v1, v2)
	}
}

func NoFunc3[A any, B any, C any](def1 A, def2 B, def3 C, f func(A, B, C) (any, error)) NoFunc {
	return func(s *Scope) (any, error) {
		v1, v2, v3, err := StepAndCast3(s, def1, def2, def3)
		if err != nil {
			return nil, err
		}
		return f(v1, v2, v3)
	}
}

func NoFunc4[A any, B any, C any, D any](def1 A, def2 B, def3 C, def4 D, f func(A, B, C, D) (any, error)) NoFunc {
	return func(s *Scope) (any, error) {
		v1, v2, v3, v4, err := StepAndCast4(s, def1, def2, def3, def4)
		if err != nil {
			return nil, err
		}
		return f(v1, v2, v3, v4)
	}
}
