package rbtree

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
)

type visualCase struct {
	name string

	buildBefore func(tr *RBT[int])
	op          func(tr *RBT[int])

	beforeDrawing string
	afterDrawing  string
}

func runVisualCase(t *testing.T, tc visualCase) {
	t.Helper()

	t.Run(tc.name, func(t *testing.T) {
		tr := NewRBT[int]()
		tc.buildBefore(tr)
		before := drawTree(tr)
		if got, want := strings.Join(before, "\n"), strings.TrimPrefix(tc.beforeDrawing, "\n"); got != want {
			t.Fatalf("before tree does not match its drawing:\ngot:\n%s\nwant:\n%s", got, want)
		}

		tc.op(tr)

		after := drawTree(tr)
		if got, want := strings.Join(after, "\n"), strings.TrimPrefix(tc.afterDrawing, "\n"); got != want {
			t.Fatalf("after tree does not match its drawing:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})
}

func drawTree(tr *RBT[int]) []string {
	if tr.Root() == nil {
		return nil
	}
	return BuildLines(tr.Root(), SHAPE_TREE, 2, true, false)
}

func TestLeftRotate(t *testing.T) {
	t.Parallel()

	cases := []visualCase{
		{
			name: "x=1 has only a right child",

			buildBefore: func(tr *RBT[int]) {
				x := &RBTNode[int]{parent: tr.tnil, left: tr.tnil, right: tr.tnil, value: 1}
				y := &RBTNode[int]{parent: x, left: tr.tnil, right: tr.tnil, value: 2}
				x.right = y
				tr.root = x
				tr.size = 2
			},

			op: func(tr *RBT[int]) {
				leftRotate(tr, tr.root)
			},

			beforeDrawing: `
1
 \
  2`,
			afterDrawing: `
  2
 /
1`,
		},
		{
			name: "x=2 and its right child both have two children",

			buildBefore: func(tr *RBT[int]) {
				x := &RBTNode[int]{parent: tr.tnil, value: 2}
				a := &RBTNode[int]{parent: x, left: tr.tnil, right: tr.tnil, value: 1}
				y := &RBTNode[int]{parent: x, value: 4}
				b := &RBTNode[int]{parent: y, left: tr.tnil, right: tr.tnil, value: 3}
				c := &RBTNode[int]{parent: y, left: tr.tnil, right: tr.tnil, value: 5}
				x.left, x.right = a, y
				y.left, y.right = b, c
				tr.root = x
				tr.size = 5
			},

			op: func(tr *RBT[int]) {
				leftRotate(tr, tr.root)
			},

			beforeDrawing: `
  2
 / \
1   4
   / \
  3   5`,
			afterDrawing: `
    4
   / \
  2   5
 / \
1   3`,
		},
	}

	for _, tc := range cases {
		runVisualCase(t, tc)
	}
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
