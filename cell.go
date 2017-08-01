package sod

type Cell struct {
	Val []Value
	pz  *Puzzle
	crd Coord
}

func NewBlankCell(size int) *Cell {
	itm := new(Cell)
	itm.Val = make([]Value, size)
	return itm
}

func (pz *Puzzle) NewCell(crd Coord) *Cell {
	size := pz.Len()
	itm := new(Cell)
	itm.pz = pz
	itm.crd = crd
	itm.Val = make([]Value, size)
	return itm
}
func NewCellValues(vals []Value) *Cell {
	itm := NewBlankCell(len(vals))
	for i, vl := range vals {
		itm.Val[i] = vl
	}
	return itm
}
func (cl Cell) String() string {
	ret_str := "["
	csl := ""
	for _, v := range cl.Values() {
		ret_str += csl + v.String()
		csl = ","
	}
	ret_str += "]"
	return ret_str
}
func (cl Cell) Values() []Value {
	return cl.Val
}
func (src Cell) Copy(dst *Cell) {
	dst.pz = src.pz
	dst.crd = src.crd
	len := copy(dst.Val, src.Val)
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
func (cl Cell) NotValues(nv []Value) []Value {
	ret_arr := make([]Value, 0, cl.Len()-len(nv))
	for _, val := range cl.Values() {
		if !containsValue(val, nv) {
			ret_arr = append(ret_arr, val)
		}
	}
	return ret_arr
}
func (cl Cell) Exist(val Value) bool {
	return cl.FindVal(val) < len(cl.Val)
}
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

func (cl Cell) Len() int {
	return len(cl.Val)
}
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
func (cl *Cell) SetVal(val Value) {
	cl.Val = []Value{val}
}
func (ln Cell) others(val Value) []Value {
	ret_val := make([]Value, 0)
	for _, vl := range ln.Values() {
		if vl != val {
			ret_val = append(ret_val, vl)
		}
	}
	return ret_val
}
