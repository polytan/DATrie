package trie

import "fmt"

type _BC struct {
	b		int		//base array
	c		int		//check array
	n		int		//the visit count, used for fast locating valid position
}

type bcArray struct {
	array	[]_BC	//the underline array
	fast	int		//the fast locating list
	slow	int		//the slow locating list
}

const BCARRAY_FREE_FAST int = 0
const BCARRAY_FREE_SLOW int = 1
const BCARRAY_MIN_SIZE int = 256

func newBCArray(size int, reserve int) *bcArray {
	if size < reserve || size < BCARRAY_MIN_SIZE {
		panic("make array with no enough size: " + string(size))
	}
	
	arr := &bcArray{
		array: make([]_BC, size),
		fast: BCARRAY_FREE_FAST,
		slow: BCARRAY_FREE_SLOW,
	}
	
	//fast free list
	arr.array[BCARRAY_FREE_FAST].b = BCARRAY_FREE_FAST
	arr.array[BCARRAY_FREE_FAST].c = BCARRAY_FREE_FAST
	
	//slow free list
	arr.array[BCARRAY_FREE_SLOW].b = BCARRAY_FREE_SLOW
	arr.array[BCARRAY_FREE_SLOW].c = BCARRAY_FREE_SLOW
	
	arr.linkFreeItem(reserve, int(size))
	return arr
}

func (a *bcArray) extend(size int) {
	if size <= BCARRAY_MIN_SIZE {
		panic("no enough size: " + string(size))
	}
	
	sBefore := len(a.array)
	a.array = append(a.array, make([]_BC, size)...)
	sAfter := len(a.array)
	
	a.linkFreeItem(sBefore, sAfter)
//	fmt.Println("Extend array to : ", sAfter)
}

//link all free items into fast locating list
func (a *bcArray) linkFreeItem(s, e int) {
	for i := s; i < e; i ++ {
		a.array[i].b = -(i + 1)
		a.array[i].c = i - 1
		a.array[i].n = 0
	}
	
	h := a.fast
	
	//link start node s with current list
	a.array[a.array[h].c].b = -s
	a.array[s].c = a.array[h].c
	
	//link end node e with current list
	a.array[e - 1].b = -h
	a.array[h].c = e - 1
	
	return
}

//take the ith item from free list
func (a *bcArray) take(i int) {
	n := -a.array[i].b
	p := a.array[i].c
	
	a.array[n].c = p
	a.array[p].b = -n
	
	a.array[i].b = 0
	a.array[i].c = 0
}
func (a *bcArray) setBase(i int, b int) { a.array[i].b = b }
func (a *bcArray) setCheck(i int, c int) { a.array[i].c = c }
func (a *bcArray) setPayload(i int, n int) { a.array[i].n = n }
func (a *bcArray) getBase(i int) int { return a.array[i].b }
func (a *bcArray) getCheck(i int) int { return a.array[i].c }
func (a *bcArray) getPayload(i int) int { return a.array[i].n }

func (a *bcArray) nextPos(i int, l int) int {
	next := -a.array[i].b
	if next <= l {
		return next
	}
	
	if a.array[next].n > 150 {
		pos := a.nextPos(next, l)
		
		if a.slow != l {
			//pick the node from fast list
			n := -a.array[next].b
			p := a.array[next].c
			
			a.array[n].c = p
			a.array[p].b = -n
			
			//add the node to slow list
			a.array[next].b = -a.slow
			a.array[next].c = a.array[a.slow].c
			a.array[a.array[a.slow].c].b = -next
			a.array[a.slow].c = next
		}
		
		next = pos
	} else {
		a.array[next].n += 1
	}
	
	return next
}

func (a *bcArray) searchPosition(items []byte, min byte, max byte, start int) int {	
	pos := -1
	for i := a.nextPos(start, start); i > start; i = a.nextPos(i, start) {
		ok := true
		p := i - int(min)
		
		if p + int(max) >= len(a.array) {
			continue	//this location exceeds max size
		}
		
		for _, k := range items {
			if a.array[p + int(k)].b >= 0 {
				ok = false
				break
			}
		}
		
		if ok {
			pos = p
			break
		}
	}
	
	return pos
}

func (a *bcArray) dumpFree() {
	var i int = a.fast
	
	fmt.Println("Dump first 10 fast free nodes: ")
	for k := 0; k < 10; k ++ {
		fmt.Printf("node[%d], b:%d, c:%d, n:%d\n", 
								i, a.array[i].b, a.array[i].c, a.array[i].n)
		if a.array[i].b == a.fast {
			break
		}
		
		i = -a.array[i].b
	}
	
	i = a.slow
	fmt.Println("Dump first 10 slow free nodes: ")
	for k := 0; k < 10; k ++ {
		fmt.Printf("node[%d], b:%d, c:%d, n:%d\n",
								i, a.array[i].b, a.array[i].c, a.array[i].n)
		if a.array[i].b == a.slow {
			break
		}
		
		i = -a.array[i].b
	}
}

