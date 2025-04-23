package sorted

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	keys        = []string{"a", "c", "b", "z", "x"}
	sorted_keys = []string{"a", "b", "c", "x", "z"}
)

func isStringLess(x, y string) bool {
	return x < y
}

func TestSortedMapCreation(t *testing.T) {
	assertT := assert.New(t)

	sm := NewSortedMap[string, int](isStringLess)
	assertT.NotNil(sm)
	assertT.Equal(0, sm.Len())
}

func TestSortedMapGetPut(t *testing.T) {
	assertT := assert.New(t)

	sm := NewSortedMap[string, int](isStringLess)
	sm.Put("a", 1)
	sm.Put("c", 2)
	sm.Put("b", 3)

	assertT.Equal(3, sm.Len())

	v, ok := sm.Get("b")
	assertT.True(ok)
	assertT.Equal(3, v)

	_, ok = sm.Get("z")
	assertT.False(ok)
}

func TestSortedMapKeysExtraction(t *testing.T) {
	assertT := assert.New(t)

	sm := NewSortedMap[string, int](isStringLess)
	for _, k := range keys {
		sm.Put(k, 1)
	}

	mapKeys := sm.Keys()
	for i, k := range mapKeys {
		assertT.Equal(sorted_keys[i], k)
	}
}

func TestSortedMapIterator(t *testing.T) {
	assertT := assert.New(t)

	sm := NewSortedMapEx[string, int](len(keys), isStringLess)
	for i, k := range keys {
		sm.Put(k, (i*3+13)%5)
	}

	it := sm.Iterator()
	for i := 0; it.HasNext(); i++ {
		k, _ := it.Next()
		assertT.Equal(sorted_keys[i], k)
	}
}

func TestSortedMapCompute(t *testing.T) {
	assertT := assert.New(t)

	sm := NewSortedMap[string, int](isStringLess)
	sm.Put("a", 1)
	sm.Put("b", 2)
	sm.Put("c", 3)

	sm.Compute("b", func(k string, v int) int {
		return v * 2
	})

	v, ok := sm.Get("b")
	assertT.True(ok)
	assertT.Equal(4, v)

	sm.Compute("z", func(k string, v int) int {
		return 100
	})

	v, ok = sm.Get("z")
	assertT.True(ok)
	assertT.Equal(100, v)
}

func TestSortedMapPutAll(t *testing.T) {
	assertT := assert.New(t)

	sm1 := NewSortedMap[string, int](isStringLess)
	sm1.Put("a", 1)
	sm1.Put("c", 2)
	sm1.Put("b", 3)
	sm2 := NewSortedMap[string, int](isStringLess)
	sm2.Put("a", 0)

	assertT.Equal(1, sm2.Len())
	sm2.PutAll(sm1)
	assertT.Equal(3, sm2.Len())
	val, ok := sm2.Get("a")
	assertT.True(ok)
	assertT.Equal(1, val)
	val, _ = sm2.Get("b")
	assertT.Equal(3, val)
	val, _ = sm2.Get("c")
	assertT.Equal(2, val)
}

func TestSortedMapRemove(t *testing.T) {
	assertT := assert.New(t)

	sm := NewSortedMap[string, int](isStringLess)
	sm.Put("a", 1)
	sm.Put("b", 2)
	sm.Put("c", 3)

	sm.Remove("b")
	assertT.Equal(2, sm.Len())
	_, ok := sm.Get("b")
	assertT.False(ok)
	_, ok = sm.Get("a")
	assertT.True(ok)
	_, ok = sm.Get("c")
	assertT.True(ok)

	sm.Remove("z")
	assertT.Equal(2, sm.Len())
	_, ok = sm.Get("a")
	assertT.True(ok)
	_, ok = sm.Get("c")
	assertT.True(ok)
}

func TestBinSearch(t *testing.T) {
	assertT := assert.New(t)

	assertT.Equal(0, binSearch(sorted_keys, "a", isStringLess))
	assertT.Equal(1, binSearch(sorted_keys, "b", isStringLess))
	assertT.Equal(2, binSearch(sorted_keys, "c", isStringLess))
	assertT.Equal(3, binSearch(sorted_keys, "x", isStringLess))
	assertT.Equal(4, binSearch(sorted_keys, "z", isStringLess))

	assertT.Equal(-4, binSearch(sorted_keys, "k", isStringLess))
	assertT.Equal(-5, binSearch(sorted_keys, "y", isStringLess))
	assertT.Equal(-1, binSearch(sorted_keys, "A", isStringLess))
	assertT.Equal(-6, binSearch(sorted_keys, "|", isStringLess))
}
