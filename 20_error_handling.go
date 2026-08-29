package main

import (
	"errors"
	"fmt"
)

// error is just an interface with one method
// type error interface { Error() string }

// simplest way to create an error
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// fmt.Errorf - build an error message with values inside it
func findUser(id int) (string, error) {
	users := map[int]string{1: "Alice", 2: "Bob"}
	name, ok := users[id]
	if !ok {
		return "", fmt.Errorf("user %d not found", id)
	}
	return name, nil
}

// sentinel error - a fixed error value other code can check against
var ErrNotFound = errors.New("not found")

func findItem(id int) error {
	if id != 1 {
		return ErrNotFound
	}
	return nil
}

// wrapping - keep the original error attached with %w
func loadConfig() error {
	err := findItem(5)
	if err != nil {
		return fmt.Errorf("loadConfig failed: %w", err)
	}
	return nil
}

// custom error type - a struct that implements Error()
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Msg: "must not be negative"}
	}
	return nil
}



func main() {

	//basic error check - the standard Go pattern
	// result, err := divide(10, 0)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// } else {
	// 	fmt.Println("result:", result)
	// }

	// // //the error case
	// _, err = divide(10, 0)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }

	// //error with a formatted message
	// name, err := findUser(1)
	// fmt.Println(name, err) // user 5 not found

	// //comparing against a sentinel error
	err := findItem(5)
	// if errors.Is(err, ErrNotFound) {
	// 	fmt.Println("sentinel matched: item missing")
	// }

	// // //wrapped error still matches with errors.Is
	// err = loadConfig()
	// fmt.Println(err)                         // loadConfig failed: not found
	// fmt.Println(errors.Is(err, ErrNotFound)) // true, even though wrapped

	// //custom error type
	err = validateAge(-5)
	fmt.Println(err)

	// //errors.As - pull the concrete type back out of an error
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		fmt.Println("field:", valErr.Field)
	}

	// //nil error means "no problem"
	err = validateAge(20)
	fmt.Println(err == nil) // true


	// //multiple errors combined into one (Go 1.20+)
	err1 := errors.New("first problem")
	err2 := errors.New("second problem")
	combined := errors.Join(err1, err2)
	fmt.Println(combined)

}
