package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const maxMemoryOutput = 20

type program struct {
	memory     []int
	runPos     int
	terminated bool
}

// newProgram reads in the input data in the form of a single CSV string
func newProgram(inputData string) (*program, error) {
	memoryStr := strings.Split(strings.TrimSpace(inputData), ",")
	memoryData := make([]int, len(memoryStr))

	var err error
	for i, memValStr := range memoryStr {
		memoryData[i], err = strconv.Atoi(memValStr)
		if err != nil {
			return nil, err
		}
	}

	return &program{memory: memoryData, runPos: 0}, nil
}

// run executes the int code currently stored in the provided memory
func (p *program) run() error {
	fmt.Printf(logFormat, "BEGIN: ")

	for !p.terminated {
		p.print()
		fmt.Println()

		// Check current position is in memory
		if p.outOfBounds() {
			return errors.Errorf("memory out of bounds: pos %d", p.runPos)
		}

		// Read current operation
		op, err := readOp(p)
		if err != nil {
			return err
		}

		// Apply operation
		op.Apply(p)

		// Move to next operation
		op.Next(p)
	}
	return nil
}

// print will output the entire memory contents if the program is small ( < 20 addresses)
func (p *program) print() {
	// Only print out the current program if its small
	if len(p.memory) < maxMemoryOutput {
		data := make([]string, len(p.memory))
		for i, entry := range p.memory {
			data[i] = strconv.Itoa(entry)
			if i == p.runPos {
				data[i] = "*" + data[i]
			}
		}
		fmt.Printf("%v", data)
	}
}

// addrOutOfBounds detects if the provided pointer is out of bounds
func (p *program) addrOutOfBounds(addr int) bool {
	return addr < 0 || addr >= len(p.memory)
}

// outOfBounds detects if program is currently out of bounds
func (p *program) outOfBounds() bool {
	return p.addrOutOfBounds(p.runPos)
}

// containsOp detects if op is contained in current program
func (p *program) containsOp(opSize int) bool {
	return p.runPos+opSize < len(p.memory)
}

// read current memory value at run position
func (p *program) read() int {
	return p.readOffset(0)
}

// read memory value ar run position with offset
func (p *program) readOffset(offset int) int {
	return p.readAddr(p.runPos + offset)
}

// read memory value at absolute position
func (p *program) readAddr(addr int) int {
	return p.memory[addr]
}

// store memory value ar run position with offset
func (p *program) storeAtAddr(addr int, val int) {
	p.memory[addr] = val
}

// generate string representation of a memory address in the format of value [index]
func (p *program) stringAddr(addr int) string {
	return fmt.Sprintf("%d[%d]", p.readAddr(addr), addr)
}

// dumps memory out into a CSV
func (p *program) dumpMemory() string {
	dump := make([]string, len(p.memory))
	for i, memVal := range p.memory {
		dump[i] = strconv.Itoa(memVal)
	}
	return strings.Join(dump, ",")
}
