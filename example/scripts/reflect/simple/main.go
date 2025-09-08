package reflect_example

import (
	"errors"
	"fmt"
	"reflect"
)

// Person represents a simple struct with various field types
type Person struct {
	Name    string  `json:"name" validate:"required"`
	Age     int     `json:"age" validate:"min=0,max=150"`
	Height  float64 `json:"height" validate:"min=0"`
	IsAlive bool    `json:"is_alive"`
}

// GetFieldNames returns all field names of a struct as a string slice
func GetFieldNames(s interface{}) ([]string, error) {
	// Check if the interface is a pointer and get its element
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ensure we're working with a struct
	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or a pointer to a struct")
	}

	t := v.Type()
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Name
	}

	return fieldNames, nil
}

// GetFieldTags returns all tags for a specific tag key of a struct
func GetFieldTags(s interface{}, tagKey string) (map[string]string, error) {
	// Check if the interface is a pointer and get its element
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ensure we're working with a struct
	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or a pointer to a struct")
	}

	t := v.Type()
	tags := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag, ok := field.Tag.Lookup(tagKey); ok {
			tags[field.Name] = tag
		}
	}

	return tags, nil
}

// SetField sets the value of a field in a struct by name
func SetField(s interface{}, fieldName string, value interface{}) error {
	// Check if the interface is a pointer and get its element
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("input must be a pointer to a struct")
	}
	v = v.Elem()

	// Ensure we're working with a struct
	if v.Kind() != reflect.Struct {
		return errors.New("input must be a pointer to a struct")
	}

	// Find the field by name
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s does not exist", fieldName)
	}

	// Check if the field is settable
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set (it may be unexported)", fieldName)
	}

	// Convert the value to the field's type
	val := reflect.ValueOf(value)
	if field.Type() != val.Type() {
		return fmt.Errorf("value type %s does not match field type %s", val.Type(), field.Type())
	}

	// Set the value
	field.Set(val)
	return nil
}

// DumpStructValues prints all field names and their values
func DumpStructValues(s interface{}) error {
	// Check if the interface is a pointer and get its element
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ensure we're working with a struct
	if v.Kind() != reflect.Struct {
		return errors.New("input must be a struct or a pointer to a struct")
	}

	t := v.Type()
	fmt.Println("Struct:", t.Name())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		fmt.Printf("  %s: %v\n", field.Name, value)
	}

	return nil
}
