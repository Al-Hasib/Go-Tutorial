package main

import (
	"fmt"
	"os/exec"
)

// the real material for this lesson is calc/table_driven_test.go - read it
// alongside this file. this just runs it and shows the real output.
func main() {

	// ---------- WHY TABLE-DRIVEN ----------

	//go has no built-in parameterized-test or assertion library, so the
	//idiomatic answer to "I want to run the same test with 6 different
	//inputs" is: build a slice of cases (a "table"), then loop over it with
	//t.Run. adding a 7th case means adding one line, not a whole new
	//function.
	fmt.Println("running the table-driven tests in calc/table_driven_test.go:")
	fmt.Println()
	run("go", "test", "-v", "-run", "TestDivideTable|TestIsPalindromeTable", "./calc/...")

	fmt.Println()
	fmt.Println("--- what to notice ---")
	fmt.Println("each table row becomes its own named subtest (e.g. TestDivideTable/division_by_zero_errors)")
	fmt.Println("go test -run TestDivideTable/division_by_zero_errors would run JUST that one row")
	fmt.Println("a failing row points at itself by name - no need to guess which of 6 cases broke")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Println("(command exited non-zero:", err, ")")
	}
}
