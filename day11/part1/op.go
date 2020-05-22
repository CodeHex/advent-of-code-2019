package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"

type opCode int

const (
	OpCodeAdd           opCode = 1
	OpCodeMultiply      opCode = 2
	OpCodeInput         opCode = 3
	OpCodeOutput        opCode = 4
	OpCodeJumpIfTrue    opCode = 5
	OpCodeJumpIfFalse   opCode = 6
	OpCodeLessThan      opCode = 7
	OpCodeEqual         opCode = 8
	OpCodeShiftRelative opCode = 9
	OpCodeHalt          opCode = 99
)

type instruction interface {
	Apply(c *computer)
}

func extractModeData(opcode int64) []int {
	modeData := opcode / 100
	var modes []int
	digits := fmt.Sprintf("%d", modeData)
	for i := len(digits) - 1; i >= 0; i-- {
		mode, _ := strconv.Atoi(string(digits[i]))
		modes = append(modes, mode)
	}
	return modes
}

func readOp(c *computer) (instruction, error) {
	opValue, err := c.read()
	if err != nil {
		return nil, errors.Wrap(err, "unable to read op code")
	}

	code := opValue % 100
	modes := extractModeData(opValue)

	switch opCode(code) {
	case OpCodeAdd:
		return newAddOp(c, modes)
	case OpCodeMultiply:
		return newMultiplyOp(c, modes)
	case OpCodeInput:
		return newInputOp(c, modes)
	case OpCodeOutput:
		return newOutputOp(c, modes)
	case OpCodeJumpIfTrue:
		return newJumpOp(c, modes, true)
	case OpCodeJumpIfFalse:
		return newJumpOp(c, modes, false)
	case OpCodeLessThan:
		return newLessThanOp(c, modes)
	case OpCodeEqual:
		return newEqualOp(c, modes)
	case OpCodeShiftRelative:
		return newShiftRelativeOp(c, modes)
	case OpCodeHalt:
		return newHaltOp(c)
	default:
		return nil, errors.Errorf("unrecognized op code %d", opValue)
	}
}

// ---- OPERATIONS ----

type param struct {
	val  int64
	mode int
}

type basicOp struct {
	params []param
}

func newBasicOp(c *computer, paramSize int, inputModes []int) (basicOp, error) {
	completeModes := make([]int, paramSize)
	for m, mode := range inputModes {
		completeModes[m] = mode
	}

	params := make([]param, paramSize)
	for i := 0; i < paramSize; i++ {
		p, err := c.read()
		if err != nil {
			return basicOp{}, errors.Wrapf(err, "unable to read param %d", i)
		}
		if completeModes[i] == PostionMode && c.addrOutOfBounds(int(p)) {
			return basicOp{}, errors.Errorf("param %d references memory out of bounds, param value %d", i, p)
		}
		if completeModes[i] == RelativeMode && c.addrOutOfBounds(int(p)+c.relativeBase) {
			return basicOp{}, errors.Errorf("param %d references memory out of bounds,  param value %d rel value %d", i, p, c.relativeBase)
		}
		params[i] = param{val: p, mode: completeModes[i]}
	}
	return basicOp{params: params}, nil
}

// ---- binary op ----

type binaryOp struct {
	basicOp
	operator  func(x, y int64) int64
	logFormat string
}

func newBinaryOp(logFormat string, c *computer, modes []int, operator func(x, y int64) int64) (binaryOp, error) {
	op, err := newBasicOp(c, 3, modes)
	return binaryOp{basicOp: op, operator: operator, logFormat: logFormat}, err
}

func (b binaryOp) Apply(c *computer) {
	arg1 := c.readMode(b.params[0])
	arg2 := c.readMode(b.params[1])
	result := b.operator(arg1, arg2)
	c.storeAtAddr(b.params[2], result)
	c.logf(b.logFormat+"\n", arg1, arg2, result)
}

// ---- Add Op ----

type addOp struct{ binaryOp }

func add(x, y int64) int64 { return x + y }

func newAddOp(c *computer, modes []int) (addOp, error) {
	op, err := newBinaryOp("ADD : %d + %d = %d", c, modes, add)
	return addOp{binaryOp: op}, err
}

// ---- Multiply Op ----

type multiplyOp struct{ binaryOp }

func multiply(x, y int64) int64 { return x * y }

func newMultiplyOp(c *computer, modes []int) (multiplyOp, error) {
	op, err := newBinaryOp("MULT: %d * %d = %d", c, modes, multiply)
	return multiplyOp{binaryOp: op}, err
}

// ---- Halt Op ----

type haltOp struct{ basicOp }

func newHaltOp(c *computer) (haltOp, error) {
	op, err := newBasicOp(c, 0, nil)
	return haltOp{basicOp: op}, err
}

func (h haltOp) Apply(c *computer) {
	c.terminated = true
	c.logOutf("%sHALT%s\n", Red, Reset)
}

// ---- Input Op ----

type inputOp struct{ basicOp }

func newInputOp(c *computer, modes []int) (inputOp, error) {
	op, err := newBasicOp(c, 1, modes)
	return inputOp{basicOp: op}, err
}

func (i inputOp) Apply(c *computer) {
	var input int64
	switch {
	case c.inScanner != nil:
		c.logOutf("%sENTER INPUT: %s", Blue, Reset)
		c.inScanner.Scan()
		scannedInput, err := strconv.ParseInt(c.inScanner.Text(), 10, 64)
		if err != nil {
			fmt.Printf("unable to parse input, %s", err.Error())
			os.Exit(1)
		}
		input = scannedInput
	case c.inChan != nil:
		c.logOutf("%sIN : %d%s\n", Green, input, Reset)
		input = <-c.inChan
		c.logOutf("%sIN : %d%s\n", Green, input, Reset)
	default:
		fmt.Printf("computer has not input method")
		os.Exit(1)
	}
	c.storeAtAddr(i.params[0], input)
}

// ---- Output Op ----

type outputOp struct{ basicOp }

func newOutputOp(c *computer, modes []int) (outputOp, error) {
	op, err := newBasicOp(c, 1, modes)
	return outputOp{basicOp: op}, err
}

func (o outputOp) Apply(c *computer) {
	out := c.readMode(o.params[0])
	c.outputs = append(c.outputs, out)
	if c.outChan != nil {
		c.outChan <- out
	}
	c.logOutf("%sOUT : %d%s\n", Green, out, Reset)
}

// ---- Jump Op ----

type jumpOp struct {
	basicOp
	jumpWhen bool
}

func newJumpOp(c *computer, modes []int, jumpWhen bool) (jumpOp, error) {
	op, err := newBasicOp(c, 2, modes)
	return jumpOp{basicOp: op, jumpWhen: jumpWhen}, err
}

func (j jumpOp) Apply(c *computer) {
	arg := c.readMode(j.params[0])
	argIsNonZero := arg != 0
	if argIsNonZero == j.jumpWhen {
		c.insPtr = int(c.readMode(j.params[1]))
		c.logf("JUMP: to %d (arg %d)\n", c.insPtr, arg)
		return
	}
	c.logf("CONT: (arg %d)\n", arg)
}

// ---- Less than Op ----

type lessThanOp struct{ binaryOp }

func lessThan(x, y int64) int64 {
	if x < y {
		return 1
	}
	return 0
}

func newLessThanOp(c *computer, modes []int) (lessThanOp, error) {
	op, err := newBinaryOp("LESS: %d < %d (%d)", c, modes, lessThan)
	return lessThanOp{binaryOp: op}, err
}

// ---- Greater than Op ----

type equalOp struct{ binaryOp }

func equal(x, y int64) int64 {
	if x == y {
		return 1
	}
	return 0
}

func newEqualOp(c *computer, modes []int) (lessThanOp, error) {
	op, err := newBinaryOp("EQ  : %d == %d (%d)", c, modes, equal)
	return lessThanOp{binaryOp: op}, err
}

// ---- Shift Relative Op ----

type shiftRelativeOp struct{ basicOp }

func newShiftRelativeOp(c *computer, modes []int) (shiftRelativeOp, error) {
	op, err := newBasicOp(c, 1, modes)
	return shiftRelativeOp{basicOp: op}, err
}

func (s shiftRelativeOp) Apply(c *computer) {
	shift := int(c.readMode(s.params[0]))
	c.logf("REL : %d->%d\n", c.relativeBase, c.relativeBase+shift)
	c.relativeBase += shift
}
