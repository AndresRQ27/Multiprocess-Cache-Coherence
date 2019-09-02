package main

import (
	"sync"
	"testing"
)

func TestMemory_MemoryRead(t *testing.T) {
	type fields struct {
		MemoryMap map[int]string
		Mux       *sync.Mutex
	}
	type args struct {
		memCell int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
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
		Mux       *sync.Mutex
	}
	type args struct {
		memCell  int
		memValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
}
