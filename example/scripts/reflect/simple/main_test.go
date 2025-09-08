package reflect_example

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetFieldNames(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
		wantErr  bool
	}{
		{
			name:     "Valid struct",
			input:    Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			expected: []string{"Name", "Age", "Height", "IsAlive"},
			wantErr:  false,
		},
		{
			name:     "Valid struct pointer",
			input:    &Person{Name: "Jane", Age: 25, Height: 165.0, IsAlive: true},
			expected: []string{"Name", "Age", "Height", "IsAlive"},
			wantErr:  false,
		},
		{
			name:     "Non-struct input",
			input:    "Not a struct",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldNames(tt.input)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result if no error expected
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("GetFieldNames() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestGetFieldTags(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		tagKey   string
		expected map[string]string
		wantErr  bool
	}{
		{
			name:   "JSON tags",
			input:  Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			tagKey: "json",
			expected: map[string]string{
				"Name":    "name",
				"Age":     "age",
				"Height":  "height",
				"IsAlive": "is_alive",
			},
			wantErr: false,
		},
		{
			name:   "Validate tags",
			input:  Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			tagKey: "validate",
			expected: map[string]string{
				"Name":   "required",
				"Age":    "min=0,max=150",
				"Height": "min=0",
			},
			wantErr: false,
		},
		{
			name:     "Non-existent tag",
			input:    Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			tagKey:   "nonexistent",
			expected: map[string]string{},
			wantErr:  false,
		},
		{
			name:     "Non-struct input",
			input:    "Not a struct",
			tagKey:   "json",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldTags(tt.input, tt.tagKey)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result if no error expected
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("GetFieldTags() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestSetField(t *testing.T) {
	tests := []struct {
		name      string
		person    *Person
		fieldName string
		value     interface{}
		wantErr   bool
		expected  Person
	}{
		{
			name:      "Set Name field",
			person:    &Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			fieldName: "Name",
			value:     "Jane",
			wantErr:   false,
			expected:  Person{Name: "Jane", Age: 30, Height: 182.5, IsAlive: true},
		},
		{
			name:      "Set Age field",
			person:    &Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			fieldName: "Age",
			value:     25,
			wantErr:   false,
			expected:  Person{Name: "John", Age: 25, Height: 182.5, IsAlive: true},
		},
		{
			name:      "Set IsAlive field",
			person:    &Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			fieldName: "IsAlive",
			value:     false,
			wantErr:   false,
			expected:  Person{Name: "John", Age: 30, Height: 182.5, IsAlive: false},
		},
		{
			name:      "Field does not exist",
			person:    &Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			fieldName: "Weight",
			value:     75.5,
			wantErr:   true,
			expected:  Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
		},
		{
			name:      "Type mismatch",
			person:    &Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
			fieldName: "Age",
			value:     "25", // String instead of int
			wantErr:   true,
			expected:  Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetField(tt.person, tt.fieldName, tt.value)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result if no error expected
			if !tt.wantErr {
				if !reflect.DeepEqual(*tt.person, tt.expected) {
					t.Errorf("After SetField() = %+v, want %+v", *tt.person, tt.expected)
				}
			}
		})
	}

	// Test with non-pointer
	t.Run("Non-pointer input", func(t *testing.T) {
		person := Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true}
		err := SetField(person, "Name", "Jane")
		if err == nil {
			t.Errorf("SetField() should return error for non-pointer input")
		}
	})
}

func TestDumpStructValues(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	person := Person{Name: "John", Age: 30, Height: 182.5, IsAlive: true}
	err := DumpStructValues(person)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check results
	if err != nil {
		t.Errorf("DumpStructValues() error = %v", err)
	}

	// Check that the output contains all the field values
	expectedContents := []string{
		"Struct: Person",
		"Name: John",
		"Age: 30",
		"Height: 182.5",
		"IsAlive: true",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("DumpStructValues() output does not contain %q", expected)
		}
	}

	// Test error case
	err = DumpStructValues("not a struct")
	if err == nil {
		t.Errorf("DumpStructValues() should return error for non-struct input")
	}
}
