package expiry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assertT := assert.New(t)

	s := newSet[string]()

	s.add("Hello")
	assertT.Equal(1, s.size())
	assertT.True(s.contains("Hello"))

	s.add("Hello")
	assertT.Equal(1, s.size())
	assertT.True(s.contains("Hello"))

	s.add("World")
	assertT.Equal(2, s.size())
	assertT.True(s.contains("World"))

	s.remove("Hello")
	assertT.Equal(1, s.size())
	assertT.False(s.contains("Hello"))
	assertT.True(s.contains("World"))
}
