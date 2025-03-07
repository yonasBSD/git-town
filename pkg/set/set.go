package set

import (
	"cmp"
	"maps"
	"slices"
)

// Set implements a simple generic Set implementation.
// The zero value is not valid. Use set.New to create new instances.
type Set[T cmp.Ordered] map[T]struct{}

func New[T cmp.Ordered](values ...T) Set[T] {
	result := Set[T]{}
	result.Add(values...)
	return result
}

func (self Set[T]) Add(values ...T) {
	for _, value := range values {
		self[value] = struct{}{}
	}
}

func (self Set[T]) AddSet(other Set[T]) {
	self.Add(other.Values()...)
}

func (self Set[T]) Values() []T {
	result := slices.Collect(maps.Keys(self))
	slices.Sort(result)
	return result
}
