package expiry

import (
	"sync"
	"sync/atomic"
	_ "unsafe"
)

// "Upgradable" R/W lock "inspired" by https://upstash.com/blog/upgradable-rwlock-for-go
type upgradableRWMutex struct {
	w            sync.Mutex   // held if there are pending writers
	writerSem    uint32       // semaphore for writers to wait for completing readers
	readerSem    uint32       // semaphore for readers to wait for completing writers
	readerCount  atomic.Int32 // number of pending readers
	readerWait   atomic.Int32 // number of departing readers
	upgraded     bool         // whether write lock was upgraded
	upgradedRead bool         // whether read lock was upgraded
}

//go:linkname semaphoreAcquire sync.runtime_Semacquire
func semaphoreAcquire(s *uint32)

//go:linkname semaphoreRelease sync.runtime_Semrelease
func semaphoreRelease(s *uint32, handoff bool, skipframes int)

const rwmutexMaxReaders = 1 << 30

// -- Convenience wrappers ---

// Locks for reading while executing function
func (rw *upgradableRWMutex) readAtomically(f func()) {
	rw.rLock()
	defer rw.rUnlock()
	f()
}

func (rw *upgradableRWMutex) writeAtomically(f func()) {
	rw.lock()
	defer rw.unlock()
	f()
}

// Locks first for reading function execution allowing optional lock upgrade
// with `upgradableRWMutex.upgradeWLock` at some later time
func (rw *upgradableRWMutex) maybeLockForWriting(f func()) {
	rw.upgradeRLock()
	defer rw.upgradableRUnlock()
	f()
}

// -- Upgradable functionality ---

// First, resolve competition with other writers.
// Disallow writers to acquire the lock
func (rw *upgradableRWMutex) upgradeRLock() {
	rw.w.Lock()
	rw.upgradedRead = true
}

// Upgrade current R-locks to W-lock
func (rw *upgradableRWMutex) upgradeWLock() {
	rw.upgraded = true
	// Announce to readers there is a pending writer.
	r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
	// Wait for active readers.
	if r != 0 && rw.readerWait.Add(r) != 0 {
		semaphoreAcquire(&rw.writerSem)
	}
}

// Undoes upgrade of r-locks and unlocksthe R-lock
func (rw *upgradableRWMutex) upgradableRUnlock() {
	rw.upgradedRead = false
	if rw.upgraded {
		rw.upgraded = false
		rw.unlock()
	} else {
		rw.w.Unlock()
	}
}

// -- Standard functionality of `sync.RWMutex` (no racing checks though) --

// Locks rw for writing - standard implementation.
func (rw *upgradableRWMutex) lock() {
	// First, resolve competition with other writers.
	rw.w.Lock()
	// Announce to readers there is a pending writer.
	r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
	// Wait for active readers.
	if r != 0 && rw.readerWait.Add(r) != 0 {
		semaphoreAcquire(&rw.writerSem)
	}
}

// Unlocks the lock for writing - standard implementation.
func (rw *upgradableRWMutex) unlock() {
	// Announce to readers there is no active writer.
	r := rw.readerCount.Add(rwmutexMaxReaders)
	// Unblock blocked readers, if any.
	for i := 0; i < int(r); i++ {
		semaphoreRelease(&rw.readerSem, false, 0)
	}
	// Allow other writers to proceed.
	rw.w.Unlock()
}

// Locks for reading - standard implementation
func (rw *upgradableRWMutex) rLock() {
	if rw.readerCount.Add(1) < 0 {
		// A writer is pending, wait for it.
		semaphoreAcquire(&rw.readerSem)
	}
}

// Undoes a single rLock call - standard implementation
func (rw *upgradableRWMutex) rUnlock() {
	if r := rw.readerCount.Add(-1); r < 0 {
		// Outlined slow-path to allow the fast-path to be inlined
		// A writer is pending.
		if rw.readerWait.Add(-1) == 0 {
			// The last reader unblocks the writer.
			semaphoreRelease(&rw.writerSem, false, 1)
		}
	}
}
