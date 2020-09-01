package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// generate phases
	signal := parseInput(input)
	phases := generatePhases(signal)

	for i := 1; i <= 100; i++ {
		signal = applyPhases(signal, phases)
	}

	fmt.Printf("after 100 phases - %s\n", signalString(signal[0:8]))
}

func parseInput(input string) []int {
	digits := make([]int, len(input))
	var err error
	for i, char := range input {
		digits[i], err = strconv.Atoi(string(char))
		if err != nil {
			panic(errors.Wrap(err, "unable to parse input"))
		}
	}
	return digits
}

var basePhase = []int{0, 1, 0, -1}

func generatePhases(signal []int) [][]int {
	size := len(signal)
	phases := make([][]int, size)
	for signalIndex := range signal {
		phase := make([]int, size+1)
		for i := 0; i < size+1; i++ {
			phase[i] = basePhase[(i/(signalIndex+1))%4]
		}
		phase = phase[1:]
		phases[signalIndex] = phase
	}
	return phases
}

func applyPhases(signal []int, phases [][]int) []int {
	result := make([]int, len(signal))

	for i, phase := range phases {
		sum := 0
		for j, phaseVal := range phase {
			sum += signal[j] * phaseVal
		}
		result[i] = sum % 10
		if result[i] < 0 {
			result[i] = result[i] * -1
		}
	}
	return result
}

func signalString(signal []int) string {
	out := ""
	for _, digit := range signal {
		out += strconv.Itoa(digit)
	}
	return out
}
