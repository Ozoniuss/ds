package main

import (
	"slices"
)

type SortedList struct {
	elems []int
}

func (s *SortedList) Insert(x int) {
	s.elems = append(s.elems, x)
	slices.Sort(s.elems)
}
