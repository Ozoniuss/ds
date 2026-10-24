# RBTree

Provides a red black tree implementation in Go. Oriented towards competitive programming or the implementation of other data structures such as ordered sets and ordered multisets.

The library exposes a minimal public interface. In general, the way it thinks about modelling trees looks something like below. Note that this is not an actual interface exposed by the rbtree package, primarily because there isn't really a reason to abstract it under an interface, but also for some technical reasons (see appendix). My intention with the snippet below is to rather present it as documentation.

```go
// Tree represents the possible operations on binary search trees. Various tree
// types (e.g. regular BST, balanced BST, red black tree etc.) implement this
// interface.
//
// To match semantics of list and map Go containers, a tree that is initialized
// directly via the type itself should be a valid tree with size 0.
type Tree[T cmp.Ordered] interface {
	// Retrieve the Root of this tree. Returns nil for a tree that had no nodes
	// inserted to it.
	Root() Node[T]
	// Return the number of nodes in the tree.
	Size() int
	// Count returns the number of elements with the given value present in the
	// tree.
	Count(value T) int
	// InsertUnique adds an element to the tree, if the value wasn't previously
	// inserted. 
	InsertUnique(value T)
	// InsertEqual adds an element to the tree regardless whether the value was
	// previously inserted or not. A new node should be allocated for that value.
	InsertEqual(value T)
	// Delete a value from the binary search tree. Implementations should
    // guarantee the value no longer exists or its count was decreased by 1 in
    // the case of  trees accepting multiple values, once this method was
    // called. Otherwise, the method must panic.
	Delete(value T)
}

// Node represents a node in a binary search tree. It is read-only.
//
// All writes to the tree happen via the Tree object, since writes defined on
// Nodes don't have a clear expected behavior. Consider the following:
//
//	  8
//	 / \
//	4   9
//
// If writes were allowed on Nodes, writing "6" via the node "9" should not
// place it as a child of "9" since that would violate the BST property of the
// entire tree. Deletes are even more confusing: it may or may not delete the
// node on which the operation is called on. Therefore, a write via any node
// that is part of the tree should be treated the same way as a write via the
// root node. For me, it's simpler to think of the write operation as "write a
// value to the tree" rather than "write a value to this node that is part of a
// tree".
// 
// Similarly, retrieving the tree size and root via the node has an ambiguous 
// meaning. These values can either be the same for all nodes part of a tree,
// but could also mean the size and root of the tree object rooted at that node.
// Having two abstractions makes the intent clear and doesn't require
// distinguishing between the "root" node and other nodes.
type Node[T cmp.Ordered] interface {
	// Value returns the value stored in the node.
	Value() T
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

### What are the problems with the above interface?

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

There are ways to work around it. For example, you can express your interfaces as follows. To be honest, I did not find it obvious to understand why:

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

### Did you consider other tree libraries?

Yes, but I built mine for several reasons. Main one was that I wanted a minimal API with the shape that I liked and similar semantics to `list` and `map` types (`nil` panics, an initialized object is valid even when empty). I wanted to have full control over how the API evolves in the future. Other reasons were having a pure Go implementation with no dependencies, a mechanism for bundling the complete code in a single file, no logging from the library, writing my own tests and benchmarks and just the fun of implementing it myself.

Note that I'm not against providing more functionality for the tree library. For example, formatting trees is something I found very useful (see past [README](https://github.com/Ozoniuss/ds/blob/445a3540e2aba1f7e839ce6e146ff06c0a76a34a/tree/README.md)). But I wanted to package those as extensions (currently work in progress) and not have it part of the core library.

Here are some other implementations:

- [sortedset](https://github.com/wangjia184/sortedset) -  I didn't like the API, test suite is minimal and has other dependencies. I found a few others that seem to be mirrors, like [sortedset-go](https://github.com/zavitax/sortedset-go) and another [sortedset](https://github.com/matthew-hartman/sortedset).
- [go-datastructures](https://github.com/Workiva/go-datastructures) - Provides much more than I need, has dependencies, also not really happy with the tree interface.
- [go-set](https://github.com/hashicorp/go-set/tree/main) - This one is from Hashicorp and is a bit closer to what I wanted, but the public interface still has methods I didn't want (e.g. [Intersect](https://github.com/hashicorp/go-set/blob/main/treeset.go#L390C22-L390C31)), it also has dependencies and I wanted more test coverage.

Honestly, I also haven't really researched in depth. There may be some that I missed which would have matched my requirements.
