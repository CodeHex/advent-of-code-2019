package main

import (
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

			testProgram, err := newProgram(tc.input)
			require.NoError(t, err)

			err = testProgram.run()
			require.NoError(t, err)

			assert.Equal(t, tc.expOutput, testProgram.dumpMemory())
		})
	}
}
