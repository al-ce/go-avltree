# go-avltree

A golang [AVL tree](https://en.wikipedia.org/wiki/AVL_tree) implementation.
Accepts `Ordered` (integer, float, and string) types as defined by the [constraints](https://pkg.go.dev/golang.org/x/exp/constraints) package.

## Example

```go
package main

import (
	"fmt"

	avl "github.com/al-ce/go-avltree"
)

func main() {
	tree := avl.NewAvlTree[int]()

	for i := 1; i <= 10; i++ {
		tree.Add(i)
	}

	tree.Remove(5)
	tree.Contains(5) // false
	tree.IsEmpty()   // false

	minVal, _ := tree.GetMin() // returns error if tree is empty
	fmt.Println(minVal)  // 1

	fmt.Println(tree.Size()) // 9

	fmt.Println("Traverse (get slice)")
	inorder := tree.InOrderTraverse() // returns []T
	for _, val := range inorder {
		fmt.Printf("%v ", val) // [1 2 3 4 6 7 8 9 10]
	}
	fmt.Println()

	// iter.Next() returns (val T, index int)
	// When index == -1, the iterator has reached the end of the tree
	iter := tree.NewIterator()

	fmt.Println("Iterator:")
	for {
		val, index := iter.Next()
		if index == -1 {
			break
		}
		fmt.Printf("%v ", val) // [1 2 3 4 6 7 8 9 10]
	}
	fmt.Println()

	tree.Clear()
	tree.IsEmpty() // true

	// String example
	stringTree := avl.NewAvlTree[string]()
	p := []string{"tahini", "za'atar", "chickpeas"}
	for _, v := range p {
		stringTree.Add(v)
	}
	ordered := stringTree.InOrderTraverse()
	fmt.Println(ordered) // [chickpeas tahini za'atar]
}
```
