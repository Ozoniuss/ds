# RBTree

Provides a red black tree implementation in Go. Oriented towards competitive programming or the implementation of other data structures such as ordered sets and ordered multisets.

The library exposes a minimal public interface. In general, the way it thinks about modelling trees looks something like below. Note that this is not an actual interface exposed by the rbtree package, primarily because there isn't really a reason to abstract it under an interface, but also for some technical reasons (see appendix).

```go
// Tree represents the possible operations on binary search trees. Various tree
// types (e.g. regular BST, balanced BST, red black tree etc.) implement this
// interface.
//
// Calling any method on a nil tree should panic.
//
// To match semantics of list and map Go containers, a tree that is initialized
// directly via the type itself should be a valid tree with size 0.
type Tree[T cmp.Ordered] interface {
	// Retrieve the Root of this tree. Returns nil for a tree that had no nodes
	// inserted to it.
	Root() Node[T]
	// Return the number of nodes in the tree.
	Size() int
	// Count returns the number of elements with value `value` present in the
	// tree. Implementations that require unique values will either return 0
	// or 1.
	//
	// Use this method if you need to determine whether a value belongs to the
	// tree or not.
	Count(value T) int
	// Insert a value to the binary search tree. Implementations that require
	// storing unique values will return an error if that value already exists.
	// Callers may choose to ignore the error if they just want to ensure the
	// value is present in the tree.
	Insert(value T) error
	// Delete a value from the binary search tree. Implementations that allow
	// storing multiple values of the same type should only remove one occurence
	// of the value.
	// Callers may choose to ignore the error if they just want to ensure the
	// value is deleted from a tree supporting only unique values.
	Delete(value T) error
}

// Node represents a node in a binary search tree. It is read-only.
//
// Nodes are not meant to be modified directly as all operations are expected
// to be performed using the Tree object. This is to prevent producing a tree
// that does not satisfy the binary search tree properties. The primary reason
// it exists in the first place is to allow defining custom traversal algorithms
// on the tree.
//
// Calling any method on a nil node should panic.
type Node[T cmp.Ordered] interface {
	// Value returns the value stored in the node.
	Value() T
	// Count returns the number of elements equal to the node's `value` are
	// present in the tree. Since the nodes are always bound to only one tree
	// (i.e. nodes are not created without a tree directly or transitively
	// pointing to them), implementations that require unique values will always
	// return 1.
	Count() int
	// Parent returns the parent node. This should return nil for the root node
	// and a non-nil value for all other nodes.
	//
	// This method should be used to check if an individual node is the root of
	// the tree.
	Parent() Node[T]
	// Left returns the left child of the node. This should return nil if there
	// is no left child.
	Left() Node[T]
	// Right returns the right child of the node. This should return nil if there
	// is no right child.
	Right() Node[T]
}
```

## Installation

Install with

```bash
go get github.com/Ozoniuss/ds
```

## Appendix

- What are the problems with the above interface?

The primary issue is that Go does not support [covariant types](https://go.dev/doc/faq#covariant_types). When designing the `RBTree` type I decided that I want the public methods to actually return the types themselves and not an interface, e.g.

```Go
// Root() returns the root of the tree.
func (t *RBT[T]) Root() *RBTNode[T] {
	panicIfNilRBT(t)

	// tnil must not be exposed externally
	if t.root == t.tnil {
		return nil
	}

	return t.root
}
```

But, this does not implement `Tree` as defined above.

There are ways to work around it. For example, you can express your interfaces as follows. To be honest, I did not find find it obvious to understand why:

```Go
type Node[T cmp.Ordered, N any] interface {
	Value() T
	Count() int
	Parent() N
	Left() N
	Right() N
}

type Tree[T cmp.Ordered, N Node[T, N]] interface {
	Root() N
	Size() int
	Count(value T) int
	Insert(value T) error
	Delete(value T) error
}

func _[T cmp.Ordered, N Node[T, N]]() {
	var _ Tree[T, *RBTNode[T]] = (*RBT[T])(nil)
}
```

I think it becomes a bit more clear once you understand that you don't actually need to supply a type argument for `N` at all when creating the node (but you still need to supply it when creating the tree). The example below makes it easies to see why: `Point` returns itself, so it acts as a type argument too. The difference from the `Node[T cmp.Ordered, N any]` interface is that you still need to instantiate `T` for the element types, and having multiple type parameters just looks a bit more confusing overall.

```Go
type Cloner[C any] interface {
	Clone() C
}

// Point does not need
type Point struct{ X, Y int }
func (p Point) Clone() Point { return Point{p.X, p.Y} }

func CloneAll[C Cloner[C]](xs []C) []C {
	out := make([]C, 0, len(xs))
	for _, x := range xs {
		out = append(out, x.Clone())
	}
	return out
}
```

I recommend this article about [self-referential interfaces](https://go.dev/blog/generic-interfaces) (funnily enough, this article also models trees similarly when it comes to the defined types) which goes a bit more in depth. Still, I wasn't happy with the fact that you needed to supply a node type. After all, nodes are really more of an internal aspect of a tree. For a while I even considered not having them part of the interface. I wasn't happy with how the interface looked in general and also not convinced there was any benefit for having an extra layer of indirection.

- Did you consider other tree libraries?

Yes, but I built mine for several reasons.
