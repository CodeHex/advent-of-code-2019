package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type computer struct {
	label         string
	inScanner     *bufio.Scanner
	inChan        <-chan int64
	outChan       chan<- int64
	memory        []int64
	insPtr        int
	relativeBase  int
	terminated    bool
	outputs       []int64
	disableLog    bool
	disableOutLog bool
}

const PostionMode = 0
const AbsoluteMode = 1
const RelativeMode = 2

const memorySize = 100_000

// newComputer reads in the input data in the form of a single CSV string and uses a scanner to read input
func newComputer(inputData string, scanner *bufio.Scanner) (*computer, error) {
	memory, err := parseMemoryInput(inputData)
	if err != nil {
		return nil, err
	}
	return &computer{memory: memory, insPtr: 0, inScanner: scanner}, nil
}

// newChannelComputer reads in the input data in the form of a single CSV string and uses channels to read input and send output
func newChannelComputer(inputData string, in <-chan int64, out chan<- int64) (*computer, error) {
	memory, err := parseMemoryInput(inputData)
	if err != nil {
		return nil, err
	}
	return &computer{memory: memory, insPtr: 0, inChan: in, outChan: out}, nil
}

func parseMemoryInput(inputData string) ([]int64, error) {
	memoryStr := strings.Split(strings.TrimSpace(inputData), ",")
	memoryData := make([]int64, memorySize)

	var err error
	for i, memValStr := range memoryStr {
		memoryData[i], err = strconv.ParseInt(memValStr, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return memoryData, nil
}

// makes an independent copy of the computer
func (c *computer) clone() *computer {
	newMem := make([]int64, len(c.memory))
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
	if c.outChan != nil {
		close(c.outChan)
	}
	return nil
}

// addrOutOfBounds detects if the provided pointer is out of bounds
func (c *computer) addrOutOfBounds(addr int) bool {
	return addr < 0 || addr >= len(c.memory)
}

// read current memory value at run position and advance
func (c *computer) read() (int64, error) {
	if c.addrOutOfBounds(c.insPtr) {
		return -1, errors.New("memory out of bounds")
	}
	result := c.memory[c.insPtr]
	c.insPtr++
	return result, nil
}

// read memory value based on mode
func (c *computer) readMode(p param) int64 {
	switch p.mode {
	case PostionMode:
		return c.memory[p.val]
	case AbsoluteMode:
		return p.val
	case RelativeMode:
		return c.memory[int(p.val)+c.relativeBase]
	default:
		return c.memory[p.val]
	}

}

// store memory value ar run position with offset
func (c *computer) storeAtAddr(p param, val int64) {
	if p.mode == RelativeMode {
		c.memory[int(p.val)+c.relativeBase] = val
	} else {
		c.memory[p.val] = val
	}
}

// dumps memory out into a CSV
func (c *computer) dumpMemory() string {
	dump := make([]string, len(c.memory))
	for i, memVal := range c.memory {
		dump[i] = strconv.FormatInt(memVal, 10)
	}
	return strings.Join(dump, ",")
}

func (c *computer) logOutf(format string, args ...interface{}) {
	if c.disableOutLog {
		return
	}

	if c.label != "" {
		format = fmt.Sprintf("[%s] %s", c.label, format)
	}
	fmt.Printf(format, args...)
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
