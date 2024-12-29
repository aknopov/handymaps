package ordered

// Implementation of map with predicatable iteration order that follows insertion order.
//
// Iteration order is not affected if a key is re-inserted into the map.
type OrderedMap[K comparable, V any] struct {
	backMap     map[K]V
	orderedKeys []K
	zero        V
}

type orderedMapIterator[K comparable, V any] struct {
	om  *OrderedMap[K, V]
	idx int
}

// Creates ordered map with the specified capacity.
func NewOrderedMapEx[K comparable, V any](capacity int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		backMap:     make(map[K]V, capacity),
		orderedKeys: make([]K, 0, capacity),
	}
}

// Creates ordered map with the default capacity.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return NewOrderedMapEx[K, V](0)
}

// Returns length of the map
func (sml *OrderedMap[K, V]) Len() int {
	return len(sml.orderedKeys)
}

// Returns value for the specified key. If key isn't present, returns false inthe second return value.
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, ok := om.backMap[key]
	return value, ok
}

// Associates the specified value with the specified key
func (om *OrderedMap[K, V]) Put(key K, value V) {
	if _, ok := om.backMap[key]; !ok {
		om.orderedKeys = append(om.orderedKeys, key)
	}
	om.backMap[key] = value
}

// Computes value for the specified key. If key is not present, compute function received "zero" value.
func (om *OrderedMap[K, V]) Compute(key K, compute func(K, V) V) V {
	value, ok := om.backMap[key]
	if !ok {
		value = om.zero
	}
	value = compute(key, value)
	om.Put(key, value)
	return value
}

// Returns a list the map keys in the order they were inserted.
func (om *OrderedMap[K, V]) Keys() []K {
	return om.orderedKeys
}

// Creates oterator for the map.
func (om *OrderedMap[K, V]) Iterator() *orderedMapIterator[K, V] {
	return &orderedMapIterator[K, V]{om: om, idx: 0}
}

// Checks if there are more elements to iterate.
func (it *orderedMapIterator[K, V]) HasNext() bool {
	return it.idx < len(it.om.orderedKeys)
}

// Returns the next key-value pair.
func (it *orderedMapIterator[K, V]) Next() (K, V) {
	key := it.om.orderedKeys[it.idx]
	value := it.om.backMap[key]
	it.idx++
	return key, value
}
