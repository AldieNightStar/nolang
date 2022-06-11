package nolang

import (
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
