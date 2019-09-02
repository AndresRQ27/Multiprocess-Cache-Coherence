package main

import (
	"sync"
	"time"
)

//Memory - struct that abstracts memory usage with mutex
type Memory struct {
	MemoryMap map[int]string
	Mux *sync.Mutex
}

//MemoryRead - function that reads the memoryMap, according to the cellNumber, with a delay
func (mem Memory) MemoryRead(memCell int) string {
	mem.Mux.Lock()
	defer mem.Mux.Unlock()

	time.Sleep(5000 * time.Millisecond) //Sleeps for 5 seconds
	return mem.MemoryMap[memCell]
}

//MemoryWrite - function that writes a newValue in the memoryMap, according to the cellNumber, with a delay
func (mem Memory) MemoryWrite(memCell int, memValue string) {
	mem.Mux.Lock()
	defer mem.Mux.Unlock()

	time.Sleep(5000 * time.Millisecond) //Sleeps for 5 seconds

	mem.MemoryMap[memCell] = memValue
	return
}

//NewMemory - constructor of a empty memory
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

	memory := Memory{
		MemoryMap: memMap,
		Mux:       &sync.Mutex{},
	}
	return &memory
}
