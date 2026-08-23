package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// the 11_deployment/ folder covers shipping a Go service as a container.
// this file covers the other common path: a single, self-contained,
// cross-compiled BINARY - no container, no runtime dependency at all, just
// the executable copied straight onto a target machine.
func main() {

	// ---------- CROSS-COMPILATION: NO EMULATOR OR TARGET MACHINE NEEDED ----------

	//GOOS/GOARCH tell the Go compiler what to build FOR, independent of
	//what you're building ON - this Windows machine can produce a Linux
	//ARM64 binary directly, no cross-toolchain install, no Docker, no VM
	targets := []struct {
		goos, goarch, filename string
	}{
		{"linux", "amd64", "server-linux-amd64"},
		{"linux", "arm64", "server-linux-arm64"},
		{"darwin", "arm64", "server-darwin-arm64"}, // Apple Silicon Mac
		{"windows", "amd64", "server-windows-amd64.exe"},
	}

	outDir := filepath.Join("11_deployment", "dist")
	must(os.MkdirAll(outDir, 0755))
	defer os.RemoveAll(outDir) // build artifacts - cleaned up after this lesson, not meant to live in the repo

	for _, t := range targets {
		outPath := filepath.Join(outDir, t.filename)

		//-ldflags="-s -w" strips debug symbols/DWARF info - smaller binary,
		//no effect on runtime behavior, just harder to run a debugger against later
		cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", outPath, "./11_deployment")
		cmd.Env = append(os.Environ(), "GOOS="+t.goos, "GOARCH="+t.goarch, "CGO_ENABLED=0")

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s/%s -> BUILD FAILED: %s\n", t.goos, t.goarch, out)
			continue
		}

		info, statErr := os.Stat(outPath)
		must(statErr)
		fmt.Printf("%s/%s -> %s (%.1f MB)\n", t.goos, t.goarch, t.filename, float64(info.Size())/1024/1024)
	}

	// ---------- WHAT THIS BUYS YOU, AND WHAT IT DOESN'T ----------

	fmt.Println("\n--- notice ---")
	fmt.Println("every binary above was built HERE, on Windows, with no target machine or emulator -")
	fmt.Println("that's what makes Go's cross-compilation unusually simple compared to most languages.")
	fmt.Println()
	fmt.Println("CGO_ENABLED=0 above matters: a build that uses cgo (some database drivers, some crypto")
	fmt.Println("bindings) needs a real C cross-toolchain for the TARGET platform to cross-compile - a")
	fmt.Println("pure-Go dependency tree (like every lesson in this repo) avoids that entirely.")
	fmt.Println()
	fmt.Println("this approach skips containers, but you lose what a container gives you for free: a")
	fmt.Println("consistent OS/filesystem/set of installed libraries around the binary. it fits best")
	fmt.Println("for a static binary with no external runtime dependencies - exactly this repo's shape.")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
