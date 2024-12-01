package avl

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

type Node[T constraints.Ordered] struct {
	value  T
	left   *Node[T]
	right  *Node[T]
	parent *Node[T]
	height int
}

type AvlTree[T constraints.Ordered] struct {
	root *Node[T]
}

func (node *Node[T]) balanceFactor() int {
	leftHeight, rightHeight := -1, -1
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	return rightHeight - leftHeight
}

// %% Public methods %%

func (tree *AvlTree[T]) PrintTree(node *Node[T]) {
	if node == nil {
		return
	}
	tree.PrintTree(node.left)
	fmt.Println(node.value)
	tree.PrintTree(node.right)
}

func NewAvlTree[T constraints.Ordered]() *AvlTree[T] {
	return &AvlTree[T]{root: nil}
}

func (tree *AvlTree[T]) Add(value T) {
	newNode := Node[T]{value: value, height: 0}
	if tree.root == nil {
		tree.root = &newNode
		return
	}

	var parent *Node[T]
	next := tree.root
	for next != nil {
		parent = next
		if value < next.value {
			next = next.left
		} else {
			next = next.right
		}
	}

	if value < parent.value {
		parent.left = &newNode
	} else {
		parent.right = &newNode
	}
	newNode.parent = parent

	for parent != nil {
		tree.rebalance(parent)
		parent = parent.parent
	}
}

// %% Private methods %%

func (node *Node[T]) rotateLeft() *Node[T] {
	child := node.right
	node.right = child.left
	if node.right != nil {
		node.right.parent = node
	}
	child.left = node
	node.parent = child
	node.updateHeight()
	child.updateHeight()
	return child
}

func (node *Node[T]) rotateRight() *Node[T] {
	child := node.left
	node.left = child.right
	if node.left != nil {
		node.left.parent = node
	}
	child.right = node
	node.parent = child
	node.updateHeight()
	child.updateHeight()
	return child
}

func (node *Node[T]) updateHeight() {
	if node == nil {
		return
	}
	leftHeight, rightHeight := -1, -1
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	node.height = int(math.Max(float64(leftHeight), float64(rightHeight))) + 1
}

func (tree *AvlTree[T]) inorderTraverse(node *Node[T], queue *[]T) {
	if node == nil {
		return
	}
	tree.inorderTraverse(node.left, queue)
	*queue = append(*queue, node.value)
	tree.inorderTraverse(node.right, queue)
}

func (tree *AvlTree[T]) getTreeValues(node *Node[T]) []T {
	if node == nil {
		return []T{}
	}
	values := []T{}
	tree.inorderTraverse(tree.root, &values)
	return values
}

func (tree *AvlTree[T]) rebalance(node *Node[T]) {
	nodeBalance := node.balanceFactor()
	if math.Abs(float64(nodeBalance)) <= 1 {
		node.updateHeight()
		return
	}
	nodeParent := node.parent
	var newSubtreeRoot *Node[T]

	if nodeBalance < -1 {
		if node.left.balanceFactor() > 0 {
			node.left = node.left.rotateLeft()
			node.left.parent = node
		}
		newSubtreeRoot = node.rotateRight()
	} else {
		if node.right.balanceFactor() < 0 {
			node.right = node.right.rotateRight()
			node.right.parent = node
		}
		newSubtreeRoot = node.rotateLeft()
	}
	newSubtreeRoot.parent = nodeParent
	tree.replaceChild(nodeParent, node, newSubtreeRoot)
}

func (tree *AvlTree[T]) replaceChild(parent *Node[T], child *Node[T], replacement *Node[T]) {
	if parent == nil {
		tree.root = replacement
	} else if parent.left == child {
		parent.left = replacement
	} else {
		parent.right = replacement
	}
}
