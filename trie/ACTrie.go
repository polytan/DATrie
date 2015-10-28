package trie

import (
//	"fmt"
	"sort"
)

type _AC struct {
	failure	*_T
	output	[]string
}

type ACTrie struct {
	nodeStruct	*Trie
	nodeValue	[]_AC
}

func NewACTrie() *ACTrie {
	return &ACTrie{
		nodeStruct: NewTrie(),
		nodeValue: []_AC{},
	}
}

func (this *ACTrie) reset() {
	this.nodeStruct = NewTrie()
	this.nodeValue = []_AC{}
}

func (_ *ACTrie) g(s *_T, a byte) *_T {
	if c, ok := s.childrens[a]; ok {
		return c
	}
	
	return nil
}

func (this *ACTrie) f(s *_T) *_T {
	return this.nodeValue[s.number].failure
}

func (this *ACTrie) out(s *_T) []string {
	return this.nodeValue[s.number].output
}

func (this *ACTrie) sf(s *_T, t *_T) {
	this.nodeValue[s.number].failure = t
}


func (this *ACTrie) BuildTrie(arr []string) {
//	fmt.Println("===== START BUILD ACTRIE =====")
	this.reset()
	
	sort.Strings(arr)
	
//	fmt.Println("Phase 1: build trie")
//	fmt.Println("step 1. construt trie")
	this.nodeStruct = BuildTrie(arr)
	
//	fmt.Println("step 2. set initial output and failure table")
	this.nodeValue = make([]_AC, this.nodeStruct.sequence + 1)
	for i := 0; i < len(this.nodeValue); i ++ {
		p := &this.nodeValue[i]				//pick reference instead of value
		p.failure = this.nodeStruct.root	//all node point to root
		p.output = []string{}
	}

//	fmt.Println("step 3. set output string set")
	for _, str := range arr {
		n := this.nodeStruct.searchNode(str)
		
		p := &this.nodeValue[n.number]
		p.output = append(p.output, str)
	}
	
//	fmt.Println("step 4. complete the goto function for root, not required")

//	fmt.Println("Phase 2: fill failure table")
//	fmt.Println("step 1. queue push initial items")
	Q := []*_T{}
	for _, q := range this.nodeStruct.root.childrens {
		Q = append(Q, q)
	}
	
//	fmt.Println("step 2. queue process to fill failure table")
	for ; len(Q) > 0; {
		r := Q[0]
		Q = Q[1:]
		
		for a, u := range r.childrens {
			Q = append(Q, u)
			v := this.f(r)
			
			for {
				if tmp := this.g(v, a); tmp != nil {
					break
				}
				
				if v = this.f(v); v == this.nodeStruct.root {
					break
				}
			}
			
			if va := this.g(v, a); va == nil && v == this.nodeStruct.root {
				this.sf(u, this.nodeStruct.root)
			} else {
				this.sf(u, va)
			}
			
			p := &this.nodeValue[u.number]
			p.output = append(p.output, this.out(this.f(u))...)
		}
	}
//	fmt.Println("===== END BUILD ACTRIE =====")
}

func (this *ACTrie) SearchTrie(str string) []string {
	out := []string{}

	q := this.nodeStruct.root
	for _, m := range []byte(str) {
		for {
			if this.g(q, m) != nil {
				break
			}

			if q = this.f(q); q == this.nodeStruct.root {
				break
			}
		}
		
		q = this.g(q, m)
		if q == nil {
			q = this.nodeStruct.root
			continue
		}

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
	return this.nodeStruct.Len()
}

func (this *ACTrie) Search(str string) bool {
	return this.nodeStruct.Search(str)
}

func (this *ACTrie) SearchDebug(str string) bool {
	return this.nodeStruct.SearchDebug(str)
}

func (this *ACTrie) Tokenize(str string) []string {
	return this.SearchTrie(str)
}