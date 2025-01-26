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
	listeners   *set[Listener[K, V]]
	zeroVal     V
	once
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

// Convenience wrapper for Listener interface
type ListenerWarapper struct {
	f func(ev EventType, key string, val int, err error)
}

func (lw *ListenerWarapper) Listen(ev EventType, key string, val int, err error) {
	lw.f(ev, key, val, err)
}

// Creates ExpiryMap with default field values - unlimited capacity without entries expiry
func NewExpiryMap[K comparable, V any]() *ExpiryMap[K, V] {
	var deflt V
	ret := ExpiryMap[K, V]{
		backMap:     make(map[K]V),
		maxCapacity: -1,
		ttl:         1<<63 - 1,
		loader:      func(key K) (V, error) { return deflt, errors.New("loader not defined") },
		listeners:   newSet[Listener[K, V]](),
	}
	return &ret
}

// Modifies max capacity of the map. If adding new entry exceeds map capacity, the oldest entry is evicted.
func (em *ExpiryMap[K, V]) WithMaxCapacity(maxCapacity int) *ExpiryMap[K, V] {
	em.doAtomically(func() {
		em.maxCapacity = maxCapacity
	})
	return em
}

// Modifes map entries time-to-live period
func (em *ExpiryMap[K, V]) Expirefter(ttl time.Duration) *ExpiryMap[K, V] {
	em.doAtomically(func() {
		em.ttl = ttl
	})
	return em
}

// Modifes map's loader that provides values for a new  key
func (em *ExpiryMap[K, V]) WithLoader(loader func(key K) (V, error)) *ExpiryMap[K, V] {
	em.doAtomically(func() {
		em.loader = loader
	})
	return em
}

// Adds listener to ExpiryMap events. The listeners are executed in a synchronous mode in order of their insretion.
func (em *ExpiryMap[K, V]) AddListener(listener Listener[K, V]) *ExpiryMap[K, V] {
	em.doAtomically(func() {
		em.listeners.add(listener)
	})
	return em
}

func (em *ExpiryMap[K, V]) RemoveListener(listener Listener[K, V]) *ExpiryMap[K, V] {
	em.doAtomically(func() {
		em.listeners.remove(listener)
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
