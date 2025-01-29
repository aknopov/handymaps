package expiry

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	iters   = 10
	sleepMs = 50
)

func TestNoUpgradeWithWriteLock(t *testing.T) {
	assertT := assert.New(t)

	c := atomic.Int32{}
	rw := upgradableRWMutex{}
	rw.lock()
	go func() {
		rw.upgradeRLock()
		defer rw.upgradableRUnlock()
		c.Add(1)
	}()
	time.Sleep(sleepMs * time.Millisecond)
	assertT.Equal(int32(0), c.Load())
	rw.unlock()
}

func TestCanNotWriteLockWithUpgraded(t *testing.T) {
	assertT := assert.New(t)

	c := atomic.Int32{}
	rw := upgradableRWMutex{}
	rw.upgradeRLock()

	rw.upgradeWLock()
	go func() {
		rw.lock()
		defer rw.unlock()
		c.Add(1)
	}()
	time.Sleep(sleepMs * time.Millisecond)
	assertT.Equal(int32(0), c.Load())
	rw.upgradableRUnlock()
}

func TestMutipleReadersBeforeReadUpgrade(t *testing.T) {
	rw := upgradableRWMutex{}

	rw.rLock()
	defer rw.rUnlock()

	rw.rLock()
	defer rw.rUnlock()

	rw.upgradeRLock()
	defer rw.upgradableRUnlock()
}

func TestMutipleReadersAfterReadUpgrade(t *testing.T) {
	rw := upgradableRWMutex{}
	rw.upgradeRLock()

	waitGroup := sync.WaitGroup{}
	for i := 0; i < iters; i++ {
		waitGroup.Add(1)
		go func() {
			rw.rLock()
			defer rw.rUnlock()
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	rw.upgradableRUnlock()
}

func TestNoReadsAfterUpgrade(t *testing.T) {
	assertT := assert.New(t)

	c := atomic.Int32{}
	rw := upgradableRWMutex{}
	rw.upgradeRLock()

	rw.upgradeWLock()
	go func() {
		rw.rLock()
		defer rw.rUnlock()
		c.Add(1)
	}()
	time.Sleep(sleepMs * time.Millisecond)
	assertT.Equal(int32(0), c.Load())
	rw.upgradableRUnlock()
}

func TestNoDoubleUpgrade(t *testing.T) {
	assertT := assert.New(t)

	c := atomic.Int32{}
	rw := upgradableRWMutex{}

	rw.rLock()
	defer rw.rUnlock()

	rw.rLock()
	defer rw.rUnlock()

	rw.upgradeRLock()
	go func() {
		rw.upgradeRLock()
		defer rw.upgradableRUnlock()
		c.Add(1)
	}()
	time.Sleep(sleepMs * time.Millisecond)
	assertT.Equal(int32(0), c.Load())
	rw.upgradableRUnlock()
}

func TestReadAtomically(t *testing.T) {
	rw := upgradableRWMutex{}

	waitGroup := sync.WaitGroup{}
	for i := 0; i < iters; i++ {
		waitGroup.Add(1)
		go rw.readAtomically(func() {
			waitGroup.Done()
		})
	}
	waitGroup.Wait()
}

func TestWriteAtomically(t *testing.T) {
	assertT := assert.New(t)

	c := 0 // Da - simple var!
	rw := upgradableRWMutex{}

	waitGroup := sync.WaitGroup{}
	for i := 0; i < iters; i++ {
		waitGroup.Add(1)
		go rw.writeAtomically(func() {
			c++
			waitGroup.Done()
		})
	}
	waitGroup.Wait()
	assertT.Equal(iters, c)
}

func TestMaybeLockForWriting(t *testing.T) {
	assertT := assert.New(t)

	c := 0 // Da - simple var!
	rw := upgradableRWMutex{}

	rw.maybeLockForWriting(func() {
		rw.upgradeWLock()
		c = 1
	})

	assertT.Equal(1, c)
}
