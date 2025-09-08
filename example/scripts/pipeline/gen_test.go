package main

import "testing"

func TestGen(t *testing.T) {
	tests := []struct {
		name        string
		inputValues []int
		wantSum     int
	}{
		{
			name:        "base case",
			inputValues: []int{1: 2},
			wantSum:     6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			g := gen(tt.inputValues...)
			// create channel
			testChan := make(chan int)

			// execute
			testChan <- tt.wantSum

			// assert
			if g != testChan {
				t.Errorf("unexpected sum, got %d, wanted %d", o.sum, tt.wantSum)
			}
		})
	}
}
