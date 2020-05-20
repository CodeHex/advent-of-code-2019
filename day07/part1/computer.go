package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type computer struct {
	label      string
	inScanner  *bufio.Scanner
	inChan     <-chan int
	outChan    chan<- int
	memory     []int
	insPtr     int
	terminated bool
	outputs    []int
	disableLog bool
}

// newComputer reads in the input data in the form of a single CSV string
func newComputer(inputData string, scanner *bufio.Scanner) (*computer, error) {
	memory, err := parseMemoryInput(inputData)
	if err != nil {
		return nil, err
	}
	return &computer{memory: memory, insPtr: 0, inScanner: scanner}, nil
}

// newComputer reads in the input data in the form of a single CSV string
func newChannelComputer(inputData string, in <-chan int, out chan<- int) (*computer, error) {
	memory, err := parseMemoryInput(inputData)
	if err != nil {
		return nil, err
	}
	return &computer{memory: memory, insPtr: 0, inChan: in, outChan: out}, nil
}

func parseMemoryInput(inputData string) ([]int, error) {
	memoryStr := strings.Split(strings.TrimSpace(inputData), ",")
	memoryData := make([]int, len(memoryStr))

	var err error
	for i, memValStr := range memoryStr {
		memoryData[i], err = strconv.Atoi(memValStr)
		if err != nil {
			return nil, err
		}
	}
	return memoryData, nil
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
		if c.addrOutOfBounds(c.insPtr) {
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
	if c.inScanner == nil {
		return -1, errors.New("computer has no scanner")
	}
	c.inScanner.Scan()
	return strconv.Atoi(c.inScanner.Text())
}

// addrOutOfBounds detects if the provided pointer is out of bounds
func (c *computer) addrOutOfBounds(addr int) bool {
	return addr < 0 || addr >= len(c.memory)
}

// read current memory value at run position and advance
func (c *computer) read() (int, error) {
	if c.addrOutOfBounds(c.insPtr) {
		return -1, errors.New("memory out of bounds")
	}
	result := c.memory[c.insPtr]
	c.insPtr++
	return result, nil
}

// read memory value based on mode
func (c *computer) readMode(p param) int {
	if p.mode {
		return p.val
	}
	return c.memory[p.val]
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

func (c *computer) logf(format string, args ...interface{}) {
	if c.disableLog {
		return
	}

	if c.label != "" {
		format = fmt.Sprintf("[%s] %s", c.label, format)
	}
	fmt.Printf(format, args...)
}
