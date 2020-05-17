package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExampleWirePathsForManhattenDistance(t *testing.T) {
	tt := map[string]struct {
		inputWire1 string
		inputWire2 string
		expDist    int
	}{
		"q example": {inputWire1: "R8,U5,L5,D3", inputWire2: "U7,R6,D4,L4", expDist: 6},
		"example 1": {inputWire1: "R75,D30,R83,U83,L12,D49,R71,U7,L72", inputWire2: "U62,R66,U55,R34,D71,R55,D58,R83", expDist: 159},
		"example 2": {inputWire1: "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51", inputWire2: "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7", expDist: 135},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			wire1, err := parseWirePath(tc.inputWire1)
			require.NoError(t, err)

			wire2, err := parseWirePath(tc.inputWire2)
			require.NoError(t, err)

			wire1Map := generatePointMap(wire1)
			wire2Map := generatePointMap(wire2)
			crossWires := crossPoints(wire1Map, wire2Map)
			closestPoint := closestManhattenPoint(crossWires)

			assert.Equal(t, tc.expDist, closestPoint.ManhattenDist())
		})
	}
}

func TestExampleWirePathsForSteps(t *testing.T) {
	tt := map[string]struct {
		inputWire1 string
		inputWire2 string
		expSteps   int
	}{
		"q example": {inputWire1: "R8,U5,L5,D3", inputWire2: "U7,R6,D4,L4", expSteps: 30},
		"example 1": {inputWire1: "R75,D30,R83,U83,L12,D49,R71,U7,L72", inputWire2: "U62,R66,U55,R34,D71,R55,D58,R83", expSteps: 610},
		"example 2": {inputWire1: "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51", inputWire2: "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7", expSteps: 410},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			wire1, err := parseWirePath(tc.inputWire1)
			require.NoError(t, err)

			wire2, err := parseWirePath(tc.inputWire2)
			require.NoError(t, err)

			wire1Map := generatePointMap(wire1)
			wire2Map := generatePointMap(wire2)
			crossWires := crossPoints(wire1Map, wire2Map)
			_, steps := closestStepsPoint(crossWires)

			assert.Equal(t, tc.expSteps, steps)
		})
	}
}
