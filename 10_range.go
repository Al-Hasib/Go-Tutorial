package main

import "fmt"

func main() {

	//range - iterate over elements in an array, slice, string, map, or channel

	// iterate over an array
	// arr := [5]int{1, 2, 3, 4, 5}
	// for i, v := range arr {
	// 	fmt.Println("index", i, "value", v)
	// 	// fmt.Printf("Index: %d, Value: %d\n", i, v)
	// }

	// numbers := []int{10, 20, 30}

	// for _ , value := range numbers {
	// 	fmt.Println(value)
	// }

	// // iterate over a slice
	// slice := []string{"apple", "banana", "cherry"}
	// for i, v := range slice {
	// 	fmt.Printf("Index: %d, Value: %s\n", i, v)
	// }

	// iterate over a string
	// str := "hello world"
	// for i, v := range str {
	// 	fmt.Printf("Index: %d, Value: %c\n", i, v)
	// }

	// for _, char := range "Go" {
	// 	fmt.Printf("%c\n", char)
	// }

	// // iterate over a map
	myMap := map[string]int{"one": 1, "two": 2, "three": 3}
	// for k, v := range myMap {
	// 	fmt.Printf("Key: %s, Value: %d\n", k, v)
	// }

	// get only keys from a map
	for k := range myMap {
		fmt.Println("Key:", k)
	}

	// // get only values from a map
	for _, v := range myMap {
		fmt.Println("Value:", v)
	}


}
