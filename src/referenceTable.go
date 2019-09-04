package main

import "sync"

//Table - struct that manages a reference Table of which cache has the most recent value
type Table struct {
	referenceTable map[int]string
	mux *sync.Mutex
	//TODO: insert sleep when access?
}

//TableWrite - method that writes a value in the table. Thread-safe
func (t *Table) TableWrite(tableCell int, value string)  {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.referenceTable[tableCell] = value
	return
}

//TableRead - method that reads a value in the table. Thread-safe
func (t *Table) TableRead(tableCell int) string {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.referenceTable[tableCell]
}

//TableClear - method that clears a value in the table. Thread-safe
func (t *Table) TableClear(tableCell int)  {
	t.TableWrite(tableCell, "")
	return
}

//CreateReferenceTable - function that creates an empty reference Table
func CreateReferenceTable() *Table  {
	tableMap := map[int]string{
		0: "",
		1: "",
		2: "",
		3: "",
		4: "",
		5: "",
		6: "",
		7: "",
		8: "",
		9: "",
		10: "",
		11: "",
		12: "",
		13: "",
		14: "",
		15: "",
	}
	tab := Table{
		referenceTable: tableMap,
		mux:            &sync.Mutex{},
	}
	return &tab
}