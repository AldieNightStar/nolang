package nolang

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/AldieNightStar/golexem"
)

func filterComments(toks []any) []any {
	res := make([]any, 0, 32)
	for _, tok := range toks {
		if _, ok := tok.(golexem.COMMENT); ok {
			continue
		}
		res = append(res, tok)
	}
	return res
}

func processLabelDefs(toks []any) ([]any, map[string]int) {
	labs := make(map[string]int, 64)
	arr := make([]any, 0, 32)
	pos := 0
	for _, tok := range toks {
		if etc, ok := tok.(golexem.ETC); ok {
			if strings.HasPrefix(etc.Value, ":") {
				labs[etc.Value[1:]] = pos
				continue
			}
		}
		arr = append(arr, tok)
		pos += 1
	}
	return arr, labs
}

func processLabelPointers(labels map[string]int, toks []any) []any {
	arr := make([]any, 0, 32)
	for _, tok := range toks {
		if etc, ok := tok.(golexem.ETC); ok {
			if strings.HasPrefix(etc.Value, "@") {
				labNum, labFound := labels[etc.Value[1:]]
				if !labFound {
					wrongToken := golexem.ETC(golexem.NewToken(etc.Value[1:], 0))
					wrongToken.LineNumber = etc.LineNumber
					arr = append(arr, wrongToken)
					continue
				}
				numToken := golexem.NUMBER(golexem.NewToken("", float64(labNum)))
				numToken.LineNumber = etc.LineNumber
				arr = append(arr, numToken)
				continue
			}
		}
		arr = append(arr, tok)
	}
	return arr
}

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
