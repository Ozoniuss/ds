package main

import (
	"fmt"

	"github.com/Ozoniuss/ds/rbtree"
)

//     6
//    / \
//   2   13
//  / \   \
// 1   4   19

func explore(node *rbtree.RBTNode[int]) {
	if node == nil {
		return
	}
	// do something with the parent before exploration
	fmt.Println("before", node.Value())

	explore(node.Left())
	explore(node.Right())

	// do something with the parent after exploration
	fmt.Println("after", node.Value())
}

var maxHeight = 0

func explore2(node *rbtree.RBTNode[int], currentDetph int) (nodeCount int) {
	if node == nil {
		return 0
	}
	maxHeight = max(currentDetph, maxHeight)

	leftNc := explore2(node.Left(), currentDetph+1)
	rightNc := explore2(node.Right(), currentDetph+1)

	return leftNc + rightNc + 1
}

func main() {
	t := rbtree.NewRBT[int]()

	t.Insert(2)
	t.Insert(6)
	t.Insert(13)
	t.Insert(1)
	t.Insert(4)
	t.Insert(19)

	explore(t.Root())

}
