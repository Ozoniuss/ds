package tree

import (
	"cmp"
	"testing"
)

func TestEqual(t *testing.T) {
	type testcase[T cmp.Ordered] struct {
		name  string
		t1    Tree[T]
		t2    Tree[T]
		equal bool
	}

	testcases := []testcase[int]{
		{
			name:  "nil",
			t1:    nil,
			t2:    nil,
			equal: true,
		},
		{
			name:  "empty",
			t1:    NewBST[int](),
			t2:    NewBST[int](),
			equal: true,
		},
		{
			name:  "nilAndEmpty",
			t1:    nil,
			t2:    NewBST[int](),
			equal: false,
		},
		{
			name: "nilAndTree",
			t1:   nil,
			t2: &BST[int]{
				root: &BSTNode[int]{value: 3},
				size: 1,
			},
			equal: false,
		},
		{
			name: "equal roots",
			t1: &BST[int]{
				root: &BSTNode[int]{value: 3},
				size: 1,
			},
			t2: &BST[int]{
				root: &BSTNode[int]{value: 3},
				size: 1,
			},
			equal: true,
		},
		{
			name: "equal trees",
			t1: &BST[int]{
				root: &BSTNode[int]{value: 3, left: &BSTNode[int]{
					value: 2,
				}},
				size: 2,
			},
			t2: &BST[int]{
				root: &BSTNode[int]{value: 3, left: &BSTNode[int]{
					value: 2,
				}},
				size: 2,
			},
			equal: true,
		},
		// following are ai generated, could be simplified but good
		// enough for now
		{
			name:  "differentRootValue",
			t1:    &BST[int]{root: &BSTNode[int]{value: 3}, size: 1},
			t2:    &BST[int]{root: &BSTNode[int]{value: 4}, size: 1},
			equal: false,
		},
		{
			name: "emptyAndSingle",
			t1:   NewBST[int](),
			t2: &BST[int]{
				root: &BSTNode[int]{value: 1},
				size: 1,
			},
			equal: false,
		},
		{
			name: "sameValuesDifferentStructure",
			t1: func() Tree[int] {
				root := &BSTNode[int]{value: 3}
				left := &BSTNode[int]{value: 2, parent: root}
				root.left = left
				return &BST[int]{root: root, size: 2}
			}(),
			t2: func() Tree[int] {
				root := &BSTNode[int]{value: 2}
				right := &BSTNode[int]{value: 3, parent: root}
				root.right = right
				return &BST[int]{root: root, size: 2}
			}(),
			equal: false,
		},
		{
			name: "mirrorNotEqual",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 4}
				l1 := &BSTNode[int]{value: 3, parent: r}
				l2 := &BSTNode[int]{value: 2, parent: l1}
				r.left = l1
				l1.left = l2
				return &BST[int]{root: r, size: 3}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 4}
				rr1 := &BSTNode[int]{value: 3, parent: r}
				rr2 := &BSTNode[int]{value: 2, parent: rr1}
				r.right = rr1
				rr1.right = rr2
				return &BST[int]{root: r, size: 3}
			}(),
			equal: false,
		},
		{
			name: "threeLevelEqual",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rg := &BSTNode[int]{value: 7, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				lr := &BSTNode[int]{value: 4, parent: l}
				rl := &BSTNode[int]{value: 6, parent: rg}
				rr := &BSTNode[int]{value: 8, parent: rg}
				r.left, r.right = l, rg
				l.left, l.right = ll, lr
				rg.left, rg.right = rl, rr
				return &BST[int]{root: r, size: 7}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rg := &BSTNode[int]{value: 7, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				lr := &BSTNode[int]{value: 4, parent: l}
				rl := &BSTNode[int]{value: 6, parent: rg}
				rr := &BSTNode[int]{value: 8, parent: rg}
				r.left, r.right = l, rg
				l.left, l.right = ll, lr
				rg.left, rg.right = rl, rr
				return &BST[int]{root: r, size: 7}
			}(),
			equal: true,
		},
		{
			name: "threeLevelValueMismatch",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rg := &BSTNode[int]{value: 7, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				lr := &BSTNode[int]{value: 4, parent: l}
				rl := &BSTNode[int]{value: 6, parent: rg}
				rr := &BSTNode[int]{value: 8, parent: rg}
				r.left, r.right = l, rg
				l.left, l.right = ll, lr
				rg.left, rg.right = rl, rr
				return &BST[int]{root: r, size: 7}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rg := &BSTNode[int]{value: 7, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				lr := &BSTNode[int]{value: 4, parent: l}
				rl := &BSTNode[int]{value: 6, parent: rg}
				rr := &BSTNode[int]{value: 9, parent: rg} // value differs here
				r.left, r.right = l, rg
				l.left, l.right = ll, lr
				rg.left, rg.right = rl, rr
				return &BST[int]{root: r, size: 7}
			}(),
			equal: false,
		},
		{
			name: "fiveNodeEqual",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rl := &BSTNode[int]{value: 8, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				rrl := &BSTNode[int]{value: 6, parent: rl}
				r.left, r.right = l, rl
				l.left = ll
				rl.left = rrl
				return &BST[int]{root: r, size: 5}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rl := &BSTNode[int]{value: 8, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				rrl := &BSTNode[int]{value: 6, parent: rl}
				r.left, r.right = l, rl
				l.left = ll
				rl.left = rrl
				return &BST[int]{root: r, size: 5}
			}(),
			equal: true,
		},
		{
			name: "fiveNodeSizeMismatch",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rl := &BSTNode[int]{value: 8, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				rrl := &BSTNode[int]{value: 6, parent: rl}
				r.left, r.right = l, rl
				l.left = ll
				rl.left = rrl
				return &BST[int]{root: r, size: 5}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l := &BSTNode[int]{value: 3, parent: r}
				rl := &BSTNode[int]{value: 8, parent: r}
				ll := &BSTNode[int]{value: 2, parent: l}
				r.left, r.right = l, rl
				l.left = ll
				// missing rl.left
				return &BST[int]{root: r, size: 4}
			}(),
			equal: false,
		},
		{
			name: "deepLeftEqual",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l1 := &BSTNode[int]{value: 4, parent: r}
				l2 := &BSTNode[int]{value: 3, parent: l1}
				l3 := &BSTNode[int]{value: 2, parent: l2}
				r.left = l1
				l1.left = l2
				l2.left = l3
				return &BST[int]{root: r, size: 4}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l1 := &BSTNode[int]{value: 4, parent: r}
				l2 := &BSTNode[int]{value: 3, parent: l1}
				l3 := &BSTNode[int]{value: 2, parent: l2}
				r.left = l1
				l1.left = l2
				l2.left = l3
				return &BST[int]{root: r, size: 4}
			}(),
			equal: true,
		},
		{
			name: "deepLeftRightMismatch",
			t1: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				l1 := &BSTNode[int]{value: 4, parent: r}
				l2 := &BSTNode[int]{value: 3, parent: l1}
				l3 := &BSTNode[int]{value: 2, parent: l2}
				r.left = l1
				l1.left = l2
				l2.left = l3
				return &BST[int]{root: r, size: 4}
			}(),
			t2: func() Tree[int] {
				r := &BSTNode[int]{value: 5}
				r1 := &BSTNode[int]{value: 6, parent: r}
				r2 := &BSTNode[int]{value: 7, parent: r1}
				r3 := &BSTNode[int]{value: 8, parent: r2}
				r.right = r1
				r1.right = r2
				r2.right = r3
				return &BST[int]{root: r, size: 4}
			}(),
			equal: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if Equal(tc.t1, tc.t2) != tc.equal {
				t.Errorf("trees are not equal\nt1:\n%s\nt2:\n%s\n", FormatTree(tc.t1, FormatHorizontalSquared), FormatTree(tc.t2, FormatHorizontalSquared))
			}
		})
	}
}
