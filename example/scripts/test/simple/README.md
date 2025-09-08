## simple unit tests

### Key Go Testing Concepts

Test files: Go test files are named with _test.go suffix
Test functions: Test functions start with Test and take a *testing.T parameter
Table-driven tests: A common Go pattern for testing multiple scenarios
Subtests: Using t.Run() to organize test cases
Error reporting: Using t.Errorf() to report test failures

### Running the Tests
You can run these tests with:
`go test`

For verbose output:

`go test -v`

To run a specific test:

`go test -run TestGreet`
