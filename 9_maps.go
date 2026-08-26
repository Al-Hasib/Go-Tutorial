package main

import "fmt"

func main() {

	//maps - collection of key-value pairs, unordered, dynamic size, fast lookups

	// // create a map
	// var myMap map[string]int
	// fmt.Println(myMap)
	// fmt.Println(len(myMap))
	// fmt.Println(myMap == nil)

	// // initialize a map
	myMap := make(map[string]int)
	// fmt.Println(myMap)

	// // add key-value pairs to the map
	myMap["one"] = 1
	myMap["two"] = 2
	myMap["three"] = 3
	fmt.Println(myMap)
	// fmt.Println(len(myMap))

	// // update a value in the map
	// myMap["two"] = 22
	// fmt.Println(myMap)

	// // access a value by key
	// fmt.Println(myMap["two"])

	// // check if a key exists
	// value, ok := myMap["one"]
	// if ok {
	// 	fmt.Println(value)
	// } else {
	// 	fmt.Println("Key not found")
	// }

	// // delete a key-value pair
	// delete(myMap, "two")
	// fmt.Println(myMap)

	// // iterate over a map
	// for key, value := range myMap {
	// 	fmt.Println(key, value)
	// }

	// //length of a map
	// fmt.Println(len(myMap))

	// // nested maps
	// nestedMap := make(map[string]map[string]int)
	// nestedMap["first"] = make(map[string]int)
	// nestedMap["first"]["one"] = 1
	// nestedMap["first"]["two"] = 2

	// nestedMap["second"] = make(map[string]int)
	// nestedMap["second"]["one"] = 1
	// nestedMap["second"]["two"] = 2
	// fmt.Println(nestedMap)

	// // maps with slices as values
	// sliceMap := make(map[string][]int)
	// sliceMap["numbers"] = []int{1, 2, 3, 4, 5}
	// sliceMap["letters"] = []int{65, 66, 67}
	// fmt.Println(sliceMap)
	// fmt.Println(sliceMap["numbers"])


}
