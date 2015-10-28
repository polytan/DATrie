package main

import (
	"os"
	"io"
	"bufio"
	"fmt"
	"strings"
	"time"
//	"runtime/pprof"
	"github.com/polytan/DATrie/trie"
)

func main() {
	input := "out.txt"
	search := "中华人民共和国"
//	search := "aabaacaadaaeaaf"
//	bc1 := loadTrie(input)
//	if bc1 == nil {
//		panic("Error while createing trie")
//	}
	
//	bc1.Search(search)
	
	bc2 := loadDATrie(input)
	if bc2 == nil {
		panic("Error while creating DATrie")
	}
	
	bc2.Search(search)
	
//	bc3 := loadACTrie(input)
//	if bc3 == nil {
//		panic("Error while creating ACTrie")
//	}
	
//	list3 := bc3.SearchTrie(search)
//	fmt.Println(list3)
	
//	bc4 := loadACDATrie(input)
//	if bc4 == nil {
//		panic("Error while creating ACDATrie")
//	}
	
//	list4 := bc4.SearchTrie(search)
//	fmt.Println(list4)

//	runtime.GC()
	
//	memFile, _ := os.Create("memory.log")
//	pprof.WriteHeapProfile(memFile)
	
}

func testSearch(t *trie.ACTrie, str string) {
	r := t.Search(str)
	fmt.Println("Search word ", str, ", result: ", r)
}

func readFile(file string) []string {
	start := time.Now()
	f, err := os.Open(file)
	if err != nil {
		return []string{}
	}
	defer f.Close()
	buf := bufio.NewReader(f)

	arr := []string{}
	for {
		line, err := buf.ReadString(byte('\n'))
		if err != nil && err != io.EOF {
			return []string{}
		}
		
		line = strings.Trim(line, "\r\n\t ")
		if len(line) > 0 {
			arr = append(arr, line)
		}
		
		if err == io.EOF {
			break
		}
	}
	
	end := time.Now()
	fmt.Println("Load file: ", end.Sub(start))
	
	return arr
}

func loadTrie(file string) *trie.Trie {
	var start, end time.Time
	
	arr := readFile(file)

	start = time.Now()
	t := trie.NewTrie()
	for _, str := range arr {
		if t.Search(str) {
			fmt.Println("Already exist item in trie: ", str)
		}
		t.Add(str)
	}
	end = time.Now()
	fmt.Println("Build trie: ", end.Sub(start))
	fmt.Println("Words in trie: ", t.Len())
	
	//check trie
	start = time.Now()
	for i := 0; i < 10; i ++ {
		for _, str := range arr {
			if !t.Search(str) {
				panic("fatal error while testing: " + str)
			}
		}
	}
	end = time.Now()
	fmt.Println("Validate trie: ", end.Sub(start))
	
	return t
}

func loadDATrie(file string) *trie.DATrie {
	var start, end time.Time
	
	arr := readFile(file)
	
	start = time.Now()
	//reorg when trie has so much nodes, expande bcarray to this size when no enough space
	t := trie.NewDATrie(100000, 102400)
//	for _, str := range arr {
//		t.Add(str)
//	}
//	t.Build()
	t.BuildFromStrings(arr)	//build one time
	end = time.Now()
	fmt.Println("Build DATrie: ", end.Sub(start))
	fmt.Println("Words in DATrie: ", t.Len())
	
	start = time.Now()
	for i := 0; i < 10; i ++ {
		for _, str := range arr {
			if !t.Search(str) {
				fmt.Println("fatal error while testing: " + str)
			}
		}
	}
	end = time.Now()
	fmt.Println("Validate trie: ", end.Sub(start))
	
	return t
}

func loadACTrie(file string) *trie.ACTrie {
	var start, end time.Time
	
	arr := readFile(file)
	
	start = time.Now()
	t := trie.NewACTrie()
	t.BuildTrie(arr)
	end = time.Now()
	fmt.Println("Build ACTrie: ", end.Sub(start))
	fmt.Println("Words in trie: ", t.Len())
	
	//check trie
	
	return t
}

func loadACDATrie(file string) *trie.ACDATrie {
	var start, end time.Time
	
	arr := readFile(file)
	
	start = time.Now()
	t := trie.NewACDATrie(100000, 102400)
	t.BuildFromArray(arr)
	end = time.Now()
	fmt.Println("Build ACDATrie: ", end.Sub(start))
	fmt.Println("Words in trie: ", t.Len())
	
	//check trie
	
	return t
	
}