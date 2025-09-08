package main

import (
	"fmt"
)

// Operator est une interface qui applique une opération sur deux entiers
type Operator interface {
	Apply(a, b int) int
}

type Add struct{}

func (Add) Apply(a, b int) int { return a + b }

type Multiply struct{}

func (Multiply) Apply(a, b int) int { return a * b }

func main() {
	var op Operator

	op = Add{}
	fmt.Println("3 + 4 =", op.Apply(3, 4))

	op = Multiply{}
	fmt.Println("3 * 4 =", op.Apply(3, 4))
}
