package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type opCode int

const logFormat = "%-40s"

const (
	OpCodeAdd      opCode = 1
	OpCodeMultiply opCode = 2
	OpCodeHalt     opCode = 99
)

type op interface {
	Next(p *program)
	Apply(p *program)
}

func readOp(p *program) (op, error) {
	code := opCode(p.read())

	switch code {
	case OpCodeAdd:
		return newAddOp(p)
	case OpCodeMultiply:
		return newMultiplyOp(p)
	case OpCodeHalt:
		return newHaltOp(p)
	default:
		return nil, errors.Errorf("unrecognized op code %d", code)
	}
}

// ---- OPERATIONS ----

type basicOp struct {
	opcodePos int
	opSize    int
}

func (b basicOp) Next(p *program) {
	p.runPos += b.opSize
}

// ---- Binary Op ----

type binaryOp struct {
	basicOp
	arg1Ptr   int
	arg2Ptr   int
	resultPtr int
}

func newBinaryOp(p *program) (binaryOp, error) {
	const binaryOpSize = 4
	if !p.containsOp(binaryOpSize) {
		return binaryOp{}, errors.New("op is incomplete, not enough input data")
	}

	op := binaryOp{
		basicOp:   basicOp{opcodePos: p.runPos, opSize: binaryOpSize},
		arg1Ptr:   p.readOffset(1),
		arg2Ptr:   p.readOffset(2),
		resultPtr: p.readOffset(3),
	}

	if p.addrOutOfBounds(op.arg1Ptr) {
		return binaryOp{}, errors.Errorf("first arg address invalid, memory out of bounds (%d[%d])", op.arg1Ptr, p.runPos+1)
	}

	if p.addrOutOfBounds(op.arg2Ptr) {
		return binaryOp{}, errors.Errorf("second arg address invalid, memory out of bounds (%d[%d])", op.arg2Ptr, p.runPos+2)
	}

	if p.addrOutOfBounds(op.resultPtr) {
		return binaryOp{}, errors.Errorf("result address is invalid, memory out of bounds (%d[%d])", op.resultPtr, p.runPos+3)
	}
	return op, nil
}

// ---- Add Op ----

type addOp struct{ binaryOp }

func newAddOp(p *program) (addOp, error) {
	op, err := newBinaryOp(p)
	return addOp{binaryOp: op}, err
}

func (a addOp) Apply(p *program) {
	beforeArg1, beforeArg2 := p.stringAddr(a.arg1Ptr), p.stringAddr(a.arg2Ptr)
	result := p.readAddr(a.arg1Ptr) + p.readAddr(a.arg2Ptr)
	p.storeAtAddr(a.resultPtr, result)
	log := fmt.Sprintf("ADD  : %s + %s = %s", beforeArg1, beforeArg2, p.stringAddr(a.resultPtr))
	fmt.Printf(logFormat, log)
}

// ---- Multiply Op ----

type multiplyOp struct{ binaryOp }

func newMultiplyOp(p *program) (multiplyOp, error) {
	op, err := newBinaryOp(p)
	return multiplyOp{binaryOp: op}, err
}

func (m multiplyOp) Apply(p *program) {
	beforeArg1, beforeArg2 := p.stringAddr(m.arg1Ptr), p.stringAddr(m.arg2Ptr)
	result := p.readAddr(m.arg1Ptr) * p.readAddr(m.arg2Ptr)
	p.storeAtAddr(m.resultPtr, result)
	log := fmt.Sprintf("MULT : %s * %s = %s", beforeArg1, beforeArg2, p.stringAddr(m.resultPtr))
	fmt.Printf(logFormat, log)
}

// ---- Halt Op ----

type haltOp struct{ basicOp }

func newHaltOp(p *program) (haltOp, error) {
	return haltOp{basicOp: basicOp{opcodePos: p.runPos, opSize: 1}}, nil
}

func (h haltOp) Apply(p *program) {
	p.terminated = true
	log := fmt.Sprintf("HALT : at [%d]", h.opcodePos)
	fmt.Printf(logFormat+"\n\n", log)
}
