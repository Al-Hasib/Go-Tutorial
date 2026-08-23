package main

import "fmt"

func main() {

	//operators
	// a := 10
	// b := 3

	//arithmetic operators
	// fmt.Println("Addition:", a+b)
	// fmt.Println("Subtraction:", a-b)
	// fmt.Println("Multiplication:", a*b)
	// fmt.Println("Division:", a/b)
	// fmt.Println("Modulus:", a%b)

	// comparison operators
	// fmt.Println("Equal to:", a == b)
	// fmt.Println("Not equal to:", a != b)
	// fmt.Println("Greater than:", a > b)
	// fmt.Println("Less than:", a < b)
	// fmt.Println("Greater than or equal to:", a >= b)
	// fmt.Println("Less than or equal to:", a <= b)

	//logical operators
	// fmt.Println("Logical AND:", a > 5 && b < 5)
	// fmt.Println("Logical OR:",  a > 5 || b < 5)
	// fmt.Println("Logical NOT:", !(a > 5))

	//assignment operators
	// a += 5 // a = a + 5
	// fmt.Println("After addition assignment:", a)
	// a -= 3 // a = a - 3
	// fmt.Println("After subtraction assignment:", a)
	// a *= 2 // a = a * 2
	// fmt.Println("After multiplication assignment:", a)
	// a /= 4 // a = a / 4
	// fmt.Println("After division assignment:", a)

	//increment and decrement operators
	// a++ // a = a + 1
	// fmt.Println("After increment:", a)
	// a-- // a = a - 1
	// fmt.Println("After decrement:", a)

	//bitwise operators
	e := 5  // 0101 in binary
	f := 3  // 0011 in binary
	fmt.Println("Bitwise AND:", e & f)   // 0001 in binary, which is 1
	fmt.Println("Bitwise OR:", e | f)    // 0111 in binary, which is 7
	fmt.Println("Bitwise XOR:", e ^ f)   // 0110 in binary, which is 6
	fmt.Println("Bitwise AND NOT:", e &^ f) // 0100 in binary, which is 4
	fmt.Println("Left Shift:", e << 1)    // 1010 in binary, which is 10
	fmt.Println("Right Shift:", e >> 1)   // 0010 in binary, which is 2
}
