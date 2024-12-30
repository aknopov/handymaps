# Handy Maps

repository contains a number of map collections extending standard implementation

## BiMap

Bimap is an implementation of bi-directional map with unique sets of keys and values. Iteration over map entries follows insertion order, unless duplicate keys or values were added. Example -
```go
import "github.com/aknopov/handymaps/bimap"

aBimap := bimap.NewBiMap[int, float32]()
aBimap.Put(1, 2.71828)
invMap := aBimap.Inverse()
idx, _ := invMap.GetValue(2.71828)
fmt.Printf("%d\n", idx)
```

## Ordered map

Implementation of map with predicatable iteration order that follows insertion order. Example - 
```go
import "github.com/aknopov/handymaps/ordered"

oMap := ordered.NewOrderedMapEx[rune, string](10)
for _, c := range "hello" {
    oMap.Put(c, string(c))
}
fmt.Println(oMap.Keys())
```
Output `[104 101 108 111]`

## SortedMap

Implementation of map which keys are sorted in according to comparison function. Example -
```go
import "github.com/aknopov/handymaps/sorted"

sMap := sorted.NewSortedMapEx[rune, string](10, func(a, b rune) bool { return a < b })
for _, c := range "hello" {
    sMap.Put(c, string(c))
}
fmt.Println(sMap.Keys())
```
Output - `[101 104 108 111]`
