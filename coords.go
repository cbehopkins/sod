package sod

import (
	"log"
	"strconv"
)

// Coord implements a coordinate structure
// Coordinate specified as {x,y}
// Or {col,row}
type Coord []int

func (cr Coord) String() string {
	retStr := "["
	if len(cr) != 2 {
		log.Fatal("Coord incorrect length")
	}
	retStr += "x=" + strconv.Itoa(cr.getColumn()) + " "
	retStr += "y=" + strconv.Itoa(cr.getRow())
	retStr += "]"
	return retStr
}
func (cr Coord) getRow() int {
	return cr[1]
}
func (cr Coord) getColumn() int {
	return cr[0]
}

// Eq returns true if two coords are equal
func (cr Coord) Eq(cs Coord) bool {
	if cr.getRow() != cs.getRow() {
		return false
	}
	if cr.getColumn() != cs.getColumn() {
		return false
	}
	return true
}
