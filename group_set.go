package sod

import "log"

type GroupSet []Group

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

// execute on everything in the groupset except
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
