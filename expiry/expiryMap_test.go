package expiry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	maxCapacity = 123
	ttl         = time.Duration(777000000)
)

func TestCreation(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		Expirefter(ttl)

	assertT.NotNil(t, em)
	assertT.Equal(maxCapacity, em.Capacity())
	assertT.Equal(ttl, em.ExpiringAfter())
	assertT.NotNil(em.loader)
	assertT.Equal(0, len(em.listeners))
}

func TestGet(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		Expirefter(ttl)
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
