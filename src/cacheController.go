package main

/**
	If someone is requesting a value in another cache, first update the snoop and then send
	When updating the snoop via broadcast, use memCell mod 8 to only ask in the correct cacheLine
	Use a channel switch to maintain active the broadcast channel and the other CPUs channels
 */

//CacheController - struct that manages request for cache and Memory usage in a coherent way
type CacheController struct {
	Name string
	PrivateCache map[int]*CacheLine

	BroadcastChannel *chan Broadcast

	ChannelCPU0 *chan Message
	ChannelCPU1 *chan Message
	ChannelCPU2 *chan Message
	ChannelCPU3 *chan Message

	SharedMemory *Memory
	SharedTable *Table
}

//CacheControllerRead - method that returns the need value
func (controller *CacheController) CacheControllerRead(memCell int) string  {

	cacheCell := controller.PrivateCache[memCell%CellsInCache]
	cacheData, cacheState, cacheMem := cacheCell.CacheLineRead()

	
	return "" //TODO: modify return statement
}

