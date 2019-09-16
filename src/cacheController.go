package main

import (
	"time"
	"sync"
	"fmt"
	"strconv"
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
	//Channels to listen for broadcasts
	PublicChannelCPU0 *chan Message
	PublicChannelCPU1 *chan Message
	PublicChannelCPU2 *chan Message
	PublicChannelCPU3 *chan Message
	//Channels to listen for responses
	PrivateChannelCPU0 *chan Message
	PrivateChannelCPU1 *chan Message
	PrivateChannelCPU2 *chan Message
	PrivateChannelCPU3 *chan Message

	GUIChannel *chan Message
}

//Listen - method in infinite loop that listens for requests
func (controller *CacheController) Listen() {
	myPublicChannel := controller.PublicChannelName(controller.Name)
	myPrivateChannel := controller.PrivateChannelName(controller.Name)
	for {
		select {
		case receivedMessage := <-*myPublicChannel: //Reads broadcast messages
			controller.UpdateInfo(receivedMessage)
		case receivedMessage := <-*controller.PrivateProcessor:  //Case when processor LDR or STR
			if receivedMessage.Value == "" { //If value is empty, Processor sent a LDR (read)
				controller.ProcessorRead(receivedMessage.Tag, myPublicChannel, myPrivateChannel)
			} else {
				controller.ProcessorWrite(receivedMessage.Tag, receivedMessage.Value, myPublicChannel)
			}
		}
	}
}

//ProcessorRead - method that looks for the memory value and returns it through a channel
func (controller *CacheController) ProcessorRead(memoryAddress int, myPublicChannel, myPrivateChannel *chan Message) {
	cacheOwner, cacheState, cacheTag, cacheData := controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //First, look in the cache
	msg := ""

	if cacheState == "I" { //Normal miss
		msg += "Read miss in block " + strconv.Itoa(memoryAddress) + " (I)"
		controller.ReadMiss(memoryAddress, myPublicChannel, myPrivateChannel)  //Updates the value in the cache for the correct one
		_, _, cacheTag, cacheData = controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Retrieve the new value
		*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the new value from cache to the processor

	} else if cacheState == "S" {
		if cacheTag == memoryAddress { //Read hit
			msg += "Read hit in block " + strconv.Itoa(memoryAddress) + " (S)"
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData from cache to the processor
		} else { //Allocation miss
			if cacheOwner { //If cacheOwner disappears, save the value to memory
				controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			}
			msg += "Read miss in block " + strconv.Itoa(memoryAddress) + " (Allocation-S)"
			controller.ReadMiss(memoryAddress, myPublicChannel, myPrivateChannel)  //Updates the value in the cache for the correct one
			_, _, cacheTag, cacheData = controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Retrieve the new value
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the new value from cache to the processor
		}

	} else if cacheState == "M" { //"M" always have ownership
		if cacheTag == memoryAddress { //Read hit
			msg += "Read hit in block " + strconv.Itoa(memoryAddress) + " (M)"
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData from cache to the processor
		} else { //Allocation miss
			msg += "Read miss in block " + strconv.Itoa(memoryAddress) + " (Allocation-M)"
			controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			controller.ReadMiss(memoryAddress, myPublicChannel, myPrivateChannel)  //Updates the value in the cache for the correct one
			_, _, cacheTag, cacheData = controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Retrieve the new value
			*controller.PrivateProcessor <- Message{Tag:cacheTag, Value:cacheData} //Sends the new value from cache to the processor
		}

	} else {
		fmt.Println("Invalid cacheState somewhere. Check your code")
	}

	msg += "\n"
	*controller.GUIChannel<- Message{Value:msg, CPU:controller.Name}
	return
}

//ProcessorWrite - method that writes the memory value
func (controller *CacheController) ProcessorWrite(memoryAddress int, memoryValue string, myPublicChannel *chan Message)  {
	_, cacheState, cacheTag, cacheData := controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineRead() //Get the important values from the cache
	msg := ""

	if cacheState == "I" { //Tag could not match, but invalidate where you write, not what you have
		msg += "Write miss in block " + strconv.Itoa(memoryAddress) + " (I)"
		controller.WriteMiss(memoryAddress, myPublicChannel)
		
	} else if cacheState == "S" {
		if cacheTag == memoryAddress {
			msg += "Write hit in block " + strconv.Itoa(memoryAddress) + " (Invalidate-S)"
			controller.Invalidate(memoryAddress, myPublicChannel)
		} else {
			msg += "Write miss in block " + strconv.Itoa(memoryAddress) + " (S)"
			controller.WriteMiss(memoryAddress, myPublicChannel)
		}

	} else if cacheState == "M" {
		if cacheTag == memoryAddress {
			//Do nothing, only you have it
			msg += "Write hit in block " + strconv.Itoa(memoryAddress) + " (M)"
		} else {
			msg += "Write miss in block " + strconv.Itoa(memoryAddress) + " (M)"
			controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
			controller.WriteMiss(memoryAddress, myPublicChannel)
		}

	} else {
		fmt.Println("Invalid cacheState somewhere. Check your code")
	}

	controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(true, memoryAddress, memoryValue, "M") //Write the new value to the cache
	*controller.PrivateProcessor <- Message{Tag:-1, Value:"ok"} //Sends ok to the processor
	
	msg += "\n"
	*controller.GUIChannel<- Message{Value:msg, CPU:controller.Name}
	return
}

//UpdateInfo - method that manages the received messages from MyChannel (from a broadcast)
func (controller *CacheController) UpdateInfo(message Message)  {
	cacheOwner, cacheState, cacheTag, cacheData := controller.PrivateCache[message.Tag%BlocksInCache].CacheLineRead() //Get the important values from the cache
	msg := ""

	//Responses for other threads. Use private channel
	if message.Value == "Read" { //In read, you answer always answer, either blank or the cacheData
		 askingChannel := *controller.PrivateChannelName(message.CPU) //Channel that is asking for the info

		if cacheTag == message.Tag {
			if cacheState == "S" {
				if cacheOwner { //You're the cacheOwner. Answer the call
					askingChannel <-Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData
					controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "S") //Lose ownership
					msg += "Block:"+ strconv.Itoa(message.Tag) + " SO -> S by " + message.CPU + "\n"
				} else {
					askingChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
				}

			} else if cacheState == "M" { //If M, always the cacheOwner
				askingChannel <-Message{Tag:cacheTag, Value:cacheData} //Sends the cacheData
				controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "S") //Lose ownership and change cacheState
				msg += "Block:"+ strconv.Itoa(message.Tag) + " MO -> S by " + message.CPU + "\n"

			} else if cacheState == "I" {
				askingChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
				
			} else {
				askingChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
				fmt.Println("Invalid cacheState")
			}	
		} else {
			askingChannel <-Message{Tag:cacheTag, Value:""} //Sends empty
		}

	} else if message.Value == "Write" { //In write, you update the states
		if cacheTag == message.Tag {
			if cacheState == "S" {
				if cacheOwner { //"S", but owner. So basically "M"
					*controller.GUIChannel<- Message{Value:"Memory accessed\n", CPU:controller.Name}
					controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory	
					msg += "Block:"+ strconv.Itoa(message.Tag) + " SO -> I by " + message.CPU + "\n"
				} else {
					msg += "Block:"+ strconv.Itoa(message.Tag) + " S -> I by " + message.CPU + "\n"
				}
				controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "I")	
	
			} else if cacheState == "M" {
				*controller.GUIChannel<- Message{Value:"Memory accessed\n", CPU:controller.Name}
				controller.SharedMemory.MemoryWrite(cacheTag, cacheData) //Write-back to memory
				controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "I")
				msg += "Block:"+ strconv.Itoa(message.Tag) + " M -> I by " + message.CPU + "\n"
	
			} //Exclude "I" cacheState. No need to answer or do anything	
		}//Ignore if it's not the tag

	} else if message.Value == "Invalidate" { //No write-back since is an upgrade of ownership
		controller.PrivateCache[message.Tag%BlocksInCache].CacheLineWrite(false, cacheTag, cacheData, "I")
		msg += "Block:"+ strconv.Itoa(message.Tag) + " S -> I (Invalidate) by " + message.CPU + "\n"
		
	} else {
		fmt.Println("Thread expected a message update but received an answer with cacheData.")
	}	

	*controller.GUIChannel<- Message{Value:msg, CPU:controller.Name}
	return
}

//ReadMiss - method that manages the procedure generated from a cache miss
func (controller *CacheController) ReadMiss(memoryAddress int, myPublicChannel, myPrivateChannel *chan Message)  {
	controller.Mux.Lock()
	for len(*myPublicChannel) > 0 {
		controller.UpdateInfo(<-*myPublicChannel)	
	}
	controller.BroadcastInfo(memoryAddress, "Read", controller.Name) //Tells the other CC a cache miss occured
	controller.Mux.Unlock()

	//Receive a message from every processor. Just one contains info
	receivedMessage1 := <-*myPrivateChannel
	receivedMessage2 := <-*myPrivateChannel
	receivedMessage3 := <-*myPrivateChannel

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
		*controller.GUIChannel<- Message{Value:"Memory accessed\n", CPU:controller.Name}
		memoryValue := controller.SharedMemory.MemoryRead(memoryAddress)
		controller.PrivateCache[memoryAddress%BlocksInCache].CacheLineWrite(
			true, memoryAddress, memoryValue, "S")
	}
	return
}

//WriteMiss - method that broadcast a WriteMiss to the other processors
func (controller *CacheController) WriteMiss(memoryAddress int, myPublicChannel *chan Message)  {
	controller.Mux.Lock()
	for index := 0; index < len(*myPublicChannel); index++ {
		controller.UpdateInfo(<-*myPublicChannel)	
	}
	controller.BroadcastInfo(memoryAddress, "Write", controller.Name) //Tells the other CC a cache miss occured
	controller.Mux.Unlock()
}

//Invalidate - method that broadcast to invalidate a given memoryAddress if in cache
func (controller *CacheController) Invalidate(memoryAddress int, myPublicChannel *chan Message)  {
	controller.Mux.Lock()
	if len(*myPublicChannel) < 1 { //There are unread messages. Handle this first
		controller.UpdateInfo(<-*myPublicChannel)
	}
	controller.BroadcastInfo(memoryAddress, "Invalidate", controller.Name) //Tells the other CC a cache miss occured
	controller.Mux.Unlock()
}

//BroadcastInfo - method that broadcast to all CPUs except myself. Penalty receive from the usage
func (controller *CacheController) BroadcastInfo(cacheTag int, value, cpu string) {
	//Sends 1 message to each of the 3 processors

	if controller.Name != "CPU0" {*controller.PublicChannelCPU0 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	
	if controller.Name != "CPU1" {*controller.PublicChannelCPU1 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	
	if controller.Name != "CPU2" {*controller.PublicChannelCPU2 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	
	if controller.Name != "CPU3" {*controller.PublicChannelCPU3 <-Message{Tag:cacheTag, Value:value, CPU:cpu}}	

	time.Sleep(2 * Clock) //Penalization time for using the bus
	return
}

//PublicChannelName - returns the correct channel depending on the given name. Useful when listening to answers or sending answers
func (controller *CacheController) PublicChannelName(publicChannelName string) *chan Message {
	switch publicChannelName {
	case "CPU0": 
		return controller.PublicChannelCPU0
	case "CPU1":
		return controller.PublicChannelCPU1
	case "CPU2":
		return controller.PublicChannelCPU2
	case "CPU3":
		return controller.PublicChannelCPU3
	default:
		return nil //Should never get here, unless you fuck up
	}
}

//PrivateChannelName - returns the correct channel depending on the given name. Useful when listening to answers or sending answers
func (controller *CacheController) PrivateChannelName(privateChannelName string) *chan Message {
	switch privateChannelName {
	case "CPU0": 
		return controller.PrivateChannelCPU0
	case "CPU1":
		return controller.PrivateChannelCPU1
	case "CPU2":
		return controller.PrivateChannelCPU2
	case "CPU3":
		return controller.PrivateChannelCPU3
	default:
		return nil //Should never get here, unless you fuck up
	}
}