package main

import (
	"fmt"

	"github.com/AldieNightStar/nolang"
)

func main() {
	// args := []string{"sample.r"}
	// args := os.Args[1:]
	args := []string{"sample.r"}
	if len(args) < 1 {
		fmt.Println("Usage:\n\tnolang file")
		return
	}
	scope := nolang.LoadFile(args[0])
	scope.Mem["reset"] = nolang.NoFunc(func(s *nolang.Scope) (any, error) {
		s.Pos = 0
		return nil, nil
	})
	err := scope.Run()
	if err != nil {
		fmt.Println(err)
	}
}
