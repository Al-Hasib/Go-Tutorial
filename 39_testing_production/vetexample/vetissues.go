// Package vetexample deliberately contains a few real mistakes that the Go
// COMPILER accepts (this file builds cleanly) but `go vet` flags as
// suspicious - that gap is exactly what vet exists to close.
package vetexample

import (
	"fmt"
	"sync"
)

// BadPrintf has a format/argument mismatch. Printf's format argument isn't
// type-checked by the compiler (it just takes ...any), so this compiles -
// but vet's "printf" check statically verifies %-verbs against the actual
// argument types and catches the mismatch here: %d expects a number, this
// passes a string.
func BadPrintf() {
	name := "Alice"
	fmt.Printf("count: %d\n", name)
}

// Locker embeds a Mutex directly, by value.
type Locker struct {
	mu    sync.Mutex
	value int
}

// passByValue takes a Locker BY VALUE - which copies its Mutex too. Each
// copy then has its own separate, useless lock protecting nothing shared.
// vet's "copylocks" check flags passing (or returning, or assigning) a
// value containing a Mutex like this.
func passByValue(l Locker) {
	l.mu.Lock()
	l.value++
	l.mu.Unlock()
}

func UseLocker() {
	l := Locker{}
	passByValue(l) // <- this copy is what vet flags
}

// afterReturn has a statement that can never execute. the Go COMPILER's
// own check here is narrower than you might expect: it only looks at
// whether a function's FINAL statement terminates (return/panic/etc) - it
// doesn't trace reachability through the whole body. that's why the
// trailing `return 0` below is required just to satisfy the compiler, even
// though the line above it can never run either way. vet's "unreachable"
// check does trace real control flow, and catches the dead Println line
// that the compiler happily ignored.
func afterReturn() int {
	return 42
	fmt.Println("this line can never run")
	return 0
}
