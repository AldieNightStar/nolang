package main

import (
	"fmt"
	"os"

	"github.com/AldieNightStar/nolang"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage:\n\tnolang file")
		return
	}
	scope := nolang.LoadFile(args[0])
	err := scope.Run()
	if err != nil {
		fmt.Println(err)
	}
}
