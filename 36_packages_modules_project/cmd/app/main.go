// ---------- PROJECT STRUCTURE ----------
//
// this file lives at cmd/app/main.go, not just main.go in the project root.
// that's a common Go convention (NOT enforced by the compiler, unlike
// internal/): put each runnable program under cmd/<binary-name>/main.go,
// so the module can hold several independent binaries side by side, e.g.
//   cmd/app/main.go       -> builds the "app" binary
//   cmd/worker/main.go    -> would build a separate "worker" binary
// everything those binaries share (business logic, helpers) lives OUTSIDE
// cmd/, in ordinary packages like mathutil/ and stringutil/ - see this
// project's full layout:
//
//   36_packages_modules_project/
//     go.mod                    <- module definition (see notes below)
//     go.sum                    <- dependency checksums (see notes below)
//     cmd/app/main.go           <- this file: the entry point
//     internal/mathutil/        <- only importable from inside this module
//     stringutil/               <- importable from anywhere
//
// for a small project, even this is arguably more structure than needed -
// the Go team's own advice is to start flat (everything in one package)
// and only split into cmd/, internal/, etc. once a real reason shows up.
// this project deliberately shows the "grown a little" shape since that's
// what you'll see in most real-world Go repos.
package main

import (
	"fmt"

	"example.com/packagesdemo/internal/mathutil"
	"example.com/packagesdemo/stringutil"

	"github.com/google/uuid" // a real third-party dependency - see notes below
)

func main() {

	// ---------- PACKAGES ----------

	//each import path above maps to a directory. "mathutil" and "stringutil"
	//are the package NAMES (declared with `package mathutil` in their .go
	//files) - by convention the last element of the import path matches the
	//package name, which is why we call functions as mathutil.Add(...), not
	//by the full import path.
	fmt.Println("mathutil.Add(2, 3):", mathutil.Add(2, 3))
	fmt.Println("mathutil.Multiply(4, 5):", mathutil.Multiply(4, 5))

	cfg := mathutil.NewConfig(2, "prices")
	fmt.Println("cfg.Precision (exported field):", cfg.Precision)
	fmt.Println("cfg.Label() (getter for unexported field):", cfg.Label())

	//mathutil.roundToEven and cfg.label are NOT reachable from here - this
	//package only sees what mathutil chose to export:
	//   mathutil.roundToEven(4) // compile error: undefined (unexported)
	//   cfg.label               // compile error: undefined (unexported)

	fmt.Println("stringutil.Reverse(\"hello\"):", stringutil.Reverse("hello"))
	fmt.Println("stringutil.Shout(\" go \"):", stringutil.Shout(" go "))

	// ---------- GO MODULES ----------

	//go.mod (one directory up from cmd/, at this project's root) is what
	//turns this directory tree into a MODULE - a versioned, importable unit.
	//its first line, "module example.com/packagesdemo", is the import-path
	//prefix every package below inherits: that's why mathutil's full import
	//path is "example.com/packagesdemo/internal/mathutil", matching its
	//folder path under the module root.
	//
	//go.mod also pins a Go language version (the "go 1.xx" line), and lists
	//every external dependency this module needs, each with an exact
	//version - that list is what makes builds reproducible on any machine.
	//
	//this module was created with:
	//   go mod init example.com/packagesdemo
	//run once, right when a project starts.

	// ---------- DEPENDENCIES ----------

	//uuid is a real, tiny, widely used dependency - it was added by running:
	//   go get github.com/google/uuid@latest
	//that command did three things:
	// 1. downloaded the package source (cached locally)
	// 2. added a `require github.com/google/uuid v1.6.0` line to go.mod
	// 3. recorded cryptographic checksums for it in go.sum
	//
	//go.sum isn't for humans to read - it's there so a rebuild (by you, a
	//teammate, or CI) can verify the downloaded code is byte-for-byte the
	//same as what you tested against, not something swapped out later.
	//that's what "@v1.6.0" (semantic versioning: major.minor.patch) pins.
	//
	//`go mod tidy` keeps go.mod/go.sum accurate - it adds anything your code
	//imports but go.mod is missing, and removes entries nothing imports
	//anymore. running it after adding or removing imports is standard habit.
	id := uuid.New()
	fmt.Println("generated uuid (from a real dependency):", id)

	fmt.Println("done")
}
