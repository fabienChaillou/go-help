package tree

import (
	"reflect"
	"testing"
)

func TestInOrderTraversal(t *testing.T) {
	// Arbre :
	//       4
	//     /   \
	//    2     6
	//   / \   / \
	//  1   3 5   7
	root := &Node{
		Value: 4,
		Left: &Node{
			Value: 2,
			Left:  &Node{Value: 1},
			Right: &Node{Value: 3},
		},
		Right: &Node{
			Value: 6,
			Left:  &Node{Value: 5},
			Right: &Node{Value: 7},
		},
	}

	var result []int
	InOrderTraversal(root, func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3, 4, 5, 6, 7}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InOrderTraversal() = %v; want %v", result, expected)
	}
}
