package sod

import (
	"log"
	"testing"
)

func exampleChain() Chain {
	chain := make(Chain, 0)
	return chain
}
func TestChainFul0(t *testing.T) {
	pz, _ := samplePuzzle2()
	pz.Solve()
	// This puzzle has some twins in row 7
	//  [row=7
	//   (0)[5,9]      (1)[2,5]    (2)[1]      (3)[7]    (4)[4]  (5)[3,6]  (6)[8]      (7)[3,6]    (8)[2,5]
	// location 1 and 8 of this row must have either of 2 & 5, therefore 0 cannot have 5 in it.

	log.Println(pz)

	interestingGroup := NewGroup(0, pz)
	for i := 0; i < 9; i++ {
		crd := Coord{i, 7}
		interestingGroup.Add(crd)
	}

	modified := pz.twinWorkGroup(*interestingGroup)
	if modified {
		log.Println("Modified success")

	} else {
		log.Fatal("Unmodified", pz)
	}
}

// as above, but using the general solver
func TestChainFul1(t *testing.T) {
	pz, _ := samplePuzzle2()
	pz.Solve()
	// This puzzle has some twins in row 7
	//  [row=7
	//   (0)[5,9]      (1)[2,5]    (2)[1]      (3)[7]    (4)[4]  (5)[3,6]  (6)[8]      (7)[3,6]    (8)[2,5]
	// location 1 and 8 of this row must have either of 2 & 5, therefore 0 cannot have 5 in it.

	log.Println(pz)

	modified := pz.TwinSolver()
	if modified {
		//pz.TrySolver()
		log.Println("Modified success", pz)

	} else {
		log.Fatal("Unmodified", pz)
	}
}
func TestChainFul2(t *testing.T) {
	pz, _ := samplePuzzle2()
	pz.Solve()
	err := pz.SolveAll()
	if err == nil {
		//pz.TrySolver()
		log.Println("Success", pz)

	} else {
		log.Fatal("Unmodified", pz, err)
	}
}
func TestChain0(t *testing.T) {
	// In this we want to spot the 1->2->3->1 chain
	// Without being tripped up by {1,5}
	testValues := [][]int{
		{1, 2},
		{2, 3},
		{1, 3},
		{4, 5},
		{4, 6},
		{1, 5},
	}
	tv := NewChain(testValues)
	log.Println(tv.SearchChain())
}
func TestChain1(t *testing.T) {
	// In this we want to spot the 1->2->3->1 chain
	// Without being tripped up by {1,5}
	testValues := [][]int{
		{1, 2},
		{2, 3},
		{1, 3},
		{4, 5},
		{4, 6},
		{4, 5},
	}
	tv := NewChain(testValues)
	log.Println(tv.SearchChain())
}
