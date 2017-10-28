package sod

import (
	"errors"
	"log"
)

// Group coords together into managable lumps
type Group struct {
	crdList []Coord
	pz      *Puzzle
}

func (gr Group) String() string {
	retStr := "["
	spacer := ""
	for _, itm := range gr.crdList {
		retStr += spacer + itm.String()
		spacer = " "
	}
	retStr += "]"
	return retStr
}

// Items - the items (list of coords) contained in the group
func (gr Group) Items() []Coord {
	return gr.crdList
}

// Len of the group
func (gr Group) Len() int {
	return len(gr.crdList)
}

// NewGroup return a new group
func NewGroup(sz int, pz *Puzzle) *Group {
	itm := new(Group)
	itm.crdList = make([]Coord, 0, sz)
	itm.pz = pz
	return itm
}

// GetCel from a coord
func (gr Group) GetCel(crd Coord) *Cell {
	return gr.pz.GetCel(crd)
}
func (gr Group) valCheck(value Value) bool {
	// Check if a value appears in all the supplied cells
	for _, cellCrd := range gr.Items() {
		cel := gr.GetCel(cellCrd)
		if cel.FindVal(value) == cel.Len() {
			return false
		}
	}
	return true
}

// SelfCheck a group for consistency
func (gr Group) SelfCheck() error {
	// Many rules associated with what can be in a group.
	// let's self check to make sure we don't violate them.
	pz := gr.pz
	if pz == nil {
		log.Fatal("Uninital,ized puzzle")
	}
	pzLen := pz.Len()

	// First off for a solved cell, a number should only appear once
	setMap := make(map[Value]struct{})
	for i := 0; i < pzLen; i++ {
		setMap[Value(i+1)] = struct{}{}
	}
	for _, crd := range gr.Items() {
		cl := pz.GetCel(crd)
		if cl.Len() == 1 {
			val := cl.Val[0]
			_, ok := setMap[val]
			if ok {
				delete(setMap, val)
			} else {
				return errors.New("value already deleted")
			}
		}
	}

	// Now check that for unsolved cells
	// A value in the unsolved cell
	// still exists in the waiting to be solved list
	for _, crd := range gr.Items() {
		cl := pz.GetCel(crd)
		if cl.Len() > 1 {
			// This is a cell that it could be one of many values
			for _, val := range cl.Val {
				// For each of the possible values
				// set map by this point has had solved values removed
				_, ok := setMap[val]
				if !ok {
					// check it's not one that's marked as solved
					log.Printf("Range:%v\nValue:%v\n", gr, val)
					log.Fatal(pz)
					return errors.New("Set missing self value")
				}
			}
		}
	}

	return nil
}

// RmChain takes a list of chains as a job list
// Each chain is one job
// Each chain is a list of links (fake cells)
// We operate on cells not in the list
// Removing the values in the link from those cells not in the list
func (gr Group) RmChain(resultCh []Chain) {
	if len(resultCh) > 0 {
		//log.Println("Removing Chain", result_ch, gr)
	}
	pz := gr.pz
	// Each chain contains a list of Links(cells)
	for _, chain := range resultCh {
		// That share a pair of numbers.
		// The pair in each of these can be removed from every other cell in the group
		valuesToRemove := make(map[Value]struct{})
		coordToIgnore := make([]Coord, 0, len(chain))

		var link Cell
		for _, link = range chain {
			coordToIgnore = append(coordToIgnore, link.crd)
			//log.Println("Sacred Coord",link.coord)
			for _, val := range link.Values() {
				v := val
				//log.Println("Removable Val",v)
				valuesToRemove[v] = struct{}{}
			}
		}
		rmList := make([]Value, 0, len(valuesToRemove))
		rmFunc := func(crd Coord) bool {
			for _, cti := range coordToIgnore {
				if cti.Eq(crd) {
					//log.Println("Skipping cord because it is us")
					return true
				}
			}
			for va := range valuesToRemove {
				if pz.ValExist(va, crd) {
					//log.Printf("Removing value %v, from crd %v\n", va,crd)
					rmList = append(rmList, va)
				}
			}
			pz.RemoveVals(rmList, crd)
			//pz.RemoveVals(va, crd)
			// Revert the slice to no items, full capacity
			rmList = rmList[0:0]
			return true
		}

		gr.ExAll(rmFunc)
	}
}

// RmLinks - This will remove the values in the links
// from the coordinates the link point to
// Does it as a batch operation for efficiency
func (gr Group) RmLinks(resultCh []Chain) {
	pz := gr.pz
	for _, chain := range resultCh {
		var link Cell
		for _, link = range chain {
			crd := link.crd
			vals := link.Val
			pz.RemoveVals(vals, crd)
		}
	}
}

// ExOthers - For a group
// Go through all the coords in the group, except the one supplied
// and run the function
// On that coordinate
func (gr Group) ExOthers(co Coord, todo func(Coord) bool) {
	lFunc := func(crd Coord) bool {
		if !co.Eq(crd) {
			if !todo(crd) {
				return false
			}
		}
		// Run on all values do not early abort
		return true
	}

	gr.ExAll(lFunc)
}

// ExAll for each coord in the group run the specified function
func (gr Group) ExAll(todo func(Coord) bool) {
	for _, crd := range gr.Items() {
		// Abort processing on false
		if !todo(crd) {
			return
		}
	}
}

// Add a coord to the group
func (gr *Group) Add(co Coord) {
	gr.crdList = append(gr.crdList, co)
}
