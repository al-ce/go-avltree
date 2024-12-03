package avl

import (
	"slices"
	"testing"
)

func rangeWithSteps(start, end, step int) []int {
	var result []int
	for i := start; i <= end; i += step {
		result = append(result, i)
	}
	return result
}

var cases = [][]int{
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
	{-1, -2, -3},          // Negative values
	{-5, -3, 1, 3, 5},     // Negative and positive values
	{50, 40, 60, 30, 70, 20, 80, 45},
	{50, 40, 60, 30, 70, 20, 80, 15},
	{50, 40, 60, 30, 70, 20, 80, 35},
	{50, 40, 60, 30, 70, 20, 80, 25},
	rangeWithSteps(-9, 16, 2),
	rangeWithSteps(0, 34, 3),
}

func deepCopyTestCases(original [][]int) [][]int {
	sliceCopy := make([][]int, len(original))
	for i, slice := range original {
		sliceCopy[i] = make([]int, len(slice))
		copy(sliceCopy[i], slice)
	}
	return sliceCopy
}

func populateTree(values []int) *AvlTree[int] {
	tree := NewAvlTree[int]()
	for _, value := range values {
		tree.Add(value)
	}
	return tree
}

// Tests the AVL tree with integer values, covering all basic rotation cases
func TestIntegerTree(t *testing.T) {
	cases = deepCopyTestCases(cases)

	for i, testCase := range cases {
		tree := populateTree(testCase)
		actual := tree.InorderTraverse(tree.root, nil)
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

		actual := tree.InorderTraverse(tree.root, nil)
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
		{1.1, 2.2, 3.3}, // Ordered
		{3.3, 2.2, 1.1}, // Reversed
		{2.2, 1.1, 3.3}, // Mixed
	}

	for i, testCase := range cases {
		tree := NewAvlTree[float64]()

		for _, value := range testCase {
			tree.Add(value)
		}

		actual := tree.InorderTraverse(tree.root, nil)
		expected := slices.Clone(testCase)
		slices.Sort(expected)

		if !slices.Equal(actual, expected) {
			t.Errorf("Test case %d: tree.Add(...) = %v; want %v", i, actual, expected)
		}
	}
}

// Test Contains method, indicating whether a value exists in the AVL tree
func TestContains(t *testing.T) {
	cases = deepCopyTestCases(cases)
	for _, testCase := range cases {
		values := testCase
		tree := populateTree(values)
		for _, v := range values {
			actual := tree.Contains(v)
			expected := true
			if actual != expected {
				t.Errorf("Test contains value: tree.Contains(%v) = %v; want %v", v, actual, expected)
			}
		}
	}
}

func TestDoesNotContain(t *testing.T) {
	tree := populateTree([]int{1, 2, 3})
	actual := tree.Contains(4)
	expected := false
	if actual != expected {
		t.Errorf("Test contains value: tree.Contains(4) = %v; want %v", actual, expected)
	}
}

// Test removing a value that does not exist in the AVL tree
func TestRemoveNonexistingValue(t *testing.T) {
	values := []int{1, 2, 3}
	tree := populateTree(values)
	actual := tree.Remove(0)
	expected := tree.Contains(0)
	if actual != expected {
		t.Errorf("Test remove non-existing value: tree.Remove(0) = %v; want %v", actual, expected)
	}
}

// Test removing a value that does not exist in the AVL tree
func TestRemoveValues(t *testing.T) {
	for i, testCase := range cases {
		for j, v := range testCase {
			tree := populateTree(testCase)

			// Successful remove returns true, so negate this to check against
			// `Contains(value) == false`

			actual := !tree.Remove(v)
			expected := tree.Contains(v)

			if actual != expected {
				t.Errorf("Test case %d.%d: tree.Remove(%v) = %v; want %v", i+1, j, v, actual, expected)
			}

			// Ensure order was maintained during removal
			actualValues := tree.InorderTraverse(tree.root, nil)
			expectedValues := slices.Clone(actualValues)
			slices.Sort(expectedValues)
			if !slices.Equal(actualValues, expectedValues) {
				t.Errorf("Test case %d.%d: tree.Remove(%v) = %v; want %v", i+1, j, v, actual, expected)
			}

		}
	}
}

// Test removing multiple values until the tree is empty
func TestRemoveMultipleValues(t *testing.T) {
	for i, testCase := range cases {
		tree := populateTree(testCase)
		for j, v := range testCase {

			// Successful remove returns true, so negate this to check against
			// `Contains(value) == false`
			actual := !tree.Remove(v)
			expected := tree.Contains(v)

			if actual != expected {
				t.Errorf("Test case %d.%d: tree.Remove(%v) = %v; want %v", i+1, j, v, actual, expected)
			}
		}
	}
}
