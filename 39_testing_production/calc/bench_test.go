package calc

import (
	"strconv"
	"strings"
	"testing"
)

// a benchmark function's name starts with Benchmark, takes *testing.B, and
// calls the code under test exactly b.N times - the testing framework picks
// b.N by running the loop for increasing counts until the timing is stable.
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(2, 3)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	//building this once outside the loop, so the benchmark measures
	//IsPalindrome itself, not string repetition too
	input := strings.Repeat("ab", 500) + "c" + strings.Repeat("ba", 500) // not a palindrome, worst case for the early-exit check

	b.ResetTimer() // in case setup above took measurable time, don't count it
	for i := 0; i < b.N; i++ {
		IsPalindrome(input)
	}
}

// BenchmarkIsPalindromeSizes uses b.Run to compare the SAME function across
// several input sizes in one benchmark run - handy for seeing how a
// function scales, not just its speed at one fixed size.
func BenchmarkIsPalindromeSizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		input := strings.Repeat("a", size)
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs() // also report allocations/op, not just time/op
			for i := 0; i < b.N; i++ {
				IsPalindrome(input)
			}
		})
	}
}
