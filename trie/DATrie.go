package trie

import (
	"sort"
	"fmt"
)

var DATRIE_START_LOC int = 256
var DATRIE_HEAD_LOC int = DATRIE_START_LOC - 1

type reference struct {
	len		uint16
	min		byte
	max		byte
	base	int
	ref		*_T
}
func newReference(len uint16, min byte, max byte, ref *_T) *reference {
	return &reference{
		len: len,
		min: min,
		max: max,
		base: 0,
		ref: ref,
	}
}

type byReference []*reference
func (r byReference) Len() int				{return len(r)}
func (r byReference) Swap(i, j int)		{r[i], r[j] = r[j], r[i]}
func (r byReference) Less(i, j int) bool {
	if r[i].len != r[j].len {
		return r[i].len > r[j].len
	}
	
	return r[i].ref.number < r[j].ref.number
}


type DATrie struct {
	reorgSize	int
	expBCSize	int
	arrayItems	int
	states		int
	array		*bcArray
	addTrie		*Trie
	list		[]string
}

func NewDATrie(reorgSize, expBCSize int) *DATrie {
	if reorgSize <= 0 || expBCSize <= 0 {
		panic("Negative parameter");
	}
	
	return &DATrie{
		reorgSize: reorgSize,
		expBCSize: expBCSize,
		arrayItems: 0,
		states: 0,
		array: newBCArray(expBCSize, DATRIE_START_LOC),
		addTrie: NewTrie(),
		list: []string{},
	}
}

func (this *DATrie) reset() {
	this.arrayItems = 0
	this.states = 0
	this.array = newBCArray(this.expBCSize, DATRIE_START_LOC)
	this.addTrie = NewTrie()
}

func (this *DATrie) insert(r *reference) int {
	bs := []byte{}
	for b, _ := range r.ref.childrens {
		bs = append(bs, b)
	}
	
	pos := this.array.searchPosition(bs, r.min, r.max, BCARRAY_FREE_FAST)
	if pos < 0 {
		if r.len < 3 || r.max - r.min < 3 {
			pos = this.array.searchPosition(bs, r.min, r.max, BCARRAY_FREE_SLOW)
		}
	}
	
	if pos < 0 {	//no suitable position
		this.array.extend(this.expBCSize)
		
		pos = this.array.searchPosition(bs, r.min, r.max, BCARRAY_FREE_SLOW)
		
		if pos < 0 {
			pos = this.array.searchPosition(bs, r.min, r.max, BCARRAY_FREE_FAST)
		}
		
		if pos < 0 {
			this.array.dumpFree()
			panic("Could have suitable position at here!")
		}
	}
	
	for k, _ := range r.ref.childrens {
		this.array.take(pos + int(k))	//take the position now
	}
	
	r.base = pos
	r.ref.payload = r
	return pos
}

func (this *DATrie) updateList() {
	lines := []string{}
	for _, str := range this.list {
		if this.Search(str) {
			lines = append(lines, str)
		}
	}
	
	this.list = lines
}

func (this *DATrie) buildFromTrie(_trie *Trie) {
	this.reset()
	
	//get non-leave nodes
	refs := []*reference{}
	_trie.WideVisit(func(n *_T) bool {
		if n.isLeaf() {
			return true	//continue
		}
		
		var min, max byte
		var length uint16
		min = 255
		max = 0
		length = uint16(len(n.childrens))
		
		for k, _ := range n.childrens {
			if min > k {
				min = k
			}
			if max < k {
				max = k
			}
		}
		
		refs = append(refs, newReference(length, min, max, n))
		
		return true
	})

	//sort the reference, we build DATrie using this order
	sort.Sort(byReference(refs))
	
	//insert siblings into BCArray, update check first
	for _, r := range refs {
		this.insert(r)
	}
	
	//now update base and check
	//1. set B/C/N of root node
	this.array.setBase(DATRIE_HEAD_LOC, _trie.root.payload.(*reference).base)
	this.array.setCheck(DATRIE_HEAD_LOC, -DATRIE_HEAD_LOC)
	this.array.setValue(DATRIE_HEAD_LOC, _trie.root.number)
	
	//2. set B/C/N of other nodes
	_trie.root.payload = DATRIE_HEAD_LOC
	_trie.WideVisit(func(n *_T) bool {
		if n.isLeaf() {	
			return true	//continue
		}
		
		ploc := n.payload.(int)		//the parent location in array
		basep := this.array.getBase(ploc)		//the base of loc
		for k, v := range n.childrens {
			idx := basep + int(k)
			
			if v.isLeaf() {
				this.array.setBase(idx, 0)
			} else {
				this.array.setBase(idx, v.payload.(*reference).base)
			}
			
			if v.isWord() {
				this.array.setCheck(idx, -ploc)
			} else {
				this.array.setCheck(idx, ploc)
			}
			
			v.payload = idx
			this.array.setValue(idx, v.number)
		}
		
		return true
	})
	
	this.arrayItems = _trie.Len()
	this.states = _trie.sequence
	return
}

func (this *DATrie) reOrg() {
	fmt.Println("Reorg with ", len(this.list), " items")
	if len(this.list) <= 0 {
		return
	}
	
	//sort the list
	sort.Strings(this.list)
	
	//build trie node
	_trie := BuildTrie(this.list)
	
	//build DATrie from trie tree
	this.buildFromTrie(_trie)
	
	return
}

func abs(elem int) int {
	if elem < 0 {
		return -elem
	}
	
	return elem
}

func (this *DATrie) prefix(str string) int {
	if this.arrayItems <= 0 {	//no item
		return -1
	}
	
	var i int = DATRIE_HEAD_LOC
	for _, b := range []byte(str) {
		base := this.array.getBase(i)
		
		loc := base + int(b)
		if abs(this.array.getCheck(loc)) != i {
			return -1	//prefix does not exist
		}
		
		i = loc
	}
	
	return i
}

func (this *DATrie) searchArray(str string) bool {
	loc := this.prefix(str)
	
	return loc > 0 && this.array.getCheck(loc) < 0
}

func (this *DATrie) delArrayItem(str string) {
	loc := this.prefix(str)
	
	if loc > 0 && this.array.getCheck(loc) < 0 {
		this.arrayItems --
		this.array.setCheck(loc, -this.array.getCheck(loc))	//we can add it back
	}
}

func (this *DATrie) Add(str string) bool {
	if this.addTrie.Search(str) {		//already exist in trie
		return true
	}
	
	if pre := this.prefix(str); pre > 0 {	//the path is already in array
		if c := this.array.getCheck(pre); c > 0 {
			this.arrayItems ++
			this.array.setCheck(pre, -c)
		}
		
		return true
	}
	
	this.list = append(this.list, str)
	this.addTrie.Add(str)
	
	if this.addTrie.Len() >= this.reorgSize {
		this.updateList()
		this.reOrg()
	}
	
	return true
}

func (this *DATrie) Del(str string) bool {
	if !this.Search(str) {
		return true
	}
	
	if this.addTrie.Search(str) {
		return this.addTrie.Del(str)
	}
	
	this.delArrayItem(str)
	
	return true
}

func (this *DATrie) Search(str string) bool {
	return this.searchArray(str) || this.addTrie.Search(str)
}

func (this *DATrie) SearchDebug(str string) bool {
	if this.addTrie.SearchDebug(str) {
		return true
	}
	
	if this.arrayItems <= 0 {
		return false
	}
	
	var i int = DATRIE_HEAD_LOC
	
	fmt.Printf("nil ==> base: %d, check: %d\n", this.array.getBase(i), this.array.getCheck(i))
	for _, b := range []byte(str) {
		base := this.array.getBase(i)

		loc := base + int(b)
		if abs(this.array.getCheck(loc)) != i {
			return false
		}
		
		i = loc
		fmt.Printf("%d ==> base: %d, check: %d\n", b, this.array.getBase(i), this.array.getCheck(i))
	}
	
	return this.array.getCheck(i) < 0
}


func (this *DATrie) Build() { this.reOrg() }

func (this *DATrie) BuildFromStrings(lines []string) {
	this.list = lines
	this.reOrg()
}

func (this *DATrie) Len() int {
	return this.arrayItems + this.addTrie.Len()
}
