package main

import "fmt"

//CacheLine - struct that describes the values of each cache line (state, mem, data)
type CacheLine struct {
	State, MemCell int
	Data string
}

//CacheLineRead - method that returns the values inside the CacheLine struct
func (line *CacheLine) CacheLineRead() (string, string, int) {
	var state string

	switch line.State {
		case 0: {
			state = "M"
		}
		case 1: {
			state = "S"
		}
		case 2: {
			state = "I"
		}
		default: {
			fmt.Println("Invalid value given")
			state = ""
		}
	}

	return line.Data, state, line.MemCell
}

//CacheLineWrite - method that writes the given parameters into the cacheLine
func (line *CacheLine) CacheLineWrite(memCell int, data, state string)  {

	switch state {
		case "M": {
			line.State = 0
		}
		case "S": {
			line.State = 1
		}
		case "I": {
			line.State = 2
		}
		default: {
			fmt.Println("Invalid state given")
			line.State = -1
		}
	}

	line.Data = data
	line.MemCell = memCell
	return
}

/**
	Interface Stringer that prints the values of the CacheLine
 */
func (line *CacheLine) String() string {
	return fmt.Sprintf("State: %v, MemCell: %v, Data: %v", line.State, line.MemCell, line.Data)
}

//ClearCacheLine - method that resets its value to default
func (line *CacheLine) ClearCacheLine() {
	line.State = -1
	line.MemCell = -1
	line.Data = ""
	return
}

//EmptyCacheLine - Constructor of CacheLine that initialize its value to default
func EmptyCacheLine() *CacheLine {
	return &CacheLine{
		State: -1,
		MemCell: -1,
		Data:    "",
	}
}