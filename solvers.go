package sod

import (
	"log"
	"math/rand"
	"sync"
)

// SolveAll run all algos and solve
func (pz Puzzle) SolveAll() error {

	pz.Solve()
	res := pz.Solved()

	if res == ErrZeroCell {
		return res
	}
	if res == nil {
		log.Println("Solved with simple solver")
		return res
	}

	log.Println("Trying complex solver")
	// This keeps experimenting until solved - very slow!
	err := pz.TrySolver()
	if err != nil {
		return err
	}
	res = pz.Solved()

	if res == ErrZeroCell {
		return res
	}
	if res == nil {
		log.Println("Solved with complex solver")
		return res
	}

	log.Println("Not solvable")
	return nil

}

// Solve with the quick solvers
func (pz Puzzle) Solve() {
	run := true
	for run {
		run = false
		for _, candidate := range pz.Coords() {
			rt := pz.GroupEliminate(candidate)
			run = run || rt
		}
		rt := pz.LoneItems()
		run = run || rt
	}
	pz.TwinSolver()
}

// masterCoord represents the master coord function
// for a supplied coordinate
// go through the values in the cell and return true if
// the possible isn't in there
func (pz Puzzle) masterCoord(crd Coord, currentPossible Value) bool {
	cell := pz.GetCel(crd)
	values := cell.Values()
	for _, v := range values {
		//log.Println("Checking", v, currentPossible, crd)
		if v == currentPossible {
			return false
		}
	}
	return true
}

// GroupEliminate for a group eliminate what we can
// For a given candidate look at each value in the candidate
// Then see if this value appears in any of the group cells
// If it does not then that is the only cell it can appear in
// we can therefore run a set on that candiddate cell
// to that value
func (pz Puzzle) GroupEliminate(candidate Coord) (valueModified bool) {
	//log.Println("GroupEliminate called on", candidate)
	// Candidate must not exist in group
	// To make other code easier, don't error, just move on
	if pz.GetCel(candidate).Len() <= 1 {
		//log.Println("Nothing to do")
		return false
	}

	allGroups := pz.allGroupSets()

	for _, currentPossible := range pz.GetCel(candidate).Values() {
		//log.Println("Checking if value is in another cell", currentPossible)
		var badValue Value
		// This is the group itterator
		// It's called once per coord as we go through the group
		coordFunc := func(crd Coord) bool {
			if !pz.masterCoord(crd, currentPossible) {
				// If this value appears in another cell
				// Then there is no point doing anything else
				// In this group
				//log.Println("Matched", v, currentPossible, crd)
				badValue = currentPossible
				return false
			}
			return true
		}
		// This is the master funciton that is called once per group
		mastFunc := func(gr Group) bool {
			// run the coordFunc on every cell except the candidate
			gr.ExOthers(candidate, coordFunc)

			if badValue == 0 {
				// We haven't found a bad value
				pz.SetVal(currentPossible, candidate)
				pz.IncreaseDifficuly()
				valueModified = true
				//log.Println("Setting value, because it does not appear elsewhere in group", currentPossible, candidate, badValue)
				//log.Println(gr)
			}
			return true
		}
		// run mastFunc on ever cell in every group
		allGroups.ExAll(mastFunc)
	}
	return
}

type mapItems map[Value]Coord

func newMapItemsInit(puzSize int) mapItems {
	checkMap := make(mapItems)
	checkMap.init(puzSize)
	return checkMap
}
func (mi mapItems) init(puzSize int) {
	for i := 0; i < puzSize; i++ {
		mi[Value(i)] = nil
	}
}
func (mi mapItems) checkValues(crd Coord, values []Value) {
	for _, v := range values {
		tmp, ok := mi[v]
		if ok {
			if tmp == nil {
				// This is the first find
				mi[v] = crd
			} else {
				// This is the second time we've found the value
				delete(mi, v)
			}
		}
	}
}

// LoneItems work on a single group
// Look through a group and see if a number only appears in one cell of the group.
// If so then that must be the value for that Cell
func (pz Puzzle) LoneItems() (valueModified bool) {
	puzSize := pz.Len()
	checkMap := newMapItemsInit(puzSize)

	// If the value existsin only 1 co-ord then we've found it
	coordFunc := func(crd Coord) bool {
		checkMap.checkValues(crd, pz.GetCel(crd).Values())
		return true
	}
	mastFunc := func(gr Group) bool {
		checkMap.init(puzSize)
		gr.ExAll(coordFunc)
		for value, coordinate := range checkMap {
			if coordinate != nil && (pz.GetCel(coordinate).Len() > 1) {
				//log.Printf("For group %v\n%v only occurs in %v\n", gr, value, coordinate)
				pz.SetVal(value, coordinate)
				pz.IncreaseDifficuly()
				valueModified = true
			}
		}
		return true
	}
	pz.allGroupSets().ExAll(mastFunc)
	return
}

// Needs a pointer receiver so that NewCrdCnt can build a list to the sources
func (pz *Puzzle) cntMapFunc(crd Coord, cell *Cell, values []Value, cntMap intMap) {

	// Populating a map from count of items in cells to list
	valLen := len(values)
	itm, ok := cntMap[valLen]
	if !ok {
		itm = NewCrdCnt(crd, pz)
	} else {
		itm.Add(crd)
	}
	cntMap.Add(valLen, itm)
}

// Needs a pointer receiver so that NewCrdCnt can build a list to the sources
func (pz *Puzzle) valMapFunc(crd Coord, cell *Cell, values []Value, valMap valueMap) {

	// Populating a map from value to to cells that occurs in
	for _, value := range values {
		itm, ok := valMap[value]
		if !ok {
			itm = NewCrdCnt(crd, pz)
		} else {
			itm.Add(crd)
		}
		valMap.Add(value, itm)
	}
}
func (pz Puzzle) buildMaps(crd Coord, cntMap intMap, valMap valueMap) bool {
	cell := pz.GetCel(crd)
	values := cell.Values()
	pz.cntMapFunc(crd, cell, values, cntMap)
	pz.valMapFunc(crd, cell, values, valMap)
	return true
}

// For one co-ord Count , make a chain of loops
func (pz Puzzle) mkChain(input *CrdCnt) (resultCh []Chain) {
	chain := make(Chain, 0)
	for _, crd := range input.LocList.Items() {
		cel := pz.GetCel(crd)
		nl := NewCellValues(cel.Values())
		nl.crd = crd
		chain = append(chain, *nl)
	}
	resultCh = chain.SearchChain()
	return resultCh
}
func (pz Puzzle) cntExamine(cm intMap) (resultCh, delCh []Chain) {
	// Look at the cells that have 2 items in them
	// Is there another cell that has the same paiting
	// So called Naked Twins
	resultCh = make([]Chain, 0)
	delCh = make([]Chain, 0)

	val, ok := cm[2]
	if ok {
		//log.Println("These are two entry cells", val)
		resultCh = pz.mkChain(val)
		//if len(result_ch)>0 {log.Println("Rx things to delete", result_ch)}
	}
	return
}
func (pz Puzzle) valGrind(value Value, cells Group, vm valueMap) bool {
	var modified bool
	// Grind through each of the values in a cell (that's not the one we're given)
	// looking for another value that appears in one of the cells this appears in
	for _, crd := range cells.Items() {
		modified = modified || pz.GetCel(crd).grindCrd(value, cells, vm)
	}
	return modified
}

func (pz Puzzle) valExamine(vm valueMap) bool {
	var modified bool
	// If a value appears in only 2 cells
	// together with another value that only appears in the same 2 cells
	// It's a hidden twin
	// Therefore any other values in those cells can be removed
	for value, crdCount := range vm {
		if crdCount.Cnt == 1 {
			//log.Println("The value appears twice", value)
			tmp := pz.valGrind(value, *crdCount.LocList, vm)
			modified = tmp || modified
		}
	}
	return modified
}
func (pz Puzzle) twinWorkGroup(gr Group) bool {
	var modified bool
	// This is confusing, so let's be totally clear
	// This will be a map from count of items in cells to list of coords
	// i.e. 1 -> Cells that have one ite, 2-> cells that have 2 items
	var cntMap intMap
	// This is a map from value to cells that that value occurs in
	// i.e 1 -> Cells that contain the value 1
	var valMap valueMap

	//log.Printf("Examining Group %v\n", pz.StringGroup(gr))
	// Clear the maps before each group
	cntMap = newIntMap()
	valMap = newValueMap()
	// Populate the maps
	// We only want to parse through the maps once as it is heavy duty
	lCoordFunc := func(crd Coord) bool {
		return pz.buildMaps(crd, cntMap, valMap)
	}
	gr.ExAll(lCoordFunc)
	// Now examine the maps
	rmChain, _ := pz.cntExamine(cntMap)
	if len(rmChain) > 0 {
		gr.RmChain(rmChain)
		modified = true
	}

	// return true if we modified anything
	tmpMod := pz.valExamine(valMap)
	return modified || tmpMod

}

// TwinSolver solve by looking for twins
// For all the twin processing stuff we want:
// * The coordinates of cells with 2(,3,4) etc items in them
// * A count of how many times each number appears (and which cells this is)
func (pz *Puzzle) TwinSolver() (modified bool) {
	mastFunc := func(gr Group) bool {
		tmpMod := pz.twinWorkGroup(gr)
		modified = modified || tmpMod
		return true
	}
	pz.ExAllGroups(mastFunc)
	return modified
}

// TrySolver try solving by trying things
func (pz Puzzle) TrySolver() error {
	var maxSolvers runCount
	maxSolvers.cnt = 10
	return pz.experimentalSolver(&maxSolvers)
}
func (pz Puzzle) experimentalSolver(rc *runCount) error {

	res := pz.SelfCheck()
	if res != nil {
		return res
	}
	res = pz.Solved()

	if res != ErrUnsolved {
		return res
	}
	pz.Solve()
	if res != ErrUnsolved {
		return res
	}

	err := rc.grab()
	if err != nil {
		//log.Println(err)
		return err
	}
	// Let's be thorough about this
	cellFunc := func(crd Coord) bool {
		cell := pz.GetCel(crd)
		// Find a cell that has multiple possibilities
		numNeeded := cell.Len()
		if numNeeded > 1 {

			solutionsGrid := make([]Puzzle, numNeeded)
			for i := range solutionsGrid {
				// Take a copy of the src puzzle

				solutionsGrid[i] = pz.Duplicate()
			}
			goodSolutionsFound := 0
			solutionFound := 0
			// For each possible value
			for i, possibleValue := range cell.Val {
				// Set this cell in that copy
				solutionsGrid[i].SetVal(possibleValue, crd)

				// So we may have a solution here
				// Recursively solve as needed
				// Call ourselves as we check if solved before running
				err := solutionsGrid[i].experimentalSolver(rc)
				if err == nil {
					// This means we have a good solution
					goodSolutionsFound++
					solutionFound = i
				} else if err == errMaxSolves {
					// We cannot know if we have a good solution
					return false
				} else {
					// Something went wrong
					// This is not a possible solution
					solutionsGrid[i] = Puzzle{}
				}
			}
			if goodSolutionsFound == 0 {
				// keep looking for more solutions
				return true
			} else if goodSolutionsFound == 1 {
				// If after this thereis only one valid solution
				// We have a winner
				// Find it
				// replace the source with the solution found
				solutionsGrid[solutionFound].Copy(pz)
				pz.AddDifficulty(len(cell.Val))
				// Stop the cell itterator
				return false
			}
			// otherwise more than one possible solution
			// TBD Compare the solutions
			// If in all solutions one of the cells retains the same values
			// Then we can safely set that in the source puzzle
			// Otherwise

		}
		// Keep searching for more options
		return true
	}

	pz.ExAll(cellFunc)

	res = pz.SelfCheck()
	if res != nil {
		return res
	}
	res = pz.Solved()
	if res == nil {
		return nil
	}
	return res
}

// MaxDifficultyRand tries to create the maximum difficulty of puzzle possible
// Currently does so randomly
func (pz Puzzle) MaxDifficultyRand() (dst Puzzle) {

	puzLen := pz.Len()
	dst = *NewPuzzle()
	rndSrc := make(chan int)
	rndSrcCloser := make(chan struct{})
	go func() {
		r := rand.New(rand.NewSource(99))
		for {
			select {
			case _, ok := <-rndSrcCloser:
				if !ok {
					close(rndSrc)
					return
				}

			default:
				rndValue := r.Intn(puzLen)
				rndSrc <- rndValue
			}
		}

	}()
	var Solved bool
	for !Solved {
		// Pick a random solved cell from the source puzzle
		col := <-rndSrc
		row := <-rndSrc
		coord := Coord{col, row}
		srcCell := pz.GetCel(coord)
		if srcCell.Len() != 1 {
			continue
		}

		dstCel := dst.GetCel(coord)
		// No point working on an already solved cell
		if dstCel.Len() == 1 {
			continue
		}
		// Set the value
		dst.SetVal(srcCell.Val[0], coord)
		// Test solving it on a copy
		dstPoss := dst.Duplicate()
		log.Println("Trying with puzzle", dstPoss)
		if dstPoss.TrySolver() == nil {
			// This is a solvable puzzle
			Solved = true
		}

	}
	close(rndSrcCloser)
	<-rndSrc
	return
}

// MaxDifficulty returns a puzzle that should be the maximum possible difficulty
func (pz Puzzle) MaxDifficulty() (dst Puzzle) {

	puzLen := pz.Len()
	puzzlesToRun := make([]*Puzzle, 0)
	possibleReductions := make(map[int]*Puzzle)

	for col := 0; col < puzLen; col++ {
		for row := 0; row < puzLen; row++ {
			coord := Coord{col, row}
			// Pick a solved cell from the source puzzle
			if pz.GetCel(coord).Len() == 1 {
				possPuz := NewPuzzle()

				//log.Println("This is a possible removal", coord)
				copySolved := func(crd Coord) bool {
					if !crd.Eq(coord) {
						srcCell := pz.GetCel(crd)
						if srcCell.Len() == 1 {
							possPuz.SetVal(srcCell.Val[0], crd)
						}
					}
					return true
				}

				*possPuz.difficulty = *pz.difficulty
				possPuz.ExAll(copySolved)

				puzzlesToRun = append(puzzlesToRun, possPuz)
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(puzzlesToRun))
	resChan := make(chan *Puzzle)

	for _, possPuz := range puzzlesToRun {
		go func(tpz *Puzzle) {
			// Is this puzzle solvable?
			solvable := tpz.TrySolver()
			if solvable != nil {
				resChan <- tpz
			}
			wg.Done()
		}(possPuz)
	}
	var wgMap sync.WaitGroup
	go func() {
		wgMap.Add(1)
		for possPuz := range resChan {
			possibleReductions[possPuz.Difficulty()] = possPuz
		}
		wgMap.Done()
	}()
	wg.Wait()
	close(resChan)
	wgMap.Wait()

	maxDifficultyFound := 0
	for diffi := range possibleReductions {
		if diffi > maxDifficultyFound {
			maxDifficultyFound = diffi
		}
	}
	log.Println("Max difficulty found:", maxDifficultyFound)
	puz, ok := possibleReductions[maxDifficultyFound]
	if !ok {
		log.Println("No reduction found")
		dst = pz
	} else {
		if puz == nil {
			log.Fatal("Nill puzzle in reduction")
		}
		dst = puz.Duplicate()
	}

	return
}
