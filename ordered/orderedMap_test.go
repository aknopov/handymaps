package ordered

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreation(t *testing.T) {
	assert := assert.New(t)

	om := NewOrderedMap[string, int]()
	assert.NotNil(om)
	assert.Equal(0, om.Len())
}

func TestGetPut(t *testing.T) {
	assert := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	assert.Equal(3, om.Len())

	v, ok := om.Get("b")
	assert.True(ok)
	assert.Equal(2, v)

	_, ok = om.Get("z")
	assert.False(ok)
}

func TestKeysExtraction(t *testing.T) {
	assert := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("c", 3)
	om.Put("b", 2)
	om.Put("z", -1)
	om.Put("y", 13)

	keys := om.Keys()
	assert.Equal(5, len(keys))
	assert.Equal("a", keys[0])
	assert.Equal("c", keys[1])
	assert.Equal("b", keys[2])
	assert.Equal("z", keys[3])
	assert.Equal("y", keys[4])
}

func TestIterator(t *testing.T) {
	assert := assert.New(t)

	keys := []string{"a", "c", "b", "z", "y"}

	om := NewOrderedMapEx[string, int](len(keys))
	for i, k := range keys {
		om.Put(k, i)
	}

	it := om.Iterator()
	for it.HasNext() {
		k, v := it.Next()
		assert.Equal(keys[v], k)
	}
}

func TestCompute(t *testing.T) {
	assert := assert.New(t)

	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om.Compute("b", func(k string, v int) int {
		return v * 2
	})

	v, ok := om.Get("b")
	assert.True(ok)
	assert.Equal(4, v)

	om.Compute("z", func(k string, v int) int {
		return 100
	})

	v, ok = om.Get("z")
	assert.True(ok)
	assert.Equal(100, v)
}
