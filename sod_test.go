package sod

import (
	"log"
	"testing"
)

func TestBas0(t *testing.T) {
	testPuzzle := NewPuzzle()
	result := testPuzzle.SelfCheck()
	if result != nil {
		log.Fatal("Self check fail", result)
	}
}

// See what happens with
//http://www.sudokuessentials.com/support-files/sudoku-very-hard-1.pdf
func TestSimplify(t *testing.T) {

	testPuzzle := NewPuzzle()
	sc := func() {
		result := testPuzzle.SelfCheck()
		if result != nil {
			log.Fatal("Self check fail", result)
		}
	}
	sc()

	// Set Value,Col,Row
	// numbered from 0
	testPuzzle.setValC(3, 1, 0)
	testPuzzle.setValC(4, 3, 0)
	testPuzzle.setValC(8, 4, 0)
	testPuzzle.setValC(6, 6, 0)
	testPuzzle.setValC(9, 8, 0)
	testPuzzle.setValC(2, 4, 1)
	testPuzzle.setValC(7, 5, 1)
	testPuzzle.setValC(8, 0, 2)
	testPuzzle.setValC(3, 3, 2)
	testPuzzle.setValC(1, 1, 3)
	testPuzzle.setValC(9, 2, 3)
	testPuzzle.setValC(7, 0, 4)
	testPuzzle.setValC(8, 1, 4)
	testPuzzle.setValC(2, 5, 4)
	testPuzzle.setValC(9, 7, 4)
	testPuzzle.setValC(3, 8, 4)
	testPuzzle.setValC(4, 5, 5)
	testPuzzle.setValC(8, 6, 5)
	testPuzzle.setValC(7, 7, 5)
	testPuzzle.setValC(5, 5, 6)
	testPuzzle.setValC(6, 8, 6)
	testPuzzle.setValC(1, 3, 7)
	testPuzzle.setValC(3, 4, 7)
	testPuzzle.setValC(9, 0, 8)
	testPuzzle.setValC(2, 2, 8)
	testPuzzle.setValC(4, 4, 8)
	testPuzzle.setValC(8, 5, 8)
	testPuzzle.setValC(1, 7, 8)
	//log.Println("Row 7\n",testPuzzle)
	sc()

	if true {
		testPuzzle.Solve()
		sc()
	}
	if true {
		// This is a minimal puzzle so:
		// Give is an extra solved location
		testPuzzle.setValC(2, 0, 0)
		difficultPuz := testPuzzle.MaxDifficulty()
		//log.Println(difficultPuz)
		if difficultPuz.roughCheck(*testPuzzle) {
			// all is good
			// every solved cell in testPuzzle
			// exists as a value in difficultPuz
		} else {
			log.Fatal("Missing Cells")
		}

		// if MaxDifficulty has worked
		// there will be some cells marked as solved in testPuzzle
		// that are not solved in difficultPuz
		if testPuzzle.lessRoughCheck(difficultPuz) {
			// If this is true
			// then every solved difficultPuz cell is solved in testPuzzle
			log.Fatal("MaxDifficulty has done nothing")
		}
	}
}
func samplePuzzle0() (testPuzzle *Puzzle, expectedResult [][]Value) {
	expectedResult = [][]Value{
		{2, 3, 7, 4, 8, 1, 6, 5, 9},
		{6, 9, 4, 5, 2, 7, 1, 3, 8},
		{8, 5, 1, 3, 6, 9, 2, 4, 7},
		{4, 1, 9, 8, 7, 3, 5, 6, 2},
		{7, 8, 5, 6, 1, 2, 4, 9, 3},
		{3, 2, 6, 9, 5, 4, 8, 7, 1},
		{1, 4, 3, 2, 9, 5, 7, 8, 6},
		{5, 7, 8, 1, 3, 6, 9, 2, 4},
		{9, 6, 2, 7, 4, 8, 3, 1, 5},
	}
	//log.Println(expectedResult)
	testPuzzle = NewPuzzle()
	// Set Value,Col,Row
	// numbered from 0
	testPuzzle.setValC(3, 1, 0)
	testPuzzle.setValC(4, 3, 0)
	testPuzzle.setValC(8, 4, 0)
	testPuzzle.setValC(6, 6, 0)
	testPuzzle.setValC(9, 8, 0)
	testPuzzle.setValC(2, 4, 1)
	testPuzzle.setValC(7, 5, 1)
	testPuzzle.setValC(8, 0, 2)
	testPuzzle.setValC(3, 3, 2)
	testPuzzle.setValC(1, 1, 3)
	testPuzzle.setValC(9, 2, 3)
	testPuzzle.setValC(7, 0, 4)
	testPuzzle.setValC(8, 1, 4)
	testPuzzle.setValC(2, 5, 4)
	testPuzzle.setValC(9, 7, 4)
	testPuzzle.setValC(3, 8, 4)
	testPuzzle.setValC(4, 5, 5)
	testPuzzle.setValC(8, 6, 5)
	testPuzzle.setValC(7, 7, 5)
	testPuzzle.setValC(5, 5, 6)
	testPuzzle.setValC(6, 8, 6)
	testPuzzle.setValC(1, 3, 7)
	testPuzzle.setValC(3, 4, 7)
	testPuzzle.setValC(9, 0, 8)
	testPuzzle.setValC(2, 2, 8)
	testPuzzle.setValC(4, 4, 8)
	testPuzzle.setValC(8, 5, 8)
	testPuzzle.setValC(1, 7, 8)
	return testPuzzle, expectedResult
}
func samplePuzzle1() (testPuzzle *Puzzle, expectedResult [][]Value) {
	expectedResult = [][]Value{
		{2, 3, 7, 4, 8, 1, 6, 5, 9},
		{6, 9, 4, 5, 2, 7, 1, 3, 8},
		{8, 5, 1, 3, 6, 9, 2, 4, 7},
		{4, 1, 9, 8, 7, 3, 5, 6, 2},
		{7, 8, 5, 6, 1, 2, 4, 9, 3},
		{3, 2, 6, 9, 5, 4, 8, 7, 1},
		{1, 4, 3, 2, 9, 5, 7, 8, 6},
		{5, 7, 8, 1, 3, 6, 9, 2, 4},
		{9, 6, 2, 7, 4, 8, 3, 1, 5},
	}
	//log.Println(expectedResult)
	testPuzzle = NewPuzzle()
	// Let's set up some twins for examination
	testPuzzle.setValC(5, 3, 0)
	testPuzzle.setValC(3, 5, 0)
	testPuzzle.setValC(9, 7, 0)
	// row 1
	testPuzzle.setValC(6, 4, 1)
	testPuzzle.setValC(7, 5, 1)
	testPuzzle.setValC(1, 6, 1)
	testPuzzle.setValC(5, 7, 1)
	// row 2
	testPuzzle.setValC(5, 1, 2)
	testPuzzle.setValC(4, 2, 2)
	testPuzzle.setValC(9, 3, 2)
	testPuzzle.setValC(2, 4, 2)
	testPuzzle.setValC(1, 5, 2)
	testPuzzle.setValC(3, 7, 2)
	// row 3
	testPuzzle.setValC(4, 1, 3)
	testPuzzle.setValC(9, 2, 3)
	testPuzzle.setValC(3, 3, 3)
	testPuzzle.setValC(7, 4, 3)
	testPuzzle.setValC(2, 6, 3)
	// row 4
	testPuzzle.setValC(1, 0, 4)
	testPuzzle.setValC(3, 1, 4)
	testPuzzle.setValC(7, 7, 4)
	testPuzzle.setValC(9, 8, 4)
	// row 5
	testPuzzle.setValC(7, 1, 5)
	testPuzzle.setValC(5, 2, 5)
	testPuzzle.setValC(1, 4, 5)
	testPuzzle.setValC(9, 5, 5)
	testPuzzle.setValC(4, 6, 5)
	testPuzzle.setValC(3, 8, 5)
	// row 6
	testPuzzle.setValC(6, 3, 6)
	testPuzzle.setValC(5, 4, 6)
	testPuzzle.setValC(4, 5, 6)
	testPuzzle.setValC(3, 6, 6)
	testPuzzle.setValC(2, 7, 6)
	// row 7
	testPuzzle.setValC(6, 1, 7)
	testPuzzle.setValC(7, 3, 7)
	testPuzzle.setValC(3, 4, 7)
	testPuzzle.setValC(2, 5, 7)
	testPuzzle.setValC(9, 6, 7)
	// row 8
	testPuzzle.setValC(2, 1, 8)
	testPuzzle.setValC(1, 3, 8)
	testPuzzle.setValC(9, 4, 8)
	testPuzzle.setValC(8, 5, 8)
	return testPuzzle, expectedResult
}
func samplePuzzle2() (testPuzzle *Puzzle, expectedResult [][]Value) {
	expectedResult = [][]Value{}
	//log.Println(expectedResult)
	testPuzzle = NewPuzzle()

	// Let's set up some twins for examination
	testPuzzle.setValC(6, 0, 0)
	testPuzzle.setValC(2, 4, 0)
	testPuzzle.setValC(9, 8, 0)
	// row 1
	testPuzzle.setValC(1, 1, 1)
	testPuzzle.setValC(3, 3, 1)
	testPuzzle.setValC(7, 5, 1)
	testPuzzle.setValC(5, 7, 1)
	// row 2
	testPuzzle.setValC(3, 2, 2)
	testPuzzle.setValC(1, 6, 2)
	// row 3
	testPuzzle.setValC(9, 1, 3)
	testPuzzle.setValC(2, 7, 3)

	// row 4
	testPuzzle.setValC(2, 0, 4)
	testPuzzle.setValC(8, 3, 4)
	testPuzzle.setValC(7, 4, 4)
	testPuzzle.setValC(5, 5, 4)
	testPuzzle.setValC(3, 8, 4)
	// row 5
	testPuzzle.setValC(5, 2, 5)
	testPuzzle.setValC(1, 4, 5)
	testPuzzle.setValC(4, 6, 5)

	// row 6
	testPuzzle.setValC(7, 1, 6)
	testPuzzle.setValC(8, 4, 6)
	testPuzzle.setValC(9, 7, 6)

	// row 7
	testPuzzle.setValC(1, 2, 7)
	testPuzzle.setValC(4, 4, 7)
	testPuzzle.setValC(8, 6, 7)

	// row 8
	testPuzzle.setValC(2, 3, 8)
	testPuzzle.setValC(5, 4, 8)
	testPuzzle.setValC(9, 5, 8)
	return testPuzzle, expectedResult
}
func TestSolver(t *testing.T) {

	testPuzzle, expectedResult := samplePuzzle0()
	sc := func() {
		result := testPuzzle.SelfCheck()
		if result != nil {
			log.Fatal("Self check fail", result)
		}
	}
	sc()

	if true {
		testPuzzle.Solve()
		sc()
	}
	if true {
		testPuzzle.TrySolver()
		sc()
		if !testPuzzle.Check(expectedResult) {
			log.Fatal("Unexpected Result", testPuzzle)
		}
	}
}

func TestDuplicate(t *testing.T) {
	testPuzzle, expectedResult := samplePuzzle0()
	sc := func() {
		result := testPuzzle.SelfCheck()
		if result != nil {
			log.Fatal("Self check fail", result)
		}
	}
	sc()

	testPuzzle.Solve()
	sc()
	duplicatePuz := testPuzzle.Duplicate()
	//log.Println("Duplicate Before", duplicatePuz)
	clearFunc := func(crd Coord) bool {
		cel := duplicatePuz.GetCel(crd)
		//log.Println("Clearing",cel,crd)
		cel.SetVal(1)
		return true
	}
	duplicatePuz.ExAll(clearFunc)
	// TBD check duplicate is all 1s
	pzSize := testPuzzle.Len()
	expectedDuplicate := make([][]Value, pzSize)
	for i := 0; i < pzSize; i++ {
		tmpArr := make([]Value, pzSize)
		for j := 0; j < pzSize; j++ {
			tmpArr[j] = 1
		}
		expectedDuplicate[i] = tmpArr
	}
	if !duplicatePuz.Check(expectedDuplicate) {
		log.Fatal("Unexpected Result should be all 1s", duplicatePuz)
	}

	testPuzzle.TrySolver()
	sc()
	// This will fail if duplicate deleted main puzzle
	testPuzzle.Check(expectedResult)
}

func TestTwin(t *testing.T) {

	testPuzzle, _ := samplePuzzle1()
	result := testPuzzle.SelfCheck()
	if result != nil {
		log.Fatal("Self check fail", result)
	}

	//log.Println(testPuzzle)
	// ar this stage
	// 0,7 should be [4,5,8]
	// 0,8 [3,4,5,7]
	// after they should be just [4,5]
	value0 := testPuzzle.GetCel(Coord{0, 7}).Values()
	value1 := testPuzzle.GetCel(Coord{0, 8}).Values()
	if len(value0) != 3 {
		log.Fatal("expected", value0)
	}
	if len(value1) != 4 {
		log.Fatal("expected", value1)
	}
	testPuzzle.TwinSolver()
	value0 = testPuzzle.GetCel(Coord{0, 7}).Values()
	value1 = testPuzzle.GetCel(Coord{0, 8}).Values()
	if len(value0) != 2 {
		log.Fatal("expected", value0)
	}
	if len(value1) != 2 {
		log.Fatal("expected", value1)
	}
	//log.Println(testPuzzle)
}
