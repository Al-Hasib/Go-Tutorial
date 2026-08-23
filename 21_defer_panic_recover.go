package main

import "fmt"

// defer - runs a statement right before the surrounding function returns
// useful for cleanup: closing files, unlocking, printing "done", etc.

func simpleDefer() {
	defer fmt.Println("2. this runs last")
	fmt.Println("1. this runs first")
}

// multiple defers run in LIFO order - last deferred, first run
func multipleDefers() {
	defer fmt.Println("a")
	defer fmt.Println("b")
	defer fmt.Println("c")
	fmt.Println("running multipleDefers")
	// output: running..., then c, b, a
}

// defer's arguments are evaluated immediately, not when it runs
func deferArgs() {
	x := 1
	defer fmt.Println("deferred x was:", x) // captures 1 right now
	x = 100
	fmt.Println("current x is:", x)
}

// but a closure sees the LATEST value, since it reads the variable later
func deferClosure() {
	x := 1
	defer func() {
		fmt.Println("closure sees x as:", x) // reads x when it runs
	}()
	x = 100
}

// defer can change a named return value
func namedReturn() (result int) {
	defer func() {
		result++ // runs after "return 5", turning it into 6
	}()
	return 5
}

// classic use case - guarantee cleanup even if the function has many exits
func readFile(name string) {
	fmt.Println("opening", name)
	defer fmt.Println("closing", name) // always runs, no matter how we exit

	if name == "" {
		fmt.Println("bad name, returning early")
		return
	}
	fmt.Println("reading", name)
}

// panic - stops normal execution and starts unwinding the call stack
func step3() {
	panic("something went wrong in step3")
}

func step2() {
	fmt.Println("step2 start")
	step3()
	fmt.Println("step2 end") // never runs, panic already happened
}

func step1() {
	fmt.Println("step1 start")
	step2()
	fmt.Println("step1 end") // never runs either
}

// recover - only works inside a deferred function, stops the panic
func safeStep1() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in safeStep1:", r)
		}
	}()
	step1()
	fmt.Println("this still won't run") // panic skips the rest of safeStep1 too
}

// a helper that turns a panicking function into a normal error return
// the return value MUST be named - only then can defer write into it
// after a panic skips straight past "return err" below
func safeCall(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	f()
	return err
}

func main() {

	//basic defer
	simpleDefer()
	fmt.Println()

	//LIFO order
	multipleDefers()
	fmt.Println()

	//deferred args are frozen at the defer line
	deferArgs()
	fmt.Println()

	//closures read the variable's value at call time instead
	deferClosure()
	fmt.Println()

	//defer changing a named return value
	fmt.Println("namedReturn:", namedReturn()) // 6, not 5
	fmt.Println()

	//defer guarantees cleanup on every exit path
	readFile("data.txt")
	readFile("")
	fmt.Println()

	//panic unwinds the stack - uncomment to see the crash:
	// step1()

	//recover stops the panic from crashing the whole program
	safeStep1()
	fmt.Println("program is still running after the panic")
	fmt.Println()

	//wrapping panics as normal errors
	err := safeCall(func() {
		panic("boom")
	})
	fmt.Println("safeCall result:", err)

	err = safeCall(func() {
		fmt.Println("no panic this time")
	})
	fmt.Println("safeCall result:", err) // nil

}
