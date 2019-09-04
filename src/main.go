package main

import (
	"fmt"
)

//PowerOff - bool to keep main alive or kill every thread if true
var PowerOff = false 
//CellsInMemory - int that has the amount of memory lines in the simulated program
const CellsInMemory = 16
//CellsInCache - int that has the amount of cache lines in the simulated program
const CellsInCache = 8

func main() {
	fmt.Println("Hello World. How are you today?")
	fmt.Println("")

}
