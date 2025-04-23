package util

// Do-it-again Set
type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

func (s *Set[T]) Base() map[T]struct{} {
	return s.m
}

func (s *Set[T]) Add(value T) {
	s.m[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	delete(s.m, value)
}
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.m)
}
