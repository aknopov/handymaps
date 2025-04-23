package bimap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s", msg)
		}
	}()
	f()
}

func TestNewBiMap(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMapEx[string, int](5)

	assertT.Equal(0, aBimap.Size())
}

func TestBiMapBasics(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMap[string, int]()

	aBimap.Put("Hello", 1)
	assertT.Equal(1, aBimap.Size())

	assertT.True(aBimap.ContainsKey("Hello"))
	assertT.True(aBimap.ContainsValue(1))
	assertT.False(aBimap.ContainsKey("there!"))

	aBimap.Put("there!", 2)
	assertT.Equal(2, aBimap.Size())

	aBimap.RemoveKey("Hello")
	assertT.Equal(1, aBimap.Size())
	idx, ok := aBimap.GetValue("there!")
	assertT.True(ok)
	assertT.Equal(2, idx)
	_, ok = aBimap.GetValue("Hello")
	assertT.False(ok)

	aBimap.RemoveValue(2)
	assertT.Equal(0, aBimap.Size())
}

func TestDuplicatedEntries(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMap[string, int]()

	aBimap.Put("Hello", 1)
	assertT.Equal(1, aBimap.Size())
	v, _ := aBimap.GetValue("Hello")
	assertT.Equal(1, v)

	// NOP case
	aBimap.Put("Hello", 1)
	assertT.Equal(1, aBimap.Size())
	v, _ = aBimap.GetValue("Hello")
	assertT.Equal(1, v)

	aBimap.Put("Hello", 2)
	assertT.Equal(1, aBimap.Size())
	v, _ = aBimap.GetValue("Hello")
	assertT.Equal(2, v)

	aBimap.Put("test", 2)
	assertT.Equal(1, aBimap.Size())
	v, _ = aBimap.GetValue("test")
	assertT.Equal(2, v)
}

func TestInverse(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMap[string, int]()
	aBimap.Put("Hello", 1)
	aBimap.Put("there!", 2)
	assertT.Equal(2, aBimap.Size())

	iBimap := aBimap.Inverse()
	v, _ := iBimap.GetValue(1)
	assertT.Equal("Hello", v)
	v, _ = iBimap.GetValue(2)
	assertT.Equal("there!", v)
	assertT.Equal(2, iBimap.Size())

	iBimap.RemoveKey(1)
	assertT.Equal(1, iBimap.Size())
	assertT.Equal(2, aBimap.Size())
	v, _ = iBimap.GetValue(2)
	assertT.Equal("there!", v)
}

func TestGetKey(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMap[string, int]()
	aBimap.Put("Hello", 1)
	aBimap.Put("there!", 2)

	v, _ := aBimap.GetKey(1)
	assertT.Equal("Hello", v)
	v, _ = aBimap.GetKey(2)
	assertT.Equal("there!", v)

	v, ok := aBimap.GetKey(3)
	assertT.False(ok)
	assertT.Equal("", v)
}

func TestEquals(t *testing.T) {
	assertT := assert.New(t)

	bimap1 := NewBiMap[string, int]()
	bimap2 := NewBiMap[string, int]()
	assertT.True(bimap1.Equals(bimap2))
	assertT.True(bimap2.Equals(bimap1))

	bimap1.Put("Hello", 1)
	bimap1.Put("there!", 2)
	assertT.False(bimap1.Equals(bimap2))
	assertT.False(bimap2.Equals(bimap1))

	bimap2.Put("Hello", 1)
	assertT.False(bimap1.Equals(bimap2))
	assertT.False(bimap2.Equals(bimap1))
	bimap2.Put("there!", 2)
	assertT.True(bimap1.Equals(bimap2))
	assertT.True(bimap2.Equals(bimap1))
}

func TestPutAll(t *testing.T) {
	assertT := assert.New(t)

	bimap1 := NewBiMap[string, int]()
	bimap2 := NewBiMap[string, int]()
	bimap1.Put("Hello", 1)
	bimap1.Put("there!", 2)

	assertT.Equal(0, bimap2.Size())
	bimap2.PutAll(bimap1)
	assertT.Equal(2, bimap2.Size())
}

func TestKeysValues(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMap[string, int]()
	aBimap.Put("Hello", 1)
	aBimap.Put("there!", 2)

	assertT.Equal([]string{"Hello", "there!"}, aBimap.Keys())
	assertT.Equal([]int{1, 2}, aBimap.Values())

	// Replacement with duplicate key
	aBimap.Put("Hello", 5)
	assertT.Equal([]string{"Hello", "there!"}, aBimap.Keys())
	assertT.Equal([]int{5, 2}, aBimap.Values())

	// Replacement with duplicate value
	aBimap.Put("test", 5)
	assertT.Equal([]string{"test", "there!"}, aBimap.Keys())
	assertT.Equal([]int{5, 2}, aBimap.Values())
}

func TestIterator(t *testing.T) {
	assertT := assert.New(t)

	aBimap := NewBiMap[string, int]()
	message := []string{"Hello ", "there!", " Here ", " goes ", " the ", "test."}
	for i, v := range message {
		aBimap.Put(v, i+1)
	}

	it := aBimap.Iterator()
	i := 0
	for it.HasNext() {
		k, v := it.Next()
		assertT.Equal(message[i], k)
		assertT.Equal(i+1, v)
		i++
	}

	assertPanic(t, "No panic?", func() { it.Next() })
}

func TestCompareSlices(t *testing.T) {
	assertT := assert.New(t)

	assertT.True(cmpSlices([]int{1, 2, 3}, []int{1, 2, 3}))
	assertT.False(cmpSlices([]int{1, 2, 3}, []int{1, 2, 4}))
	assertT.False(cmpSlices([]int{1, 2, 3}, []int{1, 2}))
	assertT.False(cmpSlices([]int{1, 2, 3}, []int{1, 2, 3, 4}))
}

const nCount = 1000

func BenchmarkBiMapPut(b *testing.B) {
	aMap := NewBiMapEx[int, int](nCount)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i % nCount
		aMap.Put(k, k)
	}
}

func BenchmarkBiMapIteration(b *testing.B) {
	aMap := NewBiMapEx[int, int](nCount)

	for i := 0; i < nCount; i++ {
		aMap.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := aMap.Iterator()
		for it.HasNext() {
			it.Next()
		}
	}
}
