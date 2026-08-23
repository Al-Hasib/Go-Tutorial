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
}
