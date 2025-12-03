package util

type Set[T comparable] map[T]bool

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(value T) {
	s[value] = true
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for key := range s {
		slice = append(slice, key)
	}
	return slice
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	unionSet := s
	for key := range other {
		unionSet.Add(key)
	}
	return unionSet
}

func (s Set[T]) Intersection(other Set[T]) Set[T] {
	intersectionSet := NewSet[T]()
	for key := range s {
		if other.Contains(key) {
			intersectionSet.Add(key)
		}
	}
	return intersectionSet
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
	differenceSet := NewSet[T]()
	for key := range s {
		if !other.Contains(key) {
			differenceSet.Add(key)
		}
	}
	return differenceSet
}

func (s Set[T]) IsSubset(other Set[T]) bool {
	for key := range s {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}
