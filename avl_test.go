package avl

import (
	"fmt"
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

func assert[T comparable](a, b T, msg string, t *testing.T) {
	if a != b {
		t.Errorf("%s %v != %v", msg, a, b)
	}
}

func assertSlice[T comparable](a, b []T, msg string, t *testing.T) {
	if !slices.Equal(a, b) {
		t.Errorf("%s\nexpected: %v\ngot: %v", msg, b, a)
	}
}

var cases = [][]int{
	{},                    // Empty tree
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

// Test insertNode method, checking that insertion follows BST properties
func TestInsertNode(t *testing.T) {
	insertCases := [][]int{
		{10, 5, 15, 4, 6, 14, 16},        // Positives
		{0, -5, 5, -6, -4, 4, 6},         // Zero
		{-10, -15, -5, -16, -14, -6, -4}, // Negatives

		// Example:
		//      10
		//     /  \
		//    /    \
		//   5     15
		//  / \   /  \
		// 4   6 14   16
	}

	type sampleTree struct {
		root  int
		lsub  int
		rsub  int
		lsubl int
		lsubr int
		rsubl int
		rsubr int
	}

	for _, testCase := range insertCases {
		tree := NewAvlTree[int]()
		sample := sampleTree{
			root:  testCase[0],
			lsub:  testCase[1],
			rsub:  testCase[2],
			lsubl: testCase[3],
			lsubr: testCase[4],
			rsubl: testCase[5],
			rsubr: testCase[6],
		}
		for _, v := range testCase {
			tree.insertNode(v)
		}

		root := tree.GetRootNode()

		assert(root.value, sample.root, "insertNode (root)", t)
		assert(root.left.value, sample.lsub, "insertNode(root.left)", t)
		assert(root.right.value, sample.rsub, "insertNode(root.right)", t)
		assert(root.left.left.value, sample.lsubl, "insertNode(root.left.left)", t)
		assert(root.left.right.value, sample.lsubr, "insertNode(root.left.right)", t)
		assert(root.right.left.value, sample.rsubl, "insertNode(root.right.left)", t)
		assert(root.right.right.value, sample.rsubr, "insertNode(root.right.right)", t)

	}
}

// Test Contains method, indicating whether a value exists in the AVL tree.
func TestContains(t *testing.T) {
	cases = deepCopyTestCases(cases)
	for _, testCase := range cases {
		values := testCase
		tree := NewAvlTree[int]()

		for _, v := range values {
			tree.insertNode(v)
			assert(tree.Contains(v), true, fmt.Sprintf("tree.Contains(%v)", v), t)
		}
	}

}
func populateTree(t *testing.T, values []int) *AvlTree[int] {
	tree := NewAvlTree[int]()
	for i, v := range values {
		tree.Add(v)
		assert(tree.Contains(v), true, fmt.Sprintf("tree.Add(%v", v), t)
		assert(i+1, tree.GetSize(), "tree size after Add", t)

	}
	return tree
}

// Tests the AVL tree with integer values, covering all basic rotation cases
func TestIntegerTree(t *testing.T) {
	cases = deepCopyTestCases(cases)

	for _, testCase := range cases {
		tree := populateTree(t, testCase)
		actual := tree.InorderTraverse(tree.root, nil)
		expected := slices.Clone(testCase)
		slices.Sort(expected)
		assertSlice(actual, expected, "tree.Add(...)", t)
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

	for _, testCase := range cases {
		tree := NewAvlTree[string]()

		for _, value := range testCase {
			tree.Add(value)
			assert(tree.Contains(value), true, fmt.Sprintf("tree.Add(%v)", value), t)
		}

		actual := tree.InorderTraverse(tree.root, nil)
		expected := slices.Clone(testCase)
		slices.Sort(expected)

		assertSlice(actual, expected, "tree.Add(...)", t)
	}
}

// Tests the AVL tree with floating-point values
func TestFloatTree(t *testing.T) {
	cases := [][]float64{
		{1.1, 2.2, 3.3}, // Ordered
		{3.3, 2.2, 1.1}, // Reversed
		{2.2, 1.1, 3.3}, // Mixed
	}

	for _, testCase := range cases {
		tree := NewAvlTree[float64]()

		for _, value := range testCase {
			tree.Add(value)
			assert(tree.Contains(value), true, fmt.Sprintf("tree.Add(%v)", value), t)
		}

		actual := tree.InorderTraverse(tree.root, nil)
		expected := slices.Clone(testCase)
		slices.Sort(expected)
		assertSlice(actual, expected, "tree.Add(...)", t)
	}
}

// Test negative case for Contains method
func TestDoesNotContain(t *testing.T) {
	tree := populateTree(t, []int{1, 2, 3})
	assert(tree.Contains(4), false, "tree.Contains(4)", t)
}

// Test removing a value from the tree
func TestRemoveValues(t *testing.T) {
	for _, testCase := range cases {
		for _, v := range testCase {
			tree := populateTree(t, testCase)
			size := tree.GetSize()

			// Successful remove returns true, so negate this to check against
			// `Contains(value) == false`
			assert(tree.Remove(v), !tree.Contains(v), "tree.Remove(v)", t)
			assert(tree.GetSize(), size-1, "tree.size after Remove", t)

			// Ensure order was maintained during removal
			actualValues := tree.InorderTraverse(tree.root, nil)
			expectedValues := slices.Clone(actualValues)
			slices.Sort(expectedValues)
			assertSlice(actualValues, expectedValues, "tree.Remove(v)", t)

		}
	}
}

// Test removing a value that does not exist in the AVL tree (negative case)
func TestRemoveNonexistingValue(t *testing.T) {
	values := []int{1, 2, 3}
	tree := populateTree(t, values)
	size := tree.GetSize()
	assert(tree.Remove(0), tree.Contains(0), "tree.Remove(0)", t)
	assert(tree.GetSize(), size, "tree.size after Remove", t)
}

// Test removing multiple values until the tree is empty
func TestRemoveMultipleValues(t *testing.T) {
	for _, testCase := range cases {
		tree := populateTree(t, testCase)
		for _, v := range testCase {
			assert(tree.Remove(v), !tree.Contains(v), fmt.Sprintf("tree.Remove(%v)", v), t)
		}
		assert(tree.IsEmpty(), true, "tree.IsEmpty()", t)
		assert(tree.GetSize(), 0, "tree.size after Remove", t)
	}
}

func TestClearTree(t *testing.T) {
	testCase := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	tree := populateTree(t, testCase)
	tree.Clear()
	assert(tree.IsEmpty(), true, "tree.Clear()", t)
	assert(tree.GetSize(), 0, "tree.size after Remove", t)
}

func TestGetMinNode(t *testing.T) {
	var minValue int
	for _, testCase := range cases {
		tree := populateTree(t, testCase)
		if len(testCase) == 0 { // Empty tree case
			assert(tree.GetMinNode(), nil, "tree.GetMin()", t)
		} else {
			minValue = slices.Min(testCase)
			assert(tree.GetMinNode().value, minValue, "tree.GetMin()", t)
		}
	}
}

func TestGetMaxNode(t *testing.T) {
	var maxValue int
	for _, testCase := range cases {
		tree := populateTree(t, testCase)
		if len(testCase) == 0 { // Empty tree case
			assert(tree.GetMaxNode(), nil, "tree.GetMax()", t)
		} else {
			maxValue = slices.Max(testCase)
			assert(tree.GetMaxNode().value, maxValue, "tree.GetMax()", t)
		}
	}
}
