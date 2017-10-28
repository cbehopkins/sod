package sod

import "log"

// GroupSet is a collection of groups
type GroupSet []Group

// ExAll execute a function on everything in the group
func (gs GroupSet) ExAll(todo func(Group) bool) {
	for _, gr := range gs {
		if gr.pz == nil {
			log.Fatal("Nil puzzle at start of GroupSet ExAll")
		}
		if !todo(gr) {
			return
		}
	}
}

// ExOthers execute on everything in the groupset except
// the coordinate supplied
func (gs GroupSet) ExOthers(co Coord, todo func(Coord) bool) {
	lFunc := func(crd Coord) bool {
		if !co.Eq(crd) {
			return todo(crd)
		}
		return true
	}

	tmp := func(gr Group) bool {
		gr.ExAll(lFunc)
		return true
	}
	gs.ExAll(tmp)
}
