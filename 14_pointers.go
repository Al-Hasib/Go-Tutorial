package main

import "fmt"

func changeAge(age *int) {
	*age = 20
}

func main() {

	//pointers - keep address of a variable

	x := 10
	p := &x

	fmt.Println(p) //address
	fmt.Println(*p) //value

	age:=15
	changeAge(&age)
	fmt.Println(age)

	var pointer_variable *int
	var simple_variable int

	fmt.Println(pointer_variable) //nil
	fmt.Println(simple_variable) //0

}
