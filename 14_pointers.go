package main

import "fmt"

func ChangeAge(age int){
	age = 20
}

func changeAgePointer(age *int) {
	*age = 20
}

// // swapPointers needs pointers - if a, b were plain ints, they'd be copies
// // and the swap would have no effect on the caller's variables
func swapPointers(a, b *int) {
	*a, *b = *b, *a
}

// // arrays are copied when passed to a function, so mutating the copy
// // does not affect the caller's array - a pointer is required to mutate it
func modifyArray(arr [3]int) {
	arr[0] = 100
}

func modifyArrayPtr(arr *[3]int) {
	arr[0] = 100
}

// // slices already hold an internal pointer to their underlying array, so a
// // function receiving a slice (no "*") can still mutate the caller's data
func modifySlice(s []int) {
	s[0] = 100
}

func main() {

	//pointers - keep address of a variable

	// x := 10
	// p := &x

	// fmt.Println(p) //address
	// fmt.Println(*p) //value

	// age:=15
	// ChangeAge(age)
	// fmt.Println(age) 
	
	// changeAgePointer(&age)
	// fmt.Println(age)

	var pointer_variable *int
	var simple_variable int

	fmt.Println(simple_variable) //0
	fmt.Println(pointer_variable) //nil
	// if pointer_variable == nil{
	// 	fmt.Println("No Value")
	// }
	

	// // // nil pointer dereference - always check for nil before dereferencing,
	// // // *pointer_variable here would panic: "invalid memory address or nil pointer dereference"
	// if pointer_variable != nil {
	// 	fmt.Println(*pointer_variable)
	// } else {
	// 	fmt.Println("pointer_variable is nil!")
	// }

	// // swap two variables - only possible because swapPointers receives addresses
	// m, n := 1, 2
	// swapPointers(&m, &n)
	// fmt.Println("swapped:", m, n) // swapped: 2 1

	// // array is passed by value - modifyArray only changes its own copy
	// arr := [3]int{1, 2, 3}
	// modifyArray(arr)
	// fmt.Println("array after modifyArray:", arr) // unchanged: [1 2 3]

	// modifyArrayPtr(&arr)
	// fmt.Println("array after modifyArrayPtr:", arr) // changed: [100 2 3]

	// // slice is passed by value too, but that value already contains a pointer
	// // to the underlying array, so modifySlice can still mutate the caller's data
	// s := []int{1, 2, 3}
	// modifySlice(s)
	// fmt.Println("slice after modifySlice:", s) // changed: [100 2 3]

	// // new(T) allocates a zeroed T and returns a pointer to it - an alternative to &x
	numPtr := new(int) 
	fmt.Println(*numPtr) // 0 (zero value)
	*numPtr = 50
	fmt.Println(*numPtr) // 50

}
