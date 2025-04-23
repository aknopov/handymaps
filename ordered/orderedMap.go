// Package "ordered" implements a map with iteration order following insertion order.
// Iteration order is not affected if a key is inserted repeatedly into the map.
package ordered

// Map implementation
type OrderedMap[K comparable, V any] struct {
	backMap     map[K]V
	orderedKeys []K
	zeroVal     V
}

// OrderedMap iterator
type OrderedMapIterator[K comparable, V any] struct {
	om  *OrderedMap[K, V]
	idx int
}

// Creates an ordered map with the specified capacity.
//   - capacity - initial capacity
func NewOrderedMapEx[K comparable, V any](capacity int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		backMap:     make(map[K]V, capacity),
		orderedKeys: make([]K, 0, capacity),
	}
}

// Creates a zero-sized ordered map.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return NewOrderedMapEx[K, V](0)
}

// Returns the length of the map.
func (om *OrderedMap[K, V]) Len() int {
	return len(om.orderedKeys)
}

// Returns the value for the specified key. If the key isn't present, returns false in the second return value.
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, ok := om.backMap[key]
	return value, ok
}

// Associates the specified value with the specified key.
func (om *OrderedMap[K, V]) Put(key K, value V) {
	if _, ok := om.backMap[key]; !ok {
		om.orderedKeys = append(om.orderedKeys, key)
	}
	om.backMap[key] = value
}

// Copies all of the mappings from the specified map to this map.
func (om *OrderedMap[K, V]) PutAll(other *OrderedMap[K, V]) {
	for _, key := range other.orderedKeys {
		om.Put(key, other.backMap[key])
	}
}

// Removes the mapping for the specified key from this map if present.
//   - returns `true` if the value was removed
func (om *OrderedMap[K, V]) Remove(key K) bool {
	var ok bool
	if _, ok = om.backMap[key]; ok {
		delete(om.backMap, key)
		for i, k := range om.orderedKeys {
			if k == key {
				om.orderedKeys = append(om.orderedKeys[:i], om.orderedKeys[i+1:]...)
				break
			}
		}
	}
	return ok
}

// Computes the value for the specified key. If the key is not present, the compute function receives the "zero" value.
func (om *OrderedMap[K, V]) Compute(key K, compute func(K, V) V) V {
	value, ok := om.backMap[key]
	if !ok {
		value = om.zeroVal
	}
	value = compute(key, value)
	om.Put(key, value)
	return value
}

// Returns a list of the map keys in the order they were inserted.
func (om *OrderedMap[K, V]) Keys() []K {
	return om.orderedKeys
}

// Creates an iterator for the map.
func (om *OrderedMap[K, V]) Iterator() *OrderedMapIterator[K, V] {
	return &OrderedMapIterator[K, V]{om: om, idx: 0}
}

// Checks if there are more elements to iterate.
func (it *OrderedMapIterator[K, V]) HasNext() bool {
	return it.idx < len(it.om.orderedKeys)
}

// Returns the next key-value pair.
func (it *OrderedMapIterator[K, V]) Next() (k K, v V) {
	key := it.om.orderedKeys[it.idx]
	value := it.om.backMap[key]
	it.idx++
	return key, value
}
