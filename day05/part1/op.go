package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type opCode int

const (
	OpCodeAdd      opCode = 1
	OpCodeMultiply opCode = 2
	OpCodeInput    opCode = 3
	OpCodeOutput   opCode = 4
	OpCodeHalt     opCode = 99
)

type instruction interface {
	Apply(c *computer)
}

func readOp(c *computer) (instruction, error) {
	opValue, err := c.read()
	if err != nil {
		return nil, errors.Wrap(err, "unable to read op code")
	}

	code := opValue % 100
	modeData := opValue / 100

	var modes []bool
	digits := fmt.Sprintf("%d", modeData)
	for i := len(digits) - 1; i >= 0; i-- {
		flag := false
		if digits[i] == '1' {
			flag = true
		}
		modes = append(modes, flag)
	}

	switch opCode(code) {
	case OpCodeAdd:
		return newAddOp(c, modes)
	case OpCodeMultiply:
		return newMultiplyOp(c, modes)
	case OpCodeInput:
		return newInputOp(c, modes)
	case OpCodeOutput:
		return newOutputOp(c, modes)
	case OpCodeHalt:
		return newHaltOp(c, modes)
	default:
		return nil, errors.Errorf("unrecognized op code %d", opValue)
	}
}

// ---- OPERATIONS ----

type basicOp struct {
	params []int
	modes  []bool
}

func newBasicOp(c *computer, insSize int, inputModes []bool) (basicOp, error) {
	op := basicOp{modes: make([]bool, insSize)}
	for m, mode := range inputModes {
		op.modes[m] = mode
	}

	for i := 1; i < insSize; i++ {
		paramPtr, err := c.read()
		if err != nil {
			return basicOp{}, errors.Wrapf(err, "unable to read param %d", i)
		}
		if !op.modes[i-1] && c.addrOutOfBounds(paramPtr) {
			return basicOp{}, errors.Errorf("param %d references memory out of bounds, param value %d", i, paramPtr)
		}
		op.params = append(op.params, paramPtr)
	}
	return op, nil
}

// ---- Add Op ----

type addOp struct{ basicOp }

func newAddOp(c *computer, modes []bool) (addOp, error) {
	if len(modes) > 2 {
		return addOp{}, errors.New("too many modes provided to add op")
	}
	op, err := newBasicOp(c, 4, modes)
	return addOp{basicOp: op}, err
}

func (a addOp) Apply(c *computer) {
	arg1 := c.readMode(a.params[0], a.modes[0])
	arg2 := c.readMode(a.params[1], a.modes[1])
	c.storeAtAddr(a.params[2], arg1+arg2)
}

// ---- Multiply Op ----

type multiplyOp struct{ basicOp }

func newMultiplyOp(c *computer, modes []bool) (multiplyOp, error) {
	if len(modes) > 2 {
		return multiplyOp{}, errors.New("too many modes provided to multiply op")
	}
	op, err := newBasicOp(c, 4, modes)
	return multiplyOp{basicOp: op}, err
}

func (m multiplyOp) Apply(c *computer) {
	arg1 := c.readMode(m.params[0], m.modes[0])
	arg2 := c.readMode(m.params[1], m.modes[1])
	c.storeAtAddr(m.params[2], arg1*arg2)
}

// ---- Halt Op ----

type haltOp struct{ basicOp }

func newHaltOp(c *computer, modes []bool) (haltOp, error) {
	if len(modes) > 1 || modes[0] {
		return haltOp{}, errors.New("halt op doesn't support modes")
	}
	op, err := newBasicOp(c, 1, nil)
	return haltOp{basicOp: op}, err
}

func (h haltOp) Apply(c *computer) {
	c.terminated = true
}

// ---- Input Op ----

type inputOp struct{ basicOp }

func newInputOp(c *computer, modes []bool) (inputOp, error) {
	if len(modes) > 1 || modes[0] {
		return inputOp{}, errors.New("input op doesn't support modes")
	}
	op, err := newBasicOp(c, 2, nil)
	return inputOp{basicOp: op}, err
}

func (i inputOp) Apply(c *computer) {
	fmt.Printf("ENTER INPUT: ")
	input, err := c.input()
	if err != nil {
		fmt.Printf("unable to parse input, %s", err.Error())
		// terminate app
		c.terminated = true
		return
	}

	c.storeAtAddr(i.params[0], input)
}

// ---- Output Op ----

type outputOp struct{ basicOp }

func newOutputOp(c *computer, modes []bool) (outputOp, error) {
	if len(modes) > 1 {
		return outputOp{}, errors.New("too many modes provided to output op")
	}
	op, err := newBasicOp(c, 2, modes)
	return outputOp{basicOp: op}, err
}

func (o outputOp) Apply(c *computer) {
	fmt.Printf("OUTPUT: %d\n", c.readMode(o.params[0], o.modes[0]))
}
