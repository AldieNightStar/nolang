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
