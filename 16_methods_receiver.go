package main

import "fmt"

// methods - a function attached to a type
// the (c Circle) part before the name is the RECEIVER

type Circle struct {
	Radius float64
}

// value receiver - gets a COPY, so it can only read
func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

// pointer receiver - gets the ADDRESS, so it can change the struct
func (c *Circle) Grow(n float64) {
	c.Radius = c.Radius + n
}

// value receiver cannot change anything - this is a copy
func (c Circle) GrowWrong(n float64) {
	c.Radius = c.Radius + n
}

type Counter struct {
	Count int
}

func (c *Counter) Increment() {
	c.Count++
}

func (c *Counter) Reset() {
	c.Count = 0
}

// methods are not only for structs
// any type YOU declare in this package can have methods
type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
	return float64(c)*1.8 + 32
}

type Numbers []int

func (n Numbers) Sum() int {
	total := 0
	for _, v := range n {
		total += v
	}
	return total
}

// String() is special - fmt uses it when printing the value
type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// a normal function does the same thing, but is not tied to the type
func areaFunc(c Circle) float64 {
	return 3.14 * c.Radius * c.Radius
}

func main() {

	//calling a value receiver method
	c := Circle{Radius: 2}
	// fmt.Println(c.Area())

	// // //pointer receiver changes the original
	// c.Grow(3)
	// fmt.Println(c.Radius) // 5

	// // //value receiver cannot change the original
	// c.GrowWrong(100)
	// fmt.Println(c.Radius) // still 5

	// // //Go adds & automatically, so these two are the same
	// c.Grow(1)
	// (&c).Grow(1)
	// fmt.Println(c.Radius) // 7

	// //Go adds * automatically too, so a pointer can call both kinds
	// p := &Circle{Radius: 10}
	// fmt.Println(p.Area()) // same as (*p).Area()
	// p.Grow(5)
	// fmt.Println(p.Radius)

	// //pointer receivers are useful for changing state
	// counter := Counter{}
	// counter.Increment()
	// counter.Increment()
	// counter.Increment()
	// fmt.Println(counter.Count) // 3
	// counter.Reset()
	// fmt.Println(counter.Count) // 0

	// //method on a non-struct type
	// temp := Celsius(100)
	// fmt.Println(temp.ToFahrenheit()) // 212

	// //method on a slice type
	// nums := Numbers{1, 2, 3, 4, 5}
	// fmt.Println(nums.Sum()) // 15

	// //String() is called automatically when printing
	// point := Point{X: 3, Y: 7}
	// fmt.Println(point)   // (3, 7)
	// fmt.Println(point.X, point.Y) // fields still work normally

	// //method vs function - same result, different style
	// fmt.Println(c.Area())
	// fmt.Println(areaFunc(c))

	// // //a method can be stored in a variable
	// getArea := c.Area
	// fmt.Println(getArea())

}
