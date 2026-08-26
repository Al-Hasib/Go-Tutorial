package main

import (
    "fmt"
    "sort"
)

//first function
func greet() {
    fmt.Println("Hello!")
}

// // function with parameters
func add(a int, b int) int {
	return a + b	
}

// // function with multiple return values
func swap(a, b string) (string, string) {
	return b, a
}

// // function with slice as parameter
func sumSlice(numbers []int) int{
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

// named return values
func divideNumbers(a, b float64) (result float64, err error) {
	if b == 0 {
		err = fmt.Errorf("division by zero")
		return
	}
	result = a / b
	return
}

// function will take map of students number and return grade and ranking
func getResults(students map[string]int) (map[string]string, map[string]int) {
	grades := make(map[string]string)
	ranks := make(map[string]int)

	// Calculate grades
	for name, mark := range students {
		if mark >= 80 {
			grades[name] = "A+"
		} else if mark >= 70 {
			grades[name] = "A"
		} else if mark >= 60 {
			grades[name] = "B"
		} else if mark >= 50 {
			grades[name] = "C"
		} else if mark >= 40 {
			grades[name] = "D"
		} else {
			grades[name] = "F"
		}
	}

	// Get student names
	names := []string{}

	for name := range students {
		names = append(names, name)
	}

	// Sort names based on marks
	sort.Slice(names, func(i, j int) bool {
		return students[names[i]] > students[names[j]]
	})

	// Assign ranks
	for i, name := range names {
		ranks[name] = i + 1
	}

	return grades, ranks
}

// variadic function - takes variable number of arguments
func sum(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

// function with a function as a parameter
func applyOperation(a, b int, operation func(int, int) int) int{
	return operation(a, b)
}

//anonymous function
var multiply = func(a, b int) int {
	return a * b
}

func main() {

	//functions - block of code that performs a specific task, can take parameters and return values
	// greet()
	// greet()

	// summation := add(5, 3)
	// fmt.Println("Result:", summation)

	// swappedA, swappedB := swap("first", "second")
	// fmt.Println("Swapped:", swappedA, swappedB)

	// slice := []int{1, 2, 3, 4, 5}
	// sliceSum := sumSlice(slice)
	// fmt.Println("Slice Sum:", sliceSum)

	
	// divResult, divErr := divideNumbers(10, 0)
	// if divErr != nil {
	// 	fmt.Println("Error:", divErr)
	// } else {
	// 	fmt.Println("Division Result:", divResult)
	// }

	// // Example usage of getResults function
	// students := map[string]int{
	// 	"Alice":   55,
	// 	"Bob":     72,
	// 	"Charlie": 60,
	// 	"David":   95,
	// 	"Eve":     30,
	// }
	// grades, ranks := getResults(students)
	// fmt.Println("Grades:", grades)
	// fmt.Println("Ranks:", ranks)


	// // Example usage of sum function
	// total := sum(1, 2, 3, 4, 5,7,4)
	// fmt.Println("Total Sum:", total)


	// // // Example usage of applyOperation function
	// result := applyOperation(10, 5, func(x, y int) int {
	// 	return x + y
	// })
	// fmt.Println("Operation Result:", result)

	// // Example usage of anonymous function
	// multiplied := multiply(10, 5)
	// fmt.Println("Anonymous Function Result:", multiplied)

}
