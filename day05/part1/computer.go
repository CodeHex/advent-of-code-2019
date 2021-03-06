package main

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type computer struct {
	scanner    *bufio.Scanner
	memory     []int
	insPtr     int
	terminated bool
}

// newComputer reads in the input data in the form of a single CSV string
func newComputer(inputData string, scanner *bufio.Scanner) (*computer, error) {
	memoryStr := strings.Split(strings.TrimSpace(inputData), ",")
	memoryData := make([]int, len(memoryStr))

	var err error
	for i, memValStr := range memoryStr {
		memoryData[i], err = strconv.Atoi(memValStr)
		if err != nil {
			return nil, err
		}
	}

	return &computer{memory: memoryData, insPtr: 0, scanner: scanner}, nil
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
	}
	return nil
}

// input reads the value from the scanner
func (c *computer) input() (int, error) {
	if c.scanner == nil {
		return -1, errors.New("computer has no scanner")
	}
	c.scanner.Scan()
	return strconv.Atoi(c.scanner.Text())
}

// addrOutOfBounds detects if the provided pointer is out of bounds
func (c *computer) addrOutOfBounds(addr int) bool {
	return addr < 0 || addr >= len(c.memory)
}

// outOfBounds detects if program is currently out of bounds
func (c *computer) outOfBounds() bool {
	return c.addrOutOfBounds(c.insPtr)
}

// read current memory value at run position and advance
func (c *computer) read() (int, error) {
	if c.outOfBounds() {
		return -1, errors.New("memory out of bounds")
	}
	result := c.readMode(c.insPtr, false)
	c.insPtr++
	return result, nil
}

// read memory value based on mode
func (c *computer) readMode(val int, isImmediate bool) int {
	if isImmediate {
		return val
	}
	return c.memory[val]
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
