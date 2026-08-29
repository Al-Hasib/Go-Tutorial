package main

import "fmt"

// struct declared outside main, so we can add methods to it
type Rect struct {
	Width  int
	Height int
}

// method with value receiver - only reads
func (r Rect) Area() int {
	return r.Width * r.Height
}

// method with pointer receiver - can change the struct
func (r *Rect) Scale(n int) {
	r.Width = r.Width * n
	r.Height = r.Height * n
}

// embedded struct - Rect fields get promoted into Box
type Box struct {
	Rect
	Depth int
}

func main() {

	// structs - group different types of data together

	// name:= "Alice"
	// Age:=12
	// Passed:= true
	// marks:= []int{1,2,3}

	type Student struct {
		Name     string
		Age      int
		Passed   bool
		Marks    []int
		Subjects map[string]int
	}

	// // initialize an object
	// student := Student{
	// 	Name:   "Alice",
	// 	Age:    16,
	// 	Passed: true,
	// 	Marks:  []int{80, 90, 85},
	// 	Subjects: map[string]int{
	// 		"Math":    90,
	// 		"English": 85,
	// 	},
	// }

	// //access fields
	// fmt.Println(student.Passed)
	// fmt.Println(student.Subjects)

	// //modify fields
	// student.Age = 18
	// student.Passed = false

	// fmt.Println(student.Age)

	// //zero value - all fields get their zero value
	// var empty Student
	// fmt.Println(empty)              // { 0 false [] map[]}
	// fmt.Println(empty.Marks == nil) // true

	// //without field names - must follow the declared order
	// student2 := Student{"Bob", 17, true, nil, nil}
	// fmt.Println(student2)

	// //pointer to a struct
	// p := &student
	// p.Age = 20 // same as (*p).Age = 20
	// fmt.Println(student.Age)

	// // //new() gives a pointer to a zero struct
	// student3 := new(Student)
	// fmt.Println(*student3)
	// student3.Name = "Carol"
	// fmt.Println(*student3)

	//structs are copied on assignment
	// copyStudent := student
	// copyStudent.Name = "Changed"
	// fmt.Println(student.Name, copyStudent.Name)

	// // //printing formats
	// fmt.Printf("%v\n", student)  // values only
	// fmt.Printf("%+v\n", student) // with field names
	// fmt.Printf("%T\n", student)  // type name

	// //nested struct
	// type Address struct {
	// 	City string
	// }
	// type Employee struct {
	// 	Name    string
	// 	Age     int
	// 	Address Address
	// }
	// emp := Employee{Name: "Dave",Age: 20, Address: Address{City: "Dhaka"}}
	// fmt.Println(emp.Age)

	// //anonymous struct - no type name, used once
	// config := struct {
	// 	Host string
	// 	Port int
	// }{
	// 	Host: "localhost",
	// 	Port: 8080,
	// }
	// fmt.Println(config.Host)

	// //methods
	// rect := Rect{Width: 3, Height: 4}
	// fmt.Println(rect.Area())
	// rect.Scale(2)  
	// fmt.Println(rect)

	// //embedding - Width comes from the embedded Rect
	box := Box{
		Rect: Rect{Width: 2, Height: 3}, 
		Depth: 4,
	}
	fmt.Println(box.Width)
	fmt.Println(box.Area()) // method is promoted too

	// //comparing structs - works if all fields are comparable
	a := Rect{1, 2}
	b := Rect{1, 2}
	fmt.Println(a == b) // true

	//slice of structs
	// students := []Student{
	// 	{Name: "Alice", Age: 16},
	// 	{Name: "Bob", Age: 17},
	// }
	// for i := range students {
	// 	students[i].Age++ // use index, range gives a copy
	// }
	// fmt.Println(students)

	// //map of structs
	byName := map[string]Student{
		"Alice": {Name: "Alice", Age: 16},
	}
	// byName["Alice"].Age = 20 // not allowed
	s := byName["Alice"]
	s.Age = 20
	byName["Alice"] = s
	fmt.Println(byName["Alice"])

}
