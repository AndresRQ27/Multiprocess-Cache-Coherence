package main

//Broadcast - struct used in the broadcast channel between caches
type Broadcast struct {
	memCell int
	State   string
}

//Message - struct used in the messages exchange channels between CPU
type Message struct {
	memCell int
	Value string
}
