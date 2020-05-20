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

func TestProgramsWithInputAndOutputUsingChannels(t *testing.T) {
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
			in := make(chan int)
			out := make(chan int)

			testComp, err := newChannelComputer(tc.inputCode, in, out)
			require.NoError(t, err)

			go func() {
				err = testComp.run()
				require.NoError(t, err)
			}()

			in <- tc.input
			outValue := <-out
			assert.Equal(t, tc.output, outValue)
		})
	}
}

func TestSeriesComputer(t *testing.T) {
	tt := map[string]struct {
		inputData string
		input     int
		phases    []int
		expOutput int
	}{
		"example 1": {inputData: "3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0", input: 0, phases: []int{4, 3, 2, 1, 0}, expOutput: 43210},
		"example 2": {inputData: "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0", input: 0, phases: []int{0, 1, 2, 3, 4}, expOutput: 54321},
		"example 3": {inputData: "3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0", input: 0, phases: []int{1, 0, 4, 3, 2}, expOutput: 65210},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			testComp, err := newSeriesComputer(tc.inputData, "A", "B", "C", "D", "E")
			require.NoError(t, err)

			testComp.runAsync()
			testComp.loadPhases(tc.phases)
			testComp.input(tc.input)
			output := testComp.output()
			testComp.waitForCompletion()
			assert.Equal(t, tc.expOutput, output)
		})
	}
}

func TestFeedbackComputer(t *testing.T) {
	tt := map[string]struct {
		inputData string
		input     int
		phases    []int
		expOutput int
	}{
		"example 1": {inputData: "3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5", input: 0, phases: []int{9, 8, 7, 6, 5}, expOutput: 139629729},
		"example 2": {inputData: "3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10", input: 0, phases: []int{9, 7, 8, 5, 6}, expOutput: 18216},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			testComp, err := newFeedbackComputer(tc.inputData, "A", "B", "C", "D", "E")
			require.NoError(t, err)

			testComp.runAsync()
			testComp.loadPhases(tc.phases)
			testComp.input(tc.input)
			testComp.waitForCompletion()

			output := testComp.output()
			assert.Equal(t, tc.expOutput, output)
		})
	}
}
