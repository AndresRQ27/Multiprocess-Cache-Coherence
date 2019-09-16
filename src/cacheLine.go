package main

import "fmt"

//CacheLine - struct that describes the values of each cache line (state, mem, data)
type CacheLine struct {
	Ownership bool
	tag int
	State, Data string
}

//CacheLineRead - method that returns the values inside the CacheLine struct
func (line *CacheLine) CacheLineRead() (bool, string, int, string) {

	return line.Ownership, line.State, line.tag, line.Data
}

//CacheLineWrite - method that writes the given parameters into the cacheLine
func (line *CacheLine) CacheLineWrite(ownership bool, tag int, data, state string)  {

	line.State = state
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
	line.State = "I"
	line.tag = -1
	line.Data = "0"
	line.Ownership = false
	return
}

//EmptyCacheLine - Constructor of CacheLine that initialize its value to default
func EmptyCacheLine() *CacheLine {
	return &CacheLine{
		State: "I",
		tag: -1,
		Data: "0",
		Ownership: false,
	}
}