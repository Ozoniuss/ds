package tree

import "cmp"

func panicIfNilTree[T cmp.Ordered](t Tree[T]) {
	if t == nil {
		panic("nil tree")
	}
}

func panicIfNilNode[T cmp.Ordered](n Node[T]) {
	if n == nil {
		panic("nil node")
	}
}
