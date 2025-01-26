package expiry

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
				// UC What to do with timer job?
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
			// UC What to do with timer job?
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
			// UC What to do with timer job?
		}
	})
	return ok
}
