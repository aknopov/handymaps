// Package "expiry" implements ExpiryMap - a read-through cache which entries expire after certain period of time.
package expiry

import (
	"errors"
	"time"

	"github.com/aknopov/handymaps/pkg/ordered"
)

const (
	// Expiry value at which eviction is not performed
	Eternity  = 1<<63 - 1
	Unlimited = -1
)

type entry[V any] struct {
	val    V
	exptmr *time.Timer
}

// Implementation of a map which entries expire after certain time.
type ExpiryMap[K comparable, V any] struct {
	backMap     ordered.OrderedMap[K, entry[V]]
	maxCapacity int
	ttl         time.Duration
	loader      func(key K) (V, error)
	listeners   *set[Listener[K, V]]
	evictChan   chan K
	stopChan    chan bool
	upgradableRWMutex
}

type EventType int

// ExpiryMap events
const (
	// added by loader
	Added EventType = iota
	// expired and removed
	Expired
	// removed to ensure capacity
	Removed
	// requested by Get or Peek without invoking loader
	Requested
	// Peek didn't yield
	Missed
	// replaced
	Replaced
	// failed load
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
		backMap:     *ordered.NewOrderedMap[K, entry[V]](),
		maxCapacity: Unlimited,
		ttl:         Eternity,
		loader:      func(key K) (V, error) { return deflt, errors.New("loader not defined") },
		listeners:   newSet[Listener[K, V]](),
		evictChan:   make(chan K),
		stopChan:    make(chan bool),
	}

	go func() {
		for {
			select {
			case key := <-ret.evictChan:
				ret.writeAtomically(func() {
					ret.removeEntry(key)
				})
			case <-ret.stopChan:
				return
			}
		}
	}()

	return &ret
}

// Modifies max capacity of the map. If adding new entry exceeds map capacity, the oldest entry is evicted.
func (em *ExpiryMap[K, V]) WithMaxCapacity(maxCapacity int) *ExpiryMap[K, V] {
	em.maxCapacity = maxCapacity
	return em
}

// Modifes map entries time-to-live period
func (em *ExpiryMap[K, V]) ExpireAfter(ttl time.Duration) *ExpiryMap[K, V] {
	em.ttl = ttl
	return em
}

// Modifes map's loader that provides values for a new  key
func (em *ExpiryMap[K, V]) WithLoader(loader func(key K) (V, error)) *ExpiryMap[K, V] {
	em.loader = loader
	return em
}

// Returns map capacity
func (em *ExpiryMap[K, V]) Capacity() int {
	return em.maxCapacity
}

// Returns expiry period
func (em *ExpiryMap[K, V]) ExpireTime() time.Duration {
	return em.ttl
}

// Returns length of the map
func (em *ExpiryMap[K, V]) Len() int {
	var size int
	em.readAtomically(func() {
		size = em.backMap.Len()
	})
	return size
}
