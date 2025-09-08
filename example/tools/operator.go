package tools

import "fmt"

func run()  {

// composite literals in go
elements := []int{1, 2, 3, 4}
// type Thing struct {
//     name       string
//     generation int
//     model      string
// }
// thing := Thing{“Raspberry Pi”, 2, “B”}
// // or using explicit field names
// thing = Thing{name: “Raspberry Pi”, generation: 2, model: “B”}

// // Key must be either literal constant or constant expression so it’s illegal to write:
// f := func() int { return 1 }
// elements := []string{0: “zero”, f(): “one”}

// elements := []string{
//     0:     "zero",
//     1:     "one",
//     4 / 2: "two",
//     2:     "also two"
// }

f := func ([]int e)  {
	for _, i := range e {
		fmt.Printf("%#v\n", s)
	}
}()

f(elements)
	
}