package avl

import (
	"slices"
	"testing"
)

// Tests the AVL tree with integer values, covering all basic rotation cases
func TestIntegerTree(t *testing.T) {
	cases := [][]int{
		{1, 2, 3},             // Right-Right case
		{1, 3, 2},             // Right-Left case
		{2, 1, 3},             // Left-Left case
		{2, 3, 1},             // Left-Right case
		{3, 1, 2},             // Left-Right case
		{3, 2, 1},             // Left-Left case
		{10, 20, 30, 40, 50},  // Multiple Right-Right rotations
		{10, 20, 30, 50, 40},  // Right-Right followed by Right-Left
		{30, 20, 10, 5, 1},    // Multiple Left-Left rotations
		{30, 20, 10, 1, 5},    // Left-Left followed by Left-Right
		{5, 4, 6, 3, 7, 2, 8}, // Mixed rotations
	}

	for i, testCase := range cases {
		tree := NewAvlTree[int]()

		for _, value := range testCase {
			tree.Add(value)
		}
		actual := tree.getTreeValues(tree.root)
		expected := slices.Clone(testCase)
		slices.Sort(expected)
		if !slices.Equal(actual, expected) {
			t.Errorf("Test case %d: tree.Add(...) = %v; want %v", i, actual, expected)
		}
	}
}

// Tests the AVL tree with string values
func TestStringTree(t *testing.T) {
	cases := [][]string{
		{"chickpeas", "tahini", "za'atar"}, // Ordered
		{"za'atar", "tahini", "chickpeas"}, // Reversed
		{"tahini", "za'atar", "chickpeas"}, // Mixed
		{"a", "b", "c", "d", "e"},          // Sequence
		{"e", "d", "c", "b", "a"},          // Reversed sequence
	}

	for i, testCase := range cases {
		tree := NewAvlTree[string]()

		for _, value := range testCase {
			tree.Add(value)
		}

		actual := tree.getTreeValues(tree.root)
		expected := slices.Clone(testCase)
		slices.Sort(expected)

		if !slices.Equal(actual, expected) {
			t.Errorf("Test case %d: tree.Add(...) = %v; want %v", i, actual, expected)
		}
	}
}

// Tests the AVL tree with floating-point values
func TestFloatTree(t *testing.T) {
	cases := [][]float64{
		{1.1, 2.2, 3.3},             // Ordered
		{3.3, 2.2, 1.1},             // Reversed
		{2.2, 1.1, 3.3},             // Mixed
	}

	for i, testCase := range cases {
		tree := NewAvlTree[float64]()

		for _, value := range testCase {
			tree.Add(value)
		}

		actual := tree.getTreeValues(tree.root)
		expected := slices.Clone(testCase)
		slices.Sort(expected)

		if !slices.Equal(actual, expected) {
			t.Errorf("Test case %d: tree.Add(...) = %v; want %v", i, actual, expected)
		}
	}
}
