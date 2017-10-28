package sod

// Cell is a single cell of a puzzle
type Cell struct {
	Val []Value
	pz  *Puzzle
	crd Coord
}

// NewBlankCell returns an empty cell
func NewBlankCell(size int) *Cell {
	itm := new(Cell)
	itm.Val = make([]Value, size)
	return itm
}

// NewCell created within a puzzle
func (pz *Puzzle) NewCell(crd Coord) *Cell {
	size := pz.Len()
	itm := new(Cell)
	itm.pz = pz
	itm.crd = crd
	itm.Val = make([]Value, size)
	return itm
}

// NewCellValues setup the values for a cell
func NewCellValues(vals []Value) *Cell {
	itm := NewBlankCell(len(vals))
	for i, vl := range vals {
		itm.Val[i] = vl
	}
	return itm
}
func (cl Cell) String() string {
	retStr := "["
	csl := ""
	for _, v := range cl.Values() {
		retStr += csl + v.String()
		csl = ","
	}
	retStr += "]"
	return retStr
}

// Values returns the array of values in the cell
func (cl Cell) Values() []Value {
	return cl.Val
}

// Copy one cell to another
func (cl Cell) Copy(dst *Cell) {
	dst.pz = cl.pz
	dst.crd = cl.crd
	len := copy(dst.Val, cl.Val)
	// Trunkate the current cell as needed
	dst.Val = dst.Val[:len]
}
func containsValue(val Value, values []Value) bool {
	for _, v := range values {
		if v == val {
			return true
		}
	}
	return false
}

// NotValues - que?
func (cl Cell) NotValues(nv []Value) []Value {
	retArr := make([]Value, 0, cl.Len())
	for _, val := range cl.Values() {
		if !containsValue(val, nv) {
			retArr = append(retArr, val)
		}
	}
	return retArr
}

// Exist returns true of a value exists
func (cl Cell) Exist(val Value) bool {
	return cl.FindVal(val) < len(cl.Val)
}

// FindVal find out if a value exists in the cell
func (cl Cell) FindVal(val Value) int {
	// This could be optimised by assumption
	// that values will be sorted
	for i, v := range cl.Val {
		if v > val {
			return len(cl.Val)
		} else if v == val {
			return i
		}
	}
	return len(cl.Val)
}

// Len of the cell
func (cl Cell) Len() int {
	return len(cl.Val)
}

// RemoveVals from a cell
func (cl *Cell) RemoveVals(vals []Value) error {
	// TBD make this more efficient
	for _, val := range vals {
		err := cl.RemoveVal(val)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveVal - remove a single value from a cell
func (cl *Cell) RemoveVal(val Value) error {
	// Remove a possible value from the cell
	// without creating a new array
	location := cl.FindVal(val)
	cllen := len(cl.Val)
	if location >= len(cl.Val) {
		return ErrUnfoundValue
	}
	if location < cllen {
		copy(cl.Val[location:cllen], cl.Val[location+1:])
	}
	cl.Val = cl.Val[:cllen-1]
	return nil
}

// SetVal - set a cell to ve a single value
func (cl *Cell) SetVal(val Value) {
	cl.Val = []Value{val}
}
func (cl Cell) others(val Value) []Value {
	retVal := make([]Value, 0)
	for _, vl := range cl.Values() {
		if vl != val {
			retVal = append(retVal, vl)
		}
	}
	return retVal
}
