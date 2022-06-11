package main

import (
	"fmt"
	"os"

	"github.com/AldieNightStar/nolang"
)

func main() {
	args := getArgs(false)
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

func getArgs(debug bool) []string {
	if debug {
		return []string{"sample.r"}
	} else {
		return os.Args[1:]
	}
}
