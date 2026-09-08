package rbtree

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

// formatRBTForTests returns a compact representation of a tree that can be
// parsed back easily into the original tree. This representation is used to
// help understand what each operation on the RBT does.
func formatRBTForTests(root *RBT[int]) string {
	lines := BuildLines(root.Root(), SHAPE_SQUARE, 2, true, true)
	return strings.Join(lines, "\n")
}

// parseRBTFromTests rebuilds the tree from its string representation in tests.
// It leverages the property that we can rebuild a regular binary search tree
// by simply inserting its nodes in BFS order. While this doesn't always work
// for RBT, we can instead create the "shape" by performing regular BST insertions
// and then just apply the color to each node.
//
// The format used in tests has the property that lines will alternatively
// contain either just node labels or tree edges, therefore we only need to
// read the lines containing labels and parse them. Tests use positive integer
// values as labels and the color is encoded at the end of each number.
//
//	                      15111111[B]
//	           v---------------+---------------v
//	      10111111[R]                     26111111[R]
//	    v------+------v                 v------+------v
//	7111111[B]   12111111[B]       17111111[B]   41111111[B]
//	                   \               /
//	               13111111[R]   16111111[R]
//
// TODO: may be worth testing that the parser is actually correct too, for
// completeness.
func parseRBTFromTests(repr string) *RBT[int] {
	parts := strings.Split(repr, "\n")
	root := &RBT[int]{}
	// we actually don't care at all about edges, since inserting in BFS order
	// on a regular tree
	for i := 0; i < len(parts); i += 2 {
		numstrs := strings.FieldsSeq(parts[i])
		for n := range numstrs {
			num, err := strconv.Atoi(n[:len(n)-3])
			if err != nil {
				msg := "parseRBTFromTests: parsing number " + n
				panic(msg)
			}
			var color string
			colorstr := n[len(n)-2]
			switch colorstr {
			case 'R':
				color = _COLOR_RED
			case 'B':
				color = _COLOR_BLACK
			default:
				msg := "parseRBTFromTests: parsing color " + string(colorstr)
				panic(msg)
			}
			bstInsertNode(root, num, color)
		}
	}
	return root
}

// This test is used to prove that trees can be recreated by inserting nodes in
// BFS order in a regular tree then applying colors.
func TestBFSReinsertionReproducesTree(t *testing.T) {
	rng := rand.New(rand.NewPCG(6, 9))

	for i := range 5_000 {
		n := 1_000
		values := rng.Perm(n)

		t.Run(fmt.Sprintf("%03d/n=%d", i, n), func(t *testing.T) {
			tr := NewRBT[int]()
			for _, v := range values {
				tr.Insert(v)
			}

			reinserted := recreateFromBFS(bfsColoredValues(tr.Root()))
			if !equalTrees(tr, reinserted) {
				t.Fatalf("BFS reinsertion produced a different tree after inserting %v", values)
			}

			for _, v := range values {
				tr.Delete(v)
			}

			reinserted = recreateFromBFS(bfsColoredValues(tr.Root()))
			if !equalTrees(tr, reinserted) {
				t.Fatalf("BFS reinsertion produced a different tree after deleting %v", values)
			}
		})
	}
}

// TestRotate showcases how rotations work.
//
// Note that rotations do not change the tree colors, therefore resulting trees
// may violate the RBT properties. These are not checked in the rotation tests
// since it is the responsibility of the insertion algorithm to enfore them.
// Colors have been included to also show that rotations do not change colors.
func TestRotate(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name string

		x int
		y int

		leftRotateXresult  string
		rightRotateYresult string
	}

	cases := []testcase{
		{
			name: "x is the root and has only a right child",
			x:    1,
			y:    2,
			leftRotateXresult: `
   2[R]
   /
 1[B]`,
			rightRotateYresult: `
1[B]
  \
  2[R]`,
		},
		{
			name: "x is the root and has both children",
			x:    2,
			y:    4,

			leftRotateXresult: `
      4[B]
    v--+--v
   2[B]  5[R]
 v--+--v
1[B]  3[R]`,
			rightRotateYresult: `
   2[B]
 v--+--v
1[B]  4[B]
    v--+--v
   3[R]  5[R]`,
		},
		{
			name: "x is the left child of its parent and has only a right child",
			x:    2,
			y:    3,

			leftRotateXresult: `
      5[B]
    v--+--v
   3[B]  8[B]
   /
 2[R]`,
			rightRotateYresult: `
   5[B]
 v--+--v
2[R]  8[B]
  \
  3[B]`,
		},
		{
			name: "x is the right child of its parent and has only a right child",
			x:    8,
			y:    9,

			leftRotateXresult: `
   5[B]
 v--+--v
2[B]  9[B]
      /
    8[R]`,
			rightRotateYresult: `
   5[B]
 v--+--v
2[B]  8[R]
        \
        9[B]`,
		},
		{
			name: "x is the left child of its parent and has both children",
			x:    5,
			y:    7,

			leftRotateXresult: `
         10[B]
       v---+---v
      7[B]   15[B]
      /
    5[R]
    /
  3[B]`,
			rightRotateYresult: `
      10[B]
    v---+---v
   5[R]   15[B]
 v--+--v
3[B]  7[B]`,
		},
		{
			name: "x is the right child of its parent and has both children",
			x:    15,
			y:    17,

			leftRotateXresult: `
    10[B]
 v----+----v
1[B]     17[B]
          /
       15[R]
        /
     13[B]`,
			rightRotateYresult: `
   10[B]
 v---+---v
1[B]   15[R]
     v---+---v
   13[B]   17[B]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before := strings.TrimPrefix(tc.rightRotateYresult, "\n")
			after := strings.TrimPrefix(tc.leftRotateXresult, "\n")

			// Note that we can test both given leftRotate and rightRotate are each
			// other's inverse function.

			t.Run(fmt.Sprintf("x=%d/left rotate", tc.x), func(t *testing.T) {
				t.Parallel()

				tr := parseRBTFromTests(before)
				leftRotate(tr, findNode(tr, tc.x))
				if got := formatRBTForTests(tr); got != after {
					t.Fatalf("left rotate did not produce the expected tree:\ngot:\n%s\nwant:\n%s", got, after)
				}
			})

			t.Run(fmt.Sprintf("y=%d/right rotate", tc.y), func(t *testing.T) {
				t.Parallel()

				tr := parseRBTFromTests(after)
				rightRotate(tr, findNode(tr, tc.y))
				if got := formatRBTForTests(tr); got != before {
					t.Fatalf("right rotate did not produce the expected tree:\ngot:\n%s\nwant:\n%s", got, before)
				}
			})
		})
	}
}

// TestTransplant showcases how rbtransplant works.
//
// Note that transplant only rewires the parent's child pointer (or the tree
// root) to v and sets v's parent; u itself is left untouched. It does not
// preserve the BST or RBT properties, so resulting trees may be invalid. Colors
// have been included to show that transplant does not change them.
func TestTransplant(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name string

		u int
		v int
		// v is the tnil sentinel and the v field is ignored
		vSentinel bool

		tree             string
		transplantResult string
	}

	cases := []testcase{
		{
			name: "u is the root and v is an internal subtree",
			u:    20,
			v:    10,
			tree: `
       20[B]
     v---+---v
   10[B]   30[B]
 v---+---v
5[R]   15[R]`,
			transplantResult: `
   10[B]
 v---+---v
5[R]   15[R]`,
		},
		{
			name: "u is the left child of its parent and v is an internal subtree",
			u:    10,
			v:    15,
			tree: `
       20[B]
     v---+---v
   10[B]   30[B]
 v---+---v
5[R]   15[R]`,
			transplantResult: `
    20[B]
  v---+---v
15[R]   30[B]`,
		},
		{
			name: "u is the right child of its parent and v is an internal subtree",
			u:    30,
			v:    5,
			tree: `
       20[B]
     v---+---v
   10[B]   30[B]
 v---+---v
5[R]   15[R]`,
			transplantResult: `
      20[B]
     v--+--v
   10[B]  5[R]
 v---+---v
5[R]   15[R]`,
		},
		{
			name:      "u is the left child of its parent and v is a sentinel",
			u:         10,
			vSentinel: true,
			tree: `
       20[B]
     v---+---v
   10[B]   30[B]
 v---+---v
5[R]   15[R]`,
			transplantResult: `
20[B]
   \
  30[B]`,
		},
		{
			name:      "u is the right child of its parent and v is a sentinel",
			u:         30,
			vSentinel: true,
			tree: `
       20[B]
     v---+---v
   10[B]   30[B]
 v---+---v
5[R]   15[R]`,
			transplantResult: `
       20[B]
        /
     10[B]
   v---+---v
  5[R]   15[R]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			want := strings.TrimPrefix(tc.transplantResult, "\n")

			tr := parseRBTFromTests(strings.TrimPrefix(tc.tree, "\n"))
			v := tr.tnil
			if !tc.vSentinel {
				v = findNode(tr, tc.v)
			}
			rbtransplant(tr, findNode(tr, tc.u), v)

			if got := formatRBTForTests(tr); got != want {
				t.Fatalf("transplant did not produce the expected tree:\ngot:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func findNode(tr *RBT[int], value int) *RBTNode[int] {
	n := tr.root
	for n != tr.tnil {
		if value < n.value {
			n = n.left
		} else if value > n.value {
			n = n.right
		} else {
			return n
		}
	}
	panic(fmt.Sprintf("findNode: value %d not found in tree", value))
}

func equalTrees(a, b *RBT[int]) bool {
	var walk func(x, y *RBTNode[int]) bool
	walk = func(x, y *RBTNode[int]) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		if x.Value() != y.Value() || x.color != y.color {
			return false
		}
		return walk(x.Left(), y.Left()) && walk(x.Right(), y.Right())
	}
	return walk(a.Root(), b.Root())
}

type coloredValue struct {
	value int
	color string
}

func bfsColoredValues(root *RBTNode[int]) []coloredValue {
	if root == nil {
		return nil
	}

	var values []coloredValue
	queue := []*RBTNode[int]{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		values = append(values, coloredValue{value: n.Value(), color: n.color})
		if l := n.Left(); l != nil {
			queue = append(queue, l)
		}
		if r := n.Right(); r != nil {
			queue = append(queue, r)
		}
	}
	return values
}

func bstInsertNode(tr *RBT[int], value int, color string) {
	tr.lazyInit()

	if tr.root == tr.tnil {
		tr.root = &RBTNode[int]{
			parent: tr.tnil,
			left:   tr.tnil,
			right:  tr.tnil,
			value:  value,
			color:  color,
		}
		tr.size = 1
		return
	}

	y := tr.tnil
	x := tr.root
	z := &RBTNode[int]{value: value}
	for x != tr.tnil {
		y = x
		if z.value < x.value {
			x = x.left
		} else {
			x = x.right
		}
	}
	z.parent = y
	if z.value < y.value {
		y.left = z
	} else {
		y.right = z
	}
	z.left = tr.tnil
	z.right = tr.tnil
	z.color = color

	tr.size++
}

func recreateFromBFS(values []coloredValue) *RBT[int] {
	tr := NewRBT[int]()
	for _, v := range values {
		bstInsertNode(tr, v.value, v.color)
	}
	return tr
}
