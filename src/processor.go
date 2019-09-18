package main

import (
	"math/rand"
	"strconv"
)

//Processor - struct that contains the parameters needed to simulate random instructions
type Processor struct {
	Name string

	InstructionNumber     int
	PublicCacheController *chan Message

	Random *rand.Rand
	GUIChannel *chan Message
}

//ExecuteInstruction - method that executes the next random instrucion once at least clock has passed
func (processor *Processor) ExecuteInstruction() {
	for {
		if processor.ExecuteNextInstruction() { //If ready to execute next instruction
			switch instruction := processor.GenerateInstructions(); instruction {
			case "LDR":
				block := processor.Random.Intn(16) //Random linear distributed number between 0 and 15
				msg := "Read at " + strconv.Itoa(block) + "\n"
				*processor.GUIChannel<- Message{Value:msg, CPU:processor.Name}

				*processor.PublicCacheController <- Message{Value: "", Tag: block, CPU: processor.Name}
				<-*processor.PublicCacheController
			case "STR":
				block := processor.Random.Intn(16) //Random linear distributed number between 0 and 15
				msg := "Write at " + strconv.Itoa(block) + "\n"
				*processor.GUIChannel<- Message{Value:msg, CPU:processor.Name}

				*processor.PublicCacheController <- Message{Value: processor.Name, Tag: block, CPU: processor.Name}
				<-*processor.PublicCacheController
			default:
				break
			}
		}
	}
}

//ExecuteNextInstruction - method that verifies 1 or more cycles has passed since the last instruction
func (processor *Processor) ExecuteNextInstruction() bool {
	last := processor.InstructionNumber
	current := InstructionCounter
	if (current - last) >= 1 { //One clock has passed, execute instruction
		processor.InstructionNumber = current
		return true

	}
	return false
}

//GenerateInstructions - method that generates a instruction between STR(write), LDR(read) or default(nothing) with a normal distribution
func (processor *Processor) GenerateInstructions() string {
	randomNumber := int(processor.Random.NormFloat64()*StdDev + Mean) //From 0 to 20 normally distributed
	switch {
	case randomNumber > 11: //Numbers bigger than 15, execute STR
		return "STR" //Write
	case randomNumber < 8: //Numbers smaller than 5, execute LDR
		return "LDR" //Read
	default: //Numbers between [7,14]
		return "default"
	}
}
