package rbtree

import (
	"cmp"
)

/*
	Implementation based on "Introduction to Algorithms by Thomas H. Cormen,
	Charles E. Leiserson, Ronald L. Rivest, and Clifford Stein."
*/

const (
	_COLOR_RED   = "red"
	_COLOR_BLACK = "black"
)

// RBT is a binary search tree that satisfies the following properties:
//
// 1. Every node is either red or black.
// 2. The root is black.
// 3. Every leaf (TNIL -- see below) is black.
// 4. If a node is red, then both its children are black.
// 5. For each node, all simple paths from the node to descendant leaves contain the
// same number of black nodes.
type RBT[T cmp.Ordered] struct {
	// root must be lazily initialized when the tree is not created via its
	// constructor.
	root *RBTNode[T]
	// size excludes sentinel nodes
	size int
	// sentinel value all internal leaves and the root point to in this tree
	tnil *RBTNode[T]
}

// NewRBT returns an initialized red black tree.
func NewRBT[T cmp.Ordered]() *RBT[T] {
	t := &RBT[T]{}
	t.lazyInit()

	return t
}

// Root() returns the root of the tree.
func (t *RBT[T]) Root() *RBTNode[T] {
	// tnil must not be exposed externally
	if t.root == t.tnil {
		return nil
	}

	return t.root
}

// Size returns the number of internal nodes in the tree.
func (t *RBT[T]) Size() int {
	return t.size
}

// For searching algorithms, it helps to know the following case is possible
// when inserting 5, 1, 7, 5
//
//   5a
//  /  \
// 1    7
//     /
//   5b
//
// Or, when inserting 5,5,5
//
//   5b
//  /  \
// 5a   5c
//

// Count returns the number of times a value was added to the tree.
func (t *RBT[T]) Count(value T) int {
	if t.root == t.tnil {
		return 0
	}

	c := t.root
	first := t.tnil

	for c != t.tnil {
		if value <= c.value {
			if value == c.value {
				first = c
			}
			c = c.left
		} else {
			c = c.right
		}
	}
	c = first

	cnt := 0
	for c != t.tnil && c.value == value {
		cnt += 1
		c = treeSuccessor(t, c)
	}

	return cnt
}

// InsertUnique adds a value to the tree, but only if the tree doesn't already
// contain this value.
func (t *RBT[T]) InsertUnique(value T) {
	t.lazyInit()

	insert(t, value, true)
}

// InsertEqual adds a value to the tree, regardless if the value already exists
// or not.
func (t *RBT[T]) InsertEqual(value T) {
	t.lazyInit()

	insert(t, value, false)
}

func insert[T cmp.Ordered](t *RBT[T], value T, unique bool) {
	if t.root == t.tnil {
		t.root = &RBTNode[T]{
			parent: t.tnil,
			left:   t.tnil,
			right:  t.tnil,
			value:  value,
			color:  _COLOR_BLACK,
		}
		t.size = 1
		return
	}

	// start from the root and find which internal leaf z would become if
	// inserted
	y := t.tnil // prev
	x := t.root // curr
	z := &RBTNode[T]{
		value: value,
	}
	for x != t.tnil {
		y = x
		if z.value < x.value {
			x = x.left
		} else {
			if unique && x.value == z.value {
				return
			}
			x = x.right
		}
	}
	// at this point y is an internal leaf, and we set z as its child
	z.parent = y

	if y == t.tnil {
		panic("new node should never have a tnil parent if root already exists")
	} else if z.value < y.value {
		y.left = z
	} else {
		y.right = z
	}
	z.left = t.tnil
	z.right = t.tnil
	z.color = _COLOR_RED

	// after regular insertion, rebalance and recolor if needed
	insertFixup(t, z)
	t.size++
}

// insertFixup will perform rebalancing after insertion, when needed
func insertFixup[T cmp.Ordered](t *RBT[T], z *RBTNode[T]) {

	for z.parent.color == _COLOR_RED {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == _COLOR_RED {
				z.parent.color = _COLOR_BLACK
				y.color = _COLOR_BLACK
				z.parent.parent.color = _COLOR_RED
				z = z.parent.parent
			} else if z == z.parent.right {
				z = z.parent
				leftRotate(t, z)
			} else {
				z.parent.color = _COLOR_BLACK
				z.parent.parent.color = _COLOR_RED
				rightRotate(t, z.parent.parent)
			}
		} else {
			// mirror of the situation above
			y := z.parent.parent.left
			if y.color == _COLOR_RED {
				z.parent.color = _COLOR_BLACK
				y.color = _COLOR_BLACK
				z.parent.parent.color = _COLOR_RED
				z = z.parent.parent
			} else if z == z.parent.left {
				z = z.parent
				rightRotate(t, z)
			} else {
				z.parent.color = _COLOR_BLACK
				z.parent.parent.color = _COLOR_RED
				leftRotate(t, z.parent.parent)
			}
		}
	}
	// we can only get in this case when terminating on branch 1, if z becomes
	// the root (and the root's parent is TNIL which is black)
	t.root.color = _COLOR_BLACK
}

// DeleteOne removes a value from the tree. If there are duplicates, it will
// remove the node of the last value that was inserted.
func (t *RBT[T]) DeleteOne(value T) {
	t.lazyInit()

	// find the latest occurence of value
	z := t.root
	c := t.tnil

	for z != t.tnil {
		if value < z.value {
			z = z.left
		} else {
			if z.value == value {
				c = z
			}
			z = z.right
		}
	}
	z = c

	// nothing to delete
	if z == t.tnil {
		return
	}

	// at this point, z is an internal leaf. delete like in a regular tree, then
	// fixup to re-color and re-balance.
	//
	// y tracks the node that is moved or removed from the tree, x tracks the
	// subtree that moves to its location.
	y := z
	yorigcolor := y.color
	var x *RBTNode[T]

	if z.left == t.tnil {
		// if z has at most one internal child, we transplant (note that we may
		// transplant a sentinel child too).
		x = z.right
		rbtransplant(t, z, z.right) // z.right may be tnil
	} else if z.right == t.tnil {
		x = z.left
		rbtransplant(t, z, z.left) // z.left may be tnil
	} else {
		y = treeMinimum(t, z.right)
		yorigcolor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			rbtransplant(t, y, y.right) // y.right may be tnil
			y.right = z.right
			y.right.parent = y
		}

		rbtransplant(t, z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yorigcolor == _COLOR_BLACK {
		rbtDeleteFixup(t, x) // x may be tnil, see above what can be tnil
	}
	t.size--
}

// rbtDeleteFixup handles recoloring and rebalancing if a black node was removed
// from the tree.
func rbtDeleteFixup[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) {
	for x != t.root && x.color == _COLOR_BLACK {

		if x == x.parent.left {
			w := x.parent.right
			if w.color == _COLOR_RED {
				w.color = _COLOR_BLACK
				x.parent.color = _COLOR_RED
				leftRotate(t, x.parent)
				w = x.parent.right
			}
			if w.left.color == _COLOR_BLACK && w.right.color == _COLOR_BLACK {
				w.color = _COLOR_RED
				x = x.parent
			} else if w.right.color == _COLOR_BLACK {
				w.left.color = _COLOR_BLACK
				w.color = _COLOR_RED
				rightRotate(t, w)
				w = x.parent.right
			} else {
				w.color = x.parent.color
				x.parent.color = _COLOR_BLACK
				w.right.color = _COLOR_BLACK
				leftRotate(t, x.parent)
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == _COLOR_RED {
				w.color = _COLOR_BLACK
				x.parent.color = _COLOR_RED
				rightRotate(t, x.parent)
				w = x.parent.left
			}
			if w.right.color == _COLOR_BLACK && w.left.color == _COLOR_BLACK {
				w.color = _COLOR_RED
				x = x.parent
			} else if w.left.color == _COLOR_BLACK {
				w.right.color = _COLOR_BLACK
				w.color = _COLOR_RED
				leftRotate(t, w)
				w = x.parent.left
			} else {
				w.color = x.parent.color
				x.parent.color = _COLOR_BLACK
				w.left.color = _COLOR_BLACK
				rightRotate(t, x.parent)
				x = t.root
			}
		}
	}
	x.color = _COLOR_BLACK
}

// RBTNode represents a tree node.
type RBTNode[T cmp.Ordered] struct {
	parent   *RBTNode[T]
	left     *RBTNode[T]
	right    *RBTNode[T]
	value    T
	color    string
	sentinel bool
}

func (n *RBTNode[T]) Value() T {
	return n.value
}

func (n *RBTNode[T]) Parent() *RBTNode[T] {
	if n.parent.sentinel {
		return nil
	}
	return n.parent
}

func (n *RBTNode[T]) Left() *RBTNode[T] {
	if n.left.sentinel {
		return nil
	}
	return n.left
}

func (n *RBTNode[T]) Right() *RBTNode[T] {
	if n.right.sentinel {
		return nil
	}
	return n.right
}

func newSentinel[T cmp.Ordered]() *RBTNode[T] {
	return &RBTNode[T]{
		color:    _COLOR_BLACK,
		sentinel: true,
	}
}

// lazyInit is used to allow creating a RBT through the type directly (i.e. as
// t := &RBT{}).
func (t *RBT[T]) lazyInit() {
	if t.tnil == nil {
		t.tnil = newSentinel[T]()
		t.root = t.tnil
	}
}

// leftRotate brings the right child of x to its position. This operation preserves
// the binary search tree property.
//
// Note that this method moves x to its (original) left child's position.
//
// It is assumed y is not tnil.
//
//	    y                                         x
//	   / \        <---left  rotate(T,x)---       / \
//	  x   c                                     a   y
//	 / \           ---right rotate(T,y)--->        / \
//	a   b                                         b   c
func leftRotate[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) {
	y := x.right
	if y == t.tnil {
		panic("leftRotate: y must not be tnil")
	}

	// set the right node for x and its parents to x, which completes the
	// leaves on left side
	x.right = y.left
	if y.left != t.tnil {
		y.left.parent = x
	}

	// y becomes x
	y.parent = x.parent
	if x.parent == t.tnil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	// x becomes left child of y
	y.left = x
	x.parent = y
}

// rightRotate brings the left child of y to its position. This operation preserves
// the binary search tree property.
//
// Note that this method moves y to its (original) right child's position.
//
// It is assumed x is not tnil.
//
//	    y                                         x
//	   / \        <---left  rotate(T,x)---       / \
//	  x   c                                     a   y
//	 / \           ---right rotate(T,y)--->        / \
//	a   b                                         b   c
func rightRotate[T cmp.Ordered](t *RBT[T], y *RBTNode[T]) {
	x := y.left
	if x == t.tnil {
		panic("rightRotate: x must not be tnil")
	}

	// set the left node for y and its parents to y, which completes the
	// leaves on right side
	y.left = x.right
	if x.right != t.tnil {
		x.right.parent = y
	}

	// x becomes y
	x.parent = y.parent
	if y.parent == t.tnil {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	// y becomes right child of x
	x.right = y
	y.parent = x
}

// rbtransplant replaces one subtree u with another subtree v. Note that this
// only changes v's references, u remains unchanged.
func rbtransplant[T cmp.Ordered](t *RBT[T], u *RBTNode[T], v *RBTNode[T]) {
	if u == nil || v == nil {
		panic("rbtransplant: called on nil node")
	}
	// note that transplant may be called on a sentinel destination
	if u.sentinel {
		panic("rbtransplant: u (source) must not be a sentinel")
	}
	// u is root, v becomes root
	if u.parent == t.tnil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

// treeMinimum finds the minimum value (the leftmost internal node) of the
// subtree rooted at t.
func treeMinimum[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) *RBTNode[T] {
	for x.left != t.tnil {
		x = x.left
	}
	return x
}

// treeSuccessor returns the next node in the tree, or t.tnil if the current
// node is the rightmost internal node
func treeSuccessor[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) *RBTNode[T] {
	if x.right != t.tnil {
		return treeMinimum(t, x.right)
	}
	y := x.parent
	for y != t.tnil && x == y.right {
		x = y
		y = y.parent
	}
	return y
}
