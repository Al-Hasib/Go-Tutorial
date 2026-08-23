package main

import "fmt"

func main(){

	//arrays 
	// fixed length, same type, contiguous memory allocation
	// var arr [5]int
	// fmt.Println(arr)


	// arr[0] = 1
	// arr[1] = 2
	// arr[2] = 3
	// arr[3] = 4
	// arr[4] = 5
	// fmt.Println(arr)

	// //array literal
	// arr2 := [5]int{1,2,3,4}
	// fmt.Println(arr2)

	// //array literal with ellipsis
	// arr3 := [...]int{1,2,3,4,5,6,7,8,9,10}
	// fmt.Println(arr3)

	// // //array literal with ellipsis and string
	// arr4 := [...]string{"a","b","c","d","e"}
	// // fmt.Println(arr4)

	// // // accessing array elements
	// // fmt.Println(arr4[0])
	// // fmt.Println(arr4[1])

	// // // changing array elements
	// arr4[0] = "z"
	// fmt.Println(arr4)

	// // //looping through arrays
	// for i := 0; i < len(arr4); i++ {
	// 	fmt.Println(arr4[i])
	// }

	// //appending arrays
	// arr5 := [...]int{1,2,3,4,5}
	// arr6 := [...]int{6,7,8,9,10}
	// arr7 := append(arr5[:], arr6[:]...)
	// fmt.Println(arr7)
	// arr7 = append(arr7, 11,12,13)
	// fmt.Println(arr7)

	// //deleting elements from an array
	arr8 := [...]int{1,2,3,4,5}
	arr9 := append(arr8[:2], arr8[3:]...)
	fmt.Println(arr9)






}
