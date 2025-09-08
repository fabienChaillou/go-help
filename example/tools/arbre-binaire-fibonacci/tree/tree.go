package tree

// Définition du noeud de l'arbre
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// Parcours in-order (gauche - racine - droite)
func InOrderTraversal(n *Node, visit func(int)) {
	if n == nil {
		return
	}
	InOrderTraversal(n.Left, visit)
	visit(n.Value)
	InOrderTraversal(n.Right, visit)
}
