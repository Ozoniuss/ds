package tree

import "cmp"

// Tree represents the possible operations on binary search trees. Various tree
// types (e.g. regular BST, balanced BST, red black tree etc.) implement this
// interface.
//
// Calling any method on a nil tree should panic.
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
// that does not satisfy the binary search tree property. The primary reason
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
