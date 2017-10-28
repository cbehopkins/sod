package sod

import "strconv"

// CrdCnt Coordinate Count
type CrdCnt struct {
	Cnt     int
	LocList *Group
}

func (cc CrdCnt) String() string {
	retStr := "CrdCnt[\n"
	retStr += "Cnt:" + strconv.Itoa(cc.Cnt+1) + ","
	retStr += cc.LocList.String()
	retStr += "\n]"
	return retStr
}

// NewCrdCnt New coordinate count
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

// Add a coordinate to the list
func (cc *CrdCnt) Add(crd Coord) {
	cc.Cnt++
	cc.set(crd)
}
