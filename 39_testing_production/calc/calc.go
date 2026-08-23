// Package calc holds a few small functions with real edge cases - just
// enough logic to make testing them meaningful. Its test files
// (unit_test.go, table_driven_test.go, bench_test.go) live right here in
// the same directory, since `go test` operates per-package-directory.
package calc

import "errors"

// Add is the simplest possible case - useful for showing the shape of a test
// before anything else gets in the way.
func Add(a, b int) int {
	return a + b
}

// Divide has a real edge case (dividing by zero) worth testing on purpose.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// IsPalindrome reads the same forwards and backwards - simple logic, but
// with enough edge cases (empty string, single character, mixed case) to
// make a good table-driven test.
func IsPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

// Max is deliberately left WITHOUT a test - see the Test Coverage lesson,
// which uses that gap to show what a coverage report actually points at.
func Max(nums ...int) (int, error) {
	if len(nums) == 0 {
		return 0, errors.New("Max of an empty list")
	}
	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max, nil
}
