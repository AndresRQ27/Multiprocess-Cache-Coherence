package main

import (
	"sync"
	"time"
)

//Memory - struct that abstracts Memory usage with mutex
type Memory struct {
	MemoryMap map[int]string
	Mux sync.Mutex
}

//MemoryRead - method that reads the memoryMap, according to the cellNumber, with a delay
func (mem *Memory) MemoryRead(memoryAddress int) string {
	mem.Mux.Lock()
	defer mem.Mux.Unlock()
	
	time.Sleep(5 * Clock) //Penalization time for using the bus
	return mem.MemoryMap[memoryAddress]
}

//MemoryWrite - method that writes a newValue in the memoryMap, according to the cellNumber, with a delay
func (mem *Memory) MemoryWrite(memoryAddress int, memValue string) {
	mem.Mux.Lock()
	defer mem.Mux.Unlock()

	time.Sleep(5 * Clock) //Penalization time for using the bus

	mem.MemoryMap[memoryAddress] = memValue
	return
}

//NewMemory - constructor of a empty Memory
func NewMemory() *Memory {
	memMap := map[int]string{
		0: "",
		1: "",
		2: "",
		3: "",
		4: "",
		5: "",
		6: "",
		7: "",
		8: "",
		9: "",
		10: "",
		11: "",
		12: "",
		13: "",
		14: "",
		15: "",
	}

	mem := Memory{
		MemoryMap: memMap,
	}
	return &mem
}
