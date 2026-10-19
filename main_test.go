package main

import (
	"slices"
	"testing"
)

func TestSortedListInsertions(t *testing.T) {

	elems := []int{5, 1, 3, 6, 8, 3, 24, 65}
	elemsSorted := slices.Clone(elems)
	slices.Sort(elemsSorted)

	s := SortedList{}
	for idx, e := range elems {
		s.Insert(e)
		elemCnt := idx + 1

		// test that you have the expected elements after every inserion

		if len(s.elems) != elemCnt {
			t.Fatalf("expected %d elements, got %d", elemCnt, len(s.elems))
		}

		for i := range elemCnt {
			if s.elems[i] != elemsSorted[i] {
				t.Fatalf("expected %d at position %d when elem count is %d, got %d", elemsSorted[i], i, s.elems[i])
			}
		}
	}
}

func TestSortedListInsertionsConstraint(t *testing.T) {
	elems := []int{5, 1, 3, 6, 8, 3, 24, 65}
	s := SortedList{}
	for _, e := range elems {
		s.Insert(e)

		// elements are always considered sorted if the size is 0 or 1
		if len(elems) < 2 {
			continue
		}
		cur := s.elems[0]
		for i := 1; i < len(s.elems); i++ {
			if s.elems[i] < cur {
				t.Fatalf("elements at positions %d and %d are not sorted", i-1, i)
			}
			cur = s.elems[i]
		}
	}
}
