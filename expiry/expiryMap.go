package expiry

import (
	"github.com/aknopov/handymaps/ordered"
)

func (em *ExpiryMap[K, V]) notifyListeners(ev EventType, key K, val V, err error) {
	for f := range em.listeners.m {
		f.Listen(ev, key, val, err)
	}
}

func (em *ExpiryMap[K, V]) removeOldest() {
	it := em.backMap.Iterator()
	if it.HasNext() {
		key, val := it.Next()
		em.backMap.Remove(key)
		em.notifyListeners(Removed, key, val, nil)
	}
}

// Returns a value associated with the given key. It can invoke `load` function if entry is not present in the map.
func (em *ExpiryMap[K, V]) Get(key K) (V, error) {
	var err error
	var val V
	var ok bool
	em.doAtomically(func() {
		if val, ok = em.backMap.Get(key); !ok {
			val, err = em.loader(key)
			if err == nil {
				for em.backMap.Len() >= em.maxCapacity {
					em.removeOldest()
				}
				em.backMap.Put(key, val)
				em.notifyListeners(Added, key, val, nil)
				// UC What to do with timer job?
			} else {
				em.notifyListeners(Failed, key, em.zeroVal, err)
			}
		} else {
			em.notifyListeners(Requested, key, val, nil)
		}
	})
	return val, err
}

// Returns the value associated to the given key. In contrast to `Get()` this method does not trigger the loader.
func (em *ExpiryMap[K, V]) Peek(key K) (V, bool) {
	val, ok := em.backMap.Get(key)
	if ok {
		em.notifyListeners(Requested, key, val, nil)
	}
	return val, ok
}

// Returns `true`, if there is a mapping for the specified key.
func (em *ExpiryMap[K, V]) ContainsKey(key K) bool {
	_, ok := em.backMap.Get(key)
	return ok
}

// Replaces synchronously the entry for a key only if currently mapped to some value.
//
//   - return `true` if value was replaced
func (em *ExpiryMap[K, V]) Replace(key K, val V) bool {
	var ok bool
	em.doAtomically(func() {
		if _, ok = em.backMap.Get(key); ok {
			em.backMap.Put(key, val)
			em.notifyListeners(Replaced, key, val, nil)
			// UC What to do with timer job?
		}
	})
	return ok
}

// Removes the mapping for a key from the cache if it is present.
//
//   - return `true` if value was removed
func (em *ExpiryMap[K, V]) Remove(key K) bool {
	var ok bool
	em.doAtomically(func() {
		ok = em.backMap.Remove(key)
		// UC What to do with timer job?
	})
	return ok
}

// Clears the cache.
func (em *ExpiryMap[K, V]) Clear() {
	em.doAtomically(func() {
		em.backMap = *ordered.NewOrderedMap[K, V]()
		// UC What to do with timer jobs?
	})
}
