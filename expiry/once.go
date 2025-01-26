package expiry

import (
	"sync"
	"sync/atomic"
)

type once struct {
	m    sync.Mutex
	done uint32
}

// Following is "inspired" by https://github.com/matryer/resync/tree/master
//
// The method does double-check locking and releases lock upon completion
func (o *once) doAtomically(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.unlockWithReset()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *once) unlockWithReset() {
	o.m.Unlock()
	atomic.StoreUint32(&o.done, 0)
}
