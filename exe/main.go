package main

import (
	"fmt"
	"github.com/cbehopkins/sod"
	"log"
)

func main() {
	testPuzzle := sod.NewPuzzle()
	fmt.Println("Blank Puzzle")
	fmt.Println(testPuzzle)
	tc := sod.Coord{1, 2}
	testPuzzle.SetVal(3, tc)
	tc = sod.Coord{1, 3}
	testPuzzle.SetVal(4, tc)
	tc = sod.Coord{5, 7}
	testPuzzle.SetVal(6, tc)
	fmt.Println(testPuzzle)
	result := testPuzzle.SelfCheck()
	if result != nil {
		log.Fatal("Self check fail", result)
	}
}
