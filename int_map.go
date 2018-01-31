package sod

type intMap map[int]*CrdCnt

func newIntMap() intMap {
	itm := new(intMap)
	*itm = make(map[int]*CrdCnt)
	return *itm
}

func (vm intMap) Add(v int, itm *CrdCnt) {
	vm[v] = itm
}
func (vm intMap) Delete(v int) {
	delete(vm, v)
}
func (vm intMap) Get(v int) *CrdCnt {
	return vm[v]
}
