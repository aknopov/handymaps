package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assertT := assert.New(t)

	s := NewSet[string]()

	s.Add("Hello")
	assertT.Equal(1, s.Size())
	assertT.True(s.Contains("Hello"))

	s.Add("Hello")
	assertT.Equal(1, s.Size())
	assertT.True(s.Contains("Hello"))

	s.Add("World")
	assertT.Equal(2, s.Size())
	assertT.True(s.Contains("World"))

	s.Remove("Hello")
	assertT.Equal(1, s.Size())
	assertT.False(s.Contains("Hello"))
	assertT.True(s.Contains("World"))
}
