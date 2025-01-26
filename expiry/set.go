package expiry

// Do-it-again set
type set[T comparable] struct {
	m map[T]struct{}
}

func newSet[T comparable]() *set[T] {
	return &set[T]{m: make(map[T]struct{})}
}

func (s *set[T]) add(value T) {
	s.m[value] = struct{}{}
}

func (s *set[T]) remove(value T) {
	delete(s.m, value)
}
func (s *set[T]) contains(value T) bool {
	_, ok := s.m[value]
	return ok
}

func (s *set[T]) size() int {
	return len(s.m)
}
