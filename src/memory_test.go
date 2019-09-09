package main

import (
	"sync"
	"testing"
	"gotest.tools/assert"
)

func TestMemory_MemoryRead(t *testing.T) {
	type fields struct {
		MemoryMap map[int]string
		Mux       sync.Mutex
	}
	testField := &fields{
		MemoryMap: map[int]string{0: "CPU0", 1: "CPU1", 2: "CPU2", 3: "CPU3",},
	}
	type args struct {
		memCell int
	}
	tests := []struct {
		name   string
		fields *fields
		args   args
		want   string
	}{
		{name: "CPU0", fields: testField, args: args{memCell: 0}, want: "CPU0",},
		{name: "CPU0", fields: testField, args: args{memCell: 1}, want: "CPU1",},
		{name: "CPU0", fields: testField, args: args{memCell: 2}, want: "CPU2",},
		{name: "CPU0", fields: testField, args: args{memCell: 3}, want: "CPU3",},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := Memory{
				MemoryMap: tt.fields.MemoryMap,
				Mux:       tt.fields.Mux,
			}
			if got := mem.MemoryRead(tt.args.memCell); got != tt.want {
				t.Errorf("Memory.MemoryRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_MemoryWrite(t *testing.T) {
	type fields struct {
		MemoryMap map[int]string
		Mux       sync.Mutex
	}
	testField := &fields{
		MemoryMap: map[int]string{0: "CPU0", 1: "CPU1", 2: "CPU2", 3: "CPU3",},
	}
	type args struct {
		memCell  int
		memValue string
	}
	tests := []struct {
		name   string
		fields *fields
		args   args
	}{
		{name: "CPU0", fields: testField, args: args{memCell: 0, memValue: "CPU1"},},
		{name: "CPU0", fields: testField, args: args{memCell: 1, memValue: "CPU2"},},
		{name: "CPU0", fields: testField, args: args{memCell: 2, memValue: "CPU3"},},
		{name: "CPU0", fields: testField, args: args{memCell: 3, memValue: "CPU0"},},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := Memory{
				MemoryMap: tt.fields.MemoryMap,
				Mux:       tt.fields.Mux,
			}
			mem.MemoryWrite(tt.args.memCell, tt.args.memValue)
		})
	}
	assert.Equal(t, testField.MemoryMap[0], "CPU1")
	assert.Equal(t, testField.MemoryMap[1], "CPU2")
	assert.Equal(t, testField.MemoryMap[2], "CPU3")
	assert.Equal(t, testField.MemoryMap[3], "CPU0")
}
