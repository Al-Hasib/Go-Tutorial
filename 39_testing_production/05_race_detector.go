package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// the real material is raceexample/counter.go + counter_test.go - an
// UnsafeCounter with a genuine data race, and a SafeCounter (Mutex-protected)
// that doesn't have one. read those alongside this file.
func main() {

	// ---------- WHY A NORMAL go test ISN'T ENOUGH ----------

	//a data race doesn't reliably crash or produce a visibly wrong result -
	//it just means the outcome is UNDEFINED. run the unsafe counter's test
	//normally (no -race) a few times, and it usually just... passes,
	//because the race rarely loses an update in practice on a given machine.
	//that's the actual danger: it can hide for a long time.
	fmt.Println("plain `go test` on the unsafe counter (usually passes - the race is still there, just not caught):")
	fmt.Println()
	runGoTest("-run", "TestUnsafeCounter", "./raceexample/...")

	// ---------- -race: WATCH EVERY MEMORY ACCESS, DON'T GUESS ----------

	fmt.Println("\nnow with -race, which instruments every memory access to PROVE a race exists,")
	fmt.Println("instead of relying on the outcome happening to look wrong:")
	fmt.Println()
	runRaceTest("TestUnsafeCounter")

	fmt.Println("\nsame treatment for the Mutex-protected SafeCounter - should stay clean under -race:")
	fmt.Println()
	runRaceTest("TestSafeCounter")

	fmt.Println("--- notice ---")
	fmt.Println("-race reports the exact two goroutines and lines that raced, plus their creation")
	fmt.Println("stacks - a real race is far easier to diagnose from that report than from guessing")
	fmt.Println("based on an occasionally-wrong number.")
	fmt.Println()
	fmt.Println("-race requires cgo (a C compiler) since it uses a C-based runtime library internally -")
	fmt.Println("that's a real, common gap on a fresh Windows machine with no C toolchain installed;")
	fmt.Println("this lesson falls back to a `golang` Docker container (which has one) when needed.")
}

// runRaceTest tries -race with the local toolchain first, and only falls
// back to Docker if the local one genuinely lacks a C compiler for it -
// most Linux/macOS setups won't need the fallback at all.
func runRaceTest(testName string) {
	out, err := exec.Command("go", "test", "-race", "-run", testName, "./raceexample/...").CombinedOutput()

	if err != nil && strings.Contains(string(out), "requires cgo") {
		fmt.Println("(local Go toolchain has no C compiler available for -race - running the same")
		fmt.Println(" test inside a `golang` Docker image instead, which ships one)")
		dockerOut, dockerErr := runInDocker(testName)
		if dockerErr != nil && strings.Contains(dockerErr.Error(), "executable file not found") {
			fmt.Println("(docker isn't available either - install a C toolchain, e.g. mingw-w64 on Windows,")
			fmt.Println(" or run this on Linux/macOS/WSL, to see -race work locally)")
			return
		}
		out, err = dockerOut, dockerErr
	}

	fmt.Print(string(out))
	if err != nil {
		//a non-zero exit is the EXPECTED, correct outcome for the unsafe
		//counter - -race deliberately fails the test run it catches a race in
		fmt.Println("(command exited non-zero:", err, ")")
	}
}

func runInDocker(testName string) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("docker", "run", "--rm",
		"-v", cwd+":/app",
		"-w", "/app",
		"golang:latest",
		"go", "test", "-race", "-run", testName, "./raceexample/...")
	return cmd.CombinedOutput()
}

func runGoTest(args ...string) {
	out, err := exec.Command("go", append([]string{"test"}, args...)...).CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Println("(command exited non-zero:", err, ")")
	}
}
