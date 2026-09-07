package main

import (
	"fmt"

	"github.com/Ozoniuss/ds/rbtree"
)

var ttyColorReset = "\033[0m"
var ttyRed = "\033[31m"

func main() {

	t1 := rbtree.NewRBT[int]()
	t1.Insert(1)
	t1.Insert(3)
	t1.Insert(4)
	t1.Insert(2)

	t4 := rbtree.NewRBT[int]()

	t4.Insert(26111111)
	t4.Insert(17111111)
	t4.Insert(41111111)
	t4.Insert(15111111)
	t4.Insert(10111111)
	t4.Insert(16111111)
	t4.Insert(7111111)
	t4.Insert(12111111)
	t4.Insert(13111111)
	// t4.Insert(14)
	// t4.Insert(3111111)

	// t4.Insert(2111111)
	// t4.Insert(1911111)
	// t4.Insert(2311111)
	// t4.Insert(2411111)
	// t4.Insert(2511111)
	// t4.Insert(2011111)
	// t4.Insert(470011111)
	// t4.Insert(3011111)
	// t4.Insert(2811111)
	// t4.Insert(3811111)
	// t4.Insert(3511111)
	// t4.Insert(3911111)

	lines := rbtree.BuildLines(t4.Root(), rbtree.SHAPE_SQUARE, 2, true, true)
	for _, l := range lines {
		fmt.Println(l)
	}
	lines = rbtree.BuildLines(t4.Root(), rbtree.SHAPE_TREE, 2, true, false)
	for _, l := range lines {
		fmt.Println(l)
	}
	lines = rbtree.BuildLines(t1.Root(), rbtree.SHAPE_SQUARE, 3, false, true)
	for _, l := range lines {
		fmt.Println(l)
	}
	lines = rbtree.BuildLines(t1.Root(), rbtree.SHAPE_TREE, 2, true, false)
	for _, l := range lines {
		fmt.Println(l)
	}
}
