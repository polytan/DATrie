package trie

//import "fmt"

type _ACDA struct {
	failure	int
	output	[]string
}

type ACDATrie struct {
	nodeStruct	*DATrie
	nodeValue	[]_ACDA
}

func NewACDATrie(reorgSize, expBCSize int) *ACDATrie {
	return &ACDATrie{
		nodeStruct: NewDATrie(reorgSize, expBCSize),
		nodeValue: []_ACDA{},
	}
}

func (this *ACDATrie) reset() {
	this.nodeStruct.reset()
	this.nodeValue = []_ACDA{}
}

func (this *ACDATrie) BuildFromACTrie(acTrie *ACTrie) {
//	fmt.Println("====>Phase 2.1: build DATrie from ACTrie.Trie")
	this.reset()
	this.nodeStruct.buildFromTrie(acTrie.nodeStruct)
	
//	fmt.Println("====>Phase 2.2: copy nodeValue from ACTrie")
	this.nodeValue = make([]_ACDA, len(acTrie.nodeValue))
	for i := 0; i < len(this.nodeValue); i ++ {
		p := &this.nodeValue[i]
		v := &acTrie.nodeValue[i]
		
		p.failure = v.failure.payload.(int)
		p.output = v.output
	}
	
	return
}

func (this *ACDATrie) BuildFromArray(arr []string) {
//	fmt.Println("==>Phase 1: build ACTrie from string")
	acTrie := NewACTrie()
	acTrie.BuildTrie(arr)
	
//	fmt.Println("==>Phase 2: build DATrie from ACTrie.Trie")
	this.BuildFromACTrie(acTrie)
	
//	fmt.Println("==>Finished!")
	return
}

func (this *ACDATrie) g(s int, a byte) int {
	b := this.nodeStruct.array.getBase(s)
	
	n := b + int(a)
	c := this.nodeStruct.array.getCheck(n)
	if s == abs(c) {
		return n
	}
	
	return -1
}

func (this *ACDATrie) f(s int) int {
	i := this.nodeStruct.array.getValue(s)
	return this.nodeValue[i].failure
}

func (this *ACDATrie) out(s int) []string {
	i := this.nodeStruct.array.getValue(s)
	return this.nodeValue[i].output
}

func (this *ACDATrie) SearchTrie(str string) []string {
	out := []string{}
	
	q := DATRIE_HEAD_LOC
	for _, m := range []byte(str) {
		for {
			if this.g(q, m) != -1 {
				break
			}
			
			if q = this.f(q); q == DATRIE_HEAD_LOC {
				break
			}
		}
		
		q = this.g(q, m)
		if q == -1 {
			q = DATRIE_HEAD_LOC
			continue
		}
		
		if o := this.out(q); len(o) > 0 {
//			fmt.Println("byte: ", m, ", out: ", o)
			out = append(out, o...)
		}
	}
	
	return out
}

func (this *ACDATrie) Len() int { return this.nodeStruct.Len() }
