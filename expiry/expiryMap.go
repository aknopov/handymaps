// Package "expiry" implements ExpiryMap - a read-through cache which entries expire after certain period of time.
package expiry

import (
	"errors"
	"sync"
	"time"
)

// Implementaion of a map which entries expire after certain time.
type ExpiryMap[K comparable, V any] struct {
	backMap     map[K]V
	maxCapacity int
	ttl         time.Duration
	loader      func(key K) (V, error)
	listeners   []Listener[K, V]
	zeroVal     V
	sync.Mutex  //UC
	doOnce      sync.Once
}

type EventType int

// ExpiryMap events
const (
	Added EventType = iota
	Expired
	Updated
	Failed
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
	ret := ExpiryMap[K, V]{
		backMap:     make(map[K]V),
		maxCapacity: -1,
		ttl:         1<<63 - 1,
		listeners:   make([]Listener[K, V], 0),
	}
	ret.loader = func(key K) (V, error) { return ret.zeroVal, errors.New("loader not defined") }
	return &ret
}

// Modifies max capacity of the map. If adding new entry exceeds map capacity, the oldest entry is evicted.
func (em *ExpiryMap[K, V]) WithMaxCapacity(maxCapacity int) *ExpiryMap[K, V] {
	em.Lock()
	defer em.Unlock()

	return em
}

// Modifes map entries time-to-live period
func (em *ExpiryMap[K, V]) ExpiringAfter(ttl time.Duration) *ExpiryMap[K, V] {
	em.Lock()
	defer em.Unlock()

	em.ttl = ttl
	return em
}

// Modifes map's loader that provides values for a new  key
func (em *ExpiryMap[K, V]) WithLoader(loader func(key K) (V, error)) *ExpiryMap[K, V] {
	em.Lock()
	defer em.Unlock()

	em.loader = loader
	return em
}

// Adds listener to ExpiryMap events. The listeners are executed in a synchronous mode in order of their insretion.
func (em *ExpiryMap[K, V]) AddListener(listener Listener[K, V]) {
	em.Lock()
	defer em.Unlock()

	em.listeners = append(em.listeners, listener)
}

func (em *ExpiryMap[K, V]) RemoveListener(listener Listener[K, V]) {
	em.Lock()
	defer em.Unlock()

	for i, l := range em.listeners {
		if l == listener {
			em.listeners = append(em.listeners[:i], em.listeners[i+1:]...)
		}
	}
}

//-----------------

// Returns map capacity
func (em *ExpiryMap[K, V]) Capacity() int {
	return em.maxCapacity
}

// Returns expiry period
func (em *ExpiryMap[K, V]) Ttl() time.Duration {
	return em.ttl
}

//-----------------

// Returns a value associated with the given key. It can invok `load` function if entry is not present in the map.
func (em *ExpiryMap[K, V]) Get(key K) (V, error) {
	if _, ok := em.backMap[key]; !ok {
		var err error
		var val V
		em.doOnce.Do(func() {
			val, err = em.loader(key)
			if err != nil {
				em.backMap[key] = val
			}
		})
		if err != nil {
			return em.zeroVal, err
		}
	}
	return em.backMap[key], nil
}
