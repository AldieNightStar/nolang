package nolang

import (
	"github.com/AldieNightStar/golexem"
)

func Load(code string) *Scope {
	toks := golexem.Parse(code)
	toks = filterComments(toks)
	return NewScope(toks)
}

func LoadFile(filename string) *Scope {
	return Load(ReadFile(filename))
}
