package main

//Message - struct used in the messages exchange channels between CPU
type Message struct {
	Tag int
	Value string
	CPU string
}