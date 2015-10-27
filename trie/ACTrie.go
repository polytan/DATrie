package trie

import (
	"fmt"
	"sort"
)

type _AC struct {
	failure	*_T
	output	[]string
}

type ACTrie struct {
	data	*Trie
	payload	[]_AC
}

func NewACTrie() *ACTrie {
	return &ACTrie{
		data: NewTrie(),
		payload: []_AC{},
	}
}

func (this *ACTrie) reset() {
	this.data = NewTrie()
	this.payload = []_AC{}
}

func (_ *ACTrie) g(s *_T, a byte) *_T {
	if s == nil {
		fmt.Println("s is nil, a is ", a)
	}
	if c, ok := s.childrens[a]; ok {
		return c
	}
	
	return nil
}

func (this *ACTrie) f(s *_T) *_T {
	return this.payload[s.number].failure
}

func (this *ACTrie) out(s *_T) []string {
	return this.payload[s.number].output
}

func (this *ACTrie) sf(s *_T, t *_T) {
	this.payload[s.number].failure = t
}


func (this *ACTrie) BuildTrie(arr []string) {
	this.reset()
	
	sort.Strings(arr)
	
	//Phase 1: build trie
	//step 1. construct trie
	for _, str := range arr {
		this.data.Add(str)
	}
	fmt.Println("Phase 1: step 1. finished!")
	
	//step 2. set initial output table
	this.payload = make([]_AC, this.data.sequence + 1)
	for i := 0; i <= this.data.sequence; i ++ {
		p := &this.payload[i]	//we need to pick reference, can not pick value
		p.failure = this.data.root	//all node point to root
		p.output = []string{}
	}

	for _, str := range arr {
		n := this.data.searchNode(str)
		
		p := &this.payload[n.number]
		p.output = append(p.output, str)
	}
	fmt.Println("Phase 1: step 2. finished!")
	
	//step 3. complete the goto function for root
	//this step is not required, since our map restrict to only exist nodes

	//Phase 2: fill failure table
	Q := []*_T{}
	for _, q := range this.data.root.childrens {
		Q = append(Q, q)
	}
	fmt.Println("Phrase 2: started!")
	
	for ; len(Q) > 0; {
		r := Q[0]
		Q = Q[1:]
		
		for a, u := range r.childrens {
			Q = append(Q, u)
			v := this.f(r)
			
			for {
				tmp := this.g(v, a)
				if tmp != nil {
					break
				}
				
				if v = this.f(v); v == this.data.root {
					break
				}
			}
			
			if va := this.g(v, a); va == nil && v == this.data.root {
				this.sf(u, this.data.root)
			} else {
				this.sf(u, va)
			}
			
			p := &this.payload[u.number]
			p.output = append(p.output, this.out(this.f(u))...)
		}
	}
	fmt.Println("Phrase 2: finished!")
}

func (this *ACTrie) SearchTrie(str string) []string {
	out := []string{}

	q := this.data.root
	for _, m := range []byte(str) {
		for {
			if this.g(q, m) != nil {
				break
			}
			q = this.f(q)
		}
		
		q = this.g(q, m)

		if o := this.out(q); len(o) > 0 {
//			fmt.Println("byte: ", m, ", out: ", o)
			out = append(out, o...)
		}
	}
	
	return out
}

func (this *ACTrie) Add(str string) bool {
	panic("Add does not supported!")
}

func (this *ACTrie) Del(str string) bool {
	panic("Del does not supported!")
}

func (this *ACTrie) Len() int {
	return this.data.Len()
}

func (this *ACTrie) Search(str string) bool {
	return this.data.Search(str)
}

func (this *ACTrie) SearchDebug(str string) bool {
	return this.data.SearchDebug(str)
}

func (this *ACTrie) Tokenize(str string) []string {
	return this.SearchTrie(str)
}