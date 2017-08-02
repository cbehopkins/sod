package sod

import (
	"errors"
	"log"
	"math"
	"strconv"
)

type Puzzle struct {
	// Rows in the inner Dimension
	// Columns on outer Dimension
	Puz        [][]Cell
	difficulty *int
}

func (pz Puzzle) IncreaseDifficuly() {
	*pz.difficulty++
}
func (pz Puzzle) Difficulty() int {
	return *pz.difficulty
}
func (pz Puzzle) AddDifficulty(val int) {
	*pz.difficulty += val
}
func (pz Puzzle) String() string {
	pz_size := pz.Len()
	ret_str := "[\n"
	for row := 0; row < pz_size; row++ {
		ret_str += "  [row=" + strconv.Itoa(row) + "\n    "
		for column := 0; column < pz_size; column++ {
			ret_str += "(" + strconv.Itoa(column) + ")" + pz.Puz[column][row].String()
		}
		ret_str += "\n  ]\n"
	}
	ret_str += "\n]"
	return ret_str
}
func (pz Puzzle) StringGroup(gp Group) string {
	ret_str := "["
	for _, crd := range gp.Items() {
		cel := pz.GetCel(crd)
		ret_str += "[" + crd.String() + "," + cel.String() + "]"
	}
	ret_str += "\n]"
	return ret_str
}
func NewPuzzle() *Puzzle {
	itm := new(Puzzle)
	itm.Puz = make([][]Cell, PuzzleSize)
	tmp_arr := NewBlankCell(PuzzleSize)
	for i := 0; i < PuzzleSize; i++ {
		tmp_arr.Val[i] = Value(i + 1)
	}
	for i := 0; i < PuzzleSize; i++ {
		row := []Cell{}
		for j := 0; j < PuzzleSize; j++ {
			tmp_copy := itm.NewCell(Coord{i, j})
			copy(tmp_copy.Val, tmp_arr.Val)
			row = append(row, *tmp_copy)
		}
		itm.Puz[i] = row
	}
	itm.difficulty = new(int)
	return itm
}
func (src Puzzle) Duplicate() (dst Puzzle) {
	dst.Puz = make([][]Cell, PuzzleSize)
	for i := 0; i < PuzzleSize; i++ {
		row := []Cell{}
		for j := 0; j < PuzzleSize; j++ {
			tmp_copy := dst.NewCell(Coord{i, j})
			src.Puz[i][j].Copy(tmp_copy)
			row = append(row, *tmp_copy)
		}
		dst.Puz[i] = row
	}
	dst.difficulty = new(int)
	*dst.difficulty = *src.difficulty
	return dst
}
func (src Puzzle) Copy(dst Puzzle) {
	if len(dst.Puz) < PuzzleSize {
		log.Fatal("Bumy", dst)
	}
	for i := 0; i < PuzzleSize; i++ {
		if len(dst.Puz[i]) < PuzzleSize {
			log.Fatal("Bum", dst)
		}
		for j := 0; j < PuzzleSize; j++ {
			src.Puz[i][j].Copy(&dst.Puz[i][j])
		}
	}
	*dst.difficulty = *src.difficulty
}

func (pz Puzzle) SetValC(value Value, col, row int) {
	crd := Coord{col, row}
	pz.SetVal(value, crd)
}

// Set a Cell in the puzzle to a known value
// Going through and updating other unsolved cells
func (pz Puzzle) SetVal(value Value, co Coord) {
	cl := pz.GetCel(co)
	// Set the value for the specified cell
  if cl.Len() == 1 {
    if !cl.Exist(value) {
      log.Fatal("Told to set a value that is not an option")
    }
    // Early abort
    //return
  }
	cl.SetVal(value)
  pz.runRemove(value,co)
}

func (pz Puzzle) runRemove(value Value, co Coord) {
	rmFunc := func(cr Coord) bool {
		pz.RemoveVal(value, cr)
		return true
	}

	// For everything in this row, remove the value
	pz.RowCoords(co.GetRow()).ExOthers(co, rmFunc)
	// For everything in this column, remove the value
	pz.ColCoords(co.GetColumn()).ExOthers(co, rmFunc)
	// For everything in this neighbourhood, remove the value
	pz.NeighCoords(co).ExOthers(co, rmFunc)
}

// Test if the puzzle is consistent
func (pz *Puzzle) SelfCheck() error {
	var result error
	if pz == nil {
		log.Fatal("Nil puzzle at start of self check")
	}
	lFunc := func(gr Group) bool {

		if gr.pz == nil {
			log.Fatal("Nil puzzle in self check")
		}
		result = gr.SelfCheck()
		if result != nil {
			return true
		} else {
			return false
		}
	}
	pz.ExAllGroups(lFunc)

	return nil
}
func (pz Puzzle) ValExist(value Value, co Coord) bool {
	cel := pz.GetCel(co)
	return cel.Exist(value)
}

// For one cell in a puzzle
// Remove the value from a particular cell co-ord
func (pz Puzzle) RemoveVal(value Value, co Coord) {
	cel := pz.GetCel(co)
	err := cel.RemoveVal(value)
	if err == nil {
		if cel.Len() == 1 {
			//pz.SetVal(cel.Val[0], co)
      pz.runRemove(cel.Val[0],co)
		}
	}
}
func (pz Puzzle) RemoveVals(values []Value, co Coord) {
	cel := pz.GetCel(co)
	err := cel.RemoveVals(values)
	if err == nil {
		if cel.Len() == 1 {
			//pz.SetVal(cel.Val[0], co)
      pz.runRemove(cel.Val[0],co)
      }
	}
}
func (pz Puzzle) Len() int {
	return len(pz.Puz)
}

// Generate all the coords for the puzzle
func (pz Puzzle) Coords() []Coord {
	returnArray := make([]Coord, 0, PuzzleSize*PuzzleSize)
	pz_len := pz.Len()

	for i := 0; i < pz_len; i++ {
		for j := 0; j < pz_len; j++ {
			new_coord := []int{i, j}
			returnArray = append(returnArray, new_coord)
		}
	}
	return returnArray
}

// Generate the Coords for all the Neighbourhoods
func (pz *Puzzle) AllNeighCoords() GroupSet {
	pz_len := pz.Len()
	pz_div := int(math.Sqrt(float64(pz_len)))
	ret_array := make([]Group, pz_len)
	i_select := 0
	i_cnt := pz_div
	for i := 0; i < pz_len; i++ {
		ret_array[i] = *NewGroup(pz_len, pz)
	}
	for i := 0; i < pz_len; i++ {
		if i_cnt == 0 {
			i_select++
			i_cnt = pz_div
		}
		if i_cnt > 0 {
			i_cnt--
		}
		j_select := 0
		j_cnt := pz_div
		for j := 0; j < pz_len; j++ {
			if j_cnt == 0 {
				j_select++
				j_cnt = pz_div
			}
			if j_cnt > 0 {
				j_cnt--
			}
			// TBD Optimize dividers
			select_num := (i_select * pz_div) + (j_select)
			new_coord := []int{i, j}
			//fmt.Printf("%v,%v is in %v\nis:%v,js:%v\n",i,j,select_num,i_select,j_select)
			ret_array[select_num].Add(new_coord)
		}
	}
	return ret_array
}
func (pz Puzzle) AllGroupSets() GroupSet {
	ret_set := pz.AllNeighCoords()
	ret_set = append(ret_set, pz.AllColCoords()...)
	ret_set = append(ret_set, pz.AllRowCoords()...)
	return ret_set
}

// Get all the coordinate of all cells in the neighbourhood
func (pz *Puzzle) NeighCoords(crd Coord) Group {
	pz_len := pz.Len()
	pz_div := int(math.Sqrt(float64(pz_len)))
	ret_array := NewGroup(pz_len, pz)

	// I want to round down to the previous multiple
	startRow := ((crd.GetRow() / pz_div) * pz_div)
	startCol := ((crd.GetColumn() / pz_div) * pz_div)
	stopRow := startRow + pz_div
	stopCol := startCol + pz_div

	for row := startRow; row < stopRow; row++ {
		for col := startCol; col < stopCol; col++ {
			new_coord := []int{col, row}
			ret_array.Add(new_coord)

		}
	}
	return *ret_array
}

// Generate the coords for a specifies row
func (pz *Puzzle) RowCoords(row_num int) Group {
	pz_len := pz.Len()
	ret_array := NewGroup(0, pz)
	for i := 0; i < pz_len; i++ {
		new_coord := Coord{i, row_num}
		ret_array.Add(new_coord)
	}
	return *ret_array
}

// Generate the Coords for a specified column
func (pz *Puzzle) ColCoords(col_num int) Group {
	pz_len := pz.Len()
	ret_array := NewGroup(0, pz)
	for i := 0; i < pz_len; i++ {
		new_coord := Coord{col_num, i}
		ret_array.Add(new_coord)
	}
	return *ret_array
}

func (pz *Puzzle) AllRowCoords() GroupSet {
	if pz == nil {
		log.Fatal("Nil puzzle at start of AllRowCoords")
	}
	pz_len := pz.Len()
	ret_array := make([]Group, 0, pz_len)
	for i := 0; i < pz_len; i++ {
		ret_array = append(ret_array, pz.RowCoords(i))
	}
	return ret_array
}
func (pz *Puzzle) AllColCoords() GroupSet {
	pz_len := pz.Len()
	ret_array := make([]Group, 0, pz_len)
	for i := 0; i < pz_len; i++ {
		ret_array = append(ret_array, pz.ColCoords(i))
	}
	return ret_array
}
func (pz *Puzzle) ExAllGroups(lFunc func(Group) bool) {

	if pz == nil {
		log.Fatal("Nil puzzle at start of ExAllGroups")
	}
	pz.AllRowCoords().ExAll(lFunc)
	pz.AllColCoords().ExAll(lFunc)
	pz.AllNeighCoords().ExAll(lFunc)
}
func (pz *Puzzle) ExAll(lFunc func(Coord) bool) {
	for _, crd := range pz.Coords() {
		result := lFunc(crd)
		// Keep going while true
		if !result {
			break
		}
	}
}
func (pz Puzzle) GetVal(crd Coord) []Value {
	cell := pz.GetCel(crd)
	return cell.Val
}
func (pz Puzzle) GetCel(crd Coord) *Cell {
	if len(crd) != 2 {
		log.Fatal("Not a 2d coord")
	}
	xc := crd[0]
	yc := crd[1]

	return &pz.Puz[xc][yc]
}

var ErrUnsolved = errors.New("Unsolved, but possible")
var ErrZeroCell = errors.New("Unsolved as zero options for cell")

func (pz Puzzle) Solved() error {
	for _, crd := range pz.Coords() {
		cel := pz.GetCel(crd)
		if cel.Len() == 0 {
			//log.Fatal("Zero length cell???")
			return ErrZeroCell
		}
		if cel.Len() > 1 {
			return ErrUnsolved
		}
	}
	return nil
}

func (pz Puzzle) Check(result [][]Value) bool {
	pz_len := pz.Len()
	for i := 0; i < pz_len; i++ {
		for j := 0; j < pz_len; j++ {
			cord := Coord{i, j}
			if pz.GetCel(cord).Val[0] != result[i][j] {
				return false
			}
		}
	}
	return true
}
func (pz Puzzle) RoughCheck(ref Puzzle) bool {
	// Check all solved values in result are in
	// the pz to be checked
	pz_len := pz.Len()
	for i := 0; i < pz_len; i++ {
		for j := 0; j < pz_len; j++ {
			ref_cord := Coord{i, j}
			ref_cel := ref.GetCel(ref_cord)
			if ref_cel.Len() == 1 {
				chk_cel := pz.GetCel(ref_cord)
				if chk_cel.FindVal(ref_cel.Val[0]) == chk_cel.Len() {
					// The value does not exist
					return false
				}
			}
		}
	}
	return true
}
func (pz Puzzle) LessRoughCheck(ref Puzzle) bool {
	// Check all solved values in result are solved in pz
	pz_len := pz.Len()
	for i := 0; i < pz_len; i++ {
		for j := 0; j < pz_len; j++ {
			ref_cord := Coord{i, j}
			ref_cel := ref.GetCel(ref_cord)
			if ref_cel.Len() == 1 {
				chk_cel := pz.GetCel(ref_cord)
				if chk_cel.Len() != 1 {
					// the pz cell is not solved
					return false
				}
				if chk_cel.Val[0] != ref_cel.Val[0] {
					// The value is incorrect
					return false
				}
			}
		}
	}
	return true
}

func (pz Puzzle) Load(src [][]Value) {
	pz_len := pz.Len()
	for i := 0; i < pz_len; i++ {
		for j := 0; j < pz_len; j++ {
			if src[i][j] != 0 {
				cord := Coord{i, j}
				pz.SetVal(Value(src[i][j]), cord)
			}
		}
	}
}
