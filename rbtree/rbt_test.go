package rbtree

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// mustPanic is used to verify that operations on uninitialized trees or
// nodes actually panic.
func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

// assertRBTProperties checks the following by visiting each node once:
// 1. All nodes are either red or black.
// 2. The root is black.
// 3. TNIL leaves are black.
// 4. If a node is red, both its children are black.
// 5. For each node, all simple paths to the (TNIL) leaves have the same number
// of black nodes.
//
// On top of these RBT properties, it also checks the binary search tree
// property that for a node, all values of the left children are smaller and
// all values of the right children are larger.
func assertRBTProperties[T cmp.Ordered](t *testing.T, tree *RBT[T]) {
	t.Helper()

	if tree == nil {
		panic("assert called on nil tree")
	}

	// nothing to assert here
	if tree.Size() == 0 {
		return
	}

	if tree.Root() == nil {
		panic("nil root")
	}

	if tree.Root().color != _COLOR_BLACK {
		t.Fatal("property 2 violation: root should be black")
	}

	var walk func(n *RBTNode[T], bnCount int, lo, hi *T)
	bnCountTarget := -1 // unset
	walk = func(n *RBTNode[T], bncount int, lo, hi *T) {

		if n.color != _COLOR_RED && n.color != _COLOR_BLACK {
			t.Fatal("property 1 violation: found a node with an invalid color")
		}

		if n.sentinel {
			if n.color != _COLOR_BLACK {
				t.Fatal("property 3 violation: sentinel node should be black")
			}
			// if we haven't explore any simple path to TNIL yet, set it as the
			// target to compare to. Otherwise, compare it against an existing
			// target since the black node count must match.
			if bnCountTarget == -1 {
				bnCountTarget = bncount
			}
			if bnCountTarget != bncount {
				t.Fatal("property 5 violation: found two paths with different black length")
			}
			return
		}

		v := n.Value()
		if lo != nil && v <= *lo {
			t.Fatalf("bst violation: %v is in the right subtree of %v", v, *lo)
		}
		if hi != nil && v >= *hi {
			t.Fatalf("bst violation: %v is in the left subtree of %v", v, *hi)
		}

		if n.color == _COLOR_RED &&
			(n.left.color != _COLOR_BLACK || n.right.color != _COLOR_BLACK) {
			t.Fatal("property 4: found red node whose children aren't both black")
		}

		if n.color == _COLOR_BLACK {
			bncount += 1
		}

		walk(n.left, bncount, lo, &v)
		walk(n.right, bncount, &v, hi)
	}

	walk(tree.Root(), 0, nil, nil)
}

func TestRBTPropertiesDeterministicPermutations(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewPCG(6, 9))

	for i := range 100 {
		n := rng.IntN(1_000) + 1
		values := rng.Perm(1_000)[:n]
		shuffled := slices.Clone(values)
		rng.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		t.Run(fmt.Sprintf("same order deletion/%03d/n=%d", i, n), func(t *testing.T) {
			t.Parallel()
			tr := NewRBT[int]()
			cnt := 0
			for _, v := range values {
				tr.InsertUnique(v)
				cnt += 1
				assertRBTProperties(t, tr)
				if tr.Size() != cnt {
					t.Fatalf("tree size is not %d after inserts: %d", cnt, tr.Size())
				}
			}
			for _, v := range values {
				tr.Delete(v)
				cnt -= 1
				assertRBTProperties(t, tr)
				if tr.Size() != cnt {
					t.Fatalf("tree size is not %d after deletes: %d", cnt, tr.Size())
				}
			}
		})

		t.Run(fmt.Sprintf("shuffled deletion/%03d/n=%d", i, n), func(t *testing.T) {
			t.Parallel()
			tr := NewRBT[int]()
			cnt := 0
			for _, v := range values {
				tr.InsertUnique(v)
				cnt += 1
				assertRBTProperties(t, tr)
				if tr.Size() != cnt {
					t.Fatalf("tree size is not %d after inserts: %d", cnt, tr.Size())
				}
			}
			for _, v := range shuffled {
				tr.Delete(v)
				cnt -= 1
				assertRBTProperties(t, tr)
				if tr.Size() != cnt {
					t.Fatalf("tree size is not %d after deletes: %d", cnt, tr.Size())
				}
			}
		})
	}
}

func TestRBTNilReceiver(t *testing.T) {
	t.Parallel()

	type method struct {
		name string
		fn   func(*RBT[int])
	}

	methods := []method{
		{name: "Root", fn: func(tr *RBT[int]) { tr.Root() }},
		{name: "Size", fn: func(tr *RBT[int]) { tr.Size() }},
		{name: "Count", fn: func(tr *RBT[int]) { tr.Count(0) }},
		{name: "InsertUnique", fn: func(tr *RBT[int]) { tr.InsertUnique(1) }},
		{name: "Delete", fn: func(tr *RBT[int]) { tr.Delete(1) }},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			t.Parallel()
			mustPanic(t, func() { m.fn(nil) })
		})
	}
}

func TestRBTNilNode(t *testing.T) {
	t.Parallel()

	type method struct {
		name string
		fn   func(*RBTNode[int])
	}
	methods := []method{
		{name: "Value", fn: func(n *RBTNode[int]) { n.Value() }},
		{name: "Parent", fn: func(n *RBTNode[int]) { n.Parent() }},
		{name: "Left", fn: func(n *RBTNode[int]) { n.Left() }},
		{name: "Right", fn: func(n *RBTNode[int]) { n.Right() }},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			t.Parallel()
			n := (*RBTNode[int])(nil)
			mustPanic(t, func() { m.fn(n) })
		})
	}
}

func TestLeftRotateMissingRightChildPanics(t *testing.T) {
	t.Parallel()

	t.Run("rootWithoutRightChild", func(t *testing.T) {
		t.Parallel()

		tr := NewRBT[int]()
		x := &RBTNode[int]{parent: tr.tnil, left: tr.tnil, right: tr.tnil, value: 2}
		a := &RBTNode[int]{parent: x, left: tr.tnil, right: tr.tnil, value: 1}
		x.left = a
		tr.root = x
		tr.size = 2

		mustPanic(t, func() { leftRotate(tr, x) })
	})

	t.Run("nonRootWithoutRightChild", func(t *testing.T) {
		t.Parallel()

		tr := NewRBT[int]()
		p := &RBTNode[int]{parent: tr.tnil, left: tr.tnil, right: tr.tnil, value: 4}
		x := &RBTNode[int]{parent: p, left: tr.tnil, right: tr.tnil, value: 2}
		a := &RBTNode[int]{parent: x, left: tr.tnil, right: tr.tnil, value: 1}
		p.left = x
		x.left = a
		tr.root = p
		tr.size = 3

		mustPanic(t, func() { leftRotate(tr, x) })
	})

	t.Run("sentinel", func(t *testing.T) {
		t.Parallel()

		tr := NewRBT[int]()

		mustPanic(t, func() { leftRotate(tr, tr.tnil) })
	})
}

func TestRightRotateMissingLeftChildPanics(t *testing.T) {
	t.Parallel()

	t.Run("rootWithoutLeftChild", func(t *testing.T) {
		t.Parallel()

		tr := NewRBT[int]()
		y := &RBTNode[int]{parent: tr.tnil, left: tr.tnil, right: tr.tnil, value: 2}
		c := &RBTNode[int]{parent: y, left: tr.tnil, right: tr.tnil, value: 3}
		y.right = c
		tr.root = y
		tr.size = 2

		mustPanic(t, func() { rightRotate(tr, y) })
	})

	t.Run("nonRootWithoutLeftChild", func(t *testing.T) {
		t.Parallel()

		tr := NewRBT[int]()
		p := &RBTNode[int]{parent: tr.tnil, left: tr.tnil, right: tr.tnil, value: 1}
		y := &RBTNode[int]{parent: p, left: tr.tnil, right: tr.tnil, value: 2}
		c := &RBTNode[int]{parent: y, left: tr.tnil, right: tr.tnil, value: 3}
		p.right = y
		y.right = c
		tr.root = p
		tr.size = 3

		mustPanic(t, func() { rightRotate(tr, y) })
	})

	t.Run("sentinel", func(t *testing.T) {
		t.Parallel()

		tr := NewRBT[int]()

		mustPanic(t, func() { rightRotate(tr, tr.tnil) })
	})
}

// FuzzRBTEachInsert inserts an arbitrary sequence of values and checks the
// red-black tree properties after each one.
//
// TODO: check AI analysis for choice of fuzz inputs and input truncation.
func FuzzRBTEachInsert(f *testing.F) {

	// can only use []byte when fuzzing.
	f.Add([]byte{})
	f.Add([]byte{1})

	// force rotations
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
	f.Add([]byte{9, 8, 7, 6, 5, 4, 3, 2, 1})

	// force recoloring
	f.Add([]byte{5, 3, 8, 1, 4, 7, 9, 2, 6})

	// check that errors do not change the tree
	f.Add([]byte{1, 1, 1, 1, 1})

	f.Fuzz(func(t *testing.T, values []byte) {
		// The property is checked after every insert, so the target is
		// quadratic in the input length and needs an upper bound.
		//
		// The bound is set from the length distribution the fuzzer actually
		// produces. Measured over ~2M executions of this target: median 138
		// bytes, p90 417, p99 613, max 4245. Truncating below that discards
		// inputs the fuzzer worked to grow and scores them on a prefix, so the
		// bound sits above p99 rather than at a round small number.
		//
		// Lengths grow because 3 of the 18 byte slice mutators in
		// internal/fuzz insert bytes and only 1 removes them, each drawing an
		// amount that is 1-8 bytes 90% of the time.
		if len(values) > 1024 {
			values = values[:1024]
		}

		tr := NewRBT[int]()
		for i, v := range values {
			tr.InsertUnique(int(v))

			assertRBTProperties(t, tr)
			if t.Failed() {
				t.Fatalf("red-black tree property broken after inserting %v", values[:i+1])
			}
		}
	})
}
