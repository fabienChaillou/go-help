package process

type order struct {
	sum int
}

func (o *order) processData(values <-chan int, done chan<- struct{}) {
	for v := range values {
		o.sum += v
	}
	done <- struct{}{}
	close(done)
}
