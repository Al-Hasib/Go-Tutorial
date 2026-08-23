package main

import (
	"fmt"
	"os/exec"
)

// the real material for this lesson is calc/bench_test.go - read it
// alongside this file. this just runs it and shows the real output.
func main() {

	// ---------- WHAT A BENCHMARK FUNCTION LOOKS LIKE ----------

	//name starts with Benchmark, takes *testing.B, and calls the code under
	//test exactly b.N times - `go test` calibrates b.N itself by running
	//the loop for longer and longer until the per-op timing stabilizes.
	//benchmarks don't run during a normal `go test` - they need -bench
	//explicitly, and -run=^$ below skips all Test functions so only the
	//benchmarks execute.
	fmt.Println("running the benchmarks in calc/bench_test.go:")
	fmt.Println()
	run("go", "test", "-run=^$", "-bench=.", "-benchmem", "./calc/...")

	fmt.Println()
	fmt.Println("--- reading the output ---")
	fmt.Println("BenchmarkAdd-8        1000000000   0.45 ns/op   0 B/op   0 allocs/op")
	fmt.Println("                  |        |            |           |")
	fmt.Println("            8 = GOMAXPROCS |     time per single op |")
	fmt.Println("                     b.N (iterations actually run)  allocations per op (from -benchmem)")
	fmt.Println()
	fmt.Println("BenchmarkIsPalindromeSizes/10, /100, /1000, /10000 come from b.Run inside one")
	fmt.Println("Benchmark function - the same code, measured across several input sizes at once,")
	fmt.Println("which is how you see a function's growth pattern instead of one single number.")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Println("(command exited non-zero:", err, ")")
	}
}
