package main

import (
	"fmt"
	"time"

	"github.com/aknopov/handymaps/bimap"
	"github.com/aknopov/handymaps/expiry"
	"github.com/aknopov/handymaps/ordered"
	"github.com/aknopov/handymaps/sorted"
)

func main() {
	biMapDemo()

	orderedDemo()

	sortedDemo()

	expiryDemo()

	testPerformance()
}

func biMapDemo() {
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
}

func orderedDemo() {
	oMap := ordered.NewOrderedMapEx[rune, string](10)
	for _, c := range "hello" {
		oMap.Put(c, string(c))
	}

	for _, c := range "hello" {
		v, _ := oMap.Get(c)
		cS := string(c)
		assert(v == cS, "Value "+cS+" not found")
	}

	fmt.Println(oMap.Keys())
}

func sortedDemo() {
	sMap := sorted.NewSortedMapEx[rune, string](10, func(a, b rune) bool { return a < b })
	for _, c := range "hello" {
		sMap.Put(c, string(c))
	}
	fmt.Println(sMap.Keys())
}

func expiryDemo() {
	defaultImpl := expiry.NewExpiryMap[string, int]().
		WithLoader(func(key string) (int, error) { return len(key), nil })

	val, _ := defaultImpl.Get("Hi")
	assert(val == 2, "`Hi` expected in the map")
	val, _ = defaultImpl.Get("Hello")
	assert(val == 5, "`Hello` expected in the map")
	assert(defaultImpl.Len() == 2, "ExpectedLen=2 [1]")
	time.Sleep(500 * time.Millisecond)
	assert(defaultImpl.Len() == 2, "ExpectedLen=2 [2]")

	expiryMap := expiry.NewExpiryMap[string, int]().
		WithLoader(func(key string) (int, error) { return len(key), nil }).
		WithMaxCapacity(2).
		ExpireAfter(50 * time.Millisecond)

	val, _ = expiryMap.Get("Hi")
	assert(val == 2, "`Hi` expected in the map")
	val, _ = expiryMap.Get("Hello")
	assert(val == 5, "`Hello` expected in the map")
	assert(expiryMap.Len() == 2, "ExpectedLen=2 [1]")
	val, _ = expiryMap.Get("World!")
	assert(val == 6, "`World!` expected in the map")
	assert(expiryMap.ContainsKey("Hello"), "`Hello` expected in the map")
	assert(expiryMap.Len() == 2, "ExpectedLen=2 [1]")

	time.Sleep(100 * time.Millisecond)
	assert(expiryMap.Len() == 0, "ExpectedLen=0 [2]")
}

const nCount = 10000000

var randSeed = int(time.Now().UnixMilli())

func pseudoRand() int {
	randSeed = (randSeed*1103515245 + 12345) & 0x7fffffff
	return randSeed
}

func assert(check bool, message string) {
	if !check {
		panic(message)
	}
}

func testPerformance() {
	data := make(map[int]int)
	for i := 0; i < nCount; i++ {
		data[i] = pseudoRand() % nCount
	}

	start := time.Now().UnixMilli()
	for _, v := range data {
		assert(v < nCount, "Value out of range")
	}
	singleMapDuration := time.Now().UnixMilli() - start

	aMap := bimap.NewBiMapEx[int, int](nCount)

	for k, v := range data {
		aMap.Put(k, v)
	}

	start = time.Now().UnixMilli()
	it := aMap.Iterator()
	for it.HasNext() {
		_, v := it.Next()
		assert(v < nCount, "Value out of range")
	}
	biMapDuration := time.Now().UnixMilli() - start

	if biMapDuration < singleMapDuration {
		fmt.Println("BiMap is faster than a standard map!")
	} else {
		fmt.Println("BiMap is slower than a standard map :(")
	}
}
