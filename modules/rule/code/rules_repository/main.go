package main

import "fmt"

func main() {
	fmt.Println("Start test repository example!")

	sql := `
	SELECT id, bu
	FROM rules
	WHERE bu = $1
`
	sql2 := "SELECT id, bu FROM rules WHERE bu = $1"

	fmt.Printf("%+q\n%+q", sql, sql2)
}
