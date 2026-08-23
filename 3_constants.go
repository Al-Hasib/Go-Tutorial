package main

import "fmt"

func main() {

	//constants
	const pi float64 = 3.14159
	const port int = 8080

	const (
		appName    = "MyApp"
		version    = "1.0.0"
		maxRetries = 5
	)

	fmt.Println(pi)
	fmt.Println(port)
	fmt.Println(appName)
	fmt.Println(version)
	fmt.Println(maxRetries)

}
