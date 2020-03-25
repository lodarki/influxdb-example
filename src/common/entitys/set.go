package entitys

import (
	"sync"
)

type Set struct {
	m map[interface{}]bool
	sync.RWMutex
}

func NewSet(len int) *Set {
	if len <= 0 {
		return &Set{
			m: map[interface{}]bool{},
		}
	}
	return &Set{
		m: make(map[interface{}]bool, len),
	}
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() int {
	return len(s.m)
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) ListForInt64() []int64 {
	s.RLock()
	defer s.RUnlock()
	var list []int64
	for item := range s.m {
		list = append(list, item.(int64))
	}
	return list
}

func (s *Set) ListForInt() []int {
	s.RLock()
	defer s.RUnlock()
	var list []int
	for item := range s.m {
		list = append(list, item.(int))
	}
	return list
}

func (s *Set) ListForString() []string {
	s.RLock()
	defer s.RUnlock()
	var list []string
	for item := range s.m {
		list = append(list, item.(string))
	}
	return list
}
