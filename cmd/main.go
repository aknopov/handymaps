package main

import (
	"fmt"

	"github.com/aknopov/handymaps/bimap"
	"github.com/aknopov/handymaps/ordered"
	"github.com/aknopov/handymaps/sorted"
)

func main() {
	aBimap := bimap.NewBiMap[int, float32]()
	aBimap.Put(1, 2.71828)
	aBimap.Put(2, 3.14153)
	fmt.Println()
	fmt.Printf("A Bimap has %d elements\n", aBimap.Size())
	it := aBimap.Iterator()
	for it.HasNext() {
		k, v := it.Next()
		fmt.Printf("M[%v] = %v\n", k, v)
	}
	invMap := aBimap.Inverse()
	idx, _ := invMap.GetValue(2.71828)
	fmt.Printf("%d\n", idx)

	oMap := ordered.NewOrderedMapEx[rune, string](10)
	for _, c := range "hello" {
		oMap.Put(c, string(c))
	}
	fmt.Println(oMap.Keys())

	sMap := sorted.NewSortedMapEx[rune, string](10, func(a, b rune) bool { return a < b })
	for _, c := range "hello" {
		sMap.Put(c, string(c))
	}
	fmt.Println(sMap.Keys())
}
