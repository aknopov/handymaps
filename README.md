![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/aknopov/handymaps/go.yml)
![Coveralls](https://img.shields.io/coverallsCoverage/github/aknopov/handymaps)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/google.golang.org/handymaps.svg)](https://pkg.go.dev/github.com/aknopov/handymaps)

# Handy Maps

Library contains a number of map collections extending standard implementation

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

## ExpiryMap

ExpiryMap is a read-through cache where entries expire after a certain period of time. Implementation is thread-safe with atomic methods of interface.
See more detals in [ExpiryMap documentation](./pkg/expiry/EXPIRYMAP.md)

Example:
```go
import "github.com/aknopov/handymaps/pkg/expiry"

expiryMap := expiry.NewExpiryMap[string, int]().
    WithLoader(func(key string) (int, error) { return len(key), nil }).
    WithMaxCapacity(2).
    ExpireAfter(50 * time.Millisecond)

val, _ = expiryMap.Get("Hi")
assert(val == 2)
val, _ = expiryMap.Get("Hello")
assert(val == 5)
assert(expiryMap.Len() == 2)
val, _ = expiryMap.Get("World!")
assert(val == 6)
assert(expiryMap.Len() == 2)
assert(expiryMap.ContainsKey("Hello"))
assert(!expiryMap.ContainsKey("Hi"))

time.Sleep(100 * time.Millisecond)
assert(expiryMap.Len() == 0)
```
