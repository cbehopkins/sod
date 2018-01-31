package sod

type valueMap map[Value]*CrdCnt

func newValueMap() valueMap {
	itm := new(valueMap)
	*itm = make(map[Value]*CrdCnt)
	return *itm
}

func (vm valueMap) Add(v Value, itm *CrdCnt) {
	vm[v] = itm
}
func (vm valueMap) Delete(v Value) {
	delete(vm, v)
}
func (vm valueMap) Get(v Value) *CrdCnt {
	return vm[v]
}
