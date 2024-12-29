// Package bimap implements bi-directional map with unique sets of keys and values.
// It preserves insertion order while iterating map entries, unless duplicate keys or values were added.
package bimap

// BiMap represents bidirectional map that has unique sets of keys and values
type BiMap[K comparable, V comparable] struct {
	keys   []K
	vals   []V
	keyIdx map[K]int
	valIdx map[V]int
	noKey  K
	noVal  V
}

type BiMapIterator[K comparable, V comparable] struct {
	biMap *BiMap[K, V]
	idx   int
}

// Creates a new zero-sized BiMap
func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return NewBiMapEx[K, V](0)
}

// Creates a new BiMap
//   - capacity - initial capacity
func NewBiMapEx[K comparable, V comparable](capacity int) *BiMap[K, V] {
	keys := make([]K, 0, capacity)
	vals := make([]V, 0, capacity)
	keyIdx := make(map[K]int, capacity)
	valIdx := make(map[V]int, capacity)
	return &BiMap[K, V]{keys: keys, vals: vals, keyIdx: keyIdx, valIdx: valIdx}
}

// Provides bi-map size
func (biMap *BiMap[K, V]) Size() int {
	return len(biMap.keys)
}

// Adds or replaces an entry in bi-map. This method has no effect if the bi-map previously had key-value mapping.
//
// Note that a call to this method could cause the size of the bimap to increase by one, stay the same, or even decrease by one.
//   - key - the entry key
//   - val - the entry value
//
// Returns this bi-map
func (biMap *BiMap[K, V]) Put(key K, val V) {
	i, okKey := biMap.keyIdx[key]
	j, okVal := biMap.valIdx[val]
	switch {
	case okVal && okKey && i == j: // NOP case
		return
	case okKey: //  new key
		oldVal := biMap.vals[i]
		delete(biMap.valIdx, oldVal)
		biMap.valIdx[val] = i
		biMap.keys[i] = key
		biMap.vals[i] = val
	case okVal: //  new value
		oldKey := biMap.keys[j]
		delete(biMap.keyIdx, oldKey)
		biMap.keyIdx[key] = j
		biMap.keys[j] = key
		biMap.vals[j] = val
	case !okKey && !okVal: // new key and value
		biMap.keyIdx[key] = len(biMap.keys)
		biMap.valIdx[val] = len(biMap.vals)
		biMap.keys = append(biMap.keys, key)
		biMap.vals = append(biMap.vals, val)
	default:
		panic("Interval BiMap error - key-value mismatch")
	}
}

// Gets value by the key
//   - key - the map key
//
// Returns found value and a flag of success
func (biMap *BiMap[K, V]) GetValue(key K) (V, bool) {
	if i, ok := biMap.keyIdx[key]; ok {
		return biMap.vals[i], true
	}
	return biMap.noVal, false
}

// Gets key by the value
//   - val - value of the matching entry
//
// Returns found key and a flag of success
func (biMap *BiMap[K, V]) GetKey(val V) (K, bool) {
	if i, ok := biMap.valIdx[val]; ok {
		return biMap.keys[i], true
	}
	return biMap.noKey, false
}

// Checks if the key is present in the map
func (biMap *BiMap[K, V]) ContainsKey(key K) bool {
	_, ok := biMap.keyIdx[key]
	return ok
}

// Checks if value is present in the map
func (biMap *BiMap[K, V]) ContainsValue(value V) bool {
	_, ok := biMap.valIdx[value]
	return ok
}

// Removes entry from bi-map based on a key
//   - biMap - bi-map to update
//   - key - key of the entry bo be removed
func (biMap *BiMap[K, V]) RemoveKey(key K) {
	if i, ok := biMap.keyIdx[key]; ok {
		val := biMap.vals[i]
		biMap.removeEntry(key, val, i)
	}
}

// Removes entry from bi-map based on a value
//   - biMap - bi-map to update
//   - val - value of the entry bo be removed
func (biMap *BiMap[K, V]) RemoveValue(val V) {
	if i, ok := biMap.valIdx[val]; ok {
		key := biMap.keys[i]
		biMap.removeEntry(key, val, i)
	}
}

func (biMap *BiMap[K, V]) removeEntry(key K, val V, i int) {
	newLen := len(biMap.keys) - 1
	delete(biMap.keyIdx, key)
	delete(biMap.valIdx, val)
	biMap.keys = append(biMap.keys[:i], biMap.keys[i+1:]...)
	biMap.vals = append(biMap.vals[:i], biMap.vals[i+1:]...)
	for j := i; j < newLen; j++ {
		biMap.keyIdx[biMap.keys[j]] = j
		biMap.valIdx[biMap.vals[j]] = j
	}
}

// Creates "inverse" copy of the bitmap
func (biMap *BiMap[K, V]) Inverse() *BiMap[V, K] {
	size := biMap.Size()
	invMap := BiMap[V, K]{}
	invMap.keys = append(invMap.keys, biMap.vals...)
	invMap.vals = append(invMap.vals, biMap.keys...)
	invMap.keyIdx = make(map[V]int, size)
	invMap.valIdx = make(map[K]int, size)
	for i := 0; i < size; i++ {
		invMap.keyIdx[invMap.keys[i]] = i
		invMap.valIdx[invMap.vals[i]] = i
	}
	return &invMap
}

func cmpSlices[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// Compares bi-map to the other
//   - biMap - first bi-map to compare
//   - other - second bi-map to compare
func (biMap *BiMap[K, V]) Equals(other *BiMap[K, V]) bool {
	return cmpSlices(biMap.keys, other.keys) && cmpSlices(biMap.vals, other.vals)
}

// Copies all of the mappings from another map to this
//   - biMap - bi-map to copy to
//   - other - bi-map to copy from
func (biMap *BiMap[K, V]) PutAll(other *BiMap[K, V]) *BiMap[K, V] {
	for i := range other.keys {
		biMap.Put(other.keys[i], other.vals[i])
	}
	return biMap
}

// Returns a slice of bi-map keys
func (biMap *BiMap[K, V]) Keys() []K {
	return biMap.keys
}

// Returns a slice of bi-map values
func (biMap *BiMap[K, V]) Values() []V {
	return biMap.vals
}

// Creates iterator over BiMap
func (biMap *BiMap[K, V]) Iterator() BiMapIterator[K, V] {
	return BiMapIterator[K, V]{biMap: biMap, idx: 0}
}

// Checks if iterator can provide another entry from bi-map
func (it *BiMapIterator[K, V]) HasNext() bool {
	return it.idx < it.biMap.Size()
}

// Provides next available bi-map entry
func (it *BiMapIterator[K, V]) Next() (K, V) {
	oldIdx := it.idx
	it.idx++
	return it.biMap.keys[oldIdx], it.biMap.vals[oldIdx]
}
