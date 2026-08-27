# Tree

Provides a standardized way of working with binary search trees in Go, as well as several binary search tree implementations. Oriented towards competitive programming or the implementation of other data structures such as ordered sets and ordered multisets.

While this library includes common tree algorithms, they are offered as an extension. The core library exposes a minimal set of operations standardized across all common tree types, which can be used to implement tree algorithms.

The following interfaces abstracts Trees and Nodes, respectively:

```go
type Tree[T cmp.Ordered] interface {
	Root() Node[T]
	Size() int
	Count(value T) int
	Insert(value T) error
	Delete(value T) error
}

type Node[T cmp.Ordered] interface {
	Value() T
	Count() int
	Parent() Node[T]
	Left() Node[T]
	Right() Node[T]
}

```

## Available implementations

- `tree.BST` (regular tree)

This implementation preserves the binary search tree properties upon addition and deletion, but doesn't do any kind of rebalancing. Operations should have O(lg n) complexity on avegrage, with worst case complexity going up to O(n) (e.g. if adding consecutive numbers to the tree).

- `tree.RBT` (red black tree)

This implements a self-balancing binary search tree using a technique called "red-black tree". Because a tree of this kind has a height smaller or equal than 2lg(n+1) where n is the number of nodes, operations have O(lg n) both average and worst case complexity.

## Installation

Install with

```bash
go get github.com/Ozoniuss/ds
```

## Usage

The library is generic, allowing you to define trees on any value set that can be ordered.

```go
func main() {
    // Instantiate a regular BST with int values.
    t1 := tree.NewBST[int]()

    // Insert values in the tree
    t1.Insert(4)
    t1.Insert(1)
    t1.Insert(8)
    t1.Insert(6)
    t1.Insert(9)
    t1.Insert(5)
    t1.Insert(11)
    t1.Insert(13)

    // Print the tree
    fmt.Println(t1)
}
```

Output:

```
  4         
 / \        
1   8       
   / \      
  6   9     
 /     \    
5       11  
         \  
          13
```

```go
func main() {
    // Instantiate a regular BST with string values.
    t2 := tree.NewBST[string]()
    
    // Insert values in the tree
    t2.Insert("5sfhskfuceskjvsdnkvjkdsn")
    t2.Insert("1dbfalkfbdslkjfbadslkfbl")
    t2.Insert("3dbfalkfbdslkjfbadslkfbl")
    t2.Insert("8dsbflkjsdbfjzhklsdbfljkds")
    t2.Insert("7dsbflkjsdbfjzhklsdbfljkds")
    t2.Insert("9dsbflkjsdbfjzhklsdbfljkds")

    // Print the tree using a more compact formatting
    fmt.Println(tree.FormatTree(t2, tree.FormatHorizontalSquared))
}
```

Output:

```
                     5sfhskfuceskjvsdnkvjkdsn                                    
           ┌────────────────────┬────────────────────┐                           
1dbfalkfbdslkjfbadslkfbl                 8dsbflkjsdbfjzhklsdbfljkds              
           └┐                          ┌─────────────┬─────────────┐             
 3dbfalkfbdslkjfbadslkfbl  7dsbflkjsdbfjzhklsdbfljkds  9dsbflkjsdbfjzhklsdbfljkds
```

## Printing

The code for printing the tree horizontally is ported from @billvanyo's [tree_printer](https://github.com/billvanyo/tree_printer/tree/master) Java library, excluding the options to print multiple trees and allowing direction agnostic branches (that is, using a character like `|` to link the parent to the child). If you'd like to understand how it works, I did my best to document the printer source code.

For more examples on how the horizontal printer behaves, visit @billvanyo's [github repository](https://github.com/billvanyo/tree_printer/tree/master)

There is also a vertical tree formatter inspired from the Linux `tree` utility that I implemented myself. See the `FormatTree` options for how to specify the formatter.

> [!TIP]
> Red-black trees are printed with colored nodes. 
