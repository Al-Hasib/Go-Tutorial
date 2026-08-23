package main

import "fmt"

func main(){

	//conditions
	age:= 20

	// if age >= 18 {
	// 	fmt.Println("You are a voter.")
	// } else{
	// 	fmt.Println("You are not a voter.")
	// }

	// if - else if - else
	// if age < 13 {
	// 	fmt.Println("You are a child.")
	// } else if age >= 13 && age < 20 {
	// 	fmt.Println("You are a teenager.")
	// } else {
	// 	fmt.Println("You are an adult.")
	// }

	//switch case
	switch age {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
		fmt.Println("You are a child.")
	case 13, 14, 15, 16, 17:
		fmt.Println("You are a teenager.")
	default:
		fmt.Println("You are an adult.")
	}

	
}
