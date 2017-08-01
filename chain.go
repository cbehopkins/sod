package sod

import "log"

type Chain []Cell

func NewChain(input [][]int) *Chain {
	ret_val := make(Chain, len(input))
	for i := 0; i < len(input); i++ {
		tmpLink := NewBlankCell(len(input[i]))
		for j := 0; j < tmpLink.Len(); j++ {
			tmpLink.Val[j] = Value(input[i][j])
		}
		ret_val[i] = *tmpLink
	}
	return &ret_val
}
func (ch *Chain) Add (input Cell) {
  *ch = append(*ch, input)
}
func (ch Chain) searchLink(startOffset, offset int, start, target Value, maxDepth int) []Cell {
	// Look through the chain for where we can get to from the target Vl
	if maxDepth < 0 {
		//log.Println("Giving up, too deep")
		return []Cell{}
	}
	//log.Printf("Searching for %v, at offset %v, Depth %v", target, offset, maxDepth)
	for i, _ := range ch {
		off := (i + offset)
		if off >= len(ch) {
			return []Cell{}
		}
		//log.Println("Looking at element", off)

		link := ch[off]

		if link.Exist(target) {
			if link.Exist(start) {
				// We're back to the start of the chain
				//log.Println("Found circuit for:", link, i)
				return []Cell{link}
			}

			// The link could be a link in a chain
			othr := link.others(target)
			if len(othr) != 1 {
				log.Fatal("")
			}
			//log.Println("No match found, looking at the next link")
			tmp := ch.searchLink(startOffset, off+1, start, othr[0], maxDepth-1)
			if len(tmp) > 0 {
				return append(tmp, link)
			} else {
				return []Cell{}
			}

		}
	}
	//log.Println("Run out of search, giving up")
	return []Cell{}
}
func (ch Chain) SearchChain() []Chain {
	ret_list := make([]Chain, 0)
	maxDepth := 9
	for offset, link := range ch {
		startVl := link.Val[0]
		targetVl := link.Val[1]
		//log.Printf("New search starting at offset:%v, chain starting with %v\n",offset,startVl)
		resultList := ch.searchLink(offset, offset+1, startVl, targetVl, maxDepth)
		if len(resultList) > 0 {
			resultList = append(resultList, link)
			//log.Println("****************Success:",resultList)
			ret_list = append(ret_list, resultList)
		}
	}
	return ret_list
}
