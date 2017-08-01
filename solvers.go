package sod

import (
	"errors"
	"log"
	"math/rand"
	"sync"
)

func (pz Puzzle) Solve() {
	run := true
	for run {
		run = false
		for _, candidate := range pz.Coords() {
			run = run || pz.GroupEliminate(candidate)
		}
		run = run || pz.LoneItems()
	}
	pz.TwinSolver()
}

// For a given candidate look at each value in the candidate
// Then see if this value appears in any of the group cells
// If it does not then that is the only cell it can appear in
// we can therefore run a set on that candiddate cell
// to that value
func (pz Puzzle) GroupEliminate(candidate Coord) (value_modified bool) {

	// Candidate must not exist in group
	// To make other code easier, don't error, just move on
	if pz.GetCel(candidate).Len() <= 1 {
		return false
	}
	allGroups := pz.AllGroupSets()

	for _, currentPossible := range pz.GetCel(candidate).Values() {
		var badValue Value
		// This is the group itterator
		// It's called once per coord as we go through the group
		coordFunc := func(crd Coord) bool {
			cell := pz.GetCel(crd)
			values := cell.Values()
			for _, v := range values {
				//log.Println("Checking", v, currentPossible, crd)
				if v == currentPossible {
					// If this value appears in another cell
					// Then there is no point doing anything else
					// In this group
					//log.Println("Matched", v, currentPossible, crd)
					badValue = currentPossible
					return false
				}
			}
			return true
		}
		// This is the master funciton that is called once per group
		mastFunc := func(gr Group) bool {
			// if the group func has EVER reported false
			// then it has at some point found a match
			gr.ExOthers(candidate, coordFunc)

			// Therefore if it has set the badValue
			if badValue != 0 {
				// then we have has a false, so found a match somewhere in this group
				// Look for more solutions in more groups
			} else {
				pz.SetVal(currentPossible, candidate)
				pz.IncreaseDifficuly()
				value_modified = true
				//log.Println("Setting value, because it does not appear elsewhere in group", currentPossible, candidate, badValue)
				//log.Println(gr)
			}
			return true
		}

		allGroups.ExAll(mastFunc)
	}
	return
}

// Look through a group and see if a number only appears in one cell of the group.
// If so then that must be the value for that Cell
func (pz Puzzle) LoneItems() (value_modified bool) {
	puz_size := pz.Len()
	var check_map map[Value]Coord

	mapInit := func() {
		check_map = make(map[Value]Coord)
		for i := 0; i < puz_size; i++ {
			check_map[Value(i)] = nil
		}
	}
	// If the value existsin only 1 co-ord then we've found it
	coordFunc := func(crd Coord) bool {
		cell := pz.GetCel(crd)
		values := cell.Values()
		for _, v := range values {
			tmp, ok := check_map[v]
			if ok {
				if tmp == nil {
					// This is the first find
					check_map[v] = crd
				} else {
					// This is the second time we've found the value
					delete(check_map, v)
				}
			}
		}
		return true
	}
	mastFunc := func(gr Group) bool {
		mapInit()
		gr.ExAll(coordFunc)
		for value, coordinate := range check_map {
			if coordinate != nil && (pz.GetCel(coordinate).Len() > 1) {
				//log.Printf("For group %v\n%v only occurs in %v\n", gr, value, coordinate)
				pz.SetVal(value, coordinate)
				pz.IncreaseDifficuly()
				value_modified = true
			}
		}
		return true
	}
	allGroups := pz.AllGroupSets()
	allGroups.ExAll(mastFunc)
	return
}

type CrdCnt struct {
	Cnt     int
	LocList *Group
}

func NewCrdCnt(crd Coord, pz *Puzzle) *CrdCnt {
	itm := new(CrdCnt)
	itm.Cnt = 0
	itm.LocList = NewGroup(1, pz)
	itm.set(crd)
	return itm
}
func (cc *CrdCnt) set(crd Coord) {
	cc.LocList.Add(crd)
}
func (cc *CrdCnt) Add(crd Coord) {
	cc.Cnt++
	cc.set(crd)
}
func (pz *Puzzle) TwinSolver() {
	// For all the twin processing stuff we want:
	// * The coordinates of cells with 2(,3,4) etc items in them
	// * A count of how many times each number appears (and which cells this is)

	// This is confusing, so let's be totally clear
	// This will be a map from count of items in cells to list of coords
	// i.e. 1 -> Cells that have one ite, 2-> cells that have 2 items
	var cntMap map[int]*CrdCnt
	// This is a map from value to cells that that value occurs in
	// i.e 1 -> Cells that contain the value 1
	var valMap map[Value]*CrdCnt

	cntMapFunc := func(crd Coord, cell *Cell, values []Value) {
		// Populating a map from count of items in cells to list
		valLen := len(values)
		itm, ok := cntMap[valLen]
		if !ok {
			itm = NewCrdCnt(crd, pz)
		} else {
			itm.Add(crd)
		}
		cntMap[valLen] = itm
	}

	valMapFunc := func(crd Coord, cell *Cell, values []Value) {
		// Populating a map from value to to cells that occurs in
		for _, value := range values {
			itm, ok := valMap[value]
			if !ok {
				itm = NewCrdCnt(crd, pz)
			} else {
				itm.Add(crd)
			}
			valMap[value] = itm
		}
	}
	cntExamine := func(cm map[int]*CrdCnt) (result_ch , del_ch []Chain) {
		// Look at the cells that have 2 items in them
		// Is there another cell that has the same paiting
		// So called Naked Twins
		result_ch = make([]Chain, 0)
    del_ch = make([]Chain, 0)

		val, ok := cm[2]
		if ok {
			//log.Println("These are two entry cells", val)
			chain := make(Chain, 0)
			for _, crd := range val.LocList.Items() {
				cel := pz.GetCel(crd)
				nl := NewCellValues(cel.Values())
				nl.crd = crd
				chain = append(chain, *nl)
			}
			result_ch = append(result_ch,chain.SearchChain()...)
		}

		return 
    }

	valGrind := func(value Value, cells Group, vm map[Value]*CrdCnt) bool {
		// Grind through each of the values in a cell (that's not the one we're given)
		// looking for another value that appears in one of the cells this appears in
		for _, crd := range cells.Items() {
			cel := pz.GetCel(crd)
			for _, val := range cel.Values() {
				if val != value && vm[val].Cnt == 1 {
					// Don't examine our self
					// Only bother to check values that also only appear in 2 cells
					if cells.valCheck(val) {
						// If the value appears in all the cells we appear in
						// Then we have a match
						// Build a list of values that are not val,value
						rem_vals := cel.NotValues([]Value{val, value})

						if len(rem_vals) > 0 {
							cel.RemoveVals(rem_vals)
							//log.Printf("%v is paired with %v: Remove %v\n", val, value, rem_vals)
						}
					}

				}
			}
		}
		return true
	}
	valExamine := func(vm map[Value]*CrdCnt) {
		// If a value appears in only 2 cells
		// together with another value that only appears in the same 2 cells
		// It's a hidden twin
		// Therefore any other values in those cells can be removed
		for value, crdCount := range vm {
			if crdCount.Cnt == 1 {
				//log.Println("The value appears twice", value)
				valGrind(value, *crdCount.LocList, vm)
			}
		}
	}
	coordFunc := func(crd Coord) bool {
		cell := pz.GetCel(crd)
		values := cell.Values()
		cntMapFunc(crd, cell, values)
		valMapFunc(crd, cell, values)
		return true
	}
	mastFunc := func(gr Group) bool {
		//log.Printf("Examining Group %v\n", pz.StringGroup(gr))
		// Clear the maps before each group
		cntMap = make(map[int]*CrdCnt)
		valMap = make(map[Value]*CrdCnt)
		// Populate the maps
		// We only want to parse through the maps once as it is heavy duty
		gr.ExAll(coordFunc)
		// Now examine the maps
		rm_chain,del_chain := cntExamine(cntMap)
		gr.RmChain(rm_chain)
    gr.RmLinks(del_chain)
		valExamine(valMap)
		return true
	}
	pz.ExAllGroups(mastFunc)
}

type runCount struct {
	sync.Mutex
	cnt int
}

var ErrMaxSolves = errors.New("max solvers have been run")

func (rc *runCount) Grab() error {
	rc.Lock()
	defer rc.Unlock()
	if rc.cnt > 0 {
		rc.cnt--
		return nil
	} else {
		return ErrMaxSolves
	}
}
func (pz Puzzle) TrySolver() error {
	var maxSolvers runCount
	maxSolvers.cnt = 10
	return pz.ExperimentalSolver(&maxSolvers)
}
func (pz Puzzle) ExperimentalSolver(rc *runCount) error {

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

	err := rc.Grab()
	if err != nil {
		//log.Println(err)
		return err
	}
	// Let's be thorough about this
	cellFunc := func(crd Coord) bool {
		cell := pz.GetCel(crd)
		// Find a cell that has multiple possibilities
		num_needed := cell.Len()
		if num_needed > 1 {

			solutions_grid := make([]Puzzle, num_needed)
			for i, _ := range solutions_grid {
				// Take a copy of the src puzzle

				solutions_grid[i] = pz.Duplicate()
			}
			good_solutions_found := 0
			solution_found := 0
			// For each possible value
			for i, possibleValue := range cell.Val {
				// Set this cell in that copy
				solutions_grid[i].SetVal(possibleValue, crd)

				// So we may have a solution here
				// Recursively solve as needed
				// Call ourselves as we check if solved before running
				err := solutions_grid[i].ExperimentalSolver(rc)
				if err == nil {
					// This means we have a good solution
					good_solutions_found++
					solution_found = i
				} else if err == ErrMaxSolves {
					// We cannot know if we have a good solution
					return false
				} else {
					// Something went wrong
					// This is not a possible solution
					solutions_grid[i] = Puzzle{}
				}
			}
			if good_solutions_found == 0 {
				// keep looking for more solutions
				return true
			} else if good_solutions_found == 1 {
				// If after this thereis only one valid solution
				// We have a winner
				// Find it
				// replace the source with the solution found
				solutions_grid[solution_found].Copy(pz)
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
	} else {
		return res
	}
}

// This tries to create the maximum difficulty of puzzle possible
// Currently does so randomly
func (src Puzzle) MaxDifficultyRand() (dst Puzzle) {

	puz_len := src.Len()
	dst = *NewPuzzle()
	rndSrc := make(chan int)
	rndSrcCloser := make(chan struct{})
	go func() {
		r := rand.New(rand.NewSource(99))
		// TBD this routine will leak
		for {
			select {
			case _, ok := <-rndSrcCloser:
				if !ok {
					close(rndSrc)
					return
				}

			default:
				// TBD seed with timestamp
				rnd_value := r.Intn(puz_len)
				rndSrc <- rnd_value
			}
		}

	}()
	var Solved bool
	for !Solved {
		// Pick a random solved cell from the source puzzle
		col := <-rndSrc
		row := <-rndSrc
		coord := Coord{col, row}
		src_cell := src.GetCel(coord)
		if src_cell.Len() != 1 {
			continue
		}

		dst_cel := dst.GetCel(coord)
		// No point working on an already solved cell
		if dst_cel.Len() == 1 {
			continue
		}
		// Set the value
		dst.SetVal(src_cell.Val[0], coord)
		// Test solving it on a copy
		dst_poss := dst.Duplicate()
		log.Println("Trying with puzzle", dst_poss)
		if dst_poss.TrySolver() == nil {
			// This is a solvable puzzle
			Solved = true
		}

	}
	close(rndSrcCloser)
	<-rndSrc
	return
}
func (src Puzzle) MaxDifficulty() (dst Puzzle) {

	puz_len := src.Len()
	puzzlesToRun := make([]*Puzzle, 0)
	possibleReductions := make(map[int]*Puzzle)

	for col := 0; col < puz_len; col++ {
		for row := 0; row < puz_len; row++ {
			coord := Coord{col, row}
			// Pick a solved cell from the source puzzle
			if src.GetCel(coord).Len() == 1 {
				possPuz := NewPuzzle()

				//log.Println("This is a possible removal", coord)
				copySolved := func(crd Coord) bool {
					if !crd.Eq(coord) {
						src_cell := src.GetCel(crd)
						if src_cell.Len() == 1 {
							possPuz.SetVal(src_cell.Val[0], crd)
						}
					}
					return true
				}

				*possPuz.difficulty = *src.difficulty
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
	for diffi, _ := range possibleReductions {
		if diffi > maxDifficultyFound {
			maxDifficultyFound = diffi
		}
	}
	log.Println("Max difficulty found:", maxDifficultyFound)
	puz, ok := possibleReductions[maxDifficultyFound]
	if !ok {
		log.Println("No reduction found")
		dst = src
	} else {
		if puz == nil {
			log.Fatal("Nill puzzle in reduction")
		}
		dst = puz.Duplicate()
	}

	return
}
