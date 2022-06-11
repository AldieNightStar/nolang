package nolang

import (
	"github.com/AldieNightStar/golexem"
)

func Load(code string) *Scope {
	var toks []any
	var labels map[string]int
	toks = golexem.Parse(code)
	toks = filterComments(toks)
	toks, labels = processLabelDefs(toks)
	toks = processLabelPointers(labels, toks)
	s := NewScope(toks)
	defaultLib(s)
	return s
}

func LoadFile(filename string) *Scope {
	return Load(readFile(filename))
}
