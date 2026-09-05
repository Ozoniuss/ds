package main

import (
	"fmt"

	tree "github.com/Ozoniuss/ds/rbtree"
)

func main() {
	t4 := tree.NewRBT[int]()

	t4.Insert(26)
	t4.Insert(17)
	t4.Insert(41)
	t4.Insert(15)
	t4.Insert(10)
	t4.Insert(16)
	t4.Insert(7)
	t4.Insert(12)
	t4.Insert(13)
	// t4.Insert(14)
	t4.Insert(3)

	t4.Insert(21)
	t4.Insert(19)
	t4.Insert(23)
	t4.Insert(24)
	t4.Insert(25)
	t4.Insert(20)
	t4.Insert(4700)
	t4.Insert(30)
	t4.Insert(28)
	t4.Insert(38)
	t4.Insert(35)
	t4.Insert(39)

	lines := tree.BuildLines(t4.Root(), 2)
	for _, l := range lines {
		fmt.Println(l)
	}
}
