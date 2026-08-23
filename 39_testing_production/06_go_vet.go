package main

import (
	"fmt"
	"os/exec"
)

// the real material is vetexample/vetissues.go - three real mistakes that
// the Go COMPILER accepts but go vet flags. read it alongside this file.
func main() {

	// ---------- vet FINDS BUGS THE COMPILER DELIBERATELY DOESN'T ----------

	//the compiler's job is to reject INVALID Go. all three functions in
	//vetexample/vetissues.go are valid Go - they compile and run - but each
	//does something that's almost certainly a mistake. that gap between
	//"compiles" and "probably correct" is exactly what vet exists to narrow.
	fmt.Println("go build ./vetexample/...  (should succeed - nothing here is a compile error)")
	fmt.Println()
	run("go", "build", "./vetexample/...")
	fmt.Println("(no output above = success)")

	fmt.Println("\ngo vet ./vetexample/...")
	fmt.Println()
	run("go", "vet", "./vetexample/...")

	fmt.Println("\n--- what each line means ---")
	fmt.Println("\"Printf format %d has arg ... of wrong type string\"")
	fmt.Println("  -> BadPrintf passes a string where %d expects a number; vet checks format")
	fmt.Println("     strings against their arguments, something the compiler can't do since")
	fmt.Println("     Printf's arguments are just ...any.")
	fmt.Println()
	fmt.Println("\"passes lock by value\" / \"copies lock value\"")
	fmt.Println("  -> Locker contains a sync.Mutex; passing a Locker BY VALUE copies that Mutex too,")
	fmt.Println("     so the copy protects nothing shared with the original. vet's copylocks check")
	fmt.Println("     catches this at both the function signature and the call site.")
	fmt.Println()
	fmt.Println("\"unreachable code\"")
	fmt.Println("  -> a statement placed after an unconditional return can never run. the compiler")
	fmt.Println("     only checks whether a function's FINAL statement terminates, not the whole")
	fmt.Println("     body - vet actually traces control flow and catches dead code in the middle too.")
	fmt.Println()
	fmt.Println("`go test` runs a small set of vet checks automatically before your tests, so some of")
	fmt.Println("these would already surface without ever typing `go vet` yourself. running it")
	fmt.Println("directly (or via `go build`, which does not run vet) gives you the full check list.")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Println("(command exited non-zero:", err, ")")
	}
}
