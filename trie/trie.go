package trie

import "fmt"

type _T struct {
	number		int		//the node number
	chains		int		//how many chains after this node
	terminal	bool	//whether this node is a word
	payload		interface{}		//the payload
	childrens	map[byte]*_T	//children nodes
}

type Trie struct {
	root		*_T
	sequence	int
}

func newT(seq int) *_T {
	return &_T{
		number: seq,
		chains: 0,
		terminal: false,
		payload: nil,
		childrens: map[byte]*_T{},
	}
}

func NewTrie() *Trie {
	return &Trie{
		root: newT(0),
		sequence: 0,
	}
}

func BuildTrie(arr []string) *Trie {
	t := NewTrie()
	
	for _, str := range arr {
		t.Add(str)
	}
	
	return t
}

type visitor func(*_T) bool
type walker func(*_T) *_T

func (t *_T) debug() {
	fmt.Printf("node[%d], chains:%d, isword:%v, childrens: %d\n",
							t.number, t.chains, t.terminal, len(t.childrens))
}
func (t *_T) isLeaf() bool { return len(t.childrens) == 0 }
func (t *_T) isWord() bool { return t.terminal }
func (t *_T) words() int { return t.chains }


func (this *Trie) next() int {
	this.sequence ++
	return this.sequence
}

func deepVisit(n *_T, v visitor) bool {
	if !v(n) {
		return false
	}
	
	if n.isLeaf() {
		return true
	}
	
	for _, c := range n.childrens {
		if !deepVisit(c, v) {
			return false
		}
	}
	
	return true
}

func (this *Trie) DeepVisit(v visitor) bool {
	return deepVisit(this.root, v)
}

func (this *Trie) WideVisit(v visitor) bool {
	for queue := []*_T{this.root}; len(queue) > 0; queue = queue[1:] {
		n := queue[0]
		
		if !v(n) {
			return false
		}
		
		for _, c := range n.childrens {
			queue = append(queue, c)
		}
	}
	
	return true
}

func (this *Trie) Walk(w walker) {
	curr := this.root
	
	for {
		if curr == nil {
			break
		}
		curr = w(curr)
	}
}

func (this *Trie) Len() int {
	return this.root.chains
}

func (this *Trie) Add(str string) bool {
	if this.Search(str) {	//exist
		return true
	}
	
	curr := this.root
	curr.chains ++	//root add words
	for _, b := range []byte(str) {
		n, ok := curr.childrens[b]
		if !ok {
			n = newT(this.next())
			curr.childrens[b] = n
		}
		
		curr = n		
		curr.chains ++
	}
	curr.terminal = true
	return true
}

func (this *Trie) Del(str string) bool {
	if !this.Search(str) {	//not exist
		return true
	}
	
	curr := this.root
	curr.chains --
	for _, b := range []byte(str) {
		n, _ := curr.childrens[b]	//don't check the existance, it must exist
		if n.chains == 1 {	//the n have only this children
			delete(curr.childrens, b)
		}
		
		curr = n
		curr.chains --
	}
	curr.terminal = false
	return true
}

func (this *Trie) Search(str string) bool {
	curr := this.root
	for _, b := range []byte(str) {
		n, ok := curr.childrens[b]
		if !ok {
			return false
		}
		
		curr = n
	}
	return curr.terminal
}

func (this *Trie) searchNode(str string) *_T {
	curr := this.root
	for _, b := range []byte(str) {
		n, ok := curr.childrens[b]
		if !ok {
			return nil
		}
		
		curr = n
	}

	return curr
}

func (this *Trie) SearchDebug(str string) bool {
	curr := this.root
	
	fmt.Printf("nil ==> ")
	curr.debug()
	for _, b := range []byte(str) {
		n, ok := curr.childrens[b]
		if !ok {
			return false
		}
		
		curr = n
		
		fmt.Printf("%d ==> ", b)
		n.debug()
	}
	return curr.terminal
}

func (this *Trie) Tokenize(str string) []string {
	return []string{}
}

