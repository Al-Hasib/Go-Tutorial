package main

import (
	"fmt"
)

func countDown(n int) {
	if n == 0 {
		return
	}

	fmt.Println(n)
	countDown(n - 1)
}


func main(){

	//recursion - function that calls itself
	// Tree traversal
	// Graph traversal
	// Directory/file traversal
	// Divide-and-conquer algorithms
	// Backtracking

	
	//factorial of a number
	var factorial func(n int) int
	factorial = func(n int) int {
		if n == 0 {
			return 1
		}
		return n * factorial(n-1)
	}

	fmt.Println(factorial(5))

	countDown(5)


}
