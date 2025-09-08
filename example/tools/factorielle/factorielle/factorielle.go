package factorielle

func Factorielle(n int) int {
	if n == 0 {
		return 1
	}
	return n * Factorielle(n-1)
}

func FactorielleIterative(n int) int {
	if n == 0 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}
