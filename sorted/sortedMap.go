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
	for _, key := range other.sortedKeys { // UC binary search
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
		for i, k := range sm.sortedKeys { // UC binary search
			if k == key {
				sm.sortedKeys = append(sm.sortedKeys[:i], sm.sortedKeys[i+1:]...)
				break
			}
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
