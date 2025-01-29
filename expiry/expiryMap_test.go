package expiry

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity)
	// Default loader returns error
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

func TestPeek(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		WithLoader(func(key string) (int, error) { return len(key), nil })

	v, ok := em.Peek("Hi")
	assertT.False(ok)
	assertT.Equal(0, v)

	v, e := em.Get("Hi")
	assertT.Nil(e)
	assertT.Equal(2, v)

	v, ok = em.Peek("Hi")
	assertT.True(ok)
	assertT.Equal(2, v)
}

func TestContainsKey(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		WithLoader(func(key string) (int, error) { return len(key), nil })

	assertT.False(em.ContainsKey("Hi"))

	v, e := em.Get("Hi")
	assertT.Nil(e)
	assertT.Equal(2, v)

	assertT.True(em.ContainsKey("Hi"))
}

func TestCapacity(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(2).
		WithLoader(func(key string) (int, error) { return len(key), nil })

	assertT.Equal(2, em.Capacity())

	_, _ = em.Get("Hi")
	_, _ = em.Get("Hello")
	assertT.Equal(2, em.Len())
	assertT.True(em.ContainsKey("Hi"))
	assertT.True(em.ContainsKey("Hello"))

	_, _ = em.Get("World!")
	assertT.Equal(2, em.Len())
	assertT.True(em.ContainsKey("Hello"))
	assertT.True(em.ContainsKey("World!"))
}

func TestReplace(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		WithLoader(func(key string) (int, error) { return len(key), nil })

	v, _ := em.Get("Hi")
	assertT.Equal(2, v)
	v, _ = em.Peek("Hi")
	assertT.Equal(2, v)

	assertT.True(em.Replace("Hi", 5))
	v, _ = em.Peek("Hi")
	assertT.Equal(5, v)

	assertT.False(em.Replace("Hello", 5))
}

func TestRemove(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		WithLoader(func(key string) (int, error) { return len(key), nil })

	_, _ = em.Get("Hi")
	_, _ = em.Get("Hello")
	assertT.Equal(2, em.Len())
	assertT.True(em.ContainsKey("Hi"))
	assertT.True(em.ContainsKey("Hello"))

	assertT.True(em.Remove("Hi"))
	assertT.Equal(1, em.Len())
	assertT.False(em.ContainsKey("Hi"))
	assertT.True(em.ContainsKey("Hello"))

	assertT.False(em.Remove("Hi"))
}

func last[T any](a []T) T {
	return a[len(a)-1]
}

func penult[T any](a []T) T {
	return a[len(a)-2]
}

func assertNotification(t *testing.T, expEvent EventType, expKey string, expVal int, expErr error, event EventType, key string, val int, err error) {
	assertT := assert.New(t)

	assertT.Equal(expEvent, event)
	assertT.Equal(expKey, key)
	assertT.Equal(expVal, val)
	assertT.Equal(expErr, err)
}

func TestNotifications(t *testing.T) {
	keys := make([]string, 0)
	vals := make([]int, 0)
	events := make([]EventType, 0)
	errs := make([]error, 0)
	callback := func(ev EventType, key string, val int, err error) {
		events = append(events, ev)
		keys = append(keys, key)
		vals = append(vals, val)
		errs = append(errs, err)
	}

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(1).
		AddListener(&ListenerWarapper{callback})

	// Default loader returns error
	_, err := em.Get("Hi")
	assertNotification(t, Failed, "Hi", 0, err, last(events), last(keys), last(vals), last(errs))

	em.WithLoader(func(key string) (int, error) { return len(key), nil })

	_, _ = em.Get("Hi")
	assertNotification(t, Added, "Hi", 2, nil, last(events), last(keys), last(vals), last(errs))
	_, _ = em.Get("Hi")
	assertNotification(t, Requested, "Hi", 2, nil, last(events), last(keys), last(vals), last(errs))
	em.Peek("Hi")
	assertNotification(t, Requested, "Hi", 2, nil, last(events), last(keys), last(vals), last(errs))
	em.Peek("Hello")
	assertNotification(t, Missed, "Hello", 0, nil, last(events), last(keys), last(vals), last(errs))

	_, _ = em.Get("Hello")
	assertNotification(t, Removed, "Hi", 2, nil, penult(events), penult(keys), penult(vals), penult(errs))
	assertNotification(t, Added, "Hello", 5, nil, last(events), last(keys), last(vals), last(errs))

	_ = em.Replace("Hello", 3)
	assertNotification(t, Replaced, "Hello", 3, nil, last(events), last(keys), last(vals), last(errs))
}

func TestClear(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		WithLoader(func(key string) (int, error) { return len(key), nil })

	_, _ = em.Get("Hi")
	_, _ = em.Get("Hello")
	assertT.Equal(2, em.Len())

	em.Clear()
	assertT.Equal(0, em.Len())
}

func TestSingleStart(t *testing.T) {
	em := NewExpiryMap[string, int]()

	assert.NotPanics(t, func() { em.Start() })
	assert.Panics(t, func() { em.Start() })
}

func TestExpiry(t *testing.T) {
	assertT := assert.New(t)

	ttlE := time.Duration(10) * time.Millisecond

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		WithLoader(func(key string) (int, error) { return len(key), nil }).
		ExpireAfter(ttlE).
		Start()

	_, _ = em.Get("Hi")
	assertT.Equal(1, em.Len())
	assertT.True(em.ContainsKey("Hi"))

	_, _ = em.Get("Hello")
	assertT.Equal(2, em.Len())
	assertT.True(em.ContainsKey("Hello"))

	time.Sleep(3 * ttlE)
	assertT.Equal(0, em.Len())

	_, _ = em.Get("World!")
	assertT.Equal(1, em.Len())
	assertT.True(em.ContainsKey("World!"))
}

func BenchmarkExpiryMap(b *testing.B) {
	ttlE := time.Duration(10) * time.Millisecond
	em := NewExpiryMap[string, string]().
		WithMaxCapacity(Unlimited).
		WithLoader(func(key string) (string, error) { return "value" + key, nil }).
		ExpireAfter(ttlE).
		Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iS := strconv.Itoa(i)
		_, _ = em.Get(iS)
		_, _ = em.Get(iS) // <- shouldn't triggert loading
	}
}
