// Package "expiry" implements ExpiryMap - a read-through cache which entries expire after certain period of time.
package expiry

import (
	"errors"
	"time"
)

// Implementation of a map which entries expire after certain time.
type ExpiryMap[K comparable, V any] struct {
	backMap     map[K]V
	maxCapacity int
	ttl         time.Duration
	loader      func(key K) (V, error)
	listeners   []Listener[K, V]
	zeroVal     V
	Once
}

type EventType int

// ExpiryMap events
const (
	Added EventType = iota
	Expired
	Requested
	Updated
	Failed
	Missed
)

// Listener interface to ExpiryMap events
type Listener[K comparable, V any] interface {
	// Function to receive on each ExpiryMap event
	//   - ev - event type
	//   - key - key assiosciated with the event
	//   - val - the associated value
	//   - err - optional error on failure
	Listen(ev EventType, key K, val V, err error)
}

// Creates ExpiryMap with default field values - unlimited capacity without entries expiry
func NewExpiryMap[K comparable, V any]() *ExpiryMap[K, V] {
	var deflt V
	ret := ExpiryMap[K, V]{
		backMap:     make(map[K]V),
		maxCapacity: -1,
		ttl:         1<<63 - 1,
		loader:      func(key K) (V, error) { return deflt, errors.New("loader not defined") },
		listeners:   make([]Listener[K, V], 0),
	}
	return &ret
}

// Modifies max capacity of the map. If adding new entry exceeds map capacity, the oldest entry is evicted.
func (em *ExpiryMap[K, V]) WithMaxCapacity(maxCapacity int) *ExpiryMap[K, V] {
	em.DoAtomically(func() {
		em.maxCapacity = maxCapacity
	})
	return em
}

// Modifes map entries time-to-live period
func (em *ExpiryMap[K, V]) Expirefter(ttl time.Duration) *ExpiryMap[K, V] {
	em.DoAtomically(func() {
		em.ttl = ttl
	})
	return em
}

// Modifes map's loader that provides values for a new  key
func (em *ExpiryMap[K, V]) WithLoader(loader func(key K) (V, error)) *ExpiryMap[K, V] {
	em.DoAtomically(func() {
		em.loader = loader
	})
	return em
}

// Adds listener to ExpiryMap events. The listeners are executed in a synchronous mode in order of their insretion.
func (em *ExpiryMap[K, V]) AddListener(listener Listener[K, V]) *ExpiryMap[K, V] {
	em.DoAtomically(func() {
		em.listeners = append(em.listeners, listener)
	})
	return em
}

func (em *ExpiryMap[K, V]) RemoveListener(listener Listener[K, V]) *ExpiryMap[K, V] {
	em.DoAtomically(func() {
		for i, l := range em.listeners {
			if l == listener {
				em.listeners = append(em.listeners[:i], em.listeners[i+1:]...)
				break
			}
		}
	})
	return em
}

// Returns map capacity
func (em *ExpiryMap[K, V]) Capacity() int {
	return em.maxCapacity
}

// Returns expiry period
func (em *ExpiryMap[K, V]) ExpiringAfter() time.Duration {
	return em.ttl
}

// Returns length of the map
func (em *ExpiryMap[K, V]) Len() int {
	return len(em.backMap)
}

//-----------------

func (em *ExpiryMap[K, V]) notifyListeners(ev EventType, key K, val V, err error) {
	for _, l := range em.listeners {
		l.Listen(ev, key, val, err)
	}
}

// Returns a value associated with the given key. It can invoke `load` function if entry is not present in the map.
func (em *ExpiryMap[K, V]) Get(key K) (V, error) {
	var err error
	var val V
	em.DoAtomically(func() {
		if _, ok := em.backMap[key]; !ok {
			val, err = em.loader(key)
			if err == nil {
				em.backMap[key] = val
				em.notifyListeners(Added, key, val, nil)
			} else {
				em.notifyListeners(Failed, key, em.zeroVal, err)
			}
		} else {
			val = em.backMap[key]
			em.notifyListeners(Requested, key, val, nil)
		}
	})
	return val, err
}

// Returns the value associated to the given key. In contrast to `Get()` this method does not trigger the loader.
func (em *ExpiryMap[K, V]) Peek(key K) (V, bool) {
	val, ok := em.backMap[key]
	if ok {
		em.notifyListeners(Requested, key, val, nil)
	} else {
		em.notifyListeners(Missed, key, em.zeroVal, nil)
	}
	return val, ok
}

// Returns `true`, if there is a mapping for the specified key.
func (em *ExpiryMap[K, V]) ContainsKey(key K) bool {
	_, ok := em.backMap[key]
	return ok
}

// Replaces synchronously the entry for a key only if currently mapped to some value.
//
//   - return `true` if value was replaced
func (em *ExpiryMap[K, V]) Replace(key K, val V) bool {
	var ok bool
	em.DoAtomically(func() {
		if _, ok = em.backMap[key]; ok {
			em.backMap[key] = val
			em.notifyListeners(Updated, key, val, nil)
		}
	})
	return ok
}

// Removes the mapping for a key from the cache if it is present.
//
//   - return `true` if value was removed
func (em *ExpiryMap[K, V]) Remove(key K, val V) bool {
	var ok bool
	em.DoAtomically(func() {
		if _, ok = em.backMap[key]; ok {
			delete(em.backMap, key)
		}
	})
	return ok
}
