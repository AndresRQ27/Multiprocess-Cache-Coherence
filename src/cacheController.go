package main

import (
	"time"
	"sync"
	"fmt"
)

/**
	If someone is requesting a value in another cache, first update the snoop and then send
	When updating the snoop via broadcast, use memoryAddress mod 8 to only ask in the correct cacheLine
	Use a channel switch to maintain active the broadcast channel and the other CPUs channels
 */

//CacheController - struct that manages request for cache and Memory usage in a coherent way
type CacheController struct {
	Name string
	PrivateCache map[int]*CacheLine
	PrivateProcessor *chan Message

	SharedMemory *Memory

	Mux *sync.Mutex
	ChannelCPU0 *chan Message
	ChannelCPU1 *chan Message
	ChannelCPU2 *chan Message
	ChannelCPU3 *chan Message
}

//Listen - method in infinite loop that listens for requests
func (controller *CacheController) Listen() {
	myChannel := controller.ChannelName(controller.Name);
	for {
		select {
		case receivedMessage := <-*myChannel:
			controller.UpdateInfo(receivedMessage)
		case receivedMessage := <-*controller.PrivateProcessor:  //Case when processor LDR or STR
			if receivedMessage.Value == "" { //If value is empty, Processor sent a LDR (read)
				controller.ProcessorRead(receivedMessage.Tag, myChannel)
			} else {
				controller.ProcessorWrite(receivedMessage.Tag, receivedMessage.Value, myChannel)
			}
		}
	}
}

//ProcessorRead - method that looks for the memory value and returns it through a channel
func (controller *CacheController) ProcessorRead(memoryAddress int, myChannel *chan Message) {
	cacheOwner, cacheState, cacheTag, cacheData := controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //First, look in the cache

	if cacheState == "I" { //Normal miss
		controller.ReadMiss(memoryAddress, myChannel)  //Updates the value in the cache for the correct one
		_, _, cacheTag, cacheData = controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Retrieve the new value
		*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the new value from cache to the processor

	} else if cacheState == "S" {
		if cacheTag == memoryAddress {
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData from cache to the processor
		} else { //Allocation miss
			if cacheOwner { //If cacheOwner disappears, save the value to memory
				controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			}
			controller.ReadMiss(memoryAddress, myChannel)  //Updates the value in the cache for the correct one
			_, _, cacheTag, cacheData = controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Retrieve the new value
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the new value from cache to the processor
		}

	} else if cacheState == "M" {
		if cacheTag == memoryAddress {
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData from cache to the processor
		} else { //Allocation miss
			controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			controller.ReadMiss(memoryAddress, myChannel)  //Updates the value in the cache for the correct one
			_, _, cacheTag, cacheData = controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Retrieve the new value
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the new value from cache to the processor
		}

	} else {
		fmt.Println("Invalid cacheState somewhere. Check your code")
	}
	return
}

//ProcessorWrite - method that writes the memory value
func (controller *CacheController) ProcessorWrite(memoryAddress int, memoryValue string, myChannel *chan Message)  {
	_, cacheState, cacheTag, cacheData := controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Get the important values from the cache
	
	if cacheState == "I" { //Tag could not match, but invalidate where you write, not what you have
		controller.WriteMiss(memoryAddress, myChannel)
		
	} else if cacheState == "S" {
		if cacheTag == memoryAddress {
			controller.Invalidate(memoryAddress, myChannel)
		} else {
			controller.WriteMiss(memoryAddress, myChannel)
		}

	} else if cacheState == "M" {
		if cacheTag == memoryAddress {
			//Do nothing, only you have it
		} else {
			controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			controller.WriteMiss(memoryAddress, myChannel)
		}

	} else {
		fmt.Println("Invalid cacheState somewhere. Check your code")
	}

	controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(true, memoryAddress, memoryValue, "M") //Write the new value to the cache
	*controller.PrivateProcessor <- Message{Tag:-1, Value:"ok"} //Sends ok to the processor
	return
}

//UpdateInfo - method that manages the received messages from MyChannel (from a broadcast)
func (controller *CacheController) UpdateInfo(message Message)  {
	cacheOwner, cacheState, cacheTag, cacheData := controller.PrivateCache[message.Tag%BlocksInCache].CacheLineRead() //Get the important values from the cache
	
	if message.Value == "Read" { //In read, you answer always answer, either blank or the cacheData
		 questionChannel := *controller.ChannelName(message.CPU) //Channel that is asking for the info

		if cacheTag == message.Tag {

			if cacheState == "S" {
				if cacheOwner { //You're the cacheOwner. Answer the call
					questionChannel <-Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData
					controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "S") //Lose ownership
				} else {
					questionChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
				}

			} else if cacheState == "M" { //If M, always the cacheOwner
				questionChannel <-Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData
				controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "S") //Lose ownership and change cacheState

			} else if cacheState == "I" {
				questionChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
				
			} else {
				questionChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
				fmt.Println("Invalid cacheState")
			}	
		} else {
			questionChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
		}

	} else if message.Value == "Write" { //In write, you update the states

		if cacheState == "S" {
			if cacheOwner { //"S", but owner. So basically "M"
				controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory	
			}
			controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "I")	

		} else if cacheState == "M" {
			controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "I")

		} //Exclude "I" cacheState. No need to answer or do anything

	} else if message.Value == "Invalidate" { //No write-back since is an upgrade of ownership
		controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "I")
		
	} else {
		fmt.Println("Thread expected a message update but received an answer with cacheData.")
	}	
	return
}

//ReadMiss - method that manages the procedure generated from a cache miss
func (controller *CacheController) ReadMiss(memoryAddress int, myChannel *chan Message)  {
	controller.Mux.Lock()
	if len(*myChannel) < 1 { //There are unread messages. Handle this first
		controller.UpdateInfo(<-*myChannel)
	}
	controller.BroadcastInfo(memoryAddress, "Read", controller.Name) //Tells the other CC a cache miss occured
	controller.Mux.Unlock()

	//Receive a message from every processor. Just one contains info
	receivedMessage1 := <-*myChannel
	receivedMessage2 := <-*myChannel
	receivedMessage3 := <-*myChannel

	if receivedMessage1.Value != "" {//Check if the first message has the value
		controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(
			true, receivedMessage1.Tag, receivedMessage1.Value, "S")
	} else if receivedMessage2.Value != "" {//Check if the second message has the value
		controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(
			true, receivedMessage2.Tag, receivedMessage2.Value, "S")
	} else if receivedMessage3.Value != "" {//Check if the third message has the value
		controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(
			true, receivedMessage3.Tag, receivedMessage3.Value, "S")
	} else {//Go get the message from memory
		memoryValue := controller.SharedMemory.MemoryRead(memoryAddress)
		controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(
			true, memoryAddress, memoryValue, "S")
	}
	return
}

//WriteMiss - method that broadcast a WriteMiss to the other processors
func (controller *CacheController) WriteMiss(memoryAddress int, myChannel *chan Message)  {
	controller.Mux.Lock()
	if len(*myChannel) < 1 { //There are unread messages. Handle this first
		controller.UpdateInfo(<-*myChannel)
	}
	controller.BroadcastInfo(memoryAddress, "Write", controller.Name) //Tells the other CC a cache miss occured
	controller.Mux.Unlock()
}

//Invalidate - method that broadcast to invalidate a given memoryAddress if in cache
func (controller *CacheController) Invalidate(memoryAddress int, myChannel *chan Message)  {
	controller.Mux.Lock()
	if len(*myChannel) < 1 { //There are unread messages. Handle this first
		controller.UpdateInfo(<-*myChannel)
	}
	controller.BroadcastInfo(memoryAddress, "Invalidate", controller.Name) //Tells the other CC a cache miss occured
	controller.Mux.Unlock()
}

//BroadcastInfo - method that broadcast to all CPUs except myself. Penalty receive from the usage
func (controller *CacheController) BroadcastInfo(cacheTag int, value, cpu string) {
	//Sends 1 message to each of the 3 processors

	if controller.Name != "CPU0" {*controller.ChannelCPU0 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	
	if controller.Name != "CPU1" {*controller.ChannelCPU1 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	
	if controller.Name != "CPU2" {*controller.ChannelCPU2 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	
	if controller.Name != "CPU3" {*controller.ChannelCPU3 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	

	time.Sleep(2 * Clock) //Penalization time for using the bus
	return
}

//ChannelName - returns the correct channel depending on the given name. Useful when listening to answers or sending answers
func (controller *CacheController) ChannelName(channelName string) *chan Message {
	switch channelName {
	case "CPU0": 
		return controller.ChannelCPU0
	case "CPU1":
		return controller.ChannelCPU1
	case "CPU2":
		return controller.ChannelCPU2
	case "CPU3":
		return controller.ChannelCPU3
	default:
		return nil //Should never get here, unless you fuck up
	}
}

/*TODO: Challenges with many LDRs & STRs
* Not so many CPU available to process all the broadcasted messages
*** CPU broadcasting the same message more than once: multiple sends to same CPU channel
****** Trash in the CPU channel. Trash in the broadcast channel
* TODO: do not send CPU name when responding to CPU Channel
* TODO: should only the cacheOwner do a write-back during a read miss?
*/