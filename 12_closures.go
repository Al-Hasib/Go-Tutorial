package main

import (
	"fmt"
)

func counter() func() int {
	value := 0

	return func() int {
		value++
		return value
	}
}

// takes a callback function and calls it with the computed result.
// the caller can pass a closure as the callback to capture its own surrounding variables.
func performOperation(a, b int, callback func(int)) {
	result := a + b
	callback(result)
}

func main() {
	// closure - function that can capture and remember the variables from its surrounding scope
	count := 0
	increment := func() int {
		count++
		return count
	}

	fmt.Println(increment())
	fmt.Println(increment())
	fmt.Println(increment())

	fmt.Println("Final count:", count)

	// closure with function returning another function
	counterFunc := counter()
	fmt.Println(counterFunc())
	fmt.Println(counterFunc())
	fmt.Println(counterFunc())

	// closure passed as a callback - the closure captures "label" from main's
	// scope even though it is invoked later, from inside performOperation
	label := "Sum result:"
	performOperation(5, 7, func(result int) {
		fmt.Println(label, result)
	})
}
