package main

import (
	"bufio"
	"bytes"
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
