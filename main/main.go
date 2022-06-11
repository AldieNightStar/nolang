package main

import (
	"fmt"

	"github.com/AldieNightStar/nolang"
)

func main() {
	scope := nolang.LoadFile("sample.txt").LoadDefaultLib()
	err := scope.Run()
	if err != nil {
		fmt.Println(err)
	}
}
