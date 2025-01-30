package ordered

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreation(t *testing.T) {
	assertT := assert.New(t)

	om := NewOrderedMap[string, int]()
	assertT.NotNil(t, om)
	assertT.Equal(0, om.Len())
}

func TestGetPut(t *testing.T) {
	assertT := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	assertT.Equal(3, om.Len())

	v, ok := om.Get("b")
	assertT.True(ok)
	assertT.Equal(2, v)

	_, ok = om.Get("z")
	assertT.False(ok)
}

func TestKeysExtraction(t *testing.T) {
	assertT := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("c", 3)
	om.Put("b", 2)
	om.Put("z", -1)
	om.Put("y", 13)

	keys := om.Keys()
	assertT.Equal(5, len(keys))
	assertT.Equal("a", keys[0])
	assertT.Equal("c", keys[1])
	assertT.Equal("b", keys[2])
	assertT.Equal("z", keys[3])
	assertT.Equal("y", keys[4])
}

func TestIterator(t *testing.T) {
	assertT := assert.New(t)

	keys := []string{"a", "c", "b", "z", "y"}

	om := NewOrderedMapEx[string, int](len(keys))
	for i, k := range keys {
		om.Put(k, i)
	}

	it := om.Iterator()
	for it.HasNext() {
		k, v := it.Next()
		assertT.Equal(keys[v], k)
	}
}

func TestCompute(t *testing.T) {
	assertT := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om.Compute("b", func(k string, v int) int {
		return v * 2
	})

	v, ok := om.Get("b")
	assertT.True(ok)
	assertT.Equal(4, v)

	om.Compute("z", func(k string, v int) int {
		return 100
	})

	v, ok = om.Get("z")
	assertT.True(ok)
	assertT.Equal(100, v)
}

func TestPutAll(t *testing.T) {
	assertT := assert.New(t)

	om1 := NewOrderedMap[string, int]()
	om2 := NewOrderedMap[string, int]()
	om1.Put("a", 1)
	om1.Put("b", 2)
	om1.Put("c", 3)

	assertT.Equal(0, om2.Len())
	om2.PutAll(om1)
	assertT.Equal(3, om2.Len())
	val, ok := om2.Get("a")
	assertT.True(ok)
	assertT.Equal(1, val)
	val, _ = om2.Get("b")
	assertT.Equal(2, val)
	val, _ = om2.Get("c")
	assertT.Equal(3, val)
}

func TestRemove(t *testing.T) {
	assertT := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	assertT.True(om.Remove("b"))
	assertT.Equal(2, om.Len())
	_, ok := om.Get("b")
	assertT.False(ok)
	_, ok = om.Get("a")
	assertT.True(ok)
	_, ok = om.Get("c")
	assertT.True(ok)

	assertT.False(om.Remove("z"))
	assertT.Equal(2, om.Len())
	_, ok = om.Get("a")
	assertT.True(ok)
	_, ok = om.Get("c")
	assertT.True(ok)
}
