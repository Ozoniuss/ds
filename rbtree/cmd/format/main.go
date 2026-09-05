package main

import (
	"fmt"

	tree "github.com/Ozoniuss/ds/rbtree"
)

func main() {

	t1 := tree.NewRBT[int]()
	t1.Insert(1)
	t1.Insert(3)
	t1.Insert(4)
	t1.Insert(2)

	t4 := tree.NewRBT[int]()

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
	t4.Insert(3111111)

	t4.Insert(2111111)
	t4.Insert(1911111)
	t4.Insert(2311111)
	t4.Insert(2411111)
	t4.Insert(2511111)
	t4.Insert(2011111)
	t4.Insert(470011111)
	t4.Insert(3011111)
	t4.Insert(2811111)
	t4.Insert(3811111)
	t4.Insert(3511111)
	t4.Insert(3911111)

	lines := tree.BuildLines(t4.Root(), 2, true)
	for _, l := range lines {
		fmt.Println(l)
	}
	lines = tree.BuildLines(t1.Root(), 2, true)
	for _, l := range lines {
		fmt.Println(l)
	}
}
