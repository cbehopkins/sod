package sod

import (
	"errors"
	"log"
)

type Group struct {
	crdList []Coord
	pz      *Puzzle
}
func (gr Group) String () string {
  ret_str := "["
  spacer := ""
  for _,itm := range gr.crdList {
    ret_str += spacer + itm.String()
    spacer = " "
  }
  ret_str+= "]"
  return ret_str
}
func (gr Group) Items() []Coord {
	return gr.crdList
}
func (gr Group) Len() int {
	return len(gr.crdList)
}
func NewGroup(sz int, pz *Puzzle) *Group {
	itm := new(Group)
	itm.crdList = make([]Coord, 0, sz)
	itm.pz = pz
	return itm
}
func (gr Group) GetCel(crd Coord) *Cell {
	return gr.pz.GetCel(crd)
}
func (grp Group) valCheck(value Value) bool {
	// Check if a value appears in all the supplied cells
	for _, cellCrd := range grp.Items() {
		cel := grp.GetCel(cellCrd)
		if cel.FindVal(value) == cel.Len() {
			return false
		}
	}
	return true
}
func (gr Group) SelfCheck() error {
	// Many rules associated with what can be in a group.
	// let's self check to make sure we don't violate them.
	pz := gr.pz
	if pz == nil {
		log.Fatal("Uninital,ized puzzle")
	}
	pz_len := pz.Len()

	// First off for a solved cell, a number should only appear once
	set_map := make(map[Value]struct{})
	for i := 0; i < pz_len; i++ {
		set_map[Value(i+1)] = struct{}{}
	}
	for _, crd := range gr.Items() {
		cl := pz.GetCel(crd)
		if cl.Len() == 1 {
			val := cl.Val[0]
			_, ok := set_map[val]
			if ok {
				delete(set_map, val)
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
				_, ok := set_map[val]
				if !ok {
				  // check it's not one that's marked as solved
					log.Printf("Range:%v\nValue:%v\n", gr, val)
          log.Fatal (pz)
					return errors.New("Set missing self value")
				}
			}
		}
	}

	return nil
}
// rmChain takes a list of chains as a job list
// Each chain is one job
// Each chain is a list of links (fake cells)
// We operate on cells not in the list
// Removing the values in the link from those cells not in the list
func (gr Group) RmChain(result_ch []Chain) {
	pz := gr.pz
	// Each chain contains a list of Links(cells)
	for _, chain := range result_ch {
		// That share a pair of numbers.
		// The pair in each of these can be removed from every other cell in the group
		values_to_remove := make(map[Value]struct{})
		coord_to_ignore := make([]Coord, 0, len(chain))

		var link Cell
		for _, link = range chain {
			coord_to_ignore = append(coord_to_ignore, link.crd)
			//log.Println("Sacred Coord",link.coord)
			for _, val := range link.Values() {
				v := val
				//log.Println("Removable Val",v)
				values_to_remove[v] = struct{}{}
			}
		}
		rm_list := make([]Value, 0, len(values_to_remove))
		rmFunc := func(crd Coord) bool {
			for _, cti := range coord_to_ignore {
				if cti.Eq(crd) {
					return true
				}
			}
			for va := range values_to_remove {
				if pz.ValExist(va, crd) {
					rm_list = append(rm_list, va)
				}
			}
			pz.RemoveVals(rm_list, crd)
			//log.Printf("Removing value %v, from crd %v\n", va,crd)
			//pz.RemoveVals(va, crd)
			// Revert the slice to no items, full capacity
			rm_list = rm_list[0:0]
			return true
		}

		gr.ExAll(rmFunc)
	}
}
// This will remove the values in the links
// from the coordinates the link point to
// Does it as a batch operation for efficiency
func (gr Group) RmLinks (result_ch []Chain) {
pz := gr.pz
for _, chain := range result_ch {
 var link Cell
 for _, link = range chain {
  crd := link.crd
  vals := link.Val
  pz.RemoveVals(vals,crd)
 }
}
}
// For a group
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
func (gr Group) ExAll(todo func(Coord) bool) {
	for _, crd := range gr.Items() {
		// Abort processing on false
		if !todo(crd) {
			return
		}
	}
}

func (gr *Group) Add(co Coord) {
	gr.crdList = append(gr.crdList, co)
}
