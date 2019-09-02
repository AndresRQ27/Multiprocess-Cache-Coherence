package main

import "fmt"

//SnoopProtocol - struct with the different states of the cache
type SnoopProtocol struct {
	M, S, I bool
}

//Interface Stringer - prints the values of the SnoopProtocol
func (snoop SnoopProtocol) String() string {
	if snoop.M {
		return fmt.Sprintf("Modified")
	} else if snoop.S {
		return fmt.Sprintf("Shared")
	} else if snoop.I {
		return fmt.Sprintf("Invalid")
	} else {
		return fmt.Sprintf("Something wrong happened in the SnoopProtocol.")
	}
}

//NewSnoopProtocol - Constructor of SnoopProtocol that initialize the struct in Invalid
func NewSnoopProtocol() *SnoopProtocol {
	return &SnoopProtocol{M:false, S:false, I:true}
}