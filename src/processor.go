package main

import (
	"math/rand"
)

//Processor - struct that contains the parameters needed to simulate random instructions
type Processor struct {
	Name string

	InstructionNumber     int
	PublicCacheController *chan Message

	Random *rand.Rand
}

func (processor *Processor) ExecuteInstruction() {
	for {
		if processor.ExecuteNextInstruction() { //If ready to execute next instruction
			switch instruction := processor.GenerateInstructions(); instruction {
			case "LDR":
				block := processor.Random.Intn(16) //Random linear distributed number between 0 and 15
				*processor.PublicCacheController <- Message{Value: "", Tag: block, CPU: processor.Name}
			case "STR":
				block := processor.Random.Intn(16) //Random linear distributed number between 0 and 15
				*processor.PublicCacheController <- Message{Value: processor.Name, Tag: block, CPU: processor.Name}
			default:
				break
			}
			<-*processor.PublicCacheController
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

	} else { //No clock has passed
		return false
	}
}

func (processor *Processor) GenerateInstructions() string {
	randomNumber := int(processor.Random.NormFloat64()*StdDev + Mean) //From 0 to 20 normally distributed
	switch {
	case randomNumber > 15: //Numbers bigger than 15, execute STR
		return "STR"
	case randomNumber < 5: //Numbers smaller than 5, execute LDR
		return "LDR"
	default: //Numbers between [7,14]
		return "default"
	}
}
