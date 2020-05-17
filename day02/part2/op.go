package main

import (
	"github.com/pkg/errors"
)

type opCode int

const logFormat = "%-40s"

const (
	OpCodeAdd      opCode = 1
	OpCodeMultiply opCode = 2
	OpCodeHalt     opCode = 99
)

type instruction interface {
	Next(c *computer)
	Apply(c *computer)
}

func readOp(c *computer) (instruction, error) {
	code := opCode(c.read())

	switch code {
	case OpCodeAdd:
		return newAddOp(c)
	case OpCodeMultiply:
		return newMultiplyOp(c)
	case OpCodeHalt:
		return newHaltOp(c)
	default:
		return nil, errors.Errorf("unrecognized op code %d", code)
	}
}

// ---- OPERATIONS ----

type basicOp struct {
	opcodePtr int
	opSize    int
}

func (b basicOp) Next(p *computer) {
	p.insPtr += b.opSize
}

// ---- Binary Op ----

type binaryOp struct {
	basicOp
	paramPtrs [3]int
}

func newBinaryOp(p *computer) (binaryOp, error) {
	const binaryOpSize = 4
	if !p.containsOp(binaryOpSize) {
		return binaryOp{}, errors.New("op is incomplete, not enough input data")
	}

	op := binaryOp{
		basicOp:   basicOp{opcodePtr: p.insPtr, opSize: binaryOpSize},
		paramPtrs: [3]int{p.readOffset(1), p.readOffset(2), p.readOffset(3)},
	}

	for i, ptr := range op.paramPtrs {
		if p.addrOutOfBounds(ptr) {
			return binaryOp{}, errors.Errorf("param %d address invalid, memory out of bounds (%d[%d])", i+1, ptr, p.insPtr+1+i)
		}
	}
	return op, nil
}

// ---- Add Op ----

type addOp struct{ binaryOp }

func newAddOp(c *computer) (addOp, error) {
	op, err := newBinaryOp(c)
	return addOp{binaryOp: op}, err
}

func (a addOp) Apply(c *computer) {
	c.storeAtAddr(a.paramPtrs[2], c.readAddr(a.paramPtrs[0])+c.readAddr(a.paramPtrs[1]))
}

// ---- Multiply Op ----

type multiplyOp struct{ binaryOp }

func newMultiplyOp(c *computer) (multiplyOp, error) {
	op, err := newBinaryOp(c)
	return multiplyOp{binaryOp: op}, err
}

func (m multiplyOp) Apply(c *computer) {
	c.storeAtAddr(m.paramPtrs[2], c.readAddr(m.paramPtrs[0])*c.readAddr(m.paramPtrs[1]))
}

// ---- Halt Op ----

type haltOp struct{ basicOp }

func newHaltOp(p *computer) (haltOp, error) {
	return haltOp{basicOp: basicOp{opcodePtr: p.insPtr, opSize: 1}}, nil
}

func (h haltOp) Apply(c *computer) {
	c.terminated = true
}
