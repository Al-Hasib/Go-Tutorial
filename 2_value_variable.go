package main

import "fmt"

func main(){

	// fmt.Println("Hello, World!") //string
	// fmt.Println(42)              //int, int8, int16, int32, int64
	// fmt.Println(3.14)            //float32, float64
	// fmt.Println(true)			// bool

	//types of value
// 	fmt.Printf("%T\n", "Hello") 
// 	fmt.Printf("%T\n", 42)
// 	fmt.Printf("%T\n", 3.14)
// 	fmt.Printf("%T\n", true)

//variables
	// var name string = "John"
	// var age int = 30
	// var height float64 = 5.9
	// var isStudent bool = true

	// fmt.Println(name)
	// fmt.Println(age)
	// fmt.Println(height)
	// fmt.Println(isStudent)

	//variable declaration with type inference
	// var city = "New York"
	// var score = 95.5
	// var isPassed = true

	// fmt.Println(city)
	// fmt.Println(score)
	// fmt.Println(isPassed)

	// // variables without value assignments
	// var country string
	// var population int
	// var temperature float64
	// var isRaining bool

	// fmt.Println(country)
	// fmt.Println(population)
	// fmt.Println(temperature)
	// fmt.Println(isRaining)

	// variable with short declaration
	// name := "Alice"
	// age := 25
	// height := 5.6
	// isStudent := true

	// fmt.Println(name)
	// fmt.Println(age)
	// fmt.Println(height)
	// fmt.Println(isStudent)

	// name, age, height, isStudent = "Bob", 30, 6.1, false

	// fmt.Println(name)
	// fmt.Println(age)
	// fmt.Println(height)
	// fmt.Println(isStudent)

	// multiple variable declaration
	var a, b, c int = 1, 2, 3
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

	var (
		name string = "John"
		age int = 30
		height float64 = 5.9
		isStudent bool = true
	)

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(height)
	fmt.Println(isStudent)
}