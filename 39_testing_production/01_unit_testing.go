package main

import (
	"fmt"
	"os/exec"
)

// unlike every other lesson in this repo, you don't "go run" a test - you
// write a _test.go file and run `go test`. the real material for THIS
// lesson lives in calc/unit_test.go (read that file alongside this one).
// this runner just executes the real `go test` command for you and shows
// the real output, so the lesson stays a single `go run 01_....go` away.
func main() {

	// ---------- WHAT MAKES SOMETHING A TEST ----------

	//a function is a test if: its name starts with Test, it takes exactly
	//one parameter of type *testing.T, and it lives in a file ending in
	//_test.go. no registration, no imports of a test framework beyond the
	//standard library's "testing" package - `go test` finds these by
	//convention and runs them all.
	fmt.Println("running the plain (non-table, non-benchmark) tests in calc/unit_test.go:")
	fmt.Println()
	run("go", "test", "-v", "-run", "TestAdd|TestDivide|TestDivide_ByZero|TestIsPalindrome|TestAdd_WithHelper", "./calc/...")

	fmt.Println()
	fmt.Println("--- what to notice in calc/unit_test.go ---")
	fmt.Println("t.Errorf: fails the test but keeps running the rest of the function")
	fmt.Println("t.Fatalf: fails AND stops immediately - use when later lines assume success")
	fmt.Println("t.Run:    a named SUBTEST, shows up as its own PASS/FAIL line above")
	fmt.Println("t.Helper: marks a helper so failures report the CALLER's line, not the helper's")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Println("(command exited non-zero:", err, ")")
	}
}
