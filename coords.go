package sod

import (
"strconv"
"log"
)

// Coordinate specified as {x,y}
// Or {col,row}
type Coord []int

func (cr Coord) String() string {
	ret_str := "["
	//spacer := ""
	//for _, v := range cr {
	//	ret_str += spacer + strconv.Itoa(v)
	//	spacer = " "
	//}
  if len(cr) != 2 {
    log.Fatal("Coord incorrect length")
  }
  ret_str += "x=" + strconv.Itoa(cr.GetColumn()) + " "
  ret_str += "y=" + strconv.Itoa(cr.GetRow())
	ret_str += "]"
	return ret_str
}
func (cr Coord) GetRow() int {
	return cr[1]
}
func (cr Coord) GetColumn() int {
	return cr[0]
}
func (cr Coord) Eq(cs Coord) bool {
	if cr.GetRow() != cs.GetRow() {
		return false
	}
	if cr.GetColumn() != cs.GetColumn() {
		return false
	}
	return true
}
