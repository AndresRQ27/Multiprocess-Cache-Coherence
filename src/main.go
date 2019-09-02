package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World. How are you today?")
	fmt.Println("")

	cache := NewCache()
	fmt.Println(cache)
}
