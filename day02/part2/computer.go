package main

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const maxMemoryOutput = 20

type computer struct {
	memory     []int
	insPtr     int
	terminated bool
}

// newComputer reads in the input data in the form of a single CSV string
func newComputer(inputData string) (*computer, error) {
	memoryStr := strings.Split(strings.TrimSpace(inputData), ",")
	memoryData := make([]int, len(memoryStr))

	var err error
	for i, memValStr := range memoryStr {
		memoryData[i], err = strconv.Atoi(memValStr)
		if err != nil {
			return nil, err
		}
	}

	return &computer{memory: memoryData, insPtr: 0}, nil
}

// makes an independent copy of the computer
func (c *computer) clone() *computer {
	newMem := make([]int, len(c.memory))
	copy(newMem, c.memory)
	return &computer{
		memory:     newMem,
		insPtr:     c.insPtr,
		terminated: c.terminated,
	}
}

// run executes the int code currently stored in the provided memory
func (c *computer) run() error {
	for !c.terminated {
		// Check current position is in memory
		if c.outOfBounds() {
			return errors.Errorf("memory out of bounds: pos %d", c.insPtr)
		}

		// Read current operation
		op, err := readOp(c)
		if err != nil {
			return err
		}

		// Apply operation
		op.Apply(c)

		// Move to next operation
		op.Next(c)
	}
	return nil
}

// addrOutOfBounds detects if the provided pointer is out of bounds
func (c *computer) addrOutOfBounds(addr int) bool {
	return addr < 0 || addr >= len(c.memory)
}

// outOfBounds detects if program is currently out of bounds
func (c *computer) outOfBounds() bool {
	return c.addrOutOfBounds(c.insPtr)
}

// containsOp detects if op is contained in current program
func (c *computer) containsOp(opSize int) bool {
	return c.insPtr+opSize < len(c.memory)
}

// read current memory value at run position
func (c *computer) read() int {
	return c.readOffset(0)
}

// read memory value ar run position with offset
func (c *computer) readOffset(offset int) int {
	return c.readAddr(c.insPtr + offset)
}

// read memory value at absolute position
func (c *computer) readAddr(addr int) int {
	return c.memory[addr]
}

// store memory value ar run position with offset
func (c *computer) storeAtAddr(addr int, val int) {
	c.memory[addr] = val
}

// dumps memory out into a CSV
func (c *computer) dumpMemory() string {
	dump := make([]string, len(c.memory))
	for i, memVal := range c.memory {
		dump[i] = strconv.Itoa(memVal)
	}
	return strings.Join(dump, ",")
}
