package main

import "fmt"

// interface - a set of method names, no implementation
// any type that has these methods "satisfies" the interface automatically
// (no "implements" keyword needed)

type Shape interface {
	Area() float64
}

type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// a function can accept the interface instead of one specific type
func printArea(s Shape) {
	fmt.Println(s.Area())
}

// interface with more than one method
type Animal interface {
	Sound() string
	Name() string
}

type Dog struct{}

func (d Dog) Sound() string { return "Woof" }
func (d Dog) Name() string  { return "Dog" }

type Cat struct{}

func (c Cat) Sound() string { return "Meow" }
func (c Cat) Name() string  { return "Cat" }

// the empty interface accepts a value of ANY type
func describe(v interface{}) {
	fmt.Println(v)
}

func main() {

	//Square and Rectangle are different types, but both satisfy Shape
	sq := Square{Side: 2}
	r := Rectangle{Width: 3, Height: 4}

	printArea(sq)
	printArea(r)

	//a slice of the interface type can hold different concrete types
	shapes := []Shape{sq, r}
	for _, s := range shapes {
		fmt.Println(s.Area())
	}

	//an interface variable can hold any type that satisfies it
	var shape Shape
	shape = sq
	fmt.Println(shape.Area())
	shape = r
	fmt.Println(shape.Area())

	//multiple methods
	animals := []Animal{Dog{}, Cat{}}
	for _, a := range animals {
		fmt.Println(a.Name(), "says", a.Sound())
	}

	//type assertion - get the concrete type back out of an interface
	var s Shape = Square{Side: 5}
	square, ok := s.(Square)
	fmt.Println(square, ok) // {5} true

	_, ok = s.(Rectangle)
	fmt.Println(ok) // false

	//type switch - check which concrete type it is
	checkShape := func(s Shape) {
		switch v := s.(type) {
		case Square:
			fmt.Println("square with side", v.Side)
		case Rectangle:
			fmt.Println("rectangle", v.Width, "x", v.Height)
		default:
			fmt.Println("unknown shape")
		}
	}
	checkShape(sq)
	checkShape(r)

	//empty interface - accepts anything
	describe(42)
	describe("hello")
	describe(sq)
	describe(true)

}
