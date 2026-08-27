package tree

import (
	"cmp"
	"fmt"
)

type BST[T cmp.Ordered] struct {
	root *BSTNode[T]
	size int
}

// NewBST returns an initialized binary search tree.
func NewBST[T cmp.Ordered]() *BST[T] {
	return &BST[T]{
		size: 0,
		root: nil,
	}
}

func (t *BST[T]) Root() Node[T] {
	panicIfNilTree(t)

	if t.root == nil {
		return nil
	}
	return t.root
}

func (t *BST[T]) Size() int {
	panicIfNilTree(t)

	return t.size
}

func (t *BST[T]) Count(value T) int {
	panicIfNilTree(t)

	if t.root == nil {
		return 0
	}
	c := t.root
	for c != nil {
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

func (t *BST[T]) Insert(value T) error {
	panicIfNilTree(t)

	if t.root == nil {
		t.root = &BSTNode[T]{
			parent: nil,
			left:   nil,
			right:  nil,
			value:  value,
		}
		t.size++
		return nil
	}

	var y *BSTNode[T] = nil
	c := t.root
	for c != nil {
		y = c
		if value < c.value {
			c = c.left
		} else if value > c.value {
			c = c.right
		} else {
			return fmt.Errorf("value already exists")
		}
	}
	if value < y.value {
		y.left = &BSTNode[T]{
			parent: y,
			left:   nil,
			right:  nil,
			value:  value,
		}
	} else {
		y.right = &BSTNode[T]{
			parent: y,
			left:   nil,
			right:  nil,
			value:  value,
		}
	}
	t.size++
	return nil
}

func (t *BST[T]) Delete(value T) error {
	panicIfNilTree(t)

	if t.root == nil {
		return fmt.Errorf("value not found")
	}

	// find node with that value
	z := t.root
	for z != nil {
		if value < z.value {
			z = z.left
		} else if value > z.value {
			z = z.right
		} else if value == z.value {
			break
		}
	}
	if z == nil {
		return fmt.Errorf("value not found")
	}

	if z.left == nil {
		transplant(t, z, z.right)
		t.size--
		return nil
	} else if z.right == nil {
		transplant(t, z, z.left)
		t.size--
		return nil
	}

	y := treeMinimum(z.right)
	if y.parent != z {
		transplant(t, y, y.right)
		y.right = z.right
		y.right.parent = y
	}
	transplant(t, z, y)
	y.left = z.left
	y.left.parent = y

	t.size--
	return nil
}

func (t *BST[T]) String() string {
	panicIfNilTree(t)

	return FormatTree(t, string(FormatHorizontal))
}

type BSTNode[T cmp.Ordered] struct {
	parent *BSTNode[T]
	left   *BSTNode[T]
	right  *BSTNode[T]
	value  T
}

func (n *BSTNode[T]) Parent() Node[T] {
	panicIfNilNode(n)

	if n.parent == nil {
		return nil
	}
	return n.parent
}

func (n *BSTNode[T]) Left() Node[T] {
	panicIfNilNode(n)

	if n.left == nil {
		return nil
	}
	return n.left
}

func (n *BSTNode[T]) Right() Node[T] {
	panicIfNilNode(n)

	if n.right == nil {
		return nil
	}
	return n.right
}

func (n *BSTNode[T]) Value() T {
	panicIfNilNode(n)

	return n.value
}

func (n *BSTNode[T]) Count() int {
	panicIfNilNode(n)

	return 1
}

// Tree helpers

// transplant replaces one subtree with another subtree
func transplant[T cmp.Ordered](t *BST[T], u *BSTNode[T], v *BSTNode[T]) {
	// u is root
	if u.parent == nil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

func treeMinimum[T cmp.Ordered](x *BSTNode[T]) *BSTNode[T] {
	y := x
	for y.left != nil {
		y = y.left
	}
	return y
}
