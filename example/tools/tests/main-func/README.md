I've created two Go files for you:

main.go - Contains two simple functions (Sum and Multiply) and a main function that demonstrates their usage.
main_test.go - Contains comprehensive tests for these functions, including:

Unit tests for both Sum and Multiply functions with various test cases
Benchmark tests to measure performance



To run these tests, you would:
```shell
bash# Run all tests
go test *.go

# Run tests with verbose output
go test -v

# Run benchmark tests
go test -bench=.

# Run tests with coverage
go test -cover
```
