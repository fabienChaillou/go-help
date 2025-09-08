package main

import (
	"fmt"

	"arbre-binaire/fibonacci"
)

// Définition du noeud de l'arbre
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// Parcours in-order (gauche - racine - droite)
func InOrderTraversal(n *Node) {
	if n == nil {
		return
	}
	InOrderTraversal(n.Left)
	fmt.Print(n.Value, " ")
	InOrderTraversal(n.Right)
}

func main() {
	// Construction de l'arbre :
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

	fmt.Print("In-Order Traversal: ")
	InOrderTraversal(root)
	fmt.Println()

	n := 10
	fmt.Printf("Fibonacci (récursif) de %d = %d\n", n, fibonacci.FibonacciRecursive(n))
	fmt.Printf("Fibonacci (itératif) de %d = %d\n", n, fibonacci.FibonacciIterative(n))
}
