package tree

import (
	"cmp"
	"errors"
)

/*
	Implementation based on "Introduction to Algorithms by Thomas H. Cormen,
	Charles E. Leiserson, Ronald L. Rivest, and Clifford Stein."
*/

const (
	_COLOR_RED   = "red"
	_COLOR_BLACK = "black"
)

// A red-black tree is a binary search tree with one extra bit of storage per node: its
// color, which can be either RED or BLACK. By constraining the node colors on any
// simple path from the root to a leaf, red-black trees ensure that no such path is more
// than twice as long as any other, so that the tree is approximately balanced.
//
// A red-black tree is a binary tree that satisﬁes the following red-black properties:
// 1. Every node is either red or black.
// 2. The root is black.
// 3. Every leaf (TNIL -- see below) is black.
// 4. If a node is red, then both its children are black.
// 5. For each node, all simple paths from the node to descendant leaves contain the
// same number of black nodes.
//
// A red-black tree with n internal nodes has height at most 2*lg(n+1).
//
// Typically, a tree will only contain unique values. By default, RBT enforces
// that constraint but can be configured to accept multiple values.
//
// There is a distinction between a tree object and a node object in order to
// have a clearer interface. All write operations and methods that return information
// about the entire tree are only defined on the tree object. The node object is
// read-only and only returns information about the node itself. Particularly:
//
// - All writes to the tree happen via the Tree object. Node objects are read-only.
// Consider the following tree:
//
//	  8
//	 / \
//	4   9
//
// If writes were allowed on Nodes, writing "6" via the node "9" should not
// place it as a child of "9" since that would violate the BST property of the
// entire tree. Deletes are even more confusing: it may or may not delete the
// node on which the operation is called on. Therefore, a write via any node that
// is part of the tree should be treated the same way as a write via the root node.
// It's simpler to thinkof the write operation as "write a value to the tree"
// rather than "write a value to this node that is part of a tree".
//
// - Retrieving the tree size and root via the node has an ambiguous meaning.
// These values can either be the same for all nodes part of a tree, but could
// also mean the size and root of the tree object rooted at that node. Having
// two abstractions makes the intent clear and doesn't require distinguishing
// between the "root" node and other nodes.
type RBT[T cmp.Ordered] struct {
	// Represents the root node of a tree. An uninitialized tree will set this to
	// nil, whereas an initialized tree will set this to tnil.
	root *RBTNode[T]
	// Represents the number of nodes in the tree, excluding sentinel nodes (see
	// below).
	size int
	// As a matter of convenience in dealing with boundary conditions in red-black
	// tree code, we use a single sentinel to represent nil-valued nodes. This sentinel
	// is an object with the same attributes as an ordinary node in the tree. Its
	// color attribute is BLACK, and its other attributes have no meaning so are
	// set to the default value of their type. We refer to this sentinel as "tnil".
	//
	// If the left or right child of a node is empty, we set it to tnil instead
	// of nil. Similarly, we set the parent of the root node to tnil.
	//
	// Note that all nil-valued nodes are set to the same tnil sentinel, meaning
	// that the parent of the sentinel is also unset.
	tnil *RBTNode[T]
	// Whether the tree allows duplicate values or not.
	allowDuplicates bool
}

type RBTOpts[T cmp.Ordered] func(t *RBT[T])

func WithAllowedDuplicates[T cmp.Ordered]() RBTOpts[T] {
	return func(t *RBT[T]) {
		t.allowDuplicates = true
	}
}

// NewRBT returns an initialized red black tree.
func NewRBT[T cmp.Ordered](opts ...RBTOpts[T]) *RBT[T] {
	tnil := sentinel[T]()
	t := &RBT[T]{
		size: 0,
		root: tnil,
		tnil: tnil,
	}
	for _, o := range opts {
		o(t)
	}

	// assert that options do not change root and tnil
	if t.tnil != tnil || t.root != tnil {
		panic("options must not modify root and tnil during initialization")
	}

	return t
}

// Root() returns the root of the tree.
func (t *RBT[T]) Root() Node[T] {
	panicIfNilRBT(t)

	// tnil must not be exposed externally
	if t.root == t.tnil {
		return nil
	}

	return t.root
}

// Size returns the number of nodes in the tree.
func (t *RBT[T]) Size() int {
	panicIfNilRBT(t)

	return t.size
}

// Count returns the number of times a value was added to the tree. If the tree
// is not configured to accept multiple values, this operation will always return
// either 0 or 1.
func (t *RBT[T]) Count(value T) int {
	panicIfNilRBT(t)

	if t.root == t.tnil {
		return 0
	}

	c := t.root
	for c != t.tnil {
		if value < c.value {
			c = c.left
		} else if value > c.value {
			c = c.right
		} else {
			return c.Count()
		}
	}
	return 0
}

// Insert adds a value to the tree. If the tree is not configured to accept multiple
// values, it will return an error when the same value is inserted twice.
func (t *RBT[T]) Insert(value T) error {
	panicIfNilRBT(t)

	if t.root == t.tnil {
		t.root = &RBTNode[T]{
			parent: t.tnil,
			left:   t.tnil,
			right:  t.tnil,
			value:  value,
			color:  _COLOR_BLACK,
			count:  1,
		}
		t.size = 1
		return nil
	}

	y := t.tnil
	x := t.root
	z := &RBTNode[T]{
		value: value,
		count: 1,
	}

	for x != t.tnil {
		y = x
		if z.value < x.value {
			x = x.left
		} else if z.value > x.value {
			x = x.right
		} else {
			if !t.allowDuplicates {
				return errors.New("value already exists")
			}
			x.count += 1
			return nil
		}
	}
	z.parent = y

	if y == t.tnil {
		t.root = z
	} else if z.value < y.value {
		y.left = z
	} else {
		y.right = z
	}
	z.left = t.tnil
	z.right = t.tnil
	z.color = _COLOR_RED

	insertFixup(t, z)
	t.size++

	return nil
}

func (t *RBT[T]) Delete(value T) error {
	panicIfNilRBT(t)

	if t.root == t.tnil {
		return errors.New("value not found")
	}

	// find z
	z := t.root
	for z != t.tnil {
		if value < z.value {
			z = z.left
		} else if value > z.value {
			z = z.right
		} else {
			// removing a value with count of 1 may require transplant operations
			if t.allowDuplicates && z.count > 1 {
				z.count -= 1
				return nil
			}
			break
		}
	}
	if z == t.tnil {
		return errors.New("value not found")
	}

	y := z
	yorigcolor := y.color
	var x *RBTNode[T]

	if z.left == t.tnil {
		x = z.right
		rbtransplant(t, z, z.right)
	} else if z.right == t.tnil {
		x = z.left
		rbtransplant(t, z, z.left)
	} else {
		y = treeMinimumRbt(t, z.right)
		yorigcolor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			rbtransplant(t, y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		rbtransplant(t, z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	if yorigcolor == _COLOR_BLACK {
		rbDeleteFixup(t, x)
	}
	t.size--
	return nil
}

func (t *RBT[T]) String() string {
	panicIfNilRBT(t)

	return FormatTree(t, string(FormatHorizontal))
}

// RBTNode represents a tree node.
type RBTNode[T cmp.Ordered] struct {
	parent   *RBTNode[T]
	left     *RBTNode[T]
	right    *RBTNode[T]
	value    T
	color    string
	count    int
	sentinel bool
}

func (n *RBTNode[T]) Value() T {
	panicIfNilRBTNodeOrSentinel(n)

	return n.value
}

func (n *RBTNode[T]) Count() int {
	panicIfNilRBTNodeOrSentinel(n)

	return n.count
}

func (n *RBTNode[T]) Parent() Node[T] {
	panicIfNilRBTNodeOrSentinel(n)

	if n.parent.isSentinel() {
		return nil
	}
	return n.parent
}

func (n *RBTNode[T]) Left() Node[T] {
	panicIfNilRBTNodeOrSentinel(n)

	if n.left.isSentinel() {
		return nil
	}
	return n.left
}

func (n *RBTNode[T]) Right() Node[T] {
	panicIfNilRBTNodeOrSentinel(n)

	if n.right.isSentinel() {
		return nil
	}
	return n.right
}

func (n *RBTNode[T]) isSentinel() bool {
	return n.sentinel
}

// ttycolor is used for colored terminal output.
func (n *RBTNode[T]) ttycolor() string {
	panicIfNilNode(n)

	return n.color
}

func sentinel[T cmp.Ordered]() *RBTNode[T] {
	return &RBTNode[T]{
		color:    _COLOR_BLACK,
		sentinel: true,
	}
}

// panicIfNilRBT will panic if the tree is nil (hence, the tree is uninitialized).
func panicIfNilRBT[T cmp.Ordered](n *RBT[T]) {
	if n == nil {
		panic("nil rbt")
	}
}

// panicIfNilRBTNodeOrSentinel will panic if the current node is nil or a sentinel.
// Sentinels are included because they are an implementation detail to simplify
// edge case handling in algorithms, but in essence they still represent nil
// nodes and should be exposed as nil from public methods.
//
// Note that this method is typically called for public methods, as algorithms
// use the internal representation directly.
func panicIfNilRBTNodeOrSentinel[T cmp.Ordered](n *RBTNode[T]) {
	if n == nil {
		panic("nil node")
	} else if n.sentinel {
		panic("nil node")
	}
}

func leftRotate[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) {
	y := x.right
	x.right = y.left
	if y.left != t.tnil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == t.tnil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func rightRotate[T cmp.Ordered](t *RBT[T], y *RBTNode[T]) {
	x := y.left
	y.left = x.right
	if x.right != t.tnil {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == t.tnil {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	x.right = y
	y.parent = x
}

// transplant replaces one subtree with another subtree
func rbtransplant[T cmp.Ordered](t *RBT[T], u *RBTNode[T], v *RBTNode[T]) {
	// u is root
	if u.parent == t.tnil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func treeMinimumRbt[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) *RBTNode[T] {
	for x.left != t.tnil {
		x = x.left
	}
	return x
}

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
	t.root.color = _COLOR_BLACK
}

func rbDeleteFixup[T cmp.Ordered](t *RBT[T], x *RBTNode[T]) {
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
