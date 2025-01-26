package expiry

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const iters = 5

func TestOnceAutoResetSynchronous(t *testing.T) {
	assertT := assert.New(t)

	var once Once
	var calls int

	for i := 0; i < iters; i++ {
		once.DoAtomically(func() { calls++ })
	}
	assertT.Equal(iters, calls)
}

func TestOnceAutoResetAsynchronous(t *testing.T) {
	assertT := assert.New(t)

	var once Once
	var calls int
	var wg sync.WaitGroup

	wg.Add(iters)
	for i := 0; i < iters; i++ {
		once.DoAtomically(func() { defer wg.Done(); calls++ })
	}
	wg.Wait()
	assertT.Equal(iters, calls)
}
