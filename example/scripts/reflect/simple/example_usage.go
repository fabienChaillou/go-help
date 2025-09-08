package main

import (
	"fmt"
	"log"

	"example.com/reflect_example" // Use your actual module path when implementing
)

func main() {
	// Create a person instance
	person := reflect_example.Person{
		Name:    "Alice",
		Age:     28,
		Height:  175.5,
		IsAlive: true,
	}

	// Get field names
	fieldNames, err := reflect_example.GetFieldNames(person)
	if err != nil {
		log.Fatalf("Error getting field names: %v", err)
	}
	fmt.Println("Field names:", fieldNames)

	// Get JSON tags
	jsonTags, err := reflect_example.GetFieldTags(person, "json")
	if err != nil {
		log.Fatalf("Error getting JSON tags: %v", err)
	}
	fmt.Println("JSON tags:", jsonTags)

	// Get validation tags
	validationTags, err := reflect_example.GetFieldTags(person, "validate")
	if err != nil {
		log.Fatalf("Error getting validation tags: %v", err)
	}
	fmt.Println("Validation tags:", validationTags)

	// Dump struct values
	fmt.Println("\nInitial struct values:")
	reflect_example.DumpStructValues(person)

	// Modify a field using reflection
	personPtr := &person
	err = reflect_example.SetField(personPtr, "Name", "Bob")
	if err != nil {
		log.Fatalf("Error setting field: %v", err)
	}
	err = reflect_example.SetField(personPtr, "Age", 35)
	if err != nil {
		log.Fatalf("Error setting field: %v", err)
	}

	// Display modified struct
	fmt.Println("\nModified struct values:")
	reflect_example.DumpStructValues(person)
}
