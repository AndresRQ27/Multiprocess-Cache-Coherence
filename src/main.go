package main

import (
	"fmt"
	"time"
	"sync"
)

//PowerOff - bool to keep main alive or kill every thread if true
var PowerOff = false 
//BlocksInMemory - int that has the amount of memory lines in the simulated program
const BlocksInMemory = 16
//BlocksInCache - int that has the amount of cache lines in the simulated program
const BlocksInCache = 8
//Clock - int the sets the clock of the processor in seconds
const Clock = 1 * time.Second

func main() {
	
	mainChannel := make(chan int, 6)

	mux := sync.Mutex{}

	go func() {
		fmt.Printf("Start 1\n")
		for {
			select {
			case consumer1 := <- mainChannel:
				mux.Lock()
				fmt.Printf("%d at channel 1\n", consumer1)	
				mux.Unlock()
			}
		}
	}()
	go func() {
		fmt.Printf("Start 2\n")
		for {
			select {
			case consumer2 := <- mainChannel:
				mux.Lock()
				fmt.Printf("%d at channel 2\n", consumer2)	
				mux.Unlock()
			}
		}
	}()
	go func() {
		fmt.Printf("Start 3\n")
		for {
			select {
			case consumer3 := <- mainChannel:
				mux.Lock()
				fmt.Printf("%d at channel 3\n", consumer3)	
				mux.Unlock()
			}
		}
	}()
	var other int
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Ready sender")

		for index2 := 0; index2 < 6; index2++ {
			mux.Lock()
			fmt.Println(len(mainChannel))
			if len(mainChannel) < 1 {
				for index := 0; index < 4; index++ {
					mainChannel <- other
					other++
				}
				fmt.Println("Sent")
				time.Sleep(2*time.Second)
			} else {
				select {
				case consumer3 := <- mainChannel:
					fmt.Printf("%d at channel 4\n", consumer3)
					index2--
				}
			}
			mux.Unlock()
		}
		
	}()

	time.Sleep(30 * time.Second)
}