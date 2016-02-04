package util

import (
	"fmt"
)

// An unordered collection of unique elements
type Set interface {
	// Insert v into the set
	Add(v interface{})

	// Remove v from the set
	Discard(v interface{})

	// Return true if and only if the set contains v
	Contains(v interface{}) bool

	// Return a channel from which each element in the set can be read exactly once. If the set is
	// mutated before the channel is emptied, the exact data read from the channel is undefined.
	// Deprecated: If return before iterate all elements, some resource will not be released.
	Iter() <-chan interface{}

	// Execute f(v) for every element v in set until f return false. if f mutate this set,
	// behavior is undefined.
	DoWhile(f func(interface{}) bool)

	// Return an element randomly and remove it from the set
	// if the set is not empty, the second return value is true; otherwise, return false
	Pop() (interface{}, bool)

	// Return the number of elements in the set
	Size() int

	// Returns a new Set that contains exactly the same elements as this set
	Copy() Set

	// Returns a new Set that contains elements in this set but not in s
	Diff(s Set) Set

	// Return a new Set that contains element both in this set and in s
	Intersection(s Set) Set

	// Return a new Set that contains element either in this set or in s
	Union(s Set) Set

	// Remove all elements in the set
	Clean()

	// Print elements in the set in a human-readable format
	String() string
}

func NewSet(items ...interface{}) Set {
	s := &setImpl{
		data: make(map[interface{}]struct{}),
	}
	for _, item := range items {
		s.Add(item)
	}
	return s
}

type setImpl struct {
	data map[interface{}]struct{}
}

func (s *setImpl) Add(v interface{}) {
	s.data[v] = struct{}{}
}

func (s *setImpl) Discard(v interface{}) {
	delete(s.data, v)
}

func (s *setImpl) Contains(v interface{}) bool {
	_, ok := s.data[v]
	return ok
}

func (s *setImpl) Iter() <-chan interface{} {
	iter := make(chan interface{})
	go func() {
		for key, _ := range s.data {
			iter <- key
		}
		close(iter)
	}()
	return iter
}

func (s *setImpl) DoWhile(f func(interface{}) bool) {
	for key, _ := range s.data {
		if !f(key) {
			break
		}
	}
}

func (s *setImpl) Pop() (interface{}, bool) {
	for key, _ := range s.data {
		delete(s.data, key)
		return key, true
	}
	return nil, false
}

func (s *setImpl) Size() int {
	return len(s.data)
}

func (s *setImpl) Copy() Set {
	r := NewSet()
	for v := range s.data {
		r.Add(v)
	}
	return r
}

func (s *setImpl) Diff(t Set) Set {
	r := NewSet()
	for v := range s.data {
		if !t.Contains(v) {
			r.Add(v)
		}
	}
	return r
}

func (s *setImpl) Intersection(s2 Set) Set {
	if s.Size() > s2.Size() {
		return s2.Intersection(s)
	}

	r := NewSet()
	for key, _ := range s.data {
		if s2.Contains(key) {
			r.Add(key)
		}
	}
	return r
}

func (s *setImpl) Union(s2 Set) Set {
	r := s2.Copy()
	for key, _ := range s.data {
		r.Add(key)
	}
	return r
}

func (s *setImpl) Clean() {
	s.data = make(map[interface{}]struct{})
}

func (s *setImpl) String() string {
	ret := make([]interface{}, 0, s.Size())
	for key, _ := range s.data {
		ret = append(ret, key)
	}
	return fmt.Sprintf("%v", ret)
}
