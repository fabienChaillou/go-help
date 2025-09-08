package process

import "testing"

func Test_order_processData(t *testing.T) {
	tests := []struct {
		name        string
		inputValues []int
		wantSum     int
	}{
		{
			name:        "base case",
			inputValues: []int{0, 1, 2, 3},
			wantSum:     6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			o := order{}

			// create channels
			testChan := make(chan int)
			doneChan := make(chan struct{})

			// execute
			// execute my function on a separate goroutine, so that I can move on
			// and send my data to the input channel
			go o.processData(testChan, doneChan)

			// create a loop
			for _, v := range tt.inputValues {
				testChan <- v
			}

			// close channel
			close(testChan)

			// wait for done signal
			<-doneChan

			// assert
			if o.sum != tt.wantSum {
				t.Errorf("unexpected sum, got %d, wanted %d", o.sum, tt.wantSum)
			}
		})
	}
}
