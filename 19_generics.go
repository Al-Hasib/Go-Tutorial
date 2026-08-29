package main

import "fmt"

// generics - write one function/type that works with many types
// [T any] means: T is a type parameter, "any" is its constraint

// without generics you'd need a separate function per type
func PrintSlice[T any](items []T) {
	for _, item := range items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}

// constraint: comparable lets us use == and !=
func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// a custom constraint - only these types are allowed
type Number interface {
	int | int64 | float64
}

func Sum[T Number](items []T) T {
	var total T
	for _, item := range items {
		total += item
	}
	return total
}

// two different type parameters
func Map[T, U any](items []T, f func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}

// generic struct
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}

// generic struct with two type parameters
type Pair[K, V any] struct {
	Key   K
	Value V
}

func main() {

	//works with int
	PrintSlice([]int{1, 2, 3})

	//same function works with string - no rewrite needed
	PrintSlice([]string{"a", "b", "c"})

	//Go usually infers T, but it can be written explicitly
	PrintSlice[float64]([]float64{1.1, 2.2})

	// //comparable constraint
	// fmt.Println(Contains([]int{1, 2, 3}, 2))       // true
	// fmt.Println(Contains([]string{"x", "y"}, "z")) // false

	// //custom constraint - only allows int, int64, float64
	// fmt.Println(Sum([]int{1, 2, 3}))      // 6
	// fmt.Println(Sum([]float64{1.5, 2.5})) // 4

	// //two type parameters - convert []int to []string
	// nums := []int{1, 2, 3}
	// strs := Map(nums, func(n int) string {
	// 	return fmt.Sprintf("#%d", n)
	// })
	// fmt.Println(strs)

	// //generic struct - a stack of int
	// intStack := Stack[int]{}
	// intStack.Push(10)
	// intStack.Push(20)
	// value, ok := intStack.Pop()
	// fmt.Println(value, ok) // 20 true

	// //same struct, different type - a stack of string
	// strStack := Stack[string]{}
	// strStack.Push("hello")
	// strStack.Push("world")
	// top, _ := strStack.Pop()
	// fmt.Println(top) // world

	// //generic struct with two type parameters
	// p := Pair[string, int]{Key: "age", Value: 25}
	// fmt.Println(p.Key, p.Value)

}
