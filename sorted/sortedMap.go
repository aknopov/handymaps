package sorted

import "sort"

// Map implementation with sorted keys.
type SortedMap[K comparable, V any] struct {
	backMap    map[K]V
	sortedKeys []K
	isLess     func(K, K) bool
	zeroVal    V
}

// SortedMap iterator
type SortedMapIterator[K comparable, V any] struct {
	sm  *SortedMap[K, V]
	idx int
}

// Creates a new zero-sized BiMap
func NewSortedMap[K comparable, V comparable](isLess func(K, K) bool) *SortedMap[K, V] {
	return NewSortedMapEx[K, V](0, isLess)
}

// Creates a new sorted name
//   - capacity - initial capacity
func NewSortedMapEx[K comparable, V comparable](capacity int, isLess func(K, K) bool) *SortedMap[K, V] {
	return &SortedMap[K, V]{
		backMap:    make(map[K]V, capacity),
		sortedKeys: make([]K, 0, capacity),
		isLess:     isLess,
	}
}

// Returns length of the map
func (sm *SortedMap[K, V]) Len() int {
	return len(sm.sortedKeys)
}

// Returns value for the specified key. If key isn't present, returns false inthe second return value.
func (sm *SortedMap[K, V]) Get(key K) (V, bool) {
	value, ok := sm.backMap[key]
	return value, ok
}

// Associates the specified value with the specified key.
func (sm *SortedMap[K, V]) Put(key K, value V) {
	if _, ok := sm.backMap[key]; !ok {
		sm.sortedKeys = append(sm.sortedKeys, key)
		sort.Slice(sm.sortedKeys, func(i, j int) bool {
			return sm.isLess(sm.sortedKeys[i], sm.sortedKeys[j])
		})
	}
	sm.backMap[key] = value
}

// Copies all of the mappings from the specified map to this map.
func (sm *SortedMap[K, V]) PutAll(other *SortedMap[K, V]) {
	for _, key := range other.sortedKeys {
		if _, ok := sm.backMap[key]; !ok {
			sm.sortedKeys = append(sm.sortedKeys, key)
		}
		sm.backMap[key] = other.backMap[key]
	}
	sort.Slice(sm.sortedKeys, func(i, j int) bool {
		return sm.isLess(sm.sortedKeys[i], sm.sortedKeys[j])
	})
}

// Removes the mapping for the specified key from this map if present.
func (sm *SortedMap[K, V]) Remove(key K) {
	if _, ok := sm.backMap[key]; ok {
		delete(sm.backMap, key)
		if i := binSearch(sm.sortedKeys, key, sm.isLess); i >= 0 {
			sm.sortedKeys = append(sm.sortedKeys[:i], sm.sortedKeys[i+1:]...)
		}
	}
}

// Computes value for the specified key. If key is not present, compute function received "zero" value.
func (sm *SortedMap[K, V]) Compute(key K, compute func(K, V) V) V {
	value, ok := sm.backMap[key]
	if !ok {
		value = sm.zeroVal
	}
	value = compute(key, value)
	sm.Put(key, value)
	return value
}

// Returns a list the map keys in the order they were inserted.
func (sm *SortedMap[K, V]) Keys() []K {
	return sm.sortedKeys
}

// Creates iterator for the map.
func (sm *SortedMap[K, V]) Iterator() *SortedMapIterator[K, V] {
	return &SortedMapIterator[K, V]{sm: sm, idx: 0}
}

// Checks if there are more elements to iterate.
func (it *SortedMapIterator[K, V]) HasNext() bool {
	return it.idx < len(it.sm.sortedKeys)
}

// Returns the next key-value pair.
func (it *SortedMapIterator[K, V]) Next() (k K, v V) {
	key := it.sm.sortedKeys[it.idx]
	value := it.sm.backMap[key]
	it.idx++
	return key, value
}

// Returns index of "key" in the "slice" if found. If key is not found, returns negative value X,
// such that -X-1 is the index of the largest element less than "key" (aka insertion point).
func binSearch[K comparable](slice []K, key K, isLess func(K, K) bool) int {
	low, high := 0, len(slice)
	for low < high {
		mid := (low + high) / 2
		if slice[mid] == key {
			return mid
		}
		if isLess(slice[mid], key) {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return -low - 1
}
