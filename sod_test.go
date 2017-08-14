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
	testPuzzle.SetValC(3, 1, 0)
	testPuzzle.SetValC(4, 3, 0)
	testPuzzle.SetValC(8, 4, 0)
	testPuzzle.SetValC(6, 6, 0)
	testPuzzle.SetValC(9, 8, 0)
	testPuzzle.SetValC(2, 4, 1)
	testPuzzle.SetValC(7, 5, 1)
	testPuzzle.SetValC(8, 0, 2)
	testPuzzle.SetValC(3, 3, 2)
	testPuzzle.SetValC(1, 1, 3)
	testPuzzle.SetValC(9, 2, 3)
	testPuzzle.SetValC(7, 0, 4)
	testPuzzle.SetValC(8, 1, 4)
	testPuzzle.SetValC(2, 5, 4)
	testPuzzle.SetValC(9, 7, 4)
	testPuzzle.SetValC(3, 8, 4)
	testPuzzle.SetValC(4, 5, 5)
	testPuzzle.SetValC(8, 6, 5)
	testPuzzle.SetValC(7, 7, 5)
	testPuzzle.SetValC(5, 5, 6)
	testPuzzle.SetValC(6, 8, 6)
	testPuzzle.SetValC(1, 3, 7)
	testPuzzle.SetValC(3, 4, 7)
	testPuzzle.SetValC(9, 0, 8)
	testPuzzle.SetValC(2, 2, 8)
	testPuzzle.SetValC(4, 4, 8)
	testPuzzle.SetValC(8, 5, 8)
	testPuzzle.SetValC(1, 7, 8)
	//log.Println("Row 7\n",testPuzzle)
	sc()

	if true {
		testPuzzle.Solve()
		sc()
	}
	if true {
		// This is a minimal puzzle so:
		// Give is an extra solved location
		testPuzzle.SetValC(2, 0, 0)
		difficultPuz := testPuzzle.MaxDifficulty()
		//log.Println(difficultPuz)
		if difficultPuz.RoughCheck(*testPuzzle) {
			// all is good
			// every solved cell in testPuzzle
			// exists as a value in difficultPuz
		} else {
			log.Fatal("Missing Cells")
		}

		// if MaxDifficulty has worked
		// there will be some cells marked as solved in testPuzzle
		// that are not solved in difficultPuz
		if testPuzzle.LessRoughCheck(difficultPuz) {
			// If this is true
			// then every solved difficultPuz cell is solved in testPuzzle
			log.Fatal("MaxDifficulty has done nothing")
		}
	}
}
func samplePuzzle0() (testPuzzle *Puzzle, expectedResult [][]Value) {
	expectedResult = [][]Value{
		[]Value{2, 3, 7, 4, 8, 1, 6, 5, 9},
		[]Value{6, 9, 4, 5, 2, 7, 1, 3, 8},
		[]Value{8, 5, 1, 3, 6, 9, 2, 4, 7},
		[]Value{4, 1, 9, 8, 7, 3, 5, 6, 2},
		[]Value{7, 8, 5, 6, 1, 2, 4, 9, 3},
		[]Value{3, 2, 6, 9, 5, 4, 8, 7, 1},
		[]Value{1, 4, 3, 2, 9, 5, 7, 8, 6},
		[]Value{5, 7, 8, 1, 3, 6, 9, 2, 4},
		[]Value{9, 6, 2, 7, 4, 8, 3, 1, 5},
	}
	//log.Println(expectedResult)
	testPuzzle = NewPuzzle()
	// Set Value,Col,Row
	// numbered from 0
	testPuzzle.SetValC(3, 1, 0)
	testPuzzle.SetValC(4, 3, 0)
	testPuzzle.SetValC(8, 4, 0)
	testPuzzle.SetValC(6, 6, 0)
	testPuzzle.SetValC(9, 8, 0)
	testPuzzle.SetValC(2, 4, 1)
	testPuzzle.SetValC(7, 5, 1)
	testPuzzle.SetValC(8, 0, 2)
	testPuzzle.SetValC(3, 3, 2)
	testPuzzle.SetValC(1, 1, 3)
	testPuzzle.SetValC(9, 2, 3)
	testPuzzle.SetValC(7, 0, 4)
	testPuzzle.SetValC(8, 1, 4)
	testPuzzle.SetValC(2, 5, 4)
	testPuzzle.SetValC(9, 7, 4)
	testPuzzle.SetValC(3, 8, 4)
	testPuzzle.SetValC(4, 5, 5)
	testPuzzle.SetValC(8, 6, 5)
	testPuzzle.SetValC(7, 7, 5)
	testPuzzle.SetValC(5, 5, 6)
	testPuzzle.SetValC(6, 8, 6)
	testPuzzle.SetValC(1, 3, 7)
	testPuzzle.SetValC(3, 4, 7)
	testPuzzle.SetValC(9, 0, 8)
	testPuzzle.SetValC(2, 2, 8)
	testPuzzle.SetValC(4, 4, 8)
	testPuzzle.SetValC(8, 5, 8)
	testPuzzle.SetValC(1, 7, 8)
	return testPuzzle, expectedResult
}
func samplePuzzle1() (testPuzzle *Puzzle, expectedResult [][]Value) {
	expectedResult = [][]Value{
		[]Value{2, 3, 7, 4, 8, 1, 6, 5, 9},
		[]Value{6, 9, 4, 5, 2, 7, 1, 3, 8},
		[]Value{8, 5, 1, 3, 6, 9, 2, 4, 7},
		[]Value{4, 1, 9, 8, 7, 3, 5, 6, 2},
		[]Value{7, 8, 5, 6, 1, 2, 4, 9, 3},
		[]Value{3, 2, 6, 9, 5, 4, 8, 7, 1},
		[]Value{1, 4, 3, 2, 9, 5, 7, 8, 6},
		[]Value{5, 7, 8, 1, 3, 6, 9, 2, 4},
		[]Value{9, 6, 2, 7, 4, 8, 3, 1, 5},
	}
	//log.Println(expectedResult)
	testPuzzle = NewPuzzle()
	// Let's set up some twins for examination
	testPuzzle.SetValC(5, 3, 0)
	testPuzzle.SetValC(3, 5, 0)
	testPuzzle.SetValC(9, 7, 0)
	// row 1
	testPuzzle.SetValC(6, 4, 1)
	testPuzzle.SetValC(7, 5, 1)
	testPuzzle.SetValC(1, 6, 1)
	testPuzzle.SetValC(5, 7, 1)
	// row 2
	testPuzzle.SetValC(5, 1, 2)
	testPuzzle.SetValC(4, 2, 2)
	testPuzzle.SetValC(9, 3, 2)
	testPuzzle.SetValC(2, 4, 2)
	testPuzzle.SetValC(1, 5, 2)
	testPuzzle.SetValC(3, 7, 2)
	// row 3
	testPuzzle.SetValC(4, 1, 3)
	testPuzzle.SetValC(9, 2, 3)
	testPuzzle.SetValC(3, 3, 3)
	testPuzzle.SetValC(7, 4, 3)
	testPuzzle.SetValC(2, 6, 3)
	// row 4
	testPuzzle.SetValC(1, 0, 4)
	testPuzzle.SetValC(3, 1, 4)
	testPuzzle.SetValC(7, 7, 4)
	testPuzzle.SetValC(9, 8, 4)
	// row 5
	testPuzzle.SetValC(7, 1, 5)
	testPuzzle.SetValC(5, 2, 5)
	testPuzzle.SetValC(1, 4, 5)
	testPuzzle.SetValC(9, 5, 5)
	testPuzzle.SetValC(4, 6, 5)
	testPuzzle.SetValC(3, 8, 5)
	// row 6
	testPuzzle.SetValC(6, 3, 6)
	testPuzzle.SetValC(5, 4, 6)
	testPuzzle.SetValC(4, 5, 6)
	testPuzzle.SetValC(3, 6, 6)
	testPuzzle.SetValC(2, 7, 6)
	// row 7
	testPuzzle.SetValC(6, 1, 7)
	testPuzzle.SetValC(7, 3, 7)
	testPuzzle.SetValC(3, 4, 7)
	testPuzzle.SetValC(2, 5, 7)
	testPuzzle.SetValC(9, 6, 7)
	// row 8
	testPuzzle.SetValC(2, 1, 8)
	testPuzzle.SetValC(1, 3, 8)
	testPuzzle.SetValC(9, 4, 8)
	testPuzzle.SetValC(8, 5, 8)
	return testPuzzle, expectedResult
}
func samplePuzzle2() (testPuzzle *Puzzle, expectedResult [][]Value) {
	expectedResult = [][]Value{}
	//log.Println(expectedResult)
	testPuzzle = NewPuzzle()

	// Let's set up some twins for examination
	testPuzzle.SetValC(6, 0, 0)
	testPuzzle.SetValC(2, 4, 0)
	testPuzzle.SetValC(9, 8, 0)
	// row 1
	testPuzzle.SetValC(1, 1, 1)
	testPuzzle.SetValC(3, 3, 1)
	testPuzzle.SetValC(7, 5, 1)
	testPuzzle.SetValC(5, 7, 1)
	// row 2
	testPuzzle.SetValC(3, 2, 2)
	testPuzzle.SetValC(1, 6, 2)
	// row 3
	testPuzzle.SetValC(9, 1, 3)
	testPuzzle.SetValC(2, 7, 3)

	// row 4
	testPuzzle.SetValC(2, 0, 4)
	testPuzzle.SetValC(8, 3, 4)
	testPuzzle.SetValC(7, 4, 4)
	testPuzzle.SetValC(5, 5, 4)
	testPuzzle.SetValC(3, 8, 4)
	// row 5
	testPuzzle.SetValC(5, 2, 5)
	testPuzzle.SetValC(1, 4, 5)
	testPuzzle.SetValC(4, 6, 5)

	// row 6
	testPuzzle.SetValC(7, 1, 6)
	testPuzzle.SetValC(8, 4, 6)
	testPuzzle.SetValC(9, 7, 6)

	// row 7
	testPuzzle.SetValC(1, 2, 7)
	testPuzzle.SetValC(4, 4, 7)
	testPuzzle.SetValC(8, 6, 7)

	// row 8
	testPuzzle.SetValC(2, 3, 8)
	testPuzzle.SetValC(5, 4, 8)
	testPuzzle.SetValC(9, 5, 8)
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
	pz_size := testPuzzle.Len()
	expected_duplicate := make([][]Value, pz_size)
	for i := 0; i < pz_size; i++ {
		tmp_arr := make([]Value, pz_size)
		for j := 0; j < pz_size; j++ {
			tmp_arr[j] = 1
		}
		expected_duplicate[i] = tmp_arr
	}
	if !duplicatePuz.Check(expected_duplicate) {
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
