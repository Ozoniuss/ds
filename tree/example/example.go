package main

import (
	"fmt"

	tree "github.com/Ozoniuss/ds/tree"
)

func main() {

	t1 := tree.NewRBT[int]()

	t1.Insert(4)
	t1.Insert(1)
	t1.Insert(8)
	t1.Insert(6)
	t1.Insert(9)
	t1.Insert(5)
	t1.Insert(11)
	t1.Insert(13)

	fmt.Println(tree.FormatTree(t1, tree.FormatLinuxTree))
	fmt.Println(tree.FormatTree(t1, tree.FormatHorizontal))
	fmt.Println(tree.FormatTree(t1, tree.FormatHorizontalSquared))

	t2 := tree.NewRBT[string]()
	t2.Insert("5sfhskfuceskjvsdnkvjkdsn")
	t2.Insert("1dbfalkfbdslkjfbadslkfbl")
	t2.Insert("3dbfalkfbdslkjfbadslkfbl")
	t2.Insert("8dsbflkjsdbfjzhklsdbfljkds")
	t2.Insert("7dsbflkjsdbfjzhklsdbfljkds")
	t2.Insert("9dsbflkjsdbfjzhklsdbfljkds")

	fmt.Println(tree.FormatTree(t2, tree.FormatLinuxTree))
	fmt.Println(tree.FormatTree(t2, tree.FormatHorizontal))
	fmt.Println(tree.FormatTree(t2, tree.FormatHorizontalSquared))

	t3 := tree.NewRBT[int]()

	t3.Insert(26)
	t3.Insert(17)
	t3.Insert(41)
	t3.Insert(14)
	t3.Insert(10)
	t3.Insert(16)
	t3.Insert(7)
	t3.Insert(12)
	t3.Insert(3)

	t3.Insert(21)
	t3.Insert(19)
	t3.Insert(23)
	t3.Insert(20)
	t3.Insert(47)
	t3.Insert(30)
	t3.Insert(28)
	t3.Insert(38)
	t3.Insert(35)
	t3.Insert(39)

	fmt.Println(tree.FormatTree(t3, tree.FormatLinuxTree))
	fmt.Println(tree.FormatTree(t3, tree.FormatHorizontal))
	fmt.Println(tree.FormatTree(t3, tree.FormatHorizontalSquared))

	t3.Delete(39)
	t3.Delete(35)
	t3.Delete(38)
	t3.Delete(28)
	fmt.Println(tree.FormatTree(t3, tree.FormatLinuxTree))
	fmt.Println(tree.FormatTree(t3, tree.FormatHorizontal))
	fmt.Println(tree.FormatTree(t3, tree.FormatHorizontalSquared))

}
