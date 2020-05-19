package main

import (
	"bufio"
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExamplePrograms(t *testing.T) {
	tt := map[string]struct {
		input     string
		expOutput string
	}{
		"q example": {input: "1,9,10,3,2,3,11,0,99,30,40,50", expOutput: "3500,9,10,70,2,3,11,0,99,30,40,50"},
		"example 1": {input: "1,0,0,0,99", expOutput: "2,0,0,0,99"},
		"example 2": {input: "2,3,0,3,99", expOutput: "2,3,0,6,99"},
		"example 3": {input: "2,4,4,5,99,0", expOutput: "2,4,4,5,99,9801"},
		"example 4": {input: "1,1,1,4,99,5,6,0,99", expOutput: "30,1,1,4,2,5,6,0,99"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			testComp, err := newComputer(tc.input, nil)
			require.NoError(t, err)

			err = testComp.run()
			require.NoError(t, err)

			assert.Equal(t, tc.expOutput, testComp.dumpMemory())
		})
	}
}

func TestProgramsWithInput(t *testing.T) {
	tt := map[string]struct {
		inputCode string
		inputVal  string
		expOutput string
	}{
		"simple input 1": {inputCode: "3,3,99,0", inputVal: "1", expOutput: "3,3,99,1"},
		"simple input 2": {inputCode: "3,3,99,0", inputVal: "23", expOutput: "3,3,99,23"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			buffer := bytes.NewBufferString(tc.inputVal)
			scanner := bufio.NewScanner(buffer)
			testComp, err := newComputer(tc.inputCode, scanner)
			require.NoError(t, err)

			err = testComp.run()
			require.NoError(t, err)

			assert.Equal(t, tc.expOutput, testComp.dumpMemory())
		})
	}
}

func TestImmediatePrograms(t *testing.T) {
	tt := map[string]struct {
		input     string
		expOutput string
	}{
		"simple immediate addition": {input: "1001,0,200,0,99", expOutput: "1201,0,200,0,99"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			testComp, err := newComputer(tc.input, nil)
			require.NoError(t, err)

			err = testComp.run()
			require.NoError(t, err)

			assert.Equal(t, tc.expOutput, testComp.dumpMemory())
		})
	}
}

func TestProgramsWithInputAndOutput(t *testing.T) {
	tt := map[string]struct {
		inputCode string
		input     int
		output    int
	}{
		"simple input output 1": {inputCode: "3,0,4,0,99", input: 254, output: 254},
		"simple input output 2": {inputCode: "3,0,4,0,99", input: 82, output: 82},

		"position equal false": {inputCode: "3,9,8,9,10,9,4,9,99,-1,8", input: 5, output: 0},
		"position equal true":  {inputCode: "3,9,8,9,10,9,4,9,99,-1,8", input: 8, output: 1},

		"position less than false (greater)": {inputCode: "3,9,7,9,10,9,4,9,99,-1,8", input: 23, output: 0},
		"position less than false (equal)":   {inputCode: "3,9,7,9,10,9,4,9,99,-1,8", input: 8, output: 0},
		"position less than true (less)":     {inputCode: "3,9,7,9,10,9,4,9,99,-1,8", input: 7, output: 1},

		"absolute equal false": {inputCode: "3,3,1108,-1,8,3,4,3,99", input: 5, output: 0},
		"absolute equal true":  {inputCode: "3,3,1108,-1,8,3,4,3,99", input: 8, output: 1},

		"absolute less than false (greater)": {inputCode: "3,3,1107,-1,8,3,4,3,99", input: 23, output: 0},
		"absolute less than false (equal)":   {inputCode: "3,3,1107,-1,8,3,4,3,99", input: 8, output: 0},
		"absolute less than true (less)":     {inputCode: "3,3,1107,-1,8,3,4,3,99", input: 7, output: 1},

		"position jump (0)":  {inputCode: "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", input: 0, output: 0},
		"position jump (1)":  {inputCode: "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", input: 1, output: 1},
		"position jump (>1)": {inputCode: "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", input: 999, output: 1},

		"immediate jump (0)":  {inputCode: "3,3,1105,-1,9,1101,0,0,12,4,12,99,1", input: 0, output: 0},
		"immediate jump (1)":  {inputCode: "3,3,1105,-1,9,1101,0,0,12,4,12,99,1", input: 1, output: 1},
		"immediate jump (>1)": {inputCode: "3,3,1105,-1,9,1101,0,0,12,4,12,99,1", input: 999, output: 1},

		"less equal greater test (less)": {
			inputCode: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			input:     3,
			output:    999,
		},
		"less equal greater test (equal)": {
			inputCode: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			input:     8,
			output:    1000,
		},
		"less equal greater test (greater)": {
			inputCode: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			input:     98,
			output:    1001,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			buffer := bytes.NewBufferString(strconv.Itoa(tc.input))
			scanner := bufio.NewScanner(buffer)
			testComp, err := newComputer(tc.inputCode, scanner)
			require.NoError(t, err)

			err = testComp.run()
			require.NoError(t, err)

			assert.Len(t, testComp.outputs, 1)
			assert.Equal(t, tc.output, testComp.outputs[0])
		})
	}
}
