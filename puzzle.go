package sod

import (
	"errors"
	"log"
	"math"
	"strconv"
)

// Puzzle - hum, not sure
type Puzzle struct {
	// Rows in the inner Dimension
	// Columns on outer Dimension
	Puz        [][]Cell
	difficulty *int
}

// IncreaseDifficuly increases the puzzle difficulty
// i.e. an operation has occured that means this is a more difficult puzzle
func (pz Puzzle) IncreaseDifficuly() {
	*pz.difficulty++
}

// Difficulty returns the current puizzle difficulty
func (pz Puzzle) Difficulty() int {
	return *pz.difficulty
}

// AddDifficulty add some difficulty to this puzzle
// i.e. an operation has occured that means this is a more difficult puzzle
func (pz Puzzle) AddDifficulty(val int) {
	*pz.difficulty += val
}

//for each column, return the max number of items
func (pz Puzzle) maxColItems() (retArr []int) {
	pzSize := pz.Len()
	retArr = make([]int, pzSize)
	for column := 0; column < pzSize; column++ {
		for row := 0; row < pzSize; row++ {

			tmpLen := pz.Puz[column][row].Len()
			if tmpLen > retArr[column] {
				retArr[column] = tmpLen
			}
		}
	}
	return retArr
}
func maxWidth(colWid int) int {
	retVal := 1
	for ; colWid > 1; colWid-- {
		retVal += 2
	}
	return retVal
}
func numSpaces(cnt int) string {
	retStr := ""

	for i := 0; i < cnt; i++ {
		retStr += " "
	}
	return retStr
}
func (pz Puzzle) String() string {
	pzSize := pz.Len()
	colWidths := pz.maxColItems()
	retStr := "[\n"
	for row := 0; row < pzSize; row++ {
		retStr += "  [row=" + strconv.Itoa(row) + "\n    "
		for column := 0; column < pzSize; column++ {
			tmpString := pz.Puz[column][row].String()
			maxW := maxWidth(colWidths[column])
			addWidth := 2 + maxW - (len(tmpString))
			tmpString += numSpaces(addWidth)
			retStr += "(" + strconv.Itoa(column) + ")" + tmpString
		}
		retStr += "\n  ]\n"
	}
	retStr += "\n]"
	return retStr
}

// StringGroup is like String, but for a Group
func (pz Puzzle) StringGroup(gp Group) string {
	retStr := "["
	for _, crd := range gp.Items() {
		cel := pz.GetCel(crd)
		retStr += "[" + crd.String() + "," + cel.String() + "]"
	}
	retStr += "\n]"
	return retStr
}

// NewPuzzle return a new puzzle
func NewPuzzle() *Puzzle {
	itm := new(Puzzle)
	itm.Puz = make([][]Cell, PuzzleSize)
	tmpArr := NewBlankCell(PuzzleSize)
	for i := 0; i < PuzzleSize; i++ {
		tmpArr.Val[i] = Value(i + 1)
	}
	for i := 0; i < PuzzleSize; i++ {
		row := []Cell{}
		for j := 0; j < PuzzleSize; j++ {
			tmpCopy := itm.NewCell(Coord{i, j})
			copy(tmpCopy.Val, tmpArr.Val)
			row = append(row, *tmpCopy)
		}
		itm.Puz[i] = row
	}
	itm.difficulty = new(int)
	return itm
}

// Duplicate a puzzle into a new one
func (pz Puzzle) Duplicate() (dst Puzzle) {
	dst.Puz = make([][]Cell, PuzzleSize)
	for i := 0; i < PuzzleSize; i++ {
		row := []Cell{}
		for j := 0; j < PuzzleSize; j++ {
			tmpCopy := dst.NewCell(Coord{i, j})
			pz.Puz[i][j].Copy(tmpCopy)
			row = append(row, *tmpCopy)
		}
		dst.Puz[i] = row
	}
	dst.difficulty = new(int)
	*dst.difficulty = *pz.difficulty
	return dst
}

// Copy one puzzle into another
func (pz Puzzle) Copy(dst Puzzle) {
	if len(dst.Puz) < PuzzleSize {
		log.Fatal("Bumy", dst)
	}
	for i := 0; i < PuzzleSize; i++ {
		if len(dst.Puz[i]) < PuzzleSize {
			log.Fatal("Bum", dst)
		}
		for j := 0; j < PuzzleSize; j++ {
			pz.Puz[i][j].Copy(&dst.Puz[i][j])
		}
	}
	*dst.difficulty = *pz.difficulty
}

func (pz Puzzle) setValC(value Value, col, row int) {
	crd := Coord{col, row}
	pz.SetVal(value, crd)
}

// SetVal  Set a Cell in the puzzle to a known value
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
	pz.runRemove(value, co)
}

func (pz Puzzle) runRemove(value Value, co Coord) {
	rmFunc := func(cr Coord) bool {
		pz.RemoveVal(value, cr)
		return true
	}

	// For everything in this row, remove the value
	pz.rowCoords(co.getRow()).ExOthers(co, rmFunc)
	// For everything in this column, remove the value
	pz.colCoords(co.getColumn()).ExOthers(co, rmFunc)
	// For everything in this neighbourhood, remove the value
	pz.neighCoords(co).ExOthers(co, rmFunc)
}

// SelfCheck Tests if the puzzle is consistent
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
		}
		return false
	}
	pz.ExAllGroups(lFunc)

	return nil
}

// ValExist returns true if the value exists
func (pz Puzzle) ValExist(value Value, co Coord) bool {
	cel := pz.GetCel(co)
	return cel.Exist(value)
}

// RemoveVal - For one cell in a puzzle
// Remove the value from a particular cell co-ord
func (pz Puzzle) RemoveVal(value Value, co Coord) {
	cel := pz.GetCel(co)
	err := cel.RemoveVal(value)
	if err == nil {
		if cel.Len() == 1 {
			//pz.SetVal(cel.Val[0], co)
			pz.runRemove(cel.Val[0], co)
		}
	}
}

// RemoveVals remove the specified values (if they exist) from the coord
func (pz Puzzle) RemoveVals(values []Value, co Coord) {
	cel := pz.GetCel(co)
	err := cel.RemoveVals(values)
	if err == nil {
		if cel.Len() == 1 {
			//pz.SetVal(cel.Val[0], co)
			pz.runRemove(cel.Val[0], co)
		}
	}
}

// Len of the puzzle
func (pz Puzzle) Len() int {
	return len(pz.Puz)
}

// Coords Generate all the coords for the puzzle
func (pz Puzzle) Coords() []Coord {
	returnArray := make([]Coord, 0, PuzzleSize*PuzzleSize)
	pzLen := pz.Len()

	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			newCoord := []int{i, j}
			returnArray = append(returnArray, newCoord)
		}
	}
	return returnArray
}

// Generate the Coords for all the Neighbourhoods
func (pz *Puzzle) allNeighCoords() GroupSet {
	pzLen := pz.Len()
	pzDiv := int(math.Sqrt(float64(pzLen)))
	retArray := make([]Group, pzLen)
	iSelect := 0
	iCnt := pzDiv
	for i := 0; i < pzLen; i++ {
		retArray[i] = *NewGroup(pzLen, pz)
	}
	for i := 0; i < pzLen; i++ {
		if iCnt == 0 {
			iSelect++
			iCnt = pzDiv
		}
		if iCnt > 0 {
			iCnt--
		}
		jSelect := 0
		jCnt := pzDiv
		for j := 0; j < pzLen; j++ {
			if jCnt == 0 {
				jSelect++
				jCnt = pzDiv
			}
			if jCnt > 0 {
				jCnt--
			}
			// TBD Optimize dividers
			selectNum := (iSelect * pzDiv) + (jSelect)
			newCoord := []int{i, j}
			//fmt.Printf("%v,%v is in %v\nis:%v,js:%v\n",i,j,select_num,i_select,j_select)
			retArray[selectNum].Add(newCoord)
		}
	}
	return retArray
}
func (pz Puzzle) allGroupSets() GroupSet {
	retSet := pz.allNeighCoords()
	retSet = append(retSet, pz.allColCoords()...)
	retSet = append(retSet, pz.allRowCoords()...)
	return retSet
}

// Get all the coordinate of all cells in the neighbourhood
func (pz *Puzzle) neighCoords(crd Coord) Group {
	pzLen := pz.Len()
	pzDiv := int(math.Sqrt(float64(pzLen)))
	retArray := NewGroup(pzLen, pz)

	// I want to round down to the previous multiple
	startRow := ((crd.getRow() / pzDiv) * pzDiv)
	startCol := ((crd.getColumn() / pzDiv) * pzDiv)
	stopRow := startRow + pzDiv
	stopCol := startCol + pzDiv

	for row := startRow; row < stopRow; row++ {
		for col := startCol; col < stopCol; col++ {
			newCoord := []int{col, row}
			retArray.Add(newCoord)

		}
	}
	return *retArray
}

// Generate the coords for a specifies row
func (pz *Puzzle) rowCoords(rowNum int) Group {
	pzLen := pz.Len()
	retArray := NewGroup(0, pz)
	for i := 0; i < pzLen; i++ {
		newCoord := Coord{i, rowNum}
		retArray.Add(newCoord)
	}
	return *retArray
}

// Generate the Coords for a specified column
func (pz *Puzzle) colCoords(colNum int) Group {
	pzLen := pz.Len()
	retArray := NewGroup(0, pz)
	for i := 0; i < pzLen; i++ {
		newCoord := Coord{colNum, i}
		retArray.Add(newCoord)
	}
	return *retArray
}

func (pz *Puzzle) allRowCoords() GroupSet {
	if pz == nil {
		log.Fatal("Nil puzzle at start of allRowCoords")
	}
	pzLen := pz.Len()
	retArray := make([]Group, 0, pzLen)
	for i := 0; i < pzLen; i++ {
		retArray = append(retArray, pz.rowCoords(i))
	}
	return retArray
}

// Return a set of all the column coords
func (pz *Puzzle) allColCoords() GroupSet {
	pzLen := pz.Len()
	retArray := make([]Group, 0, pzLen)
	for i := 0; i < pzLen; i++ {
		retArray = append(retArray, pz.colCoords(i))
	}
	return retArray
}

// ExAllGroups run a func on all groups
func (pz *Puzzle) ExAllGroups(lFunc func(Group) bool) {

	if pz == nil {
		log.Fatal("Nil puzzle at start of ExAllGroups")
	}
	pz.allRowCoords().ExAll(lFunc)
	pz.allColCoords().ExAll(lFunc)
	pz.allNeighCoords().ExAll(lFunc)
}

// ExAll run a func across all coords
func (pz *Puzzle) ExAll(lFunc func(Coord) bool) {
	for _, crd := range pz.Coords() {
		result := lFunc(crd)
		// Keep going while true
		if !result {
			break
		}
	}
}

// GetVal returns the value at the specified coordinate
func (pz Puzzle) GetVal(crd Coord) []Value {
	cell := pz.GetCel(crd)
	return cell.Val
}

// GetCel returns the cel at a specified coordinate
func (pz Puzzle) GetCel(crd Coord) *Cell {
	if len(crd) != 2 {
		log.Fatal("Not a 2d coord")
	}
	xc := crd[0]
	yc := crd[1]

	return &pz.Puz[xc][yc]
}

// ErrUnsolved - Unsolved, but possible
var ErrUnsolved = errors.New("Unsolved, but possible")

// ErrZeroCell - Unsolved as zero options for cell
var ErrZeroCell = errors.New("Unsolved as zero options for cell")

// Solved returns an error if unsolved
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

// Check the puzzle for consistency
func (pz Puzzle) Check(result [][]Value) bool {
	pzLen := pz.Len()
	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			cord := Coord{i, j}
			if pz.GetCel(cord).Val[0] != result[j][i] {
				log.Printf("Coord %v,%v not match %v,%v", i, j, pz.GetCel(cord).Val[0], result[j][i])
				return false
			}
		}
	}
	return true
}

// Match Check puzzles contain the same data
func (pz Puzzle) Match(ref Puzzle) bool {
	pzLen := pz.Len()
	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			refCord := Coord{i, j}
			refCel := ref.GetCel(refCord)
			pzCel := pz.GetCel(refCord)

			if refCel.Len() == pzCel.Len() {
				refValues := refCel.Values()
				puzValues := pzCel.Values()
				for i, v := range refValues {
					if puzValues[i] != v {
						return false
					}
				}
			} else {
				return false
			}
		}
	}
	return true
}
func (pz Puzzle) roughCheck(ref Puzzle) bool {
	// Check all solved values in result are in
	// the pz to be checked
	pzLen := pz.Len()
	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			refCord := Coord{i, j}
			refCel := ref.GetCel(refCord)
			if refCel.Len() == 1 {
				chkCel := pz.GetCel(refCord)
				if chkCel.FindVal(refCel.Val[0]) == chkCel.Len() {
					// The value does not exist
					return false
				}
			}
		}
	}
	return true
}
func (pz Puzzle) lessRoughCheck(ref Puzzle) bool {
	// Check all solved values in result are solved in pz
	pzLen := pz.Len()
	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			refCord := Coord{i, j}
			refCel := ref.GetCel(refCord)
			if refCel.Len() == 1 {
				chkCel := pz.GetCel(refCord)
				if chkCel.Len() != 1 {
					// the pz cell is not solved
					return false
				}
				if chkCel.Val[0] != refCel.Val[0] {
					// The value is incorrect
					return false
				}
			}
		}
	}
	return true
}

func (pz Puzzle) load(src [][]Value) {
	pzLen := pz.Len()
	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			if src[i][j] != 0 {
				cord := Coord{i, j}
				pz.SetVal(Value(src[i][j]), cord)
			}
		}
	}
}
