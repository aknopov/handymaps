package expiry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity)
	v, e := em.Get("Hi")
	assertT.NotNil(e)
	assertT.Equal(0, v)
	assertT.Equal(0, em.Len())

	em.WithLoader(func(key string) (int, error) { return len(key), nil })
	v, e = em.Get("Hi")
	assertT.Nil(e)
	assertT.Equal(2, v)
	assertT.Equal(1, em.Len())

	v, e = em.Get("Hello")
	assertT.Nil(e)
	assertT.Equal(5, v)
	assertT.Equal(2, em.Len())
}

func TesPeek(t *testing.T) {
	// UC
}

func TestContainsKey(t *testing.T) {
	// UC
}

func TestExpiry(t *testing.T) {
	// UC
}

func TestReplace(t *testing.T) {
	// UC
}

func TestRemove(t *testing.T) {
	// UC
}

func TestNotifications(t *testing.T) {
	// UC
}

func TestFailurePrpoagation(t *testing.T) {
	// UC
}


