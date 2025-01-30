package expiry

import (
	"time"
)

func (em *ExpiryMap[K, V]) notifyListeners(ev EventType, key K, val V, err error) {
	for f := range em.listeners.m {
		f.Listen(ev, key, val, err)
	}
}

func (em *ExpiryMap[K, V]) removeEntry(key K) bool {
	if val, ok := em.backMap.Get(key); ok {
		val.exptmr.Stop()
		em.backMap.Remove(key)
		em.notifyListeners(Removed, key, val.val, nil)
		return true
	}
	return false
}

func (em *ExpiryMap[K, V]) removeOldest() {
	keys := em.backMap.Keys()
	if len(keys) > 0 {
		em.removeEntry(keys[0])
	}
}

func (em *ExpiryMap[K, V]) assumeAlive() {
	if em.evictChan == nil {
		panic("The map has been discarded!")
	}
}

// Removes all entries and stops eviction loop
func (em *ExpiryMap[K, V]) Discard() {
	em.assumeAlive()

	em.Clear()
	em.writeAtomically(func() {
		em.stopChan <- true
		close(em.evictChan)
		em.evictChan = nil
	})
}

// Returns a value associated with the given key. It can invoke `load` function if entry is not present in the map.
func (em *ExpiryMap[K, V]) Get(key K) (V, error) {
	em.assumeAlive()

	var err error
	var val V
	em.maybeLockForWriting(func() {
		if ent, ok := em.backMap.Get(key); !ok {
			val, err = em.loadValue(key)
		} else {
			val = ent.val
			em.notifyListeners(Requested, key, val, nil)
		}
	})
	return val, err
}

func (em *ExpiryMap[K, V]) loadValue(key K) (V, error) {
	val, err := em.loader(key)
	if err == nil {
		em.upgradeWLock()
		for em.maxCapacity != -1 && em.backMap.Len() >= em.maxCapacity {
			em.removeOldest()
		}
		keyTimer := time.NewTimer(em.ttl)
		em.backMap.Put(key, entry[V]{val: val, exptmr: keyTimer})
		go func() {
			<-keyTimer.C
			em.evictChan <- key
		}()
		em.notifyListeners(Added, key, val, nil)
	} else {
		em.notifyListeners(Failed, key, val, err) // val has "zero" value
	}
	return val, err
}

// Returns the value associated to the given key. In contrast to `Get()` this method does not trigger the loader.
func (em *ExpiryMap[K, V]) Peek(key K) (V, bool) {
	em.assumeAlive()

	var ent entry[V]
	var ok bool
	em.readAtomically(func() {
		ent, ok = em.backMap.Get(key)
		if ok {
			em.notifyListeners(Requested, key, ent.val, nil)
		} else {
			em.notifyListeners(Missed, key, ent.val, nil)
		}
	})
	return ent.val, ok
}

// Returns `true`, if there is a mapping for the specified key.
func (em *ExpiryMap[K, V]) ContainsKey(key K) bool {
	em.assumeAlive()

	var ok bool
	em.readAtomically(func() {
		_, ok = em.backMap.Get(key)
	})
	return ok
}

// Replaces synchronously the entry for a key if present. This operationresets doesn't change the expiry time.
//
//   - return `true` if value was replaced
func (em *ExpiryMap[K, V]) Replace(key K, val V) bool {
	em.assumeAlive()

	var ok bool
	em.readAtomically(func() {
		if ent, oki := em.backMap.Get(key); oki {
			em.backMap.Put(key, entry[V]{val: val, exptmr: ent.exptmr})
			em.notifyListeners(Replaced, key, val, nil)
			ok = true
		}
	})
	return ok
}

// Removes the mapping for a key from the cache if it is present.
//
//   - return `true` if value was removed
func (em *ExpiryMap[K, V]) Remove(key K) bool {
	em.assumeAlive()

	var ok bool
	em.readAtomically(func() {
		ok = em.removeEntry(key)
	})
	return ok
}

// Clears the cache.
func (em *ExpiryMap[K, V]) Clear() {
	em.assumeAlive()

	em.writeAtomically(func() {
		keys := em.backMap.Keys()
		for _, key := range keys {
			em.removeEntry(key)
		}
	})
}

// Adds listener to ExpiryMap events. The listeners are executed in a synchronous mode in order of their insretion.
func (em *ExpiryMap[K, V]) AddListener(listener Listener[K, V]) *ExpiryMap[K, V] {
	em.assumeAlive()

	em.writeAtomically(func() {
		em.listeners.add(listener)
	})
	return em
}

// Removes listener to ExpiryMap events.
func (em *ExpiryMap[K, V]) RemoveListener(listener Listener[K, V]) *ExpiryMap[K, V] {
	em.assumeAlive()

	em.writeAtomically(func() {
		em.listeners.remove(listener)
	})
	return em
}
