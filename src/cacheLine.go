package main

import "fmt"

//CacheLine - struct that describes the values of each cache line (state, mem, data)
type CacheLine struct {
	Ownership bool
	State, tag int
	Data string
}

//CacheLineRead - method that returns the values inside the CacheLine struct
func (line *CacheLine) CacheLineRead() (bool, string, int, string) {
	var state string

	switch line.State {
	case 0: 
		state = "M"
	case 1: 
		state = "S"
	case 2: 
		state = "I"
	default: 
		fmt.Println("Invalid value given")
		state = ""
	}

	return line.Ownership, state, line.tag, line.Data
}

//CacheLineWrite - method that writes the given parameters into the cacheLine
func (line *CacheLine) CacheLineWrite(ownership bool, tag int, data, state string)  {

	switch state {
	case "M": 
		line.State = 0
	case "S": 
		line.State = 1
	case "I": 
		line.State = 2
	default: 
		fmt.Println("Invalid state given")
		line.State = 2
	}

	line.Data = data
	line.tag = tag
	line.Ownership = ownership
	return
}

/**
	Interface Stringer that prints the values of the CacheLine
 */
func (line *CacheLine) String() string {
	return fmt.Sprintf("State: %v, tag: %v, Data: %v", line.State, line.tag, line.Data)
}

//ClearCacheLine - method that resets its value to default
func (line *CacheLine) ClearCacheLine() {
	line.State = 2
	line.tag = -1
	line.Data = ""
	line.Ownership = false
	return
}

//EmptyCacheLine - Constructor of CacheLine that initialize its value to default
func EmptyCacheLine() *CacheLine {
	return &CacheLine{
		State: 2,
		tag: -1,
		Data:    "",
		Ownership: false,
	}
}