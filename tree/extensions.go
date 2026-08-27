package tree

import "cmp"

func equalSubtree[T cmp.Ordered](n1, n2 Node[T]) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	// avoid panics beyond this point
	if (n1 == nil && n2 != nil) || (n1 != nil && n2 == nil) {
		return false
	}

	return n1.Value() == n2.Value() && equalSubtree(n1.Left(), n2.Left()) && equalSubtree(n1.Right(), n2.Right())
}

// Equal returns whether two trees have identical shape and store the same set
// of values.
func Equal[T cmp.Ordered](t1, t2 Tree[T]) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	// avoid panics beyond this point
	if (t1 == nil && t2 != nil) || (t1 != nil && t2 == nil) {
		return false
	}

	if t1.Root() == nil && t2.Root() == nil {
		return true
	}

	return equalSubtree(t1.Root(), t2.Root())
}
