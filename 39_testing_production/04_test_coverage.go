package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {

	// ---------- -cover: THE HEADLINE PERCENTAGE ----------

	fmt.Println("go test -cover ./calc/...")
	fmt.Println()
	run("go", "test", "-cover", "./calc/...")

	// ---------- -coverprofile + go tool cover -func: PER-FUNCTION DETAIL ----------

	//the single percentage above tells you almost nothing useful on its
	//own - a coverprofile plus `go tool cover -func` breaks it down by
	//function, which is where coverage actually becomes actionable
	profile := filepath.Join(os.TempDir(), "calc.coverprofile")
	fmt.Println("\ngo test -coverprofile=<file> ./calc/...  then  go tool cover -func=<file>")
	fmt.Println()
	run("go", "test", "-coverprofile="+profile, "./calc/...")
	run("go", "tool", "cover", "-func="+profile)

	// ---------- WHAT THE GAP MEANS ----------

	fmt.Println("--- notice ---")
	fmt.Println("Max shows 0.0% - calc/calc.go's Max function has no test at all (on purpose, for this lesson).")
	fmt.Println("100% coverage does NOT mean bug-free - it means every LINE ran at least once, not that")
	fmt.Println("every possible input/edge case was checked. coverage finds untested code paths; it")
	fmt.Println("can't tell you whether the assertions in the tests that DO run are actually correct.")
	fmt.Println()
	fmt.Println("go tool cover -html=<file> -o coverage.html renders the same data as a browsable page")
	fmt.Println("with covered/uncovered lines highlighted directly in the source - the most useful view")
	fmt.Println("day to day, just harder to show in a terminal transcript like this one.")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Println("(command exited non-zero:", err, ")")
	}
}
