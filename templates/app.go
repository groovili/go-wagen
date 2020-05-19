package main

import "fmt"

func main() {
	app := "{{.Name}}"
	fmt.Printf("Hello %s!", app)
}
