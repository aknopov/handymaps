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
Implementation of a read-through cache where entries expire after a certain period of time.  The code is thread-safe.
The default implementation does not expire entries and has unlimited capacity. Both parameters can be customized to reach full functionality.
It requires a user-defined load function to provide a value based on a key. Expired entries are removed asynchronously.
Implementing the `Listener` interface allows tracking of map events such as adding, removing, peeking, missing entries and load failures.
Example of use:
```go
import "github.com/aknopov/handymaps/sorted"

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
assert(expiryMap.ContainsKey("Hello"))
assert(expiryMap.Len() == 2)

time.Sleep(100 * time.Millisecond)
assert(expiryMap.Len() == 0)
```
