package expiry

import (
	"fmt"
	"runtime"
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
	timeout := time.After(sleepMs * time.Millisecond)
	done := make(chan interface{})

	go func() {
		rw := upgradableRWMutex{}

		rw.rLock()
		defer rw.rUnlock()

		rw.rLock()
		defer rw.rUnlock()

		rw.upgradeRLock()
		defer rw.upgradableRUnlock()

		done <- struct{}{}
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
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

// Tests from https://go.dev/src/sync/rwmutex_test.go

func parallelReader(m *upgradableRWMutex, clocked, cunlock, cdone chan bool) {
	m.rLock()
	clocked <- true
	<-cunlock
	m.rUnlock()
	cdone <- true
}

func doTestParallelReaders(numReaders, gomaxprocs int) {
	runtime.GOMAXPROCS(gomaxprocs)
	var m upgradableRWMutex
	clocked := make(chan bool)
	cunlock := make(chan bool)
	cdone := make(chan bool)
	for i := 0; i < numReaders; i++ {
		go parallelReader(&m, clocked, cunlock, cdone)
	}
	// Wait for all parallel RLock()s to succeed.
	for i := 0; i < numReaders; i++ {
		<-clocked
	}
	for i := 0; i < numReaders; i++ {
		cunlock <- true
	}
	// Wait for the goroutines to finish.
	for i := 0; i < numReaders; i++ {
		<-cdone
	}
}

func TestParallelReaders(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	doTestParallelReaders(1, 4)
	doTestParallelReaders(3, 4)
	doTestParallelReaders(4, 2)
}

func reader(rwm *upgradableRWMutex, num_iterations int, activity *int32, cdone chan bool) {
	for i := 0; i < num_iterations; i++ {
		rwm.rLock()
		n := atomic.AddInt32(activity, 1)
		if n < 1 || n >= 10000 {
			rwm.rUnlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for i := 0; i < 100; i++ {
		}
		atomic.AddInt32(activity, -1)
		rwm.rUnlock()
	}
	cdone <- true
}

func writer(rwm *upgradableRWMutex, num_iterations int, activity *int32, cdone chan bool) {
	for i := 0; i < num_iterations; i++ {
		rwm.lock()
		n := atomic.AddInt32(activity, 10000)
		if n != 10000 {
			rwm.unlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for i := 0; i < 100; i++ {
		}
		atomic.AddInt32(activity, -10000)
		rwm.unlock()
	}
	cdone <- true
}

func HammerRWMutex(gomaxprocs, numReaders, num_iterations int) {
	runtime.GOMAXPROCS(gomaxprocs)
	// Number of active readers + 10000 * number of active writers.
	var activity int32
	var rwm upgradableRWMutex
	cdone := make(chan bool)
	go writer(&rwm, num_iterations, &activity, cdone)
	var i int
	for i = 0; i < numReaders/2; i++ {
		go reader(&rwm, num_iterations, &activity, cdone)
	}
	go writer(&rwm, num_iterations, &activity, cdone)
	for ; i < numReaders; i++ {
		go reader(&rwm, num_iterations, &activity, cdone)
	}
	// Wait for the 2 writers and all readers to finish.
	for i := 0; i < 2+numReaders; i++ {
		<-cdone
	}
}

func TestRWMutex(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	n := 1000
	if testing.Short() {
		n = 5
	}
	HammerRWMutex(1, 1, n)
	HammerRWMutex(1, 3, n)
	HammerRWMutex(1, 10, n)
	HammerRWMutex(4, 1, n)
	HammerRWMutex(4, 3, n)
	HammerRWMutex(4, 10, n)
	HammerRWMutex(10, 1, n)
	HammerRWMutex(10, 3, n)
	HammerRWMutex(10, 10, n)
	HammerRWMutex(10, 5, n)
}
