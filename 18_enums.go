package main

import "fmt"

// Go has no "enum" keyword.
// Instead: a named type + a group of constants, usually built with iota.

// iota starts at 0 inside a const block and increases by 1 each line
type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

// giving names to the values - makes them print nicely
func (d Weekday) String() string {
	names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	return names[d]
}

// skipping a value with _
type Size int

const (
	Small Size = iota // 0
	_                 // 1 skipped
	Large             // 2
)

// starting from a different number
type HTTPStatus int

const (
	StatusOK       HTTPStatus = 200
	StatusNotFound HTTPStatus = 404
	StatusError    HTTPStatus = 500
)

// iota with a formula - powers of two, common for bit flags
type Permission int

const (
	Read    Permission = 1 << iota // 1
	Write                          // 2
	Execute                        // 4
)

// plain constants without iota still work as a simple enum
type Direction int

const (
	Up    Direction = 0
	Down  Direction = 1
	Left  Direction = 2
	Right Direction = 3
)

func main() {

	//using the enum like a normal value
	today := Wednesday
	fmt.Println(today)      // Wednesday (because of String())
	fmt.Println(int(today)) // 3

	//comparing enum values
	if today == Wednesday {
		fmt.Println("midweek")
	}

	//enum in a switch
	switch today {
	case Saturday, Sunday:
		fmt.Println("weekend")
	default:
		fmt.Println("weekday")
	}

	//skipped value
	fmt.Println(Small, Large) // 0 2

	//fixed values, not from iota
	fmt.Println(StatusOK, StatusNotFound, StatusError)

	//bit flag style enum - combine with |
	perms := Read | Write
	fmt.Println(perms) // 3

	hasWrite := perms&Write != 0
	fmt.Println(hasWrite) // true

	hasExecute := perms&Execute != 0
	fmt.Println(hasExecute) // false

	//looping over all enum values
	for d := Sunday; d <= Saturday; d++ {
		fmt.Print(d, " ")
	}
	fmt.Println()

	//plain enum without iota
	dir := Left
	fmt.Println(dir) // 2 (no String() method, so it prints the number)

}
