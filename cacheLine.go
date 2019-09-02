package main

import "fmt"

//CacheLine - struct that describes the values of each cache line (state, mem, data)
type CacheLine struct {
	Snoop *SnoopProtocol
	MemCell int
	Data string
}

/**
	Interface Stringer that prints the values of the CacheLine
 */
func (cacheLine CacheLine) String() string {
	return fmt.Sprintf("Snoop: %v, MemCell: %v, Data: %v", cacheLine.Snoop, cacheLine.MemCell, cacheLine.Data)
}

//ClearCacheLine - method that resets its value to default
func (cacheLine CacheLine) ClearCacheLine() {
	cacheLine.Snoop.M = false
	cacheLine.Snoop.S = false
	cacheLine.Snoop.I = true
	cacheLine.MemCell = -1
	cacheLine.Data = ""
	return
}

//EmptyCacheLine - Constructor of CacheLine that initialize its value to default
func EmptyCacheLine() *CacheLine  {
	return &CacheLine{
		Snoop:   NewSnoopProtocol(),
		MemCell: -1,
		Data:    "",
	}
}