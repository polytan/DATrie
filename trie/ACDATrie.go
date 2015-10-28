package trie

import (
	"fmt"
)

type _ACDA struct {
	failure	int
	output	[]string
}

type ACDATrie struct {
	data	*DATrie
	payload	[]_ACDA
}

func NewACDATrie(reorgSize, expBCSize int) *ACDATrie {
	return &ACDATrie{
		data: NewDATrie(reorgSize, expBCSize),
		payload: []_ACDA{},
	}
}

func (this *ACDATrie) reset() {
	this.data.reset()
	this.payload = []_ACDA{}
}

func (this *ACDATrie) g(s int, a byte) int {
	b := this.data.array.getBase(s)
	
	n := b + int(a)
	c := this.data.array.getCheck(n)
	if s == abs(c) {
		return n
	}
	
	return -1
}

func (this *ACDATrie) f(s int) int {
	i := this.data.array.getPayload(s)
	return this.payload[i].failure
}

func (this *ACDATrie) out(s int) []string {
	i := this.data.array.getPayload(s)
	return this.payload[i].output
}

func (this *ACDATrie) sf(s int, t int) {
	i := this.data.array.getPayload(s)
	this.payload[i].failure = t
}

func (this *ACDATrie) BuildTrie(arr []string) {
	this.reset()
	
	//Phrase 1: build trie
	//step 1. construct trie
	this.data.BuildFromStrings(arr)	//build double-array trie
	fmt.Println("Phase 1: step 1. finished!")
	
	//step 2. set initial output table
	this.payload = make([]_ACDA, this.data.states + 1)
	for i := 0; i < len(this.payload); i ++ {
		p := &this.payload[i]			//
		p.failure = DATRIE_HEAD_LOC
		p.output = []string{}
	}
	
	for _, str := range arr {
		n := this.data.prefix(str)
		i := this.data.array.getPayload(n)
		
		p := &this.payload[i]
		p.output = append(p.output, str)
	}
	fmt.Println("Phase 1: step 2. finished!")
	
	//step 3. complete the goto function for root
	//this step is not required, since our map restrict to only exist nodes
	
	//Phase 2: fill failure table
	Q := []int{}
	
	var b int = this.data.array.getBase(DATRIE_HEAD_LOC)
	var i byte
	for i = 1; i != 0; i ++ {
		q := b + int(i)
		if q >= len(this.data.array.array) {
			continue
		}
		
		if q == DATRIE_HEAD_LOC {
			continue
		}
		
		c := this.data.array.getCheck(q)
		if DATRIE_HEAD_LOC == abs(c) {
			Q = append(Q, q)
		}
	}
	fmt.Println("Phrase 2: started!")
	
	for ; len(Q) > 0; {
		r := Q[0]
		Q = Q[1:]
		
		b = this.data.array.getBase(r)
		
		var a byte
		for a = 1; a != 0; a ++ {
			u := b + int(a)
			if u >= len(this.data.array.array) {
				continue
			}
			
			if u == DATRIE_HEAD_LOC {
				continue
			}
			
			c := this.data.array.getCheck(u)
			if r != abs(c) {
				continue
			}
			
			Q = append(Q, u)
			v := this.f(r)
			
//			fmt.Println("u: ", u)
			for {
				if tmp := this.g(v, a); tmp != -1 {
					break
				}
				
				if v = this.f(v); v == DATRIE_HEAD_LOC {
					break
				}
				
//				fmt.Println("v: ", v)
			}
			
			if va := this.g(v, a); va == -1 && v == DATRIE_HEAD_LOC {
				this.sf(u, DATRIE_HEAD_LOC)
			} else {
				this.sf(u, va)
			}
			
			n := this.data.array.getPayload(u)
			p := &this.payload[n]
			p.output = append(p.output, this.out(this.f(u))...)
		}
	}
	fmt.Println("Phrase 2: finished!")
}

func (this *ACDATrie) SearchTrie(str string) []string {
	out := []string{}
	
	q := DATRIE_HEAD_LOC
	for _, m := range []byte(str) {
		for {
			if this.g(q, m) != -1 {
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

func (this *ACDATrie) Len() int {
	return this.data.Len()
}
