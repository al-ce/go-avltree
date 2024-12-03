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

// %%% Node public methods %%%

// %%% Tree public methods %%%

func NewAvlTree[T constraints.Ordered]() *AvlTree[T] {
	return &AvlTree[T]{root: nil}
}

func (tree *AvlTree[T]) PrintTree(node *Node[T]) {
	if node == nil {
		return
	}
	tree.PrintTree(node.left)
	fmt.Println(node.value)
	tree.PrintTree(node.right)
}

// Returns a slice of the tree's values in-order. Appends to the provided
// pointer to a slice. If the pointer is nil, a new slice is created.
func (tree *AvlTree[T]) InorderTraverse(node *Node[T], queue *[]T) []T {
	if queue == nil {
		queue = &[]T{}
	}
	if node == nil {
		return *queue
	}
	*queue = tree.InorderTraverse(node.left, queue)
	*queue = append(*queue, node.value)
	*queue = tree.InorderTraverse(node.right, queue)
	return *queue
}

// Returns a bool indicating whether the value exists in the tree
func (tree *AvlTree[T]) Contains(value T) bool {
	return tree.getNodeByValue(value) != nil
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

// Remove a node from the tree by value lookup.
// Returns true on successful removal, false if value was not found.
func (tree *AvlTree[T]) Remove(value T) bool {
	node := tree.getNodeByValue(value)
	if node == nil { // value was not found in the tree
		return false
	}

	parent := node.parent
	var replacement *Node[T]

	// Action node is the node where the rebalancing will start
	actionNode := parent

	// Case 1: two children, replace with in-order successor, then rebalance
	if node.left != nil && node.right != nil {

		// Find in-order successor (move right once then left all the way down)
		successor := node.right
		for successor.left != nil {
			successor = successor.left
		}

		// Assign the children of the node to remove to the successor node
		successor.left = node.left
		// If the successor wasn't the right node, then we need to give it a
		// right node. Otherwise, the successor's right node will be nil
		if successor != node.right {
			// We moved all the way down to the left.
			// If the successor has a right node, put that right node in the
			// successor's current spot
			successor.parent.left = successor.right
			if successor.right != nil {
				successor.right.parent = successor.parent
			}
			// The successor now has both the node's children as its own
			successor.right = node.right
		}
		// Complete the child->parent relationship
		node.left.parent = successor
		node.right.parent = successor

		replacement = successor

		actionNode = replacement.parent
	} else {
		// Case 2: one or no children, replace with existing child
		if node.left == nil {
			replacement = node.right
		} else if node.right == nil {
			replacement = node.left
		}
	}

	tree.replaceChild(parent, node, replacement)
	if replacement != nil {
		replacement.parent = parent
	}

	// Rebalance from the parent of the node that got moved, up to the root
	for actionNode != nil {
		tree.rebalance(actionNode)
		actionNode = actionNode.parent
	}

	return true
}

// %% Private methods %%

// %%% Node private methods %%%

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

// %%% Tree private methods %%%

func (tree *AvlTree[T]) getNodeByValue(value T) *Node[T] {
	if tree.root == nil {
		return nil
	}

	node := tree.root
	for node != nil {
		if node.value == value {
			return node
		}
		if value < node.value {
			node = node.left
		} else {
			node = node.right
		}
	}
	return nil
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

func (tree *AvlTree[T]) replaceRoot(newRoot *Node[T]) {
	tree.root = newRoot
	if newRoot != nil {
		newRoot.parent = nil
	}
}

func (tree *AvlTree[T]) replaceChild(parent *Node[T], child *Node[T], replacement *Node[T]) {
	// If we are replacing the root node
	if parent == nil {
		tree.replaceRoot(replacement)
		return
	}

	if parent.left == child {
		parent.left = replacement
	} else {
		parent.right = replacement
	}
}
