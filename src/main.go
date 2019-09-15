package main

import (
	"math/rand"
	"time"
	"sync"
	"fmt"
)

//BlocksInMemory - int that has the amount of memory lines in the simulated program
const BlocksInMemory = 16

//BlocksInCache - int that has the amount of cache lines in the simulated program
const BlocksInCache = 8

//Clock - int the sets the clock of the processor in seconds
const Clock = 1 * time.Second

//InstructionCounter - counter of the number of instructions executed
var InstructionCounter = 0

//Mean - bigger = +STR instructions / smaller = +LDR instructions
const Mean = -100
//StdDev - bigger = +STD/LDR instructions / smaller = -STD/LDR instructions
const StdDev = 2
/////Default values generates random distributed numbers between 0 and 20

func main() {
	//Shared resources created
	SharedMux := sync.Mutex{}
	SharedMemory := Memory{
		MemoryMap:map[int]string{
			0:"0",
			1:"0",
			2:"0",
			3:"0",
			4:"0",
			5:"0",
			6:"0",
			7:"0",
			8:"0",
			9:"0",
			10:"0",
			11:"0",
			12:"0",
			13:"0",
			14:"0",
			15:"0",
		},
	}

	//Max buffer of the channel can hold 3 responses from needed data and 3 broadcast messages
	CPU0PublicChannel := make(chan Message, 6)
	CPU1PublicChannel := make(chan Message, 6)
	CPU2PublicChannel := make(chan Message, 6)
	CPU3PublicChannel := make(chan Message, 6)

	CPU0PrivateChannel := make(chan Message, 6)
	CPU1PrivateChannel := make(chan Message, 6)
	CPU2PrivateChannel := make(chan Message, 6)
	CPU3PrivateChannel := make(chan Message, 6)

	//CPU0 creation
	CPU0ProcessorChannel := make(chan Message)
	CPU0Rand := rand.New(rand.NewSource(0))
	
	CPU0CC := CacheController{
		Name:"CPU0",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU0ProcessorChannel,
		SharedMemory:&SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,
	}

	CPU0Processor := Processor{
		Name:"CPU0",
		InstructionNumber:0,
		PublicCacheController:&CPU0ProcessorChannel,
		Random:CPU0Rand,
	}

	//CPU1 creation
	CPU1ProcessorChannel := make(chan Message)
	CPU1Rand := rand.New(rand.NewSource(1))
	
	CPU1CC := CacheController{
		Name:"CPU1",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU1ProcessorChannel,
		SharedMemory:&SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,
	}

	CPU1Processor := Processor{
		Name:"CPU1",
		InstructionNumber:0,
		PublicCacheController:&CPU1ProcessorChannel,
		Random:CPU1Rand,
	}

	//CPU2 creation
	CPU2ProcessorChannel := make(chan Message)
	CPU2Rand := rand.New(rand.NewSource(2))
	
	CPU2CC := CacheController{
		Name:"CPU2",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU2ProcessorChannel,
		SharedMemory:&SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,
	}

	CPU2Processor := Processor{
		Name:"CPU2",
		InstructionNumber:0,
		PublicCacheController:&CPU2ProcessorChannel,
		Random:CPU2Rand,
	}

	//CPU3 creation
	CPU3ProcessorChannel := make(chan Message)
	CPU3Rand := rand.New(rand.NewSource(3))
	
	CPU3CC := CacheController{
		Name:"CPU3",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU3ProcessorChannel,
		SharedMemory:&SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,
	}

	CPU3Processor := Processor{
		Name:"CPU3",
		InstructionNumber:0,
		PublicCacheController:&CPU3ProcessorChannel,
		Random:CPU3Rand,
	}

	//Initiate all the go routines
	go CPU0CC.Listen()
	go CPU0Processor.ExecuteInstruction()

	go CPU1CC.Listen()
	go CPU1Processor.ExecuteInstruction()

	go CPU2CC.Listen()
	go CPU2Processor.ExecuteInstruction()

	go CPU3CC.Listen()
	go CPU3Processor.ExecuteInstruction()

	time.Sleep(Clock) //Sleeps for 1 second during the first start

	//Infinite loop that manages the counter
	for {
		InstructionCounter++ //Add 1 to the instruction counter
		if InstructionCounter == 15 {
			break
		}
		time.Sleep(10*Clock) //Sleeps for 1 second
		fmt.Println()
	}

	return

	/*r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	mymap := make(map[int]int)

	for i := 0; i < 10000000; i++ {
		x := int(r.NormFloat64()*float64(2) + float64(10)) //From 0 to 20 normally distributed
		_, ok := mymap[x]
		if !ok {
			mymap[x] = 1
		} else {
			mymap[x]++
		}
	}
	fmt.Println(mymap)
	return*/

	/*mainChannel := make(chan int, 6)

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

	time.Sleep(30 * time.Second)*/
}
